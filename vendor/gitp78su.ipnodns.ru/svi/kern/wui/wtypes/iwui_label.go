package wtypes

// IWuiLabel -- текстовая метка
type IWuiLabel interface {
	IWuiWidget
	// Text -- возвращает текст метки
	Text() IWuiText
}
