package wtypes

// IWuiHx -- атрибуты HTMX
type IWuiHx interface {
	// Url -- возвращает URL HTMX
	Url() IHxUrl
	// Trigger -- возвращает триггер HTMX
	Trigger() IHxTrigger
	// Target -- возвращает цель HTMX
	Target() IHxTarget
	// Swap -- политика замены элемента
	Swap() IHxSwap
	// Oob -- политика внеполосной подкачки
	Oob() IHxSwapOob
	// Vals -- словарь дополнительных значений
	Vals() IHxVals
	// String -- возвращает строку тэгов
	String() string
}
