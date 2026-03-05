package wtypes

// IHxUrlPatch -- атрибут пути HTMX
type IHxUrlPatch interface {
	// Get -- возвращает путь атрибута
	Get() string
	// Set -- устанавливает путь атрибута
	Set(string)
}
