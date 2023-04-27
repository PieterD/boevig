package sim

import (
	"fmt"
	"reflect"
	"sort"

	"github.com/PieterD/rlpoc/internal/sim/attrib"
	"github.com/PieterD/rlpoc/lib/entity"
)

type fullEvent struct {
	Id     EntityId
	Action Action
	Attrib *attrib.Event
	Parent *attrib.EventRecentParent
	Child  *attrib.EventRecentChild
}

type Crawler struct {
	store       *entity.Store
	constructor ActionConstructor
	handler     ActionHandler
	event       fullEvent
}

type ActionHandler interface {
	Apply(a Action) error
	Revert(a Action) error
}

type ActionConstructor func(name string, settable bool) (Action, error)

func NewCrawler(store *entity.Store, actionConstructor ActionConstructor, h ActionHandler) *Crawler {
	rootId := store.NewEntityId()
	event := &attrib.Event{
		Depth:  0,
		Action: "",
	}
	eventParent := &attrib.EventParentLink{
		Parent: rootId,
	}
	store.Give(rootId, event, eventParent)
	return &Crawler{
		store:       store,
		constructor: actionConstructor,
		handler:     h,
		event: fullEvent{
			Id:     rootId,
			Attrib: event,
		},
	}
}

func (c *Crawler) Add(actions ...Action) error {
	for _, action := range actions {
		redoId, redoOk, err := c.tryRedo(action)
		if err != nil {
			return fmt.Errorf("trying redo: %w", err)
		}
		if redoOk {
			if err := c.redo(redoId); err != nil {
				return fmt.Errorf("redoing action: %w", err)
			}
			return nil
		}
		if err := c.add(action); err != nil {
			return err
		}
	}
	return nil
}

func (c *Crawler) redo(id entity.Id) error {
	panic(fmt.Errorf("not implemented"))
}

func (c *Crawler) tryRedo(action Action) (entity.Id, bool, error) {
	ids, actions, err := c.RedoOptions()
	if err != nil {
		return entity.Id{}, false, fmt.Errorf("fetching redo options: %w", err)
	}
	if len(ids) == 0 {
		return entity.Id{}, false, nil
	}
	if len(actions) != len(ids) {
		return entity.Id{}, false, fmt.Errorf("actions(%d) not the same length as ids(%d)", len(actions), len(ids))
	}
	for i, redoId := range ids {
		redoAction := actions[i]
		if action == redoAction {
			return redoId, true, nil
		}
	}
	return entity.Id{}, false, nil
}

func (c *Crawler) add(action Action) error {
	newId := c.store.NewEntityId()
	newEvent := &attrib.Event{
		Depth:  c.event.Attrib.Depth + 1,
		Action: action.Name(),
	}
	parentLink := &attrib.EventParentLink{Parent: c.event.Id}
	recentParent := &attrib.EventRecentParent{ParentId: c.event.Id}
	c.store.Give(newId, newEvent, parentLink, recentParent, action)

	childLink := &attrib.EventChildLink{Child: newId}
	recentChild := &attrib.EventRecentChild{ChildId: newId}
	c.store.Give(c.event.Id, childLink, recentChild)
	if err := c.handler.Apply(action); err != nil {
		return err
	}

	c.event = fullEvent{
		Id:     newId,
		Action: action,
		Attrib: newEvent,
		Parent: recentParent,
		Child:  nil,
	}
	return nil
}

func (c *Crawler) atRoot() bool {
	if c.event.Attrib.Depth == 0 {
		return true
	}
	return false
}

func (c *Crawler) Undo() (Action, error) {
	if c.atRoot() {
		return nil, fmt.Errorf("attempted to undo root")
	}
	parentId := c.event.Parent.ParentId
	parentEvent := &attrib.Event{}
	recentParent := &attrib.EventRecentParent{}
	if c.store.First(parentId, parentEvent, recentParent) == false {
		return nil, fmt.Errorf("parent %v (from recentparent) not found", parentId)
	}
	parentAction, err := c.constructor(parentEvent.Action, true)
	if err != nil {
		return nil, fmt.Errorf("creating new action(%s): %w", parentEvent.Action, err)
	}
	if c.store.First(parentId, parentAction) == false {
		return nil, fmt.Errorf("parent %v has no Action(%T)", parentId, parentAction)
	}
	recentChild := &attrib.EventRecentChild{ChildId: c.event.Id}
	c.store.Give(parentId, recentChild)
	childAction := c.event.Action
	if err := c.handler.Revert(childAction); err != nil {
		return nil, fmt.Errorf("reverting action %T: %w", c.event.Action, err)
	}
	c.event = fullEvent{
		Id:     parentId,
		Action: parentAction,
		Attrib: parentEvent,
		Parent: recentParent,
		Child:  recentChild,
	}
	return childAction, nil
}

func (c *Crawler) RedoOptions() ([]entity.Id, []Action, error) {
	childLink := &attrib.EventChildLink{}
	if c.store.First(c.event.Id, childLink) == false {
		return nil, nil, nil
	}
	type partialEvent struct {
		Id     entity.Id
		Action Action
	}
	var pes []partialEvent
	for {
		childId := childLink.Child
		childEvent := &attrib.Event{}
		if c.store.First(childId, childEvent) == false {
			return nil, nil, fmt.Errorf("child event %v has no Event", childId)
		}
		action, err := c.constructor(childEvent.Action, true)
		if err != nil {
			return nil, nil, fmt.Errorf("constructing action %s: %w", childEvent.Action, err)
		}
		if c.store.First(childId, action) == false {
			return nil, nil, fmt.Errorf("child event %v does not have promised action %s", childId, childEvent.Action)
		}
		pes = append(pes, partialEvent{
			Id:     childId,
			Action: action,
		})
		if c.store.Next(c.event.Id, childLink) == false {
			break
		}
	}
	sort.Slice(pes, func(i, j int) bool {
		return pes[j].Id.Less(pes[i].Id)
	})
	var ids []entity.Id
	var actions []Action
	for _, pe := range pes {
		ids = append(ids, pe.Id)
		actions = append(actions, unsettableAction(pe.Action))
	}
	return ids, actions, nil
}

func unsettableAction(a Action) Action {
	v := reflect.ValueOf(a)
	if v.Kind() == reflect.Ptr {
		return v.Elem().Interface().(Action)
	}
	return a
}
