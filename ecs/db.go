package ecs

import (
	"iter"

	"github.com/PieterD/boevig/rang"
)

type DB struct {
	components     *cStoreBook
	indices        *indexBook
	activeEntities map[EntityID]struct{}
	entityID       EntityID
}

func New() *DB {
	return &DB{
		components:     newCStoreBook(),
		indices:        newIndexBook(),
		activeEntities: make(map[EntityID]struct{}),
	}
}

func (db *DB) NewEntity(components ...Component) EntityID {
	for {
		db.entityID++
		if _, ok := db.activeEntities[db.entityID]; ok {
			continue
		}
		db.activeEntities[db.entityID] = struct{}{}
		break
	}
	db.Set(db.entityID, components...)
	return db.entityID
}

func (db *DB) Remove(id EntityID) {
	delete(db.activeEntities, id)
	db.components.Remove(id)
	db.indices.RemoveAll(id)
}

func (db *DB) Get(id EntityID, componentPtrs ...Component) bool {
	return db.components.Get(id, componentPtrs...)
}

func (db *DB) Set(id EntityID, components ...Component) {
	db.components.Add(id, components...)
	for _, component := range components {
		db.indices.Set(id, component)
	}
}

func (db *DB) Unset(id EntityID, component Component) {
	db.components.RemoveComponent(id, component)
	db.indices.Remove(id, component)
}

func (db *DB) SearchComponents(componentPtrs ...Component) iter.Seq[rang.Seekable[EntityID]] {
	return db.components.All(componentPtrs...)
}

func (db *DB) SearchIndex(indexers ...Indexer) iter.Seq[rang.Seekable[EntityID]] {
	return db.indices.Search(indexers...)
}

type SearchBuilder struct {
	db   *DB
	seqs []iter.Seq[rang.Seekable[EntityID]]
}

func (db *DB) Search() *SearchBuilder {
	return &SearchBuilder{db: db}
}

func (b *SearchBuilder) Components(componentPtrs ...Component) *SearchBuilder {
	b.seqs = append(b.seqs, b.db.SearchComponents(componentPtrs...))
	return b
}

func (b *SearchBuilder) Index(indexers ...Indexer) *SearchBuilder {
	b.seqs = append(b.seqs, b.db.SearchIndex(indexers...))
	return b
}

func (b *SearchBuilder) SeekSeq(i iter.Seq[rang.Seekable[EntityID]]) *SearchBuilder {
	b.seqs = append(b.seqs, i)
	return b
}

func (b *SearchBuilder) Done() iter.Seq[EntityID] {
	if len(b.seqs) == 0 {
		return func(yield func(EntityID) bool) {
			return
		}
	}
	if len(b.seqs) == 1 {
		return rang.UnSeek(b.seqs[0])
	}
	return rang.UnSeek(rang.NewOrdered(EntityID.Less).Intersect(b.seqs...))
}
