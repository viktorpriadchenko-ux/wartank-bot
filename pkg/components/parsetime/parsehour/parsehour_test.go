package parsehour

import (
	"testing"
)

/*
	Тест для парсера времени часов
*/

// Тестер для проверки парсера времени
type tester struct {
	t  *testing.T
	ph *ParseHour
}

func TestParseHour(t *testing.T) {
	test := &tester{
		t: t,
	}
	test.create()
	test.parse()
	test.set()
	test.reset()
}

// Сброс часов в ноль
func (sf *tester) reset() {
	sf.t.Logf("=reset=\n")
	sf.ph.Reset()
	if strHour := sf.ph.String(); strHour != "00" {
		sf.t.Errorf("setGood2(): strHour(%q)!='00'\n", strHour)
	}
}

// Устанавливает целочисленное значение часов
func (sf *tester) set() {
	sf.t.Logf("=set=\n")
	sf.setBad1()
	sf.setGood1()
}
func (sf *tester) setGood1() {
	sf.t.Logf("=setGood1=\n")
	sf.ph.Set(8)
}

// Отрицательное значение часа
func (sf *tester) setBad1() {
	sf.t.Logf("=setBad1=\n")
	sf.ph.Set(-1)
}

// Устанавливает значение часов
func (sf *tester) parse() {
	sf.t.Logf("=parse=\n")
	sf.parseBad1()
	sf.parseBad2()
	sf.parseGood1()
	sf.parseGood2()
}

// Установка правильных больших часов
func (sf *tester) parseGood2() {
	sf.t.Logf("=parseGood2=\n")
	sf.ph.Parse("867")
	if strHour := sf.ph.String(); strHour != "867" {
		sf.t.Errorf("parseGood2(): strHour(%q)!='867'\n", strHour)
	}
}

// Установка правильных часов
func (sf *tester) parseGood1() {
	sf.t.Logf("=parseGood1=\n")
	sf.ph.Parse("8")
	if strHour := sf.ph.String(); strHour != "08" {
		sf.t.Errorf("parseGood1(): strHour(%q)!='08'\n", strHour)
	}
}

// Установка отрицательных часов
func (sf *tester) parseBad2() {
	sf.t.Logf("=parseBad2=\n")
	sf.ph.Parse("-1")
}

// Установка не часов
func (sf *tester) parseBad1() {
	sf.t.Logf("=parseBad1=\n")
	sf.ph.Parse("abc")
}

// Создание парсера часов
func (sf *tester) create() {
	sf.t.Logf("=create=\n")
	sf.ph = NewParseHour()
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
