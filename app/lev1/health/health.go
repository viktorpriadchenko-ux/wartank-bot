package health

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	. "wartank/app/lev0/types"
	"wartank/app/lev1/health/health_time"
	"wartank/app/lev1/repair_time"

	// "wartank/internal/components/sound"
	. "gitp78su.ipnodns.ru/svi/kern"
	. "gitp78su.ipnodns.ru/svi/kern/kc/helpers"
	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"
)

/*
	Контролирует состояние здоровья танка
*/

// Здоровье -- контроль здоровья танка
type Здоровье struct {
	ИСражениеПроцесс                         // FIXME:
	здоровьеСейчас   *health_time.HealthTime // Изменяемое здоровье танка
	здоровьеПолное   *health_time.HealthTime // Полное здоровье танка
	еслиНадо         ISafeBool               // Необходимость восстановления
	отсчётАптечка    *repair_time.RepairTime // Время до восстановления
	isEnd            ISafeBool               // Ссылка на признак конца сражения
	логин            string                  // Для поиска контрольных строк
	chTick           chan int                // Канал для ровной отправки тиков
	промаховЛогин    int                     // Счётчик неудач поиска логина
}

// НовЗдоровье -- возвращает новый *Health
func НовЗдоровье(проц ИСражениеПроцесс) *Здоровье {
	Hassert(проц != nil, "НовЗдоровье(): ИСражениеПроцесс == nil")
	логин := проц.Бот().Имя()
	сам := &Здоровье{
		ИСражениеПроцесс: проц,
		здоровьеСейчас:   health_time.NewHealthTime(),
		здоровьеПолное:   health_time.NewHealthTime(),
		еслиНадо:         NewSafeBool(),
		отсчётАптечка:    repair_time.NewRepairTime(),
		isEnd:            проц.ЕслиКонец(),
		логин:            логин,
		chTick:           make(chan int, 2),
	}
	go сам.makeTik()
	go сам.run()
	return сам
}

// Отправляет тики с заданным равным интервалом
func (сам *Здоровье) makeTik() {
	defer func() {
		close(сам.chTick)
		сам.Отменить()
	}()
	лимитАптечка := 0 // Предел времени ожидания
	отсчётАптечка := 0
	for {
		select {
		case <-сам.Контекст().Done():
			return
		default:
			if сам.ЕслиУбит() {
				return
			}
			if сам.отсчётАптечка.Получ() == отсчётАптечка {
				лимитАптечка++
			} else {
				отсчётАптечка = сам.отсчётАптечка.Получ()
				лимитАптечка = 0
			}
			if сам.отсчётАптечка.IsReady() {
				лимитАптечка = 0
			}
			if лимитАптечка > 90 {
				return
			}
		}
		сам.chTick <- 1
		time.Sleep(time.Second * 1)
		сам.отсчётАптечка.Dec()
	}
}

// Главный цикл обработки здоровья в сражении
func (сам *Здоровье) run() {
	for {
		select {
		case <-сам.Контекст().Done():
			сам.isEnd.Set()
			return
		case <-сам.chTick:
			сам.здоровьеНайти()
			сам.найтиВремяВосстановления()
			if сам.еслиНадо.Get() {
				go сам.repair()
			}
		}
	}
}

// Полное -- возвращает объект полного здоровья танка
func (сам *Здоровье) Полное() int {
	return сам.здоровьеПолное.Get()
}

// ЕслиУбит -- возвращает признак мертвичины танка
func (сам *Здоровье) ЕслиУбит() bool {
	if сам.isEnd.Get() {
		сам.Отменить()
		return true
	}
	lstBattle := сам.СписПолучить()
	for _, strOut := range lstBattle {
		if strings.Contains(strOut, `>Ваш танк подбит.`) {
			сам.здоровьеСейчас.Set(0)
			сам.isEnd.Set()
			сам.Отменить()
			return true
		}
	}
	return сам.isEnd.Get()
}

// Ищет время восстановления ремки
func (сам *Здоровье) найтиВремяВосстановления() {
	if сам.отсчётАптечка.IsReady() {
		return
	}
	var (
		strOut      string
		lstBattle   = сам.СписПолучить()
		еслиНайдено bool
		ind         int
	)
	// <a href="pve?19-14.ILinkListener-currentControl-repairLink" class="simple-but blue"><span><span>12 секунд</span></span></a>
	//
	for ind, strOut = range lstBattle {
		if strings.Contains(strOut, `ILinkListener-currentControl-repairLink" class="simple-but blue"><span><span>`) {
			еслиНайдено = true
			break
		}
	}
	if !еслиНайдено {
		return
	}
	strOut = lstBattle[ind]
	// <a href="pve?19-14.ILinkListener-currentControl-repairLink" class="simple-but blue"><span><span>12 секунд</span></span></a>
	lstTime := strings.Split(strOut, `ILinkListener-currentControl-repairLink" class="simple-but blue"><span><span>`)
	if len(lstTime) < 2 {
		return
	}
	strTime := lstTime[1]
	strTime = strings.TrimSuffix(strTime, ` секунд</span></span></a>`)
	if err := сам.отсчётАптечка.Уст(strTime); err != nil {
		return
	}
}

// Восстанавливает здоровье (~)
func (сам *Здоровье) repair() {
	var (
		strOut            string
		lstBattleOn       = сам.СписПолучить()
		еслиНайденоRepair bool
	)
	// <span>Ремкомплект</span>
	// <a href="pve?19-14.ILinkListener-currentControl-repairLink" class="simple-but blue"><span><span>Ремкомплект</span></span></a>
	for _, strOut = range lstBattleOn {
		if strings.Contains(strOut, `<span>Ремкомплект</span>`) {
			еслиНайденоRepair = true
			break
		}
	}
	if !еслиНайденоRepair {
		return
	}
	lstLink := strings.Split(strOut, `<a href="`)
	if len(lstLink) < 2 {
		return
	}
	strLink := lstLink[1]
	lstLink = strings.Split(strLink, `" class="simple-but blue"><span><span>Ремкомплект</span></span></a>`)
	strLink = "https://wartank.ru/" + lstLink[0]
	res := сам.Сеть().Get(strLink)
	if res.IsErr() {
		fmt.Println("ремонт: ошибка HTTP, пропускаем")
		return
	}
	lstBattleOn = res.Unwrap()
	сам.СтрОбновить(lstBattleOn)
	// sound.Repair()
}

// Ищет своё здоровье (~)
func (сам *Здоровье) здоровьеНайти() {
	var (
		ind         int
		strOut      string
		еслиНайдено bool
		lstBattle   = сам.СписПолучить()
	)
	if len(lstBattle) == 0 { // Принудительно обновим сражение
		сам.Обновить()
		lstBattle = сам.СписПолучить()
	}
	// <div class="small bold green1 sh_b mb10 mt5">Половина коня</div>
	for ind, strOut = range lstBattle {
		if strings.Contains(strOut, `<div class="small bold green1 sh_b mb10 mt5">`+сам.логин+`"`) {
			еслиНайдено = true
			break
		}
	}
	if !еслиНайдено {
		сам.промаховЛогин++
		fmt.Printf("здоровье: логин '%s' не найден (%d/30), строк=%d\n", сам.логин, сам.промаховЛогин, len(lstBattle))
		if len(lstBattle) > 0 {
			// Показать первые 3 строки для диагностики
			for i := 0; i < 3 && i < len(lstBattle); i++ {
				if len(lstBattle[i]) > 120 {
					fmt.Printf("  строка[%d]: %s...\n", i, lstBattle[i][:120])
				} else {
					fmt.Printf("  строка[%d]: %s\n", i, lstBattle[i])
				}
			}
		}
		if сам.промаховЛогин >= 30 { // ~30 секунд без логина — бой точно кончился
			fmt.Println("здоровье: бой окончен — логин не найден после 30 попыток")
			сам.isEnd.Set()
			сам.Отменить()
		}
		return
	}
	сам.промаховЛогин = 0 // Сброс — логин найден
	fmt.Printf("здоровье: логин найден, ind=%d\n", ind)
	ind += 11
	if ind >= len(lstBattle) {
		fmt.Printf("здоровье: ind+11=%d >= len=%d, пропускаем\n", ind, len(lstBattle))
		return
	}
	strOut = lstBattle[ind]
	strHealth := strings.TrimPrefix(strOut, `<div class="value-block lh1"><span><span>`)
	strHealth = strings.TrimSuffix(strHealth, `</span></span></div>`)
	iHealth, err := strconv.Atoi(strHealth)
	if err != nil {
		fmt.Printf("здоровье: ошибка Atoi для '%s', пропускаем (не убиваем бой)\n", strHealth)
		return // Не убиваем бой — может быть временный глюк HTML
	}
	сам.здоровьеУстановить(iHealth)
}

// здоровьеУстановить -- устанавливает текущее здоровье
func (сам *Здоровье) здоровьеУстановить(здоровье int) {
	if здоровье < 0 {
		здоровье = 0
	}
	дельта := сам.здоровьеСейчас.Get() - здоровье
	if дельта < 0 {
		дельта = 0
	}
	// Обновляем полное здоровье при первом чтении или после лечения
	if здоровье >= сам.здоровьеПолное.Get() {
		сам.здоровьеПолное.Set(здоровье)
		сам.здоровьеСейчас.Set(здоровье)
		сам.еслиНадо.Reset()
		return
	}
	// Всегда обновляем текущее здоровье
	сам.здоровьеСейчас.Set(здоровье)

	switch {
	case сам.isEnd.Get():
		сам.здоровьеСейчас.Set(0)
		сам.isEnd.Set()
		сам.Отменить()
		return
	case здоровье <= 0:
		сам.здоровьеСейчас.Set(0)
		сам.isEnd.Set()
		сам.Отменить()
		return
	default:
		// Ремкомплект при HP < 15%
		порог := сам.здоровьеПолное.Get() * 15 / 100
		if порог <= 0 {
			порог = 1
		}
		if здоровье <= порог {
			сам.еслиНадо.Set()
			fmt.Printf("здоровье: HP=%d/%d (<15%%), нужен ремкомплект\n", здоровье, сам.здоровьеПолное.Get())
		} else {
			сам.еслиНадо.Reset()
		}
		// Получили урон — маневрируем
		if дельта > 0 {
			сам.Манёвр().УстНадо()
		}
	}
}
