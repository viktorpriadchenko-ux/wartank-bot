package wtypes

// IHxSwapOob -- внеполосная подкачка (hx-swap-oob)
//
// Атрибут 'hx-swap-oob' позволяет указать, что некоторый контент в ответе должен
// быть добавлен в DOM не в целевом элементе, а в другом месте, то есть
// «вне диапазона». Это позволяет добавлять обновления к другим элементам в ответе.
//
// В одном ответе может быть несколько целей внеполосной замены
//
//	Примеры
//
// <div id="alerts" hx-swap-oob="true">
//
// Saved!
//
// </div>
//
// hx-swap-oob="beforeend:#table2"
//
// hx-swap-oob="true"
type IHxSwapOob interface {
	// Get -- получает атрибут HTMX
	Get() string
	// Set -- устанавливает атрибут HTMX
	Set(string)
	// String -- возвращает строковое представление тэга
	String() string
}
