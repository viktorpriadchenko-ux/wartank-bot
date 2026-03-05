package parsesec

import (
	"testing"
)

/*
	Тест для парсера времени секунд
*/

// Тестер для проверки парсера времени
type tester struct {
	t  *testing.T
	ph *ParseSec
}

func TestParsesec(t *testing.T) {
	test := &tester{
		t: t,
	}
	test.create()
	test.parse()
	test.set()
	test.reset()
}

// Целочисленная установка секунд
func (sf *tester) set() {
	sf.t.Logf("=set=\n")
	sf.setBad1()
	sf.setGood1()
}

func (sf *tester) setGood1() {
	sf.t.Logf("=setGood1=\n")
	sf.ph.Set(26)
	if strHour := sf.ph.String(); strHour != "26" {
		sf.t.Errorf("setGood1(): strHour(%q)!='26'\n", strHour)
	}
}

// Отрицательные секунды
func (sf *tester) setBad1() {
	sf.t.Logf("=setBad1=\n")
	sf.ph.Set(-1)
	if strHour := sf.ph.String(); strHour != "59" {
		sf.t.Errorf("setBad1(): strHour(%q)!='59'\n", strHour)
	}
}

// Сброс секунд в ноль
func (sf *tester) reset() {
	sf.t.Logf("=reset=\n")
	sf.ph.Reset()
	if strSec := sf.ph.String(); strSec != "00" {
		sf.t.Errorf("setGood2(): strSec(%q)!='00'\n", strSec)
	}
}

// Устанавливает значение секунд
func (sf *tester) parse() {
	sf.t.Logf("=parse=\n")
	sf.parseBad1()
	sf.parseBad2()
	sf.parseBad3()
	sf.parseGood1()
	sf.parseGood2()
}

// Установка правильных секунд
func (sf *tester) parseGood2() {
	sf.t.Logf("=parseGood2=\n")
	defer func() {
		if _panic := recover(); _panic != nil {
			sf.t.Errorf("parseGood2(): panic=\n\t%v\n", _panic)
		}
	}()
	sf.ph.Parse("59")
	if strHour := sf.ph.String(); strHour != "59" {
		sf.t.Errorf("parseGood2(): strHour(%q)!='867'\n", strHour)
	}
}

// Установка правильных часов
func (sf *tester) parseGood1() {
	sf.t.Logf("=parseGood1=\n")
	defer func() {
		if _panic := recover(); _panic != nil {
			sf.t.Errorf("parseGood1(): panic=\n\t%v\n", _panic)
		}
	}()
	sf.ph.Parse("8")
	if strHour := sf.ph.String(); strHour != "08" {
		sf.t.Errorf("parseGood1(): strHour(%q)!='08'\n", strHour)
	}
}

// Установка больших часов
func (sf *tester) parseBad3() {
	sf.t.Logf("=parseBad3=\n")
	sf.ph.Parse("61")
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
	sf.ph = NewParseSec()
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
