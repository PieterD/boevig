package rang

import "iter"

func Map[FROM, TO any](seq iter.Seq[FROM], f func(FROM) TO) iter.Seq[TO] {
	return func(yield func(TO) bool) {
		for v := range seq {
			newV := f(v)
			if !yield(newV) {
				return
			}
		}
	}
}

func First[T any](seq iter.Seq[T]) (T, bool) {
	var zero T
	for v := range seq {
		return v, true
	}
	return zero, false
}

func ToSlice[T any](seq iter.Seq[T]) []T {
	var vs []T
	for v := range seq {
		vs = append(vs, v)
	}
	return vs
}
