package ktypes

import (
	"reflect"

	. "gitp78su.ipnodns.ru/svi/kern/kc/helpers"
)

// Option -- результат возможно содержащий nil
type Option[T any] struct {
	val *T
}

// NewSome - полезный результат
func NewSome[T any](value T) Option[T] {
	// Для некоторых типов нужна дополнительная проверка через reflect
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Ptr, reflect.Slice, reflect.Map, reflect.Chan, reflect.Func, reflect.Interface:
		Hassert(!v.IsNil(), "NewSome[T any](): value==nil")
	}
	sf := Option[T]{val: &value}
	return sf
}

// NewNone - нет результата в ответе
func NewNone[T any]() Option[T] {
	return Option[T]{
		val: nil,
	}
}

// IsSome - проверяет, есть ли значение
func (sf *Option[T]) IsSome() bool {
	return sf.val != nil
}

// IsSome - проверяет, есть ли значение
func (sf *Option[T]) IsNone() bool {
	return sf.val == nil
}

// Unwrap - извлекает значение (паника, если None)
func (sf *Option[T]) Unwrap() T {
	Hassert(sf.val != nil, "Option[T].Unwrap(): val==nil!")
	return *sf.val
}

// UnwrapOr - возвращает значение или дефолтное
func (sf Option[T]) UnwrapOr(defaultValue T) T {
	if sf.val == nil {
		return defaultValue
	}
	return *sf.val
}

// UnwrapOrFn -- возвращает значение, если оно есть, или результат выполнения функции
func (sf *Option[T]) UnwrapOrFn(fn func() T) T {
	Hassert(fn != nil, "Result[T].UnwrapOrFn(): fn==nil")
	if sf.val == nil {
		return fn()
	}
	return *sf.val
}

// Hassert -- проверяет, что нет ошибки (с паникой)
func (sf *Option[T]) Hassert(msgFormat string, args ...any) {
	Hassert(sf.val != nil, msgFormat, args...)
}

// Assert -- проверяет, что нет ошибки (с паникой только на локальном стенде)
func (sf *Option[T]) Assert(msgFormat string, args ...any) {
	Assert(sf.val != nil, msgFormat, args...)
}
