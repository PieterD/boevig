package sim

import (
	"fmt"

	"github.com/PieterD/rlpoc/lib/entity"
)

type (
	EntityId = entity.Id
)

type State struct {
	EntityStore *entity.Store
}

func NewState(store *entity.Store) *State {
	return &State{
		EntityStore: store,
	}
}

func (s *State) Apply(action Action) error {
	return action.Apply(s)
}

func (s *State) Revert(action Action) error {
	return action.Revert(s)
}

func (s *State) Player(vs ...entity.Value) (EntityId, bool) {
	id, ok := s.EntityStore.Entities(vs...).Fix()
	if !ok {
		return entity.Id{}, false
	}
	if s.EntityStore.First(id, vs...) == false {
		panic(fmt.Errorf("entities found %v, but first did not", id))
	}
	return id, true
}
