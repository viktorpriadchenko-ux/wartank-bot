package parse_min

import (
	"testing"
)

/*
	Тест для парсера времени часов
*/

// Тестер для проверки парсера времени
type tester struct {
	t  *testing.T
	ph *ПарсерМинут
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
func (сам *tester) set() {
	сам.t.Logf("=set=\n")
	сам.setGood1()
	сам.setBad1()
}

// Кривое число минут
func (сам *tester) setBad1() {
	сам.t.Logf("=setBad1=\n")
	defer func() {
		if _panic := recover(); _panic == nil {
			сам.t.Fatalf("setBad1(): panic==nil")
		}
		if strHour := сам.ph.String(); strHour != "08" {
			сам.t.Errorf("setBad1(): strHour(%q)!='08'\n", strHour)
		}
	}()
	сам.ph.УстМин(60)
}

func (сам *tester) setGood1() {
	сам.t.Logf("=setGood1=\n")
	сам.ph.УстМин(8)
	if strHour := сам.ph.String(); strHour != "08" {
		сам.t.Errorf("setGood1(): strHour(%q)!='08'\n", strHour)
	}
}

// Сброс часов в ноль
func (сам *tester) reset() {
	сам.t.Logf("=reset=\n")
	сам.ph.Сброс()
	if strHour := сам.ph.String(); strHour != "00" {
		сам.t.Errorf("reset(): strHour(%q)!='00'\n", strHour)
	}
}

// Устанавливает значение минут
func (сам *tester) parse() {
	сам.t.Logf("=parse=\n")
	сам.parseBad1()
	сам.parseBad2()
	сам.parseGood1()
}

// Установка правильных минут
func (сам *tester) parseGood1() {
	сам.t.Logf("=parseGood1=\n")
	сам.ph.Уст("8")
	if strHour := сам.ph.String(); strHour != "08" {
		сам.t.Errorf("parseGood1(): strHour(%q)!='08'\n", strHour)
	}
}

// Установка отрицательных минут
func (сам *tester) parseBad2() {
	сам.t.Logf("=parseBad2=\n")
	сам.ph.Уст("-1")
}

// Установка не минут
func (сам *tester) parseBad1() {
	сам.t.Logf("=parseBad1=\n")
	сам.ph.Уст("abc")
}

// Создание парсера минут
func (сам *tester) create() {
	сам.t.Logf("=create=\n")
	сам.ph = НовПарсерМинут()
	if сам.ph == nil {
		сам.t.Errorf("create(): parseHour==nil\n")
	}
	if мин := сам.ph.Получ(); мин != 0 {
		сам.t.Errorf("create(): мин(%v)!=0\n", мин)
	}
	if стрМин := сам.ph.String(); стрМин != "00" {
		сам.t.Errorf("create(): стрМин(%q)!='00'\n", стрМин)
	}
}
