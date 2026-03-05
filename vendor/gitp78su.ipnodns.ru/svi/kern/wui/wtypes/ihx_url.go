package wtypes

// IHxUrl -- атрибут метода HTMX
type IHxUrl interface {
	// Method -- возвращает метод атрибута
	Method() IHxUrlMethod
	// Patch -- возвращает путь атрибута
	Patch() IHxUrlPatch
	// String -- возвращает полное значение
	String() string
}
