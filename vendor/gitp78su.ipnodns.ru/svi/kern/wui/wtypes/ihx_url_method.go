package wtypes

// IHxUrlMethod -- атрибут метода HTMX
//
// hx-get    // READ
// hx-post   // CREATE (универсальный, по умолчанию)
// hx-patch  // UPDATE
// hx-put    // UPDATE PARTIAL
// hx-delete // DELETE
type IHxUrlMethod interface {
	// Get -- возвращает метод атрибута
	Get() string
	// Set -- устанавливает метод атрибута
	Set(string)
}
