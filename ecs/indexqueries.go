package ecs

import (
	"iter"

	"github.com/PieterD/boevig/rang"
)

type Indexer interface {
	search(book *indexBook) iter.Seq[rang.Seekable[EntityID]]
	apply(book *indexBook, id EntityID)
	remove(book *indexBook, id EntityID)
}

func EQ[T comparable](indexName string, value T) EqualityIndexer[T] {
	return EqualityIndexer[T]{
		IndexName: indexName,
		Value:     value,
	}
}

type EqualityIndexer[T comparable] struct {
	IndexName string
	Value     T
}

func (es EqualityIndexer[T]) search(book *indexBook) iter.Seq[rang.Seekable[EntityID]] {
	page := getPageG(book, es.IndexName, es.Value)
	return page.SeekSeqG(es.Value)
}

func (es EqualityIndexer[T]) apply(book *indexBook, id EntityID) {
	page := getPageG(book, es.IndexName, es.Value)
	page.Set(id, es.Value)
}

func (es EqualityIndexer[T]) remove(book *indexBook, id EntityID) {
	page := getPageG(book, es.IndexName, es.Value)
	page.Remove(id)
}
