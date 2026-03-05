package wtypes

// IWuiText -- текст WUI
type IWuiText interface {
	IWuiWidget
	// Get -- возвращает текст
	Get() string
	// Set -- устанавливает текст
	Set(string)
}
