package parse_sec

import (
	"testing"
)

/*
	Тест для парсера времени секунд
*/

// Тестер для проверки парсера времени
type tester struct {
	t  *testing.T
	ph *ПарсерСекунд
}

func TestParseSec(t *testing.T) {
	test := &tester{
		t: t,
	}
	test.create()
	test.parse()
	test.set()
	test.reset()
}

// Целочисленная установка секунд
func (сам *tester) set() {
	сам.t.Logf("=set=\n")
	сам.setGood1()
}

func (сам *tester) setGood1() {
	сам.t.Logf("=setGood1=\n")
	сам.ph.УстСек(26)
	if strHour := сам.ph.String(); strHour != "26" {
		сам.t.Errorf("setGood1(): strHour(%q)!='26'\n", strHour)
	}
}

// Сброс секунд в ноль
func (сам *tester) reset() {
	сам.t.Logf("=reset=\n")
	сам.ph.Сброс()
	if strSec := сам.ph.String(); strSec != "00" {
		сам.t.Errorf("setGood2(): strSec(%q)!='00'\n", strSec)
	}
}

// Устанавливает значение секунд
func (сам *tester) parse() {
	сам.t.Logf("=parse=\n")
	сам.parseBad1()
	сам.parseBad2()
	сам.parseBad3()
	сам.parseGood1()
	сам.parseGood2()
}

// Установка правильных секунд
func (сам *tester) parseGood2() {
	сам.t.Logf("=parseGood2=\n")
	defer func() {
		if _panic := recover(); _panic != nil {
			сам.t.Errorf("parseGood2(): panic=\n\t%v\n", _panic)
		}
	}()
	сам.ph.Уст("59")
	if strHour := сам.ph.String(); strHour != "59" {
		сам.t.Errorf("parseGood2(): strHour(%q)!='867'\n", strHour)
	}
}

// Установка правильных часов
func (сам *tester) parseGood1() {
	сам.t.Logf("=parseGood1=\n")
	defer func() {
		if _panic := recover(); _panic != nil {
			сам.t.Errorf("parseGood1(): panic=\n\t%v\n", _panic)
		}
	}()
	сам.ph.Уст("8")
	if strHour := сам.ph.String(); strHour != "08" {
		сам.t.Errorf("parseGood1(): strHour(%q)!='08'\n", strHour)
	}
}

// Установка больших часов
func (сам *tester) parseBad3() {
	сам.t.Logf("=parseBad3=\n")
	сам.ph.Уст("61")
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
	сам.ph = НовПарсерСекунд()
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
