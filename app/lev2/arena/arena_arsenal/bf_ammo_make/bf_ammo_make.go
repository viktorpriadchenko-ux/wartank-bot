// package bf_ammo_make -- бизнес-функция производства снарядов
package bf_ammo_make

import (
	"strings"

	. "gitp78su.ipnodns.ru/svi/kern/kc/helpers"
	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"

	"wartank/app/lev0/cons"
	. "wartank/app/lev0/types"
)

const (
	стрКумулятивы = "кумулятивы"
	стрБронебойки = "бронебойки"
	стрФугасы     = "фугасы"
	стрРемки      = "ремки"
)

// СнарядыСделать -- делает снаряды в арсенале
func СнарядыСделать(конт ILocalCtx) {
	арсенал_ := конт.Get("арсенал")

	арсенал := арсенал_.Val().(ИАренаАрсенал)
	еслиПостроено := арсенал.Состояние().Получ() == cons.РежимПостроено
	еслиОжидание := арсенал.Состояние().Получ() == cons.РежимОжидание
	if !(еслиОжидание || еслиПостроено) {
		return
	}
	приоритет(конт)
}

// ищет приоритет производства
func приоритет(конт ILocalCtx) {
	арсенал := конт.Get("арсенал").Val().(ИАренаАрсенал)
	var (
		ремки      = арсенал.Ремки().Получ()
		фугасы     = арсенал.Фугасы().Получ()
		кумули     = арсенал.Кумулятивы().Получ()
		бронебойки = арсенал.Бронебойки().Получ()
		снарядТип  string
	)
	{ // Контроль по числу снарядов. В равных долях без приоритетов
		// снарядТип = стрФугасы
		снарядТип = стрБронебойки
		if бронебойки > фугасы {
			снарядТип = стрФугасы
		}

		if фугасы > кумули {
			снарядТип = стрКумулятивы
		}
		if фугасы > 0 {
			if ремки < 120 { // Контроль ремки по времени суток и минимальному количеству ремок
				сделатьРемку(конт)
				return
			}
		}
		switch снарядТип {
		case стрФугасы: // Мало фугасов
			сделатьФугасы(конт)
		case стрКумулятивы: // Мало кумулятивов
			сделатьКумули(конт)
		case стрБронебойки: // Мало бронебойных
			сделатьБронебойки(конт)
		default:
			// log._rintf("ERRO Арсенал.сделать(): неизвестный тип арсенала(%v)", typeArmor)
		}
		арсенал.ПродуктСейчас().ИмяУст(снарядТип)
	}
}

// Создать бронебойные
func сделатьБронебойки(конт ILocalCtx) {
	арсенал := конт.Get("арсенал").Val().(ИАренаАрсенал)
	var (
		стрВых      string
		lstArsenal  = арсенал.СписПолучить()
		еслиНайдено bool
		инд         int
	)
	for инд, стрВых = range lstArsenal {
		if strings.Contains(стрВых, `<span class="green2">Бронебойный снаряд</span><br/>`) {
			еслиНайдено = true
			break
		}
	}
	if !еслиНайдено {
		return
	}
	стрВых = lstArsenal[инд+10]
	// Получить ссылку на бронебойные
	lstArmor := strings.Split(стрВых, `<a class="simple-but border" href="`)
	if len(lstArmor) <= 1 { // Тут возможно есть пустая строка
		return // Считаем, что производство уже запущено
	}
	strLink := lstArmor[1]
	lstArmor = strings.Split(strLink, `"><span><span>Начать производство</span></span></a>`)
	strLink = "https://wartank.ru/production/" + lstArmor[0]
	if res := арсенал.Сеть().Get(strLink); res.IsErr() {
		// log._rintf("ERRO ArsenalNet.makeArmor(): in update lstArsenal,  err=\n\t%v\n", err)
		return
	}
	арсенал.ПродуктСейчас().ИмяУст(стрБронебойки)
}

// Создать кумулятивные
func сделатьКумули(конт ILocalCtx) {
	арсенал := конт.Get("арсенал").Val().(ИАренаАрсенал)
	var (
		стрВых      string
		lstArsenal  = арсенал.СписПолучить()
		еслиНайдено bool
		инд         int
	)
	for инд, стрВых = range lstArsenal {
		if strings.Contains(стрВых, `<span class="green2">Кумулятивный снаряд</span><br/>`) {
			еслиНайдено = true
			break
		}
	}
	if !еслиНайдено {
		return
	}
	стрВых = lstArsenal[инд+10]
	if !strings.Contains(стрВых, `>Начать производство<`) {
		return
	}
	// Получить ссылку на кумулятив
	списКумул := strings.Split(стрВых, `<a class="simple-but border" href="`)
	strLink := списКумул[1]
	списКумул = strings.Split(strLink, `"><span><span>Начать производство</span></span></a>`)
	strLink = "https://wartank.ru/production/" + списКумул[0]
	if res := арсенал.Сеть().Get(strLink); res.IsErr() {
		// log._rintf("ERRO ArsenalNet.makeКумуль(): in make product arsenal кумуль , err=\n\t%v\n", err)
		return
	}
	арсенал.ПродуктСейчас().ИмяУст(стрКумулятивы)
}

// Создать фугасы
func сделатьФугасы(конт ILocalCtx) {
	арсенал := конт.Get("арсенал").Val().(ИАренаАрсенал)
	var (
		lstArsenal  = арсенал.СписПолучить()
		стрВых      string
		еслиНайдено bool
		инд         int
	)

	for инд, стрВых = range lstArsenal {
		if strings.Contains(стрВых, `<span class="green2">Фугасный снаряд</span><br/>`) {
			еслиНайдено = true
			break
		}
	}
	if !еслиНайдено {
		return
	}
	стрВых = lstArsenal[инд+10]
	if !strings.Contains(стрВых, `"><span><span>Начать производство</span></span></a>`) {
		return
	}
	// Получить ссылку на ремку
	списКумул := strings.Split(стрВых, `<a class="simple-but border" href="`)
	strLink := списКумул[1]
	списКумул = strings.Split(strLink, `"><span><span>Начать производство</span></span></a>`)
	strLink = "https://wartank.ru/production/" + списКумул[0]
	if res := арсенал.Сеть().Get(strLink); res.IsErr() {
		// log._rintf("ERRO ArsenalNet.makeФугас(): in make request arsenal product, err=\n\t%v\n", err)
		return
	}
	арсенал.ПродуктСейчас().ИмяУст(стрФугасы)
}

// Создать ремку. Выполняется если подходят условия
func сделатьРемку(конт ILocalCtx) {
	арсенал := конт.Get("арсенал").Val().(ИАренаАрсенал)
	var (
		стрВых      string
		еслиНайдено bool
		инд         int
	)
	lstArsenal := арсенал.Сеть().ВебВоркер().Получ("https://wartank.ru/production/Armory")
	// <span class="green2">Ремкомплект</span><br/>
	for инд, стрВых = range lstArsenal {
		if strings.Contains(стрВых, `<span class="green2">Ремкомплект</span><br/>`) {
			еслиНайдено = true
			break
		}
	}
	Hassert(еслиНайдено, "сделатьРемку():Не найдена контрольная строка ремок")
	стрВых = lstArsenal[инд+10]
	// Если статус изменился -- выйти
	if strings.Contains(стрВых, `</div></div></div></div></div></div></div></div>`) {
		return
	}
	// Получить ссылку на ремку
	// <a class="simple-but border" href="Armory?154-1.ILinkListener-productions-2-production-startProduceLink"><span><span>Начать производство</span></span></a>
	_ссылка := strings.TrimPrefix(стрВых, `<a class="simple-but border" href="`)
	_ссылка = strings.TrimSuffix(_ссылка, `"><span><span>Начать производство</span></span></a>`)
	// https://wartank.ru/production/Armory?154-1.ILinkListener-productions-2-production-startProduceLink
	ссылка := "https://wartank.ru/production/" + _ссылка
	_ = арсенал.Сеть().ВебВоркер().Получ(ссылка)
	арсенал.ПродуктСейчас().ИмяУст(стрРемки)
	арсенал.Состояние().Уст(cons.РежимРабота)
}
