// package bf_tank_stat -- бизнес-функция поиска статы танка
package bf_tank_stat

import (
	"strconv"
	"strings"

	. "gitp78su.ipnodns.ru/svi/kern"
	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"

	. "wartank/app/lev0/types"
)

// ТанкСтатПолучить -- получает стату танка
func ТанкСтатПолучить(конт ILocalCtx) {
	найтиАтаку(конт)
	найтиБроню(конт)
	найтиТочность(конт)
	найтиПрочность(конт)
	найтиМощь(конт)
}

func найтиМощь(конт ILocalCtx) {
	ангар := конт.Get("ангар").Val().(ИАренаАнгар)
	списАнгар := ангар.СписПолучить()
	if len(списАнгар) == 0 {
		ангар.Обновить()
		списАнгар = ангар.СписПолучить()
	}
	var (
		стрМощь  string
		еслиЕсть bool
	)
	// alt="Точность" title="Точность"/> Точность <span class="green2">
	for _, стр := range списАнгар {
		if strings.Contains(стр, `<img src="/images/icons/power.png?2" height="14" width="14"> <span class="green2">Танковая мощь: `) {
			стрМощь = стр
			еслиЕсть = true
			break
		}
	}
	if !еслиЕсть {
		NewLogBuf().Warn("найтиМощь(): не найдена строка мощи — пропускаем\n")
		return
	}
	// Вырезать ссылку на атаку
	// <img src="/images/icons/power.png?2" height="14" width="14"> <span class="green2">Танковая мощь: 3115</span>
	стрМощь = strings.TrimPrefix(стрМощь, `<img src="/images/icons/power.png?2" height="14" width="14"> <span class="green2">Танковая мощь: `)
	стрМощь = strings.TrimSuffix(стрМощь, `</span>`)
	if стрМощь == "" {
		NewLogBuf().Warn("найтиМощь(): строка мощи пустая — пропускаем\n")
		return
	}
	цМощь, ош := strconv.Atoi(стрМощь)
	if ош != nil {
		NewLogBuf().Warn("найтиМощь(): стрМощь(%v) не число, ош=%v\n", стрМощь, ош)
		return
	}
	танкСтата := конт.Get("танкСтат").Val().(ИТанкСтат)
	танкСтата.Мощь().Уст(цМощь)
}

func найтиПрочность(конт ILocalCtx) {
	ангар := конт.Get("ангар").Val().(ИАренаАнгар)
	списАнгар := ангар.СписПолучить()
	if len(списАнгар) == 0 {
		ангар.Обновить()
		списАнгар = ангар.СписПолучить()
	}
	var (
		стрПрочность string
		еслиЕсть     bool
	)
	// alt="Точность" title="Точность"/> Точность <span class="green2">
	for _, стр := range списАнгар {
		if strings.Contains(стр, `alt="Прочность" title="Прочность"/> Прочность <span class="green2">`) {
			стрПрочность = стр
			еслиЕсть = true
			break
		}
	}
	if !еслиЕсть {
		NewLogBuf().Warn("найтиПрочность(): не найдена строка прочности — пропускаем\n")
		return
	}
	// Вырезать ссылку на атаку
	// <img width="14" height="14" src="/images/icons/durability.png?1" alt="Прочность" title="Прочность"/> Прочность <span class="green2">797</span><br/>
	стрПрочность = strings.TrimPrefix(стрПрочность, `<img width="14" height="14" src="/images/icons/durability.png?1" alt="Прочность" title="Прочность"/> Прочность <span class="green2">`)
	стрПрочность = strings.TrimSuffix(стрПрочность, `</span><br/>`)
	if стрПрочность == "" {
		NewLogBuf().Warn("найтиПрочность(): строка прочности пустая — пропускаем\n")
		return
	}
	цПрочность, ош := strconv.Atoi(стрПрочность)
	if ош != nil {
		NewLogBuf().Warn("найтиПрочность(): стрПрочность(%v) не число, ош=%v\n", стрПрочность, ош)
		return
	}
	танкСтата := конт.Get("танкСтат").Val().(ИТанкСтат)
	танкСтата.Прочность().Уст(цПрочность)
}

func найтиТочность(конт ILocalCtx) {
	ангар := конт.Get("ангар").Val().(ИАренаАнгар)
	списАнгар := ангар.СписПолучить()
	if len(списАнгар) == 0 {
		ангар.Обновить()
		списАнгар = ангар.СписПолучить()
	}
	var (
		стрТочность string
		еслиЕсть    bool
	)
	// alt="Точность" title="Точность"/> Точность <span class="green2">
	for _, стр := range списАнгар {
		if strings.Contains(стр, `alt="Точность" title="Точность"/> Точность <span class="green2">`) {
			стрТочность = стр
			еслиЕсть = true
			break
		}
	}
	if !еслиЕсть {
		NewLogBuf().Warn("найтиТочность(): не найдена строка точности — пропускаем\n")
		return
	}
	// Вырезать ссылку на атаку
	// <img width="14" height="14" src="/images/icons/accuracy.png?1" alt="Точность" title="Точность"/> Точность <span class="green2">833</span><br/>
	стрТочность = strings.TrimPrefix(стрТочность, `<img width="14" height="14" src="/images/icons/accuracy.png?1" alt="Точность" title="Точность"/> Точность <span class="green2">`)
	стрТочность = strings.TrimSuffix(стрТочность, `</span><br/>`)
	if стрТочность == "" {
		NewLogBuf().Warn("найтиТочность(): строка точности пустая — пропускаем\n")
		return
	}
	цТочность, ош := strconv.Atoi(стрТочность)
	if ош != nil {
		NewLogBuf().Warn("найтиТочность(): стрТочность(%v) не число, ош=%v\n", стрТочность, ош)
		return
	}
	танкСтата := конт.Get("танкСтат").Val().(ИТанкСтат)
	танкСтата.Точность().Уст(цТочность)
}

func найтиБроню(конт ILocalCtx) {
	ангар := конт.Get("ангар").Val().(ИАренаАнгар)
	списАнгар := ангар.СписПолучить()
	if len(списАнгар) == 0 {
		ангар.Обновить()
		списАнгар = ангар.СписПолучить()
	}
	var (
		стрБроня string
		еслиЕсть bool
	)
	// alt="Броня" title="Броня"/> Броня <span class="green2">
	for _, стр := range списАнгар {
		if strings.Contains(стр, `alt="Броня" title="Броня"/> Броня <span class="green2">`) {
			стрБроня = стр
			еслиЕсть = true
			break
		}
	}
	if !еслиЕсть {
		NewLogBuf().Warn("найтиБроню(): не найдена строка брони — пропускаем\n")
		return
	}
	// Вырезать ссылку на атаку
	// <img width="14" height="14" src="/images/icons/armor.png?1" alt="Броня" title="Броня"/> Броня <span class="green2">787</span><br/>
	стрБроня = strings.TrimPrefix(стрБроня, `<img width="14" height="14" src="/images/icons/armor.png?1" alt="Броня" title="Броня"/> Броня <span class="green2">`)
	стрБроня = strings.TrimSuffix(стрБроня, `</span><br/>`)
	if стрБроня == "" {
		NewLogBuf().Warn("найтиБроню(): строка брони пустая — пропускаем\n")
		return
	}
	цБроня, ош := strconv.Atoi(стрБроня)
	if ош != nil {
		NewLogBuf().Warn("найтиБроню(): стрБроня(%v) не число, ош=%v\n", стрБроня, ош)
		return
	}
	танкСтата := конт.Get("танкСтат").Val().(ИТанкСтат)
	танкСтата.Броня().Уст(цБроня)
}

func найтиАтаку(конт ILocalCtx) {
	ангар := конт.Get("ангар").Val().(ИАренаАнгар)
	списАнгар := ангар.СписПолучить()
	if len(списАнгар) == 0 {
		ангар.Обновить()
		списАнгар = ангар.СписПолучить()
	}
	var (
		стрАтака string
		еслиЕсть bool
	)
	// alt="Атака" title="Атака"/> Атака <span class="green2">
	for _, стр := range списАнгар {
		if strings.Contains(стр, `alt="Атака" title="Атака"/> Атака <span class="green2">`) {
			стрАтака = стр
			еслиЕсть = true
			break
		}
	}
	if !еслиЕсть {
		NewLogBuf().Warn("найтиАтаку(): не найдена строка атаки — пропускаем\n")
		return
	}
	// Вырезать ссылку на атаку
	// <img width="14" height="14" src="/images/icons/attack.png?1" alt="Атака" title="Атака"/> Атака <span class="green2">698</span><br/>
	стрАтака = strings.TrimPrefix(стрАтака, `<img width="14" height="14" src="/images/icons/attack.png?1" alt="Атака" title="Атака"/> Атака <span class="green2">`)
	стрАтака = strings.TrimSuffix(стрАтака, `</span><br/>`)
	if стрАтака == "" {
		NewLogBuf().Warn("найтиАтаку(): строка атаки пустая — пропускаем\n")
		return
	}
	цАтака, ош := strconv.Atoi(стрАтака)
	if ош != nil {
		NewLogBuf().Warn("найтиАтаку(): стрАтака(%v) не число, ош=%v\n", стрАтака, ош)
		return
	}
	танкСтата := конт.Get("танкСтат").Val().(ИТанкСтат)
	танкСтата.Атака().Уст(цАтака)
}
