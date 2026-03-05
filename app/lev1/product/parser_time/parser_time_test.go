package parser_time

import (
	"testing"
)

/*
	Тест для парсинга времени
*/

// Тестер для парсера времени
type tester struct {
	t    *testing.T
	pars *ПарсерВремя
}

func TestParseTime(t *testing.T) {
	test := &tester{
		t: t,
	}
	test.create()
	test.parse()
}

// Создание парсера
func (сам *tester) create() {
	сам.t.Logf("=create=\n")
	сам.pars = НовПарсерВремя()
	if сам.pars.час == nil {
		сам.t.Errorf("create(): hour==nil\n")
	}
	if сам.pars.мин == nil {
		сам.t.Errorf("create(): min==nil\n")
	}
	if сам.pars.сек == nil {
		сам.t.Errorf("create(): sec==nil\n")
	}
	if val := сам.pars.ПолучМилСек(); val != 0 {
		сам.t.Errorf("create(): valInt(%d)!=0\n", val)
	}
}

// Парсинг строки времени
func (сам *tester) parse() {
	сам.t.Logf("=parse=\n")
	сам.parseBad1()
	сам.parseSec()
	сам.parseSecBad1()
	сам.parseMin()
	сам.parseMinSecBad1()
	сам.parseMinMinBad1()
	сам.parseHour()
	сам.parseHourSecBad1()
	сам.parseHourMinBad1()
	сам.parseHourHourBad1()
}

func (сам *tester) parseHourHourBad1() {
	сам.t.Logf("=parseHourHourBad1=\n")
	сам.pars.Уст("-11:14:54")
	if val := сам.pars.ПолучМилСек(); val != 5820 {
		сам.t.Errorf("parseHourHourBad1(): valInt(%d)!=5820\n", val)
	}
	if hour := сам.pars.Час().Получ(); hour != 1 {
		сам.t.Errorf("parseHourHourBad1(): hour(%d)!=1\n", hour)
	}
	if min := сам.pars.Мин().Получ(); min != 14 {
		сам.t.Errorf("parseHourHourBad1(): min(%d)!=14\n", min)
	}
	if sec := сам.pars.Сек().Получ(); sec != 54 {
		сам.t.Errorf("parseHourHourBad1(): sec(%d)!=54\n", sec)
	}
}

func (сам *tester) parseHourMinBad1() {
	сам.t.Logf("=parseHourMinBad1=\n")
	сам.pars.Уст("01:-4:01")
	if val := сам.pars.ПолучМилСек(); val != 5820 {
		сам.t.Errorf("parseHourMinBad1(): valInt(%d)!=5820\n", val)
	}
}

func (сам *tester) parseHourSecBad1() {
	сам.t.Logf("=parseHourSecBad1=\n")
	сам.pars.Уст("01:37:a")
	if val := сам.pars.ПолучМилСек(); val != 5820 {
		сам.t.Errorf("parseHourSecBad1(): valInt(%d)!=5820\n", val)
	}
}

func (сам *tester) parseHour() {
	сам.t.Logf("=parseHour=\n")
	сам.pars.Уст("01:37:00")
	if val := сам.pars.ПолучМилСек(); val != 5820 {
		сам.t.Errorf("parseHour(): valInt(%d)!=5820\n", val)
	}
}

func (сам *tester) parseMinMinBad1() {
	сам.t.Logf("=parseMinMinBad1=\n")
	сам.pars.Уст("60:25")
	if val := сам.pars.ПолучМилСек(); val != 444 {
		сам.t.Errorf("parseMinMinBad1(): valInt(%d)!=444\n", val)
	}
}

func (сам *tester) parseMinSecBad1() {
	сам.t.Logf("=parseMinSecBad1=\n")
	сам.pars.Уст("07:-1")
	if val := сам.pars.ПолучМилСек(); val != 7*60+24 {
		сам.t.Errorf("parseMinSecBad1(): valInt(%d)!=7*60+24\n", val)
	}
}

func (сам *tester) parseMin() {
	сам.t.Logf("=parseMin=\n")
	сам.pars.Уст("07:24")
	if val := сам.pars.ПолучМилСек(); val != 7*60+24 {
		сам.t.Errorf("parseMin(): valInt(%d)!=7*60+24\n", val)
	}
}

// Слишком большие секунды
func (сам *tester) parseSecBad1() {
	сам.t.Logf("=parseSecBad1=\n")
	сам.pars.Уст("60")
	if val := сам.pars.ПолучМилСек(); val != 28 {
		сам.t.Errorf("parseSecBad1(): valInt(%d)!=28\n", val)
	}
}

func (сам *tester) parseSec() {
	сам.t.Logf("=parseSec=\n")
	сам.pars.Уст("28")
	if val := сам.pars.ПолучМилСек(); val != 28 {
		сам.t.Errorf("parseSec(): valInt(%d)!=28\n", val)
	}
}

// Нет строки для парсинга
func (сам *tester) parseBad1() {
	сам.t.Logf("=parseBad1=\n")
	сам.pars.Уст("")
}
