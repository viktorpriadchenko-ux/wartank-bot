package ktypes

import (
	"fmt"
	"reflect"

	. "gitp78su.ipnodns.ru/svi/kern/kc/helpers"
)

// Result — обёртка вокруг результата с возможной ошибкой
//
// Может быть либо только полезное значение, либо только ошибка
type Result[T any] struct {
	val T     // Полезное значение
	err error // Ошибка
}

// NewOk -- возвращает успешный Result с значением
func NewOk[T any](result T) Result[T] {
	// Для некоторых типов нужна дополнительная проверка через reflect
	v := reflect.ValueOf(result)
	switch v.Kind() {
	case reflect.Ptr, reflect.Slice, reflect.Map, reflect.Chan, reflect.Func, reflect.Interface:
		Hassert(!v.IsNil(), "NewOk(): result==nil")
	}
	sf := Result[T]{
		val: result,
	}
	return sf
}

// NewErr -- возвращает Result с ошибкой
func NewErr[T any](err error) Result[T] {
	Hassert(err != nil, "NewError(): err==nil")
	return Result[T]{
		err: err,
	}
}

// IsOk -- возвращает true, если Result содержит значение
func (sf *Result[T]) IsOk() bool {
	return sf.err == nil
}

// IsErr -- возвращает true, если Result содержит ошибку
func (sf *Result[T]) IsErr() bool {
	return sf.err != nil
}

// Unwrap -- возвращает значение, если оно есть, иначе паникует
func (sf *Result[T]) Unwrap() T {
	Hassert(sf.err == nil, "Result[T].Unwrap(): err(%v)!=nil", sf.err)
	return sf.val
}

// UnwrapOr -- возвращает значение, если оно есть, или значение по умолчанию
func (sf *Result[T]) UnwrapOr(defaultValue T) T {
	if sf.err != nil {
		return defaultValue
	}
	return sf.val
}

// UnwrapOrFn -- возвращает значение, если оно есть, или результат выполнения функции
func (sf *Result[T]) UnwrapOrFn(fn func() T) T {
	Hassert(fn != nil, "Result[T].UnwrapOrFn(): fn==nil")
	if sf.err != nil {
		return fn()
	}
	return sf.val
}

// Error -- возвращает ошибку, если она есть
func (sf *Result[T]) Error() error {
	return sf.err
}

// Hassert -- проверяет, что нет ошибки (с паникой)
func (sf *Result[T]) Hassert(msgFormat string, args ...any) {
	msg := fmt.Sprintf(msgFormat, args...)
	Hassert(sf.err == nil, msg+", err=\n\t%v\n", sf.err)
}

// Assert -- проверяет, что нет ошибки (с паникой только на локальном стенде)
func (sf *Result[T]) Assert(msgFormat string, args ...any) {
	msg := fmt.Sprintf(msgFormat, args...)
	Assert(sf.err == nil, msg+", err=\n\t%v\n", sf.err)
}
