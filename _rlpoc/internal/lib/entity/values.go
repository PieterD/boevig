package entity

type Value interface {
	Less(Value) bool
}
