package ecs

import "reflect"

type EntityID uint64

func (a EntityID) Less(b EntityID) bool {
	return a < b
}

type Component interface {
	hdrNewPage() cStorePage
	typ() reflect.Type
	Index() []Indexer
}

type PTRContract[T Component] interface {
	*T
	Component
}

type ComponentHeader[T Component, TP PTRContract[T]] struct{}

func (_ ComponentHeader[T, TP]) hdrNewPage() cStorePage {
	return newCStorePageG[T, TP]()
}

func (_ ComponentHeader[T, TP]) typ() reflect.Type {
	return reflect.TypeOf(*new(T))
}

func (_ ComponentHeader[T, TP]) Index() []Indexer {
	return nil
}
