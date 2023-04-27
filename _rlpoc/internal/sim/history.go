package sim

import (
	"errors"
	"fmt"
	"sort"

	"github.com/PieterD/rlpoc/lib/entity"

	"github.com/google/btree"
)

type Action interface {
	Name() string
	Apply(state *State) error
	Revert(state *State) error
	entity.Value
}

type Event struct {
	Depth    uint64
	Action   Action
	Previous EventId
}

type EventId uint64

type History struct {
	store          *entity.Store
	state          *State
	currentEventId EventId
	events         []Event
	reverse        *btree.BTree // *reverseEvent
	apply          uint64
	branches       map[string]EventId
}

func NewHistory(store *entity.Store, state *State) *History {
	return &History{
		store: store,
		state: state,
		events: []Event{
			{Action: nil, Previous: 0},
		},
		reverse: btree.New(50),
	}
}

func (h *History) Add(actions ...Action) error {
	for _, action := range actions {
		if err := h.add(action); err != nil {
			return err
		}
	}
	return nil
}

func (h *History) Undo() (Action, error) {
	if len(h.events) == 1 {
		return nil, errEmptyHistory
	}
	event := h.events[h.currentEventId]
	prevEventId := event.Previous
	action := event.Action

	if err := h.state.Revert(action); err != nil {
		return nil, fmt.Errorf("applying action to state: %w", err)
	}
	h.currentEventId = prevEventId
	return action, nil
}

func (h *History) RedoOptions() []Action {
	eventIds := h.redoableEvents()
	var actions []Action
	for _, eventId := range eventIds {
		event, ok := h.getEvent(eventId)
		if !ok {
			continue
		}
		actions = append(actions, event.Action)
	}
	return actions
}

func (h *History) Tag(name string) {
	h.branches[name] = h.currentEventId
}

func (h *History) commonAncestor(a, b EventId) (EventId, bool) {
	panic("not implemented")
}

func (h *History) unroll(to EventId) error {
	panic("not implemented")
}

func (h *History) uproll(to EventId) error {
	panic("not implemented")
}

func (h *History) Checkout(name string) error {
	currentId := h.currentEventId
	branchId, ok := h.branches[name]
	if !ok {
		return fmt.Errorf("unknown branch")
	}
	ancestor, ok := h.commonAncestor(currentId, branchId)
	if !ok {
		return fmt.Errorf("no common ancestor")
	}
	if err := h.unroll(ancestor); err != nil {
		return fmt.Errorf("unrolling: %w", err)
	}
	if err := h.uproll(branchId); err != nil {
		return fmt.Errorf("uprolling: %w", err)
	}
	return nil
}

func (h *History) lastApply() uint64 {
	la := h.apply
	h.apply++
	return la
}

func (h *History) add(action Action) error {
	redoEventId, redo := h.tryRedo(action)
	if err := h.state.Apply(action); err != nil {
		return fmt.Errorf("applying action to state: %w", err)
	}
	currentEvent, ok := h.getEvent(h.currentEventId)
	if !ok {
		return fmt.Errorf("count not find current event")
	}
	prevEventId := h.currentEventId
	newEventId := redoEventId
	if !redo {
		newEventId = EventId(len(h.events))
		h.events = append(h.events, Event{
			Depth:    currentEvent.Depth + 1,
			Action:   action,
			Previous: prevEventId,
		})
	}
	h.currentEventId = newEventId
	h.reverse.ReplaceOrInsert(&reverseEvent{Prev: prevEventId, New: newEventId, LastApply: h.lastApply()})
	return nil
}

func (h *History) getEvent(eventId EventId) (Event, bool) {
	if uint64(eventId) >= uint64(len(h.events)) {
		return Event{}, false
	}
	return h.events[eventId], true
}

func (h *History) tryRedo(action Action) (EventId, bool) {
	redoEventIds := h.redoableEvents()
	for _, redoEventId := range redoEventIds {
		redoEvent, ok := h.getEvent(redoEventId)
		if !ok {
			return 0, false
		}
		if redoEvent.Action == action {
			return redoEventId, true
		}
	}
	return 0, false
}

func (h *History) redoableEvents() []EventId {
	from := &reverseEvent{
		Prev: h.currentEventId,
		New:  0,
	}
	var reverseEvents []*reverseEvent
	h.reverse.AscendGreaterOrEqual(from, func(i btree.Item) bool {
		found, ok := i.(*reverseEvent)
		if !ok {
			panic(fmt.Errorf("invalid type in btree: want %T, got %T", found, i))
		}
		if found.Prev != h.currentEventId {
			return false
		}
		reverseEvents = append(reverseEvents, found)
		return true
	})
	sort.Slice(reverseEvents, func(i, j int) bool {
		return reverseEvents[j].LastApply < reverseEvents[i].LastApply
	})
	var eventIds []EventId
	for _, rev := range reverseEvents {
		eventIds = append(eventIds, rev.New)
	}
	return eventIds
}

type reverseEvent struct {
	Prev      EventId
	New       EventId
	LastApply uint64
}

func (a *reverseEvent) Less(item btree.Item) bool {
	b, ok := item.(*reverseEvent)
	if !ok {
		panic(fmt.Errorf("invalid type in reverse event btree: want %T, got %T", b, item))
	}
	if a.Prev < b.Prev {
		return true
	}
	if b.Prev < a.Prev {
		return false
	}
	if a.New < b.New {
		return true
	}
	return false
}

var errEmptyHistory = fmt.Errorf("history is empty")

func IsEmptyHistoryError(err error) bool {
	return errors.Is(err, errEmptyHistory)
}
