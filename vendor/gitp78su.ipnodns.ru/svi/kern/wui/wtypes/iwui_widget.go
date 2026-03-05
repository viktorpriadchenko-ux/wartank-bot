package wtypes

// IWuiWidget -- WUI виджет
type IWuiWidget interface {
	// Id -- возвращает Id виджета
	Id() string
	// Html -- возвращает HTML-представление виджета
	Html() string
}
