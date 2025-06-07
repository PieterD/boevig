package ecs

import (
	"fmt"
	"iter"
	"reflect"

	"github.com/PieterD/boevig/rang"
)

type cStorePage interface {
	Add(id EntityID, c Component)
	Remove(id EntityID)
	Get(id EntityID, cPtr Component) bool
	SeekSeq() iter.Seq[rang.Seekable[EntityID]]
}

type cStoreBook struct {
	components map[reflect.Type]cStorePage
}

func newCStoreBook() *cStoreBook {
	return &cStoreBook{components: make(map[reflect.Type]cStorePage)}
}

func (cb *cStoreBook) Add(id EntityID, components ...Component) {
	for _, component := range components {
		page := cb.getPage(component)
		page.Add(id, component)
	}
}

func (cb *cStoreBook) Remove(id EntityID) {
	for _, m := range cb.components {
		m.Remove(id)
	}
}

func (cb *cStoreBook) RemoveComponent(id EntityID, component Component) {
	cb.getPage(component).Remove(id)
}

func (cb *cStoreBook) Get(id EntityID, componentPtrs ...Component) bool {
	for _, componentPtr := range componentPtrs {
		page := cb.getPage(componentPtr)
		if !page.Get(id, componentPtr) {
			return false
		}
	}
	return true
}

func (cb *cStoreBook) All(componentPtrs ...Component) iter.Seq[rang.Seekable[EntityID]] {
	o := rang.NewOrdered[EntityID](EntityID.Less)
	return o.SeekIterator(func(start *EntityID) iter.Seq[EntityID] {
		return func(yield func(EntityID) bool) {
			var seqs []iter.Seq[rang.Seekable[EntityID]]
			for _, componentPtr := range componentPtrs {
				page := cb.getPage(componentPtr)
				seqs = append(seqs, page.SeekSeq())
			}
			for sid := range o.Intersect(seqs...) {
				if start != nil && sid.Value() < *start {
					sid.Seek(*start)
					start = nil
					continue
				}
				for _, componentPtr := range componentPtrs {
					page := cb.getPage(componentPtr)
					ok := page.Get(sid.Value(), componentPtr)
					if !ok {
						panic(fmt.Errorf("page get %v component %T: ID appeared in union, not in page", sid.Value(), componentPtr))
					}
				}
				if !yield(sid.Value()) {
					return
				}
			}
		}
	})
}

func (cb *cStoreBook) All_(componentPtrs ...Component) iter.Seq[rang.Seekable[EntityID]] {
	return func(yield func(rang.Seekable[EntityID]) bool) {
		o := rang.NewOrdered[EntityID](EntityID.Less)
		var seqs []iter.Seq[rang.Seekable[EntityID]]
		for _, componentPtr := range componentPtrs {
			page := cb.getPage(componentPtr)
			seqs = append(seqs, page.SeekSeq())
		}
		for id := range o.Intersect(seqs...) {
			for _, componentPtr := range componentPtrs {
				page := cb.getPage(componentPtr)
				ok := page.Get(id.Value(), componentPtr)
				if !ok {
					panic(fmt.Errorf("page get %v component %T: ID appeared in union, not in page", id.Value(), componentPtr))
				}
			}
			if !yield(id) {
				return
			}
		}
	}
}

func (cb *cStoreBook) getPage(component Component) cStorePage {
	t := component.typ()
	cs, ok := cb.components[t]
	if !ok {
		cs = component.hdrNewPage()
		cb.components[t] = cs
	}
	return cs
}
