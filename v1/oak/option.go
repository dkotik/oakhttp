package oak

type Option[T any] struct {
	valueOrNil T
}

func Some[T any](value T) Option[T] {
	return Option[T]{
		valueOrNil: value,
	}
}

func None[T any]() Option[T] {
	return Option[T]{valueOrNil: nil}
}
