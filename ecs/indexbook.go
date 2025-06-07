package ecs

import (
	"iter"
	"reflect"

	"github.com/PieterD/boevig/rang"
)

type indexPage interface {
	Set(id EntityID, vi any)
	Remove(id EntityID)
	SeekSeq(vi any) iter.Seq[rang.Seekable[EntityID]]
}

type indexBook struct {
	pages map[indexTuple]indexPage
}

func newIndexBook() *indexBook {
	return &indexBook{
		pages: make(map[indexTuple]indexPage),
	}
}

func (book *indexBook) Set(id EntityID, c Component) {
	for _, index := range c.Index() {
		index.apply(book, id)
	}
}

func (book *indexBook) RemoveAll(id EntityID) {
	for _, page := range book.pages {
		page.Remove(id)
	}
}

func (book *indexBook) Remove(id EntityID, c Component) {
	for _, index := range c.Index() {
		index.remove(book, id)
	}
}

func (book *indexBook) Search(params ...Indexer) iter.Seq[rang.Seekable[EntityID]] {
	var seqs []iter.Seq[rang.Seekable[EntityID]]
	for _, param := range params {
		seqs = append(seqs, param.search(book))
	}
	if len(seqs) == 0 {
		return func(yield func(rang.Seekable[EntityID]) bool) {
			return
		}
	}
	if len(seqs) == 1 {
		return seqs[0]
	}
	return rang.NewOrdered(EntityID.Less).Intersect(seqs...)
}

func (book *indexBook) getPage(indexName string, value any) indexPage {
	tup := indexTuple{
		Name: indexName,
		Type: reflect.TypeOf(value),
	}
	page, ok := book.pages[tup]
	if !ok {
		return nil
	}
	return page
}

type indexTuple struct {
	Name string
	Type reflect.Type
}

func getPageG[T comparable](book *indexBook, indexName string, value T) *indexPageG[T] {
	tup := indexTuple{
		Name: indexName,
		Type: reflect.TypeOf(value),
	}
	page, ok := book.pages[tup]
	if !ok {
		page = newIndexPageG[T]()
		book.pages[tup] = page
	}
	return page.(*indexPageG[T])
}
