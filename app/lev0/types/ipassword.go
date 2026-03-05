package types

// IPassword -- возвращает объект пароля
type IPassword interface {
	// Get -- возвращает хранимый пароль
	Get() string
	// Set -- устанавливает хранимый пароль
	Set(val string)
}
