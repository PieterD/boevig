package ecs_test

import (
	"testing"

	"github.com/PieterD/boevig/ecs"
	"github.com/PieterD/boevig/rang"
	"github.com/stretchr/testify/require"
)

// Define some components

type Player struct {
	// Magic component header with internal shenanigans
	ecs.ComponentHeader[Player, *Player]
	Name string
}

type Monster struct {
	ecs.ComponentHeader[Monster, *Monster]
	Name string
}

type Location struct {
	ecs.ComponentHeader[Location, *Location]
	Coord Coord
}

// Location has an index on its coordinate
func (c Location) Index() []ecs.Indexer {
	return []ecs.Indexer{
		// Report the current value to the index system on request.
		ecs.EQ("Location.Coord", c.Coord),
	}
}

type Coord struct {
	X int
	Y int
}

func TestDB_Example(t *testing.T) {
	db := ecs.New()
	// Create some entities
	p1 := db.NewEntity(Player{Name: "player one"},
		Location{Coord: Coord{X: 1, Y: 2}})
	m1 := db.NewEntity(Monster{Name: "bat"},
		Location{Coord: Coord{X: 1, Y: 1}})
	m2 := db.NewEntity(Monster{Name: "rat"},
		Location{Coord: Coord{X: 1, Y: 2}})
	m3 := db.NewEntity(Monster{Name: "ghost"})

	t.Run("find player, get location", func(t *testing.T) {
		var player Player
		// Done returns a sequence, get the first value
		id, ok := rang.First(db.Search().Components(&player).Done())
		require.True(t, ok)
		require.Equal(t, p1, id)
		var location Location
		// Get supports an arbitrary amount of components.
		// Only if the entity given by id has all of them,
		// does Get return true.
		ok = db.Get(id, &location)
		require.True(t, ok)
		require.Equal(t, location.Coord, Coord{X: 1, Y: 2})
	})

	t.Run("find all monsters", func(t *testing.T) {
		var ids []ecs.EntityID
		var monsters []Monster
		var monster Monster
		for id := range db.Search().Components(&monster).Done() {
			ids = append(ids, id)
			// As we loop through the sequence, &monster is automatically set
			// to the component value for this id.
			// If we provided more components, they would all be filled.
			monsters = append(monsters, monster)
		}
		expectedIds := []ecs.EntityID{m1, m2, m3}
		expectedMonsters := []Monster{
			{Name: "bat"},
			{Name: "rat"},
			{Name: "ghost"},
		}
		require.Equal(t, expectedIds, ids)
		require.Equal(t, expectedMonsters, monsters)
	})
	t.Run("find all monsters with locations", func(t *testing.T) {
		var ids []ecs.EntityID
		var monsters []Monster
		var locations []Location
		var monster Monster
		var location Location
		seq := db.
			Search().
			Components(
				&monster,
				&location,
			).
			Done()
		for id := range seq {
			ids = append(ids, id)
			monsters = append(monsters, monster)
			locations = append(locations, location)
		}
		expectedIds := []ecs.EntityID{m1, m2}
		expectedMonsters := []Monster{
			{Name: "bat"},
			{Name: "rat"},
		}
		expectedLocations := []Location{
			{Coord: Coord{X: 1, Y: 1}},
			{Coord: Coord{X: 1, Y: 2}},
		}
		require.Equal(t, expectedIds, ids)
		require.Equal(t, expectedMonsters, monsters)
		require.Equal(t, expectedLocations, locations)
	})
	t.Run("find everything on 1,2", func(t *testing.T) {
		var ids []ecs.EntityID
		// Search the coordinate index.
		seq := db.Search().Index(
			// Note how we use the same function to both report the value above,
			// and to search the value here :)
			ecs.EQ("Location.Coord", Coord{X: 1, Y: 2}),
		).Done()
		for id := range seq {
			ids = append(ids, id)
		}
		expectedIds := []ecs.EntityID{p1, m2}
		require.Equal(t, expectedIds, ids)
	})
	t.Run("find monsters on 1,2", func(t *testing.T) {
		var ids []ecs.EntityID
		var monsters []Monster
		var monster Monster
		// Search on both a component and the coordinate index.
		seq := db.Search().
			Components(
				// Multiple components are possible.
				&monster,
			).
			Index(
				// Multiple index searches are possible.
				ecs.EQ("Location.Coord", Coord{X: 1, Y: 2}),
			).
			Done()
		for id := range seq {
			ids = append(ids, id)
			monsters = append(monsters, monster)
		}
		expectedIds := []ecs.EntityID{m2}
		expectedMonsters := []Monster{
			{Name: "rat"},
		}
		require.Equal(t, expectedIds, ids)
		require.Equal(t, expectedMonsters, monsters)
	})
}
