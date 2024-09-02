package lo

type ifElse[T any] struct {
	result T
	done   bool
}

func If[T any](condition bool, result T) *ifElse[T] { //nolint:revive // chain of ifElse
	if condition {
		return &ifElse[T]{result, true}
	}

	var t T
	return &ifElse[T]{t, false}
}

func (i *ifElse[T]) ElseF(resultF func() T) T {
	if i.done {
		return i.result
	}

	return resultF()
}
