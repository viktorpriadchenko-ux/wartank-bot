package parsetime

import (
	"strings"
	"sync"

	. "gitp78su.ipnodns.ru/svi/kern/kc/helpers"

	"wartank/pkg/components/parsetime/parsehour"
	"wartank/pkg/components/parsetime/parsemin"
	"wartank/pkg/components/parsetime/parsesec"
)

/*
	Выковыривает из строки время и потокобезопасно хранит его
*/

// ParseTime -- потокобезопасный ковырятор строки времени
type ParseTime struct {
	intVal int                  // Числовое значение хранимого времени
	hour   *parsehour.ParseHour // Часы метки времени
	min    *parsemin.ParseMin   // Минуты метки времени
	sec    *parsesec.ParseSec   // Секунды метки времени

	block sync.RWMutex
}

// NewParseTime -- возвращает новый *ParseTime
func NewParseTime() *ParseTime {
	return &ParseTime{
		hour: parsehour.NewParseHour(),
		min:  parsemin.NewParseMin(),
		sec:  parsesec.NewParseSec(),
	}
}

// GetInt -- возвращает хранимоесчётчика времени
func (sf *ParseTime) Get() int {
	sf.block.RLock()
	defer sf.block.RUnlock()
	return sf.intVal
}

// Parse -- разбирает строковое представление на части
func (sf *ParseTime) Parse(strTime string) {
	sf.block.Lock()
	defer sf.block.Unlock()
	Hassert(strTime != "", "CountTime.Set(): val is empty")
	// Разбить время, перевести в секунды
	lstTime := strings.Split(strTime, ":")
	if len(lstTime) == 1 { // Только секунды
		sf.hour.Reset()
		sf.min.Reset()
		sf.sec.Parse(lstTime[0])
	}
	if len(lstTime) == 2 { // Минуты, секунды
		sf.hour.Reset()
		sf.min.Parse(lstTime[0])
		sf.sec.Parse(lstTime[1])
	}
	if len(lstTime) >= 3 { // Есть всё, возможно с левыми полями в конце
		strHour := lstTime[0]
		strMin := lstTime[1]
		strSec := lstTime[2]
		sf.sec.Parse(strSec)
		sf.min.Parse(strMin)
		sf.hour.Parse(strHour)
	}
	sf.intVal = sf.hour.Get()*3600 + sf.min.Get()*60 + sf.sec.Get()
}

// Hour -- возвращает хранимые часы
func (sf *ParseTime) Hour() *parsehour.ParseHour {
	return sf.hour
}

// Min -- возвращает хранимые минуты
func (sf *ParseTime) Min() *parsemin.ParseMin {
	return sf.min
}

// Sec -- возвращает хранимые секунды
func (sf *ParseTime) Sec() *parsesec.ParseSec {
	return sf.sec
}

// String -- возвращает хранимое время
func (sf *ParseTime) String() string {
	res := sf.hour.String() + ":" + sf.min.String() + ":" + sf.sec.String()
	return res
}
