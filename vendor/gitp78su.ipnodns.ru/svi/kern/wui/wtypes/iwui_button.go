package wtypes

// IWuiButton -- WUI-кнопка
type IWuiButton interface {
	IWuiWidget
	// Text -- возвращает текст кнопки
	Text() IWuiText
	// Click -- нажатие кнопки
	Click(map[string]string) string
	// Hx -- атрибуты HTMX
	Hx() IWuiHx
}
