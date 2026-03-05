// package bf_ammo_stat -- бизнес-функция статистики снарядов
package bf_ammo_stat

import (
	"strconv"
	"strings"

	. "gitp78su.ipnodns.ru/svi/kern/kc/helpers"
	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"

	"wartank/app/lev0/cons"
	. "wartank/app/lev0/types"
)

// СнарядыСтат -- получает статистику снарядов
func СнарядыСтат(конт ILocalCtx) {
	арсенал_ := конт.Get("арсенал")
	if арсенал_ == nil { // Может быть не построен
		return
	}
	арсенал := арсенал_.Val().(ИАренаАрсенал)

	if арсенал.Состояние().Получ() == cons.РежимНеСуществует {
		return
	}
	фугасыНайти(конт)
	бронейбойкиНайти(конт)
	кумульНайти(конт)
	ремкаНайти(конт)
}

func ремкаНайти(конт ILocalCtx) {
	арсенал := конт.Get("арсенал").Val().(ИАренаАрсенал)
	lstArsenal := арсенал.СписПолучить()
	if len(lstArsenal) == 0 {
		арсенал.Обновить()
		lstArsenal = арсенал.СписПолучить()
	}
	strOut := ""
	isFind := false
	// <span class="nwr"><img class="rico vm" src="/images/shells/repairkit.gif"/> 282</span>
	for _, стрСнаряд := range lstArsenal {
		if strings.Contains(стрСнаряд, `<span class="nwr"><img class="rico vm" src="/images/shells/repairkit.gif"/> `) {
			strOut = стрСнаряд
			isFind = true
			break
		}
	}
	Hassert(isFind, "ремкаНайти(): не найдена контрольная строка")
	strOut = strings.TrimPrefix(strOut, `<span class="nwr"><img class="rico vm" src="/images/shells/repairkit.gif"/> `)
	strOut = strings.TrimSuffix(strOut, `</span>`)
	целФугас, err := strconv.Atoi(strOut)
	Hassert(err == nil, "ремкаНайти(): strOut(%v), err=\n\t%v", strOut, err)
	арсенал.Ремки().Уст(целФугас)
}

func кумульНайти(конт ILocalCtx) {
	арсенал := конт.Get("арсенал").Val().(ИАренаАрсенал)
	lstArsenal := арсенал.СписПолучить()
	if len(lstArsenal) == 0 {
		арсенал.Обновить()
		lstArsenal = арсенал.СписПолучить()
	}
	strOut := ""
	isFind := false
	// <span class="nwr"><img class="rico vm" src="/images/shells/HollowCharge.png" alt="Кумулятивный снаряд" title="Кумулятивный снаряд"/> 7340 &nbsp;&nbsp;</span>
	for _, стрСнаряд := range lstArsenal {
		if strings.Contains(стрСнаряд, `<span class="nwr"><img class="rico vm" src="/images/shells/HollowCharge.png" alt="Кумулятивный снаряд" title="Кумулятивный снаряд"/> `) {
			strOut = стрСнаряд
			isFind = true
			break
		}
	}
	Hassert(isFind, "кумульНайти(): не найдена контрольная строка")
	strOut = strings.TrimPrefix(strOut, `<span class="nwr"><img class="rico vm" src="/images/shells/HollowCharge.png" alt="Кумулятивный снаряд" title="Кумулятивный снаряд"/> `)
	strOut = strings.TrimSuffix(strOut, ` &nbsp;&nbsp;</span>`)
	целФугас, err := strconv.Atoi(strOut)
	Hassert(err == nil, "кумульНайти(): strOut(%v), err=\n\t%v", strOut, err)
	арсенал.Кумулятивы().Уст(целФугас)
}

func бронейбойкиНайти(конт ILocalCtx) {
	арсенал := конт.Get("арсенал").Val().(ИАренаАрсенал)
	lstArsenal := арсенал.СписПолучить()
	if len(lstArsenal) == 0 {
		арсенал.Обновить()
		lstArsenal = арсенал.СписПолучить()
	}
	strOut := ""
	isFind := false
	// <span class="nwr"><img class="rico vm" src="/images/shells/ArmorPiercing.png" alt="Бронебойный снаряд" title="Бронебойный снаряд"/> 7335 &nbsp;&nbsp;</span>
	for _, стрСнаряд := range lstArsenal {
		if strings.Contains(стрСнаряд, `<span class="nwr"><img class="rico vm" src="/images/shells/ArmorPiercing.png" alt="Бронебойный снаряд" title="Бронебойный снаряд"/> `) {
			strOut = стрСнаряд
			isFind = true
			break
		}
	}
	Hassert(isFind, "бронейбойкиНайти(): не найдена контрольная строка")
	strOut = strings.TrimPrefix(strOut, `<span class="nwr"><img class="rico vm" src="/images/shells/ArmorPiercing.png" alt="Бронебойный снаряд" title="Бронебойный снаряд"/> `)
	strOut = strings.TrimSuffix(strOut, ` &nbsp;&nbsp;</span>`)
	целФугас, err := strconv.Atoi(strOut)
	Hassert(err == nil, "бронейбойкиНайти(): strOut(%v), err=\n\t%v", strOut, err)
	арсенал.Бронебойки().Уст(целФугас)
}

func фугасыНайти(конт ILocalCtx) {
	арсенал := конт.Get("арсенал").Val().(ИАренаАрсенал)
	lstArsenal := арсенал.СписПолучить()
	if len(lstArsenal) == 0 {
		арсенал.Обновить()
		lstArsenal = арсенал.СписПолучить()
	}
	strOut := ""
	isFind := false
	// <span class="nwr"><img class="rico vm" src="/images/shells/HighExplosive.png" alt="Фугасный снаряд" title="Фугасный снаряд"/> 7343 &nbsp;&nbsp;</span>
	for _, стрФугас := range lstArsenal {
		if strings.Contains(стрФугас, `<span class="nwr"><img class="rico vm" src="/images/shells/HighExplosive.png" alt="Фугасный снаряд" title="Фугасный снаряд"/> `) {
			strOut = стрФугас
			isFind = true
			break
		}
	}
	Hassert(isFind, "фугасыНайти(): не найдена контрольная строка")
	strOut = strings.TrimPrefix(strOut, `<span class="nwr"><img class="rico vm" src="/images/shells/HighExplosive.png" alt="Фугасный снаряд" title="Фугасный снаряд"/> `)
	strOut = strings.TrimSuffix(strOut, ` &nbsp;&nbsp;</span>`)
	целФугас, err := strconv.Atoi(strOut)
	Hassert(err == nil, "фугасыНайти(): strOut(%v), err=\n\t%v", strOut, err)
	арсенал.Фугасы().Уст(целФугас)
}
