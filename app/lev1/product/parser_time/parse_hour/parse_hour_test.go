package parse_hour

import (
	"testing"
)

/*
	Тест для парсера времени часов
*/

// Тестер для проверки парсера времени
type tester struct {
	t  *testing.T
	ph *ПарсерЧас
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
func (сам *tester) reset() {
	сам.t.Logf("=reset=\n")
	сам.ph.Сброс()
	if strHour := сам.ph.String(); strHour != "00" {
		сам.t.Errorf("setGood2(): strHour(%q)!='00'\n", strHour)
	}
}

// Устанавливает целочисленное значение часов
func (сам *tester) set() {
	сам.t.Logf("=set=\n")
	сам.setBad1()
	сам.setGood1()
}
func (сам *tester) setGood1() {
	сам.t.Logf("=setGood1=\n")
	сам.ph.Уст("8")
}

// Отрицательное значение часа
func (сам *tester) setBad1() {
	сам.t.Logf("=setBad1=\n")
	сам.ph.Уст("-1")
}

// Устанавливает значение часов
func (сам *tester) parse() {
	сам.t.Logf("=parse=\n")
	сам.parseBad1()
	сам.parseBad2()
	сам.parseGood1()
	сам.parseGood2()
}

// Установка правильных больших часов
func (сам *tester) parseGood2() {
	сам.t.Logf("=parseGood2=\n")
	сам.ph.Уст("867")
	if strHour := сам.ph.String(); strHour != "867" {
		сам.t.Errorf("parseGood2(): strHour(%q)!='867'\n", strHour)
	}
}

// Установка правильных часов
func (сам *tester) parseGood1() {
	сам.t.Logf("=parseGood1=\n")
	сам.ph.Уст("8")
	if strHour := сам.ph.String(); strHour != "08" {
		сам.t.Errorf("parseGood1(): strHour(%q)!='08'\n", strHour)
	}
}

// Установка отрицательных часов
func (сам *tester) parseBad2() {
	сам.t.Logf("=parseBad2=\n")
	сам.ph.Уст("-1")
}

// Установка не часов
func (сам *tester) parseBad1() {
	сам.t.Logf("=parseBad1=\n")
	сам.ph.Уст("abc")
}

// Создание парсера часов
func (сам *tester) create() {
	сам.t.Logf("=create=\n")
	сам.ph = НовПарсерЧас()
	if сам.ph == nil {
		сам.t.Errorf("create(): parseHour==nil\n")
	}
	if hour := сам.ph.Получ(); hour != 0 {
		сам.t.Errorf("create(): hour(%v)!=0\n", hour)
	}
	if strHour := сам.ph.String(); strHour != "00" {
		сам.t.Errorf("create(): strHour(%q)!='00'\n", strHour)
	}
}
