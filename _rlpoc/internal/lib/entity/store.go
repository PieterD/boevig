package entity

import (
	"fmt"
	"reflect"
)

type Store struct {
	nextEntityId Id
	attributes   map[reflect.Type]*attribStore
}

func NewStore() *Store {
	return &Store{
		nextEntityId: firstId(),
		attributes:   make(map[reflect.Type]*attribStore),
	}
}

func (s *Store) mustGetAttribStore(v Value) (settable bool, attribStore *attribStore) {
	reflectValue := reflect.ValueOf(v)
	if reflectValue.Kind() == reflect.Ptr {
		settable = true
		reflectValue = reflectValue.Elem()
	}
	typ := reflectValue.Type()
	if typ.Kind() == reflect.Ptr {
		panic(fmt.Errorf("invalid input value %T: must not be a double pointer type", v))
	}
	attribStore, ok := s.attributes[typ]
	if !ok {
		attribStore = newAttribStore()
		s.attributes[typ] = attribStore
	}
	return settable, attribStore
}

func (s *Store) NewEntityId() Id {
	id := s.nextEntityId
	s.nextEntityId = id.next()
	return id
}

func (s *Store) DestroyEntity(entityId Id) {
	for _, attribStore := range s.attributes {
		attribStore.Destroy(entityId)
	}
}

func (s *Store) Give(entityId Id, vs ...Value) {
	for _, v := range vs {
		_, attribStore := s.mustGetAttribStore(v)
		attribStore.Give(entityId, v)
	}
}

func (s *Store) Clear(v Value) {
	_, attribStore := s.mustGetAttribStore(v)
	attribStore.Clear()
}

func (s *Store) Remove(entityId Id, vs ...Value) {
	for _, v := range vs {
		_, attribStore := s.mustGetAttribStore(v)
		attribStore.Remove(entityId, v)
	}
}

func (s *Store) First(entityId Id, vs ...Value) bool {
	for _, v := range vs {
		settable, attribStore := s.mustGetAttribStore(v)
		if !settable {
			panic(fmt.Errorf("value type %T can not be set", v))
		}
		if attribStore.FirstAttrib(entityId, v) == false {
			return false
		}
	}
	return true
}

func (s *Store) Next(entityId Id, v Value) bool {
	settable, attribStore := s.mustGetAttribStore(v)
	if !settable {
		panic(fmt.Errorf("value type %T not settable", v))
	}
	if attribStore.NextAttrib(entityId, v) == false {
		return false
	}
	return true
}

func (s *Store) Entities(vs ...Value) Cursor {
	cursor := newMultiAttributeCursor(s, vs...)
	return cursor
}
