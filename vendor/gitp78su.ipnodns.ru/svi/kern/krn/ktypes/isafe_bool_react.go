package ktypes

// ISafeBoolReact -- реактивный потокобезопасный булевый признак
type ISafeBoolReact interface {
	ISafeBool
	// Add -- добавляет функцию обратного вызова
	Add(key string, fn func(bool))
	// Delete -- удаляет функцию обратного вызова по ключу
	Delete(key string)
}
