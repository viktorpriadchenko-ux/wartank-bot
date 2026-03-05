package counttime

import (
	"testing"
	"time"

	. "gitp78su.ipnodns.ru/svi/kern"
	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"
)

/*
	Тест для счётчика времени
*/

// Тестер для счётчика времени
type tester struct {
	t      *testing.T
	ct     *CountTime
	isCall ISafeBool // Признак обратного вызова
}

// Обратный вызов для счётчика времени
func (sf *tester) call() {
	<-sf.ct.chCall
}

func TestCountTime(t *testing.T) {
	test := &tester{
		t:      t,
		isCall: NewSafeBool(),
	}
	time.Sleep(time.Millisecond * 100)
	test.create()
	test.setInt()
	test.setStr()
	test.checkTick()
	test.cancel()
}

// Оменяет работу таймера
func (sf *tester) cancel() {
	sf.t.Logf("=cancel=\n")
	ct := NewCountTime()
	for len(ct.chTick) > 0 {
		<-ct.chTick
	}
	ctx := GetKernelCtx()
	ctx.Cancel()
	time.Sleep(time.Millisecond * 150)
}

// Проверяет обработчик тика
func (sf *tester) checkTick() {
	ct := NewCountTime()
	{ // Секундный тик
		ct.Parse("00:00:08")
		time.Sleep(time.Second * 1)
		ct.chTick <- 1
		time.Sleep(time.Millisecond * 20)
		if val := ct.String(); val != "00:00:08" {
			sf.t.Errorf("checkTick(): счётчик(%v)!='00:00:08'", val)
		}
	}
	{ // Проверка времени прямо сейчас
		ct.chTick <- 1
		time.Sleep(time.Millisecond * 20)
		if val := ct.String(); val != "00:00:08" {
			sf.t.Errorf("checkTick(): счётчик(%v)!='00:00:08'", val)
		}
	}
	{ // Проверка обратного вызова прямо сейчас
		strTime := time.Now().UTC().Format("15:04:05")
		ct.Parse(strTime)
		if val := ct.String(); val != strTime {
			sf.t.Errorf("checkTick(): счётчик(%v)!=%s", val, strTime)
		}
		ct.chTick <- 1
		// Выход из функции -- и есть факт обратного вызова
		sf.call()
		{ // Проверка отсутствия обратного вызова прямо сейчас
			ct.Parse("00:00:00")
			ct.chTick <- 1
			// Выход из функции -- и есть факт обратного вызова
			sf.call()
			if val := ct.Get(); val != 0 {
				sf.t.Errorf("checkTick(): счётчик(%v)!=0", val)
			}
			ct.Stop()
			time.Sleep(time.Millisecond * 50)
		}
	}
}

func (sf *tester) setStrBad1(strBad string) {
	sf.ct.Parse(strBad)
}

// Устанавливает строковое значение времени
func (sf *tester) setStr() {
	go sf.call()
	ct := NewCountTime()
	ct.Parse("") // BAD-1 пустая строка
	// BAD-2 неформатная строка
	sf.setStrBad1(":::")
	// BAD-3 кривые часы
	sf.setStrBad1("a1:02:03")
	// BAD-4 кривые минуты
	sf.setStrBad1("01:a2:03")
	// BAD-5 кривые секунды
	sf.setStrBad1("01:02:a3")
	// BAD-6 кривые только секунды
	sf.setStrBad1("a3")
	// BAD-7 кривые минуты +секунды
	sf.setStrBad1("a2:03")
	// BAD-8 кривые часы +минуты +секунды
	sf.setStrBad1("a1:02:03")
	// BAD-9 минуты +кривые секунды
	sf.setStrBad1("02:a3")
	// BAD-10 кривые минуты +секунды
	sf.setStrBad1("60:03")
	// BAD-11 кривые минуты +секунды
	sf.setStrBad1("-1:03")
	// BAD-12 минуты +кривые секунды
	sf.setStrBad1("01:60")
	// BAD-13 минуты +кривые секунды
	sf.setStrBad1("01:-1")
	// BAD-14 кривые часы +минуты + секунды
	sf.setStrBad1("-1:02:03")
	// BAD-15 кривые часы +минуты + секунды
	//sf.setStrBad1("24:02:03")
	ct.Parse("03")       // GOOD-1 секунды
	ct.Parse("02:03")    // GOOD-2 минуты секунды
	ct.Parse("01:02:03") // GOOD-3 часы минуты секунды
}

// Устанавливает число секунд для отсчёта
func (sf *tester) setInt() {
	go sf.call()
	ct := NewCountTime()
	ct.Set(-1) // Bad-1 Отрицательное число
	{          // GOOD-1
		ct.Set(8)
		if ct.parser.Hour().Get() != 0 {
			sf.t.Errorf("setInt(): GOOD-1 hour(%v)!=0", sf.ct.parser.Hour().Get())
		}
		if ct.parser.Min().Get() != 0 {
			sf.t.Errorf("setInt(): GOOD-1 min(%v)!=0", sf.ct.parser.Min().Get())
		}
		if ct.parser.Sec().Get() != 8 {
			sf.t.Errorf("setInt(): GOOD-1 sec(%v)!=8", sf.ct.parser.Sec().Get())
		}
		if strVal := ct.String(); strVal != "00:00:08" {
			sf.t.Errorf("setInt(): GOOD-1 strVal(%v)!='00:00:08'", strVal)
		}
	}
	{ // GOOD-2
		ct.Set(121)
		if ct.parser.Hour().Get() != 0 {
			sf.t.Errorf("setInt(): GOOD-2 hour(%v)!=0", sf.ct.parser.Hour().Get())
		}
		if ct.parser.Min().Get() != 2 {
			sf.t.Errorf("setInt(): GOOD-2 min(%v)!=2", sf.ct.parser.Min().Get())
		}
		if ct.parser.Sec().Get() != 1 {
			sf.t.Errorf("setInt(): GOOD-2 sec(%v)!=1", sf.ct.parser.Sec().Get())
		}
		if strVal := ct.String(); strVal != "00:02:01" {
			sf.t.Errorf("setInt(): GOOD-2 strVal(%v)!='00:02:01'", strVal)
		}
	}
	{ // GOOD-3
		ct.Set(7203)
		if ct.parser.Hour().Get() != 2 {
			sf.t.Errorf("setInt(): GOOD-3 hour(%v)!=2", sf.ct.parser.Hour().Get())
		}
		if ct.parser.Min().Get() != 0 {
			sf.t.Errorf("setInt(): GOOD-3 min(%v)!=0", sf.ct.parser.Min().Get())
		}
		if ct.parser.Sec().Get() != 3 {
			sf.t.Errorf("setInt(): GOOD-3 sec(%v)!=3", sf.ct.parser.Sec().Get())
		}
		if strVal := ct.String(); strVal != "02:00:03" {
			sf.t.Errorf("setInt(): GOOD-3 strVal(%v)!='02:00:03'", strVal)
		}
	}
}

// Правильное создание
func (sf *tester) createGood1() {
	ct := NewCountTime()
	if ct == nil {
		sf.t.Errorf("createGood1(): countTime==nil")
	}
	if val := ct.Get(); val != 0 {
		sf.t.Errorf("createGood1(): val(%v)!=0", val)
	}
}

// Создание счётчика обратного времени
func (sf *tester) create() {
	sf.createGood1()
}
