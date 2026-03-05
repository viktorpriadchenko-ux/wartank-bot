// package bf_mission_simple -- бизнес-функция забрать простые миссии
package bf_mission_simple

import (
	"strings"

	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"

	. "wartank/app/lev0/types"
)

// МиссииПростыеЗабрать -- забирает выполненные простые миссии
func МиссииПростыеЗабрать(конт ILocalCtx) {
	сражениеЗащита(конт)
	сражениеНаступление(конт)
	проведиВойну(конт)
	подряд5побед(конт)
	подряд6побед(конт)
	подряд10побед(конт)
	победаСхватка(конт)
	сделать10ресурсов(конт)
	убить3танка(конт)
	топливоДив(конт)
	upMan(конт)
}

// Проверяет награду за уничтожить 3 танка в бою
func убить3танка(конт ILocalCtx) {
	арена := конт.Get("миссии_простые").Val().(ИАрена)
	var (
		strOut      string
		lstMissions = арена.СписПолучить()
		еслиНайдено bool
		ind         int
	)
	for ind, strOut = range lstMissions {
		if strings.Contains(strOut, `Уничтожь в бою 3 танка<br/>`) {
			еслиНайдено = true
			break
		}
	}
	if !еслиНайдено {
		return
	}
	// <a class="simple-but border" href="?70-1.ILinkListener-missions-cc-0-c-awardLink"><span><span>Получить награду</span></span></a>
	ind += 23
	strOut = lstMissions[ind]
	if !strings.Contains(strOut, `>Получить награду<`) {
		return
	}
	lstLink := strings.Split(strOut, `<a class="simple-but border" href="`)
	strLink := lstLink[1]
	lstLink = strings.Split(strLink, `"><span><span>Получить награду</span></span></a>`)
	strLink = "https://wartank.ru/missions/" + lstLink[0]
	res := арена.Сеть().Get(strLink)
	if res.IsErr() {
		// log._rintf("ERRO Миссии.убить3танка(): при выполнении GET-запроса, err=\n\t%v\n", err)
		return
	}
	арена.СтрОбновить(res.Unwrap())
}

// Проверяет награду оборонительного сражения
func сражениеЗащита(конт ILocalCtx) {
	арена := конт.Get("миссии_простые").Val().(ИАрена)
	var (
		strOut      string
		еслиНайдено bool
	)
	арена.Обновить()

	списМиссия := арена.СписПолучить()
	// <a class="simple-but border" href="?23-1.ILinkListener-missions-cc-0-c-awardLink"><span><span>Получить награду</span></span></a>
	for _, strOut = range списМиссия {
		if strings.Contains(strOut, `.ILinkListener-missions-cc-0-c-awardLink`) {
			еслиНайдено = true
			break
		}
	}
	if !еслиНайдено {
		return
	}
	// <a class="simple-but border" href="?70-1.ILinkListener-missions-cc-0-c-awardLink"><span><span>Получить награду</span></span></a>
	_ссылка := strings.TrimPrefix(strOut, `<a class="simple-but border" href="`)
	_ссылка = strings.TrimSuffix(_ссылка, `"><span><span>Получить награду</span></span></a>`)
	strLink := "https://wartank.ru/missions/" + _ссылка
	res := арена.Сеть().Get(strLink)
	if res.IsErr() {
		// log._rintf("ERRO Миссии.сражениеЗащита(): при выполнении GET-запроса, err=\n\t%v\n", err)
		return
	}
	арена.СтрОбновить(res.Unwrap())
}

// Проверяет награду за одну войну
func проведиВойну(конт ILocalCtx) {
	арена := конт.Get("миссии_простые").Val().(ИАрена)
	var (
		strOut      string
		lstMissions = арена.СписПолучить()
		еслиНайдено bool
		ind         int
	)
	for ind, strOut = range lstMissions {
		if strings.Contains(strOut, `Проведи одну войну<br/>`) {
			еслиНайдено = true
			break
		}
	}
	if !еслиНайдено {
		return
	}
	// <a class="simple-but border" href="?70-1.ILinkListener-missions-cc-0-c-awardLink"><span><span>Получить награду</span></span></a>
	ind += 23
	strOut = lstMissions[ind]
	if !strings.Contains(strOut, `>Получить награду<`) {
		return
	}
	lstLink := strings.Split(strOut, `<a class="simple-but border" href="`)
	strLink := lstLink[1]
	lstLink = strings.Split(strLink, `"><span><span>Получить награду</span></span></a>`)
	strLink = "https://wartank.ru/missions/" + lstLink[0]
	res := арена.Сеть().Get(strLink)
	if res.IsErr() {
		// log._rintf("ERRO Миссии.battleWar(): при выполнении GET-запроса, err=\n\t%v\n", err)
		return
	}
	арена.СтрОбновить(res.Unwrap())
}

// Проверяет награду наступательного сражения
func сражениеНаступление(конт ILocalCtx) {
	арена := конт.Get("миссии_простые").Val().(ИАрена)
	var (
		strOut      string
		еслиНайдено bool
		ind         int
		lstMissions = арена.СписПолучить()
	)
	if len(lstMissions) == 0 {

		арена.Обновить()

		lstMissions = арена.СписПолучить()
	}
	for ind, strOut = range lstMissions {
		if strings.Contains(strOut, "Проведи одно наступательное сражение<br/>") {
			еслиНайдено = true
			break
		}
	}
	if !еслиНайдено {
		return
	}
	// <a class="simple-but border" href="?70-1.ILinkListener-missions-cc-0-c-awardLink"><span><span>Получить награду</span></span></a>
	ind += 23
	strOut = lstMissions[ind]
	if !strings.Contains(strOut, `>Получить награду<`) {
		return
	}
	lstLink := strings.Split(strOut, `<a class="simple-but border" href="`)
	strLink := lstLink[1]
	lstLink = strings.Split(strLink, `"><span><span>Получить награду</span></span></a>`)
	strLink = "https://wartank.ru/missions/" + lstLink[0]
	res := арена.Сеть().Get(strLink)
	if res.IsErr() {
		// log._rintf("ERRO Миссии.battleAttack(): при выполнении GET-запроса, err=\n\t%v\n", err)
		return
	}
	арена.СтрОбновить(res.Unwrap())
}

// Проверяет награду за схватку
func победаСхватка(конт ILocalCtx) {
	арена := конт.Get("миссии_простые").Val().(ИАрена)
	var (
		strOut      string
		lstMissions = арена.СписПолучить()
		еслиНайдено bool
		ind         int
	)
	for ind, strOut = range lstMissions {
		if strings.Contains(strOut, `Проведи одну схватку<br/>`) {
			еслиНайдено = true
			break
		}
	}
	if !еслиНайдено {
		return
	}
	// <a class="simple-but border" href="?70-1.ILinkListener-missions-cc-0-c-awardLink"><span><span>Получить награду</span></span></a>
	ind += 23
	strOut = lstMissions[ind]
	if !strings.Contains(strOut, `>Получить награду<`) {
		return
	}
	lstLink := strings.Split(strOut, `<a class="simple-but border" href="`)
	strLink := lstLink[1]
	lstLink = strings.Split(strLink, `"><span><span>Получить награду</span></span></a>`)
	strLink = "https://wartank.ru/missions/" + lstLink[0]
	res := арена.Сеть().Get(strLink)
	if res.IsErr() {
		// log._rintf("ERRO Миссии.победаСхватка(): при выполнении GET-запроса, err=\n\t%v\n", err)
		return
	}
	арена.СтрОбновить(res.Unwrap())
}

// Проверяет награду за ресурсы
func сделать10ресурсов(конт ILocalCtx) {
	арена := конт.Get("миссии_простые").Val().(ИАрена)
	var (
		strOut      string
		lstMissions = арена.СписПолучить()
		еслиНайдено bool
		ind         int
	)
	for ind, strOut = range lstMissions {
		if strings.Contains(strOut, `Произведи на базе 10 ресурсов<br/>`) {
			еслиНайдено = true
			break
		}
	}
	if !еслиНайдено {
		return
	}
	// <a class="simple-but border" href="?70-1.ILinkListener-missions-cc-0-c-awardLink"><span><span>Получить награду</span></span></a>
	ind += 23
	strOut = lstMissions[ind]
	if !strings.Contains(strOut, `>Получить награду<`) {
		return
	}
	lstLink := strings.Split(strOut, `<a class="simple-but border" href="`)
	strLink := lstLink[1]
	lstLink = strings.Split(strLink, `"><span><span>Получить награду</span></span></a>`)
	strLink = "https://wartank.ru/missions/" + lstLink[0]
	res := арена.Сеть().Get(strLink)
	if res.IsErr() {
		// log._rintf("ERRO Миссии.сделать10ресурсов(): при выполнении GET-запроса, err=\n\t%v\n", err)
		return
	}
	арена.СтрОбновить(res.Unwrap())
}

// Проверяет награду за ресурсы
func upMan(конт ILocalCtx) {
	арена := конт.Get("миссии_простые").Val().(ИАрена)
	var (
		strOut      string
		lstMissions = арена.СписПолучить()
		еслиНайдено bool
		ind         int
	)
	for ind, strOut = range lstMissions {
		if strings.Contains(strOut, `Прокачай экипаж на 100 опыта экипажа<br/>`) {
			еслиНайдено = true
			break
		}
	}
	if !еслиНайдено {
		return
	}
	// <a class="simple-but border" href="?70-1.ILinkListener-missions-cc-0-c-awardLink"><span><span>Получить награду</span></span></a>
	ind += 23
	strOut = lstMissions[ind]
	if !strings.Contains(strOut, `>Получить награду<`) {
		return
	}
	lstLink := strings.Split(strOut, `<a class="simple-but border" href="`)
	strLink := lstLink[1]
	lstLink = strings.Split(strLink, `"><span><span>Получить награду</span></span></a>`)
	strLink = "https://wartank.ru/missions/" + lstLink[0]
	res := арена.Сеть().Get(strLink)
	if res.IsErr() {
		// log._rintf("ERRO Миссии.upMan(): при выполнении GET-запроса, err=\n\t%v\n", err)
		return
	}
	арена.СтрОбновить(res.Unwrap())
	// log._rintf("INFO Миссии.upMan(): награда получена\n")
}

// Проверяет награду за топливо
func топливоДив(конт ILocalCtx) {
	арена := конт.Get("миссии_простые").Val().(ИАрена)
	var (
		strOut      string
		lstMissions = арена.СписПолучить()
		еслиНайдено bool
		ind         int
	)
	for ind, strOut = range lstMissions {
		if strings.Contains(strOut, `Получи топливо в дивизии<br/>`) {
			еслиНайдено = true
			break
		}
	}
	if !еслиНайдено {
		return
	}
	// <a class="simple-but border" href="?157-1.ILinkListener-missions-cc-0-c-awardLink"><span><span>Получить награду</span></span></a>
	ind += 19
	strOut = lstMissions[ind]
	if !strings.Contains(strOut, `>Получить награду</`) {
		return
	}
	lstLink := strings.Split(strOut, `<a class="simple-but border" href="`)
	strLink := lstLink[1]
	lstLink = strings.Split(strLink, `"><span><span>Получить награду</span></span></a>`)
	strLink = "https://wartank.ru/missions/" + lstLink[0]
	res := арена.Сеть().Get(strLink)
	if res.IsErr() {
		// log._rintf("ERRO Миссии.makeFuel(): при выполнении GET-запроса, err=\n\t%v\n", err)
		return
	}
	арена.СтрОбновить(res.Unwrap())
}

// Проверяет награду 5 боёв
func подряд5побед(конт ILocalCtx) {
	арена := конт.Get("миссии_простые").Val().(ИАрена)
	var (
		strOut      string
		lstMissions = арена.СписПолучить()
		еслиНайдено bool
		ind         int
	)
	for ind, strOut = range lstMissions {
		if strings.Contains(strOut, `Проведи 5 боев<br/>`) {
			еслиНайдено = true
			break
		}
	}
	if !еслиНайдено {
		return
	}
	// <a class="simple-but border" href="?113-1.ILinkListener-missions-cc-0-c-awardLink"><span><span>Получить награду</span></span></a>
	ind += 23
	strOut = lstMissions[ind]
	if !strings.Contains(strOut, `>Получить награду<`) {
		return
	}
	lstLink := strings.Split(strOut, `<a class="simple-but border" href="`)
	strLink := lstLink[1]
	lstLink = strings.Split(strLink, `"><span><span>Получить награду</span></span></a>`)
	strLink = "https://wartank.ru/missions/" + lstLink[0]
	res := арена.Сеть().Get(strLink)
	if res.IsErr() {
		// log._rintf("ERRO Миссии.подряд5побед(): при выполнении GET-запроса, err=\n\t%v\n", err)
		return
	}
	арена.СтрОбновить(res.Unwrap())
}

// Проверяет награду 10 боёв
func подряд10побед(конт ILocalCtx) {
	арена := конт.Get("миссии_простые").Val().(ИАрена)
	var (
		strOut      string
		lstMissions = арена.СписПолучить()
		еслиНайдено bool
		ind         int
	)
	for ind, strOut = range lstMissions {
		if strings.Contains(strOut, `Уничтожь в бою 10 танков<br/>`) {
			еслиНайдено = true
			break
		}
	}
	if !еслиНайдено {
		return
	}
	// <a class="simple-but border" href="?113-1.ILinkListener-missions-cc-0-c-awardLink"><span><span>Получить награду</span></span></a>
	ind += 23
	strOut = lstMissions[ind]
	if !strings.Contains(strOut, `>Получить награду<`) {
		return
	}
	lstLink := strings.Split(strOut, `<a class="simple-but border" href="`)
	strLink := lstLink[1]
	lstLink = strings.Split(strLink, `"><span><span>Получить награду</span></span></a>`)
	strLink = "https://wartank.ru/missions/" + lstLink[0]
	res := арена.Сеть().Get(strLink)
	if res.IsErr() {
		// log._rintf("ERRO Миссии.подряд10побед(): при выполнении GET-запроса, err=\n\t%v\n", err)
		return
	}
	арена.СтрОбновить(res.Unwrap())
}

// Проверяет награду за 6 побед подряд
func подряд6побед(конт ILocalCtx) {
	арена := конт.Get("миссии_простые").Val().(ИАрена)
	var (
		strOut      string
		lstMissions = арена.СписПолучить()
		еслиНайдено bool
		ind         int
	)
	for ind, strOut = range lstMissions {
		if strings.Contains(strOut, `Победи 6 раз подряд<br/>`) {
			еслиНайдено = true
			break
		}
	}
	if !еслиНайдено {
		return
	}
	// <a class="simple-but border" href="?113-1.ILinkListener-missions-cc-0-c-awardLink"><span><span>Получить награду</span></span></a>
	ind += 25
	strOut = lstMissions[ind]
	if !strings.Contains(strOut, `>Получить награду<`) {
		return
	}
	lstLink := strings.Split(strOut, `<a class="simple-but border" href="`)
	strLink := lstLink[1]
	lstLink = strings.Split(strLink, `"><span><span>Получить награду</span></span></a>`)
	strLink = "https://wartank.ru/missions/" + lstLink[0]
	res := арена.Сеть().Get(strLink)
	if res.IsErr() {
		// log._rintf("ERRO Миссии.подряд6побед(): при выполнении GET-запроса, err=\n\t%v\n", err)
		return
	}
	арена.СтрОбновить(res.Unwrap())
}
