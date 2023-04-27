package integration

import (
	"reflect"
	"testing"

	"github.com/PieterD/rlpoc/internal/sim/attrib"
	"github.com/PieterD/rlpoc/lib/cardinal"

	"github.com/PieterD/rlpoc/internal/sim"
	"github.com/PieterD/rlpoc/internal/sim/action"
	"github.com/PieterD/rlpoc/lib/entity"
)

func TestCrawler(t *testing.T) {
	store := entity.NewStore()
	state := sim.NewState(store)
	crawler := sim.NewCrawler(store, action.NewAction, state)

	roomId := store.NewEntityId()
	playerId := store.NewEntityId()
	walkEast := action.Walk{
		EntityId:  playerId,
		Direction: cardinal.East,
	}
	walkNorth := action.Walk{
		EntityId:  playerId,
		Direction: cardinal.North,
	}
	err := crawler.Add(
		action.SpawnRoom{
			EntityId: roomId,
			Shape:    box5x5,
		},
		action.SpawnPlayer{
			EntityId: playerId,
			Room:     roomId,
			Position: cardinal.Coord{X: 2, Y: 2},
		},
		walkEast,
	)
	if err != nil {
		t.Fatalf("failed to add actions: %v", err)
	}
	gotLocation := attrib.Location{}
	gotPlayerId, ok := state.Player(&gotLocation)
	if !ok {
		t.Fatalf("failed to get player")
	}
	if want, got := playerId, gotPlayerId; want != got {
		t.Logf("want: %#v", want)
		t.Logf("got : %#v", got)
		t.Fatalf("invalid playerid")
	}
	wantLocation := attrib.Location{
		Room:     roomId,
		Position: cardinal.Coord{X: 3, Y: 2},
	}
	if want, got := wantLocation, gotLocation; want != got {
		t.Logf("want: %#v", want)
		t.Logf("got : %#v", got)
		t.Fatalf("invalid location")
	}
	undoneAction, err := crawler.Undo()
	if err != nil {
		t.Fatalf("undoing: %v", err)
	}
	if undoneAction != walkEast {
		t.Logf("want: %#v", walkEast)
		t.Logf("got : %#v", undoneAction)
		t.Fatalf("invalid undone action")
	}
	wantLocation = attrib.Location{
		Room:     roomId,
		Position: cardinal.Coord{X: 2, Y: 2},
	}
	state.Player(&gotLocation)
	if want, got := wantLocation, gotLocation; want != got {
		t.Logf("want: %#v", want)
		t.Logf("got : %#v", got)
		t.Fatalf("invalid location")
	}
	if err := crawler.Add(walkNorth); err != nil {
		t.Fatalf("walking north: %v", err)
	}
	wantLocation = attrib.Location{
		Room:     roomId,
		Position: cardinal.Coord{X: 2, Y: 1},
	}
	state.Player(&gotLocation)
	if want, got := wantLocation, gotLocation; want != got {
		t.Logf("want: %#v", want)
		t.Logf("got : %#v", got)
		t.Fatalf("invalid location")
	}
	undoneAction, err = crawler.Undo()
	if err != nil {
		t.Fatalf("undoing action: %v", err)
	}
	if undoneAction != walkNorth {
		t.Logf("want: %#v", walkNorth)
		t.Logf("got : %#v", undoneAction)
		t.Fatalf("invalid undone action")
	}
	_, gotRedoOptions, err := crawler.RedoOptions()
	if err != nil {
		t.Fatalf("redo options: %v", err)
	}
	wantRedoOptions := []sim.Action{walkNorth, walkEast}
	if !reflect.DeepEqual(wantRedoOptions, gotRedoOptions) {
		t.Logf("want: %#v", wantRedoOptions)
		t.Logf("got : %#v", gotRedoOptions)
		t.Fatalf("invalid redo options")
	}
	if err := crawler.Add(walkEast); err != nil {
		t.Fatalf("walking east")
	}
	if _, err := crawler.Undo(); err != nil {
		t.Fatalf("undoing: %v", err)
	}
	_, gotRedoOptions, err = crawler.RedoOptions()
	if err != nil {
		t.Fatalf("redo options: %v", err)
	}
	wantRedoOptions = []sim.Action{walkEast, walkNorth}
	if !reflect.DeepEqual(wantRedoOptions, gotRedoOptions) {
		t.Logf("want: %#v", wantRedoOptions)
		t.Logf("got : %#v", gotRedoOptions)
		t.Fatalf("invalid redo options")
	}
}
