// package bf_glory_make -- бизнес-функция бой за славу
package bf_glory_make

import (
	"strings"

	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"

	. "wartank/app/lev0/types"
)

// СлаваБой -- забрать славу в бою
func СлаваБой(конт ILocalCtx) {
	стрАтака := атакаПроверить(конт)
	if стрАтака == "" {
		return
	}
	атакаНачать(конт, стрАтака)
}

// Проводит атаку на конвой
func атакаНачать(конт ILocalCtx, strOut string) {
	конвой := конт.Get("конвой").Val().(ИАренаКонвой)
	strLink := strOut
	// Можно начать разведку
	lstConvoy := конвой.Сеть().ВебВоркер().Получ(strLink)
	конвой.СтрОбновить(lstConvoy)
	начатьРазведку(конт)
	конвой.ПродуктВремяСейчас().Set("01")
}

func атакаПроверить(конт ILocalCtx) string {
	конвой := конт.Get("конвой").Val().(ИАренаКонвой)
	var (
		strOut      = ""
		еслиНайдено bool
	)
	lstConvoy := конвой.Сеть().ВебВоркер().Получ("https://wartank.ru/convoy")
	for _, strOut = range lstConvoy {
		// <div class="bot"><a class="simple-but border" w:id="findEnemy" href="convoy?50-1.ILinkListener-root-findEnemy"><span><span>Начать разведку</span></span></a></div>
		если1 := strings.Contains(strOut, `.ILinkListener-root-findEnemy"`)
		if если1 {
			_ссылка := strings.TrimPrefix(strOut, `<div class="bot"><a class="simple-but border" w:id="findEnemy" href="`)
			_ссылка = strings.TrimSuffix(_ссылка, `"><span><span>Начать разведку</span></span></a></div>`)
			strOut = "https://wartank.ru/" + _ссылка
			еслиНайдено = true
			break
		}
		если2 := strings.Contains(strOut, `<span>В БОЙ!</span>`)
		if если2 {
			lstLink := strings.Split(strOut, `<div class="bot"><a class="simple-but border" w:id="startMasking" href="`)
			if len(lstLink) == 1 {
				lstLink = strings.Split(strOut, `<div class="bot"><a class="simple-but border red" w:id="startFight" href="`)
			}
			strOut = lstLink[1]
			lstLink = strings.Split(strOut, `"><span><span>В БОЙ!</span></span></a></div>`)
			strOut = "https://wartank.ru/" + lstLink[0]
			еслиНайдено = true
			break
		}
		// <div class="bot"><a class="simple-but border" w:id="findEnemy" href="convoy?15-1.ILinkListener-root-findEnemy"><span><span>Начать разведку</span></span></a></div>
		если3 := strings.Contains(strOut, "<span>Начать разведку</span>")
		if если3 {
			_ссылка := strings.TrimPrefix(strOut, `<<div class="bot"><a class="simple-but border" w:id="findEnemy" href="`)
			_ссылка = strings.TrimSuffix(_ссылка, `"><span><span>Начать разведку</span></span></a></div>`)
			strOut = "https://wartank.ru/" + _ссылка
			еслиНайдено = true
			break
		}
		if strings.Contains(strOut, `>ОБЫЧНЫЕ<`) {
			lstLink := strings.Split(strOut, `<a href="`)
			strOut = lstLink[1]
			lstLink = strings.Split(strOut, `" class="simple-but gray"><span><span>ОБЫЧНЫЕ</span></span></a>`)
			strOut = "https://wartank.ru/" + lstLink[0]
			еслиНайдено = true
			break
		}
	}
	if !еслиНайдено { // Время ожидания
		return ""
	}
	return strOut
}

// Выполняет атаку на конвой
func начатьРазведку(конт ILocalCtx) {
	конвой := конт.Get("конвой").Val().(ИАренаКонвой)
	// Вырезать ссылку на атаку
	strOut := ""
	еслиНайдено := false
	lstConvoy := конвой.СписПолучить()
	// <div class="bot"><a class="simple-but border" w:id="findEnemy" href="convoy?50-1.ILinkListener-root-findEnemy"><span><span>Начать разведку</span></span></a></div>
	for _, strOut = range lstConvoy {
		if strings.Contains(strOut, `.ILinkListener-root-findEnemy`) {
			еслиНайдено = true
			break
		}
	}
	if !еслиНайдено { // Нечего атаковать
		return
	}
	// Атакуем конвой
	_link := strings.TrimPrefix(strOut, `<div class="bot"><a class="simple-but border" w:id="findEnemy" href="`)
	_link = strings.TrimSuffix(_link, `"><span><span>Начать разведку</span></span></a></div>`)
	// https://wartank.ru/convoy?52-1.ILinkListener-root-findEnemy
	link := "https://wartank.ru/" + _link
	{ // Выполнить атаку
		lstConvoy = конвой.Сеть().ВебВоркер().Получ(link)
		конвой.СтрОбновить(lstConvoy)
	}
}
