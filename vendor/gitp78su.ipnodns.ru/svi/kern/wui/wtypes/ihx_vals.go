package wtypes

// IHxVals -- словарь значений в элементе HTMX (hx-vals)
//
// Атрибут hx-vals позволяет добавлять параметры, которые будут отправляться с запросом AJAX.
//
// По умолчанию значением этого атрибута является список значений имени-выражения
// в формате JSON
//
//	Примеры
//
// hx-vals='{"myVal": "My Value"}'
//
// hx-vals='js:{lastKey: event.key}'
//
// hx-vals='js:{x: event.clientX, y: event.clientY}'
type IHxVals interface {
	// Get --  возвращает элемент словаря
	Get(key string) any
	// Set -- устанавливает элемент словаря
	Set(key string, val any)
	// Del -- удаляет элемент из словаря
	Del(key string)
	// Clear -- очищает весь словарь
	Clear()
	// Len -- возвращает размер словаря
	Len() int
	// String -- возвращает строковое представление тэга
	String() string
}
