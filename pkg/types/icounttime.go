package types

/*
	Интерфейс к счётчику оставшегося времени
*/

// ICountTime -- нтерфейс к счётчику оставшегося времени
type ICountTime interface {
	// Parse -- устанавливает интервал времени
	Parse(string) error
	// Set -- устанавливает интервал времени из числа секунд
	Set(int) error
	// Get -- возвращает оставшееся время
	Get() int
	// String -- возвращает стороковое представление оставшегося времени
	String() string
	// Stop -- останавливает работу отсчёта
	Stop()
	// ChanSig -- возвращает канал тиков
	ChanSig() <-chan int
}

func FakeTest() {
	// fmt._rintf("fake test\n")
}
