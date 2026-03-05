package parsetime

import (
	"testing"
)

/*
	Тест для парсинга времени
*/

// Тестер для парсера времени
type tester struct {
	t    *testing.T
	pars *ParseTime
}

func TestParseTime(t *testing.T) {
	test := &tester{
		t: t,
	}
	test.create()
	test.parse()
}

// Создание парсера
func (sf *tester) create() {
	sf.t.Logf("=create=\n")
	sf.pars = NewParseTime()
	if sf.pars.hour == nil {
		sf.t.Errorf("create(): hour==nil\n")
	}
	if sf.pars.min == nil {
		sf.t.Errorf("create(): min==nil\n")
	}
	if sf.pars.sec == nil {
		sf.t.Errorf("create(): sec==nil\n")
	}
	if val := sf.pars.Get(); val != 0 {
		sf.t.Errorf("create(): valInt(%d)!=0\n", val)
	}
}

// Парсинг строки времени
func (sf *tester) parse() {
	sf.t.Logf("=parse=\n")
	sf.parseBad1()
	sf.parseSec()
	sf.parseSecBad1()
	sf.parseMin()
	sf.parseMinSecBad1()
	sf.parseMinMinBad1()
	sf.parseHour()
	sf.parseHourSecBad1()
	sf.parseHourMinBad1()
	sf.parseHourHourBad1()
}

func (sf *tester) parseHourHourBad1() {
	sf.t.Logf("=parseHourHourBad1=\n")
	sf.pars.Parse("-11:14:54")
	if val := sf.pars.Get(); val != 5820 {
		sf.t.Errorf("parseHourHourBad1(): valInt(%d)!=5820\n", val)
	}
	if hour := sf.pars.Hour().Get(); hour != 1 {
		sf.t.Errorf("parseHourHourBad1(): hour(%d)!=1\n", hour)
	}
	if min := sf.pars.Min().Get(); min != 14 {
		sf.t.Errorf("parseHourHourBad1(): min(%d)!=14\n", min)
	}
	if sec := sf.pars.Sec().Get(); sec != 54 {
		sf.t.Errorf("parseHourHourBad1(): sec(%d)!=54\n", sec)
	}
}

func (sf *tester) parseHourMinBad1() {
	sf.t.Logf("=parseHourMinBad1=\n")
	sf.pars.Parse("01:-4:01")
	if val := sf.pars.Get(); val != 5820 {
		sf.t.Errorf("parseHourMinBad1(): valInt(%d)!=5820\n", val)
	}
}

func (sf *tester) parseHourSecBad1() {
	sf.t.Logf("=parseHourSecBad1=\n")
	sf.pars.Parse("01:37:a")
	if val := sf.pars.Get(); val != 5820 {
		sf.t.Errorf("parseHourSecBad1(): valInt(%d)!=5820\n", val)
	}
}

func (sf *tester) parseHour() {
	sf.t.Logf("=parseHour=\n")
	sf.pars.Parse("01:37:00")
	if val := sf.pars.Get(); val != 5820 {
		sf.t.Errorf("parseHour(): valInt(%d)!=5820\n", val)
	}
}

func (sf *tester) parseMinMinBad1() {
	sf.t.Logf("=parseMinMinBad1=\n")
	sf.pars.Parse("60:25")
	if val := sf.pars.Get(); val != 444 {
		sf.t.Errorf("parseMinMinBad1(): valInt(%d)!=444\n", val)
	}
}

func (sf *tester) parseMinSecBad1() {
	sf.t.Logf("=parseMinSecBad1=\n")
	sf.pars.Parse("07:-1")
	if val := sf.pars.Get(); val != 7*60+24 {
		sf.t.Errorf("parseMinSecBad1(): valInt(%d)!=7*60+24\n", val)
	}
}

func (sf *tester) parseMin() {
	sf.t.Logf("=parseMin=\n")
	sf.pars.Parse("07:24")
	if val := sf.pars.Get(); val != 7*60+24 {
		sf.t.Errorf("parseMin(): valInt(%d)!=7*60+24\n", val)
	}
}

// Слишком большие секунды
func (sf *tester) parseSecBad1() {
	sf.t.Logf("=parseSecBad1=\n")
	sf.pars.Parse("60")
	if val := sf.pars.Get(); val != 28 {
		sf.t.Errorf("parseSecBad1(): valInt(%d)!=28\n", val)
	}
}

func (sf *tester) parseSec() {
	sf.t.Logf("=parseSec=\n")
	sf.pars.Parse("28")
	if val := sf.pars.Get(); val != 28 {
		sf.t.Errorf("parseSec(): valInt(%d)!=28\n", val)
	}
}

// Нет строки для парсинга
func (sf *tester) parseBad1() {
	sf.t.Logf("=parseBad1=\n")
	sf.pars.Parse("")
}
