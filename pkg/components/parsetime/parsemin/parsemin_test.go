package parsemin

import (
	"testing"
)

/*
	Тест для парсера времени часов
*/

// Тестер для проверки парсера времени
type tester struct {
	t  *testing.T
	ph *ParseMin
}

func TestParseMin(t *testing.T) {
	test := &tester{
		t: t,
	}
	test.create()
	test.parse()
	test.set()
	test.reset()
}

// Устанавливает целочисленное значение
func (sf *tester) set() {
	sf.t.Logf("=set=\n")
	sf.setGood1()
	sf.setBad1()
}

// Кривое число минут
func (sf *tester) setBad1() {
	sf.t.Logf("=setBad1=\n")
	sf.ph.Set(60)
	if strHour := sf.ph.String(); strHour != "08" {
		sf.t.Errorf("setBad1(): strHour(%q)!='08'\n", strHour)
	}
}

func (sf *tester) setGood1() {
	sf.t.Logf("=setGood1=\n")
	sf.ph.Set(8)
	if strHour := sf.ph.String(); strHour != "08" {
		sf.t.Errorf("setGood1(): strHour(%q)!='08'\n", strHour)
	}
}

// Сброс часов в ноль
func (sf *tester) reset() {
	sf.t.Logf("=reset=\n")
	sf.ph.Reset()
	if strHour := sf.ph.String(); strHour != "00" {
		sf.t.Errorf("reset(): strHour(%q)!='00'\n", strHour)
	}
}

// Устанавливает значение минут
func (sf *tester) parse() {
	sf.t.Logf("=parse=\n")
	sf.parseBad1()
	sf.parseBad2()
	sf.parseGood1()
}

// Установка правильных минут
func (sf *tester) parseGood1() {
	sf.t.Logf("=parseGood1=\n")
	sf.ph.Parse("8")
	if strHour := sf.ph.String(); strHour != "08" {
		sf.t.Errorf("parseGood1(): strHour(%q)!='08'\n", strHour)
	}
}

// Установка отрицательных минут
func (sf *tester) parseBad2() {
	sf.t.Logf("=parseBad2=\n")
	sf.ph.Parse("-1")
}

// Установка не минут
func (sf *tester) parseBad1() {
	sf.t.Logf("=parseBad1=\n")
	sf.ph.Parse("abc")
}

// Создание парсера минут
func (sf *tester) create() {
	sf.t.Logf("=create=\n")
	sf.ph = NewParseMin()
	if sf.ph == nil {
		sf.t.Errorf("create(): parseHour==nil\n")
	}
	if hour := sf.ph.Get(); hour != 0 {
		sf.t.Errorf("create(): hour(%v)!=0\n", hour)
	}
	if strHour := sf.ph.String(); strHour != "00" {
		sf.t.Errorf("create(): strHour(%q)!='00'\n", strHour)
	}
}
