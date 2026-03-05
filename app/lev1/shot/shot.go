package shot

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	. "wartank/app/lev0/types"
	"wartank/app/lev1/shot/damage"
	"wartank/app/lev1/shot/shot_time"

	// "wartank/internal/components/sound"
	"wartank/app/lev0/alias"

	. "gitp78su.ipnodns.ru/svi/kern"
	. "gitp78su.ipnodns.ru/svi/kern/kc/helpers"
	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"
)

/*
	Исходник предоставляет выстрел со свойствами:
		- время до выстрела
		- длительность перезарядки

	Первый параметр постоянно изменяется (после выстрела восстанавливается)
	Второй параметр меняется медленно (в зависимости от количества очков после выстрела)
*/

// выстрел -- объект выстрела
type выстрел struct {
	ИСражениеПроцесс                     // FIXME:
	перезарядка      *shot_time.ShotTime // Сколько времени нужно для полной перезарядки
	урон             *damage.Damage      // Урон от выстрела с памятью
	уронВсего        alias.АУрон         // Суммарный урон
	isEnd            ISafeBool           // Признак конца сражения
	еслиБлок         ISafeBool           // Признак блокировки выстрела
	логин            string              // Логин для поиска контрольных строк
	chTick           chan int            // Тик для выстрела
	промаховПодряд   int                 // Счётчик промахов поиска кнопки атаки
}

// НовВыстрел -- возвращает новый выстрел
func НовВыстрел(проц ИСражениеПроцесс) ИВыстрел {
	Hassert(проц != nil, "НовВыстрел(): ИСражениеПроцесс == nil")
	логинТанк := проц.Бот().Имя()
	сам := &выстрел{
		ИСражениеПроцесс: проц,
		перезарядка:      shot_time.NewShotTime(),
		урон:             damage.NewDamage(),
		еслиБлок:         NewSafeBool(),
		isEnd:            проц.ЕслиКонец(),
		логин:            логинТанк,
		chTick:           make(chan int, 2),
	}
	// Атака каждые 5 секунд
	сам.перезарядка.Set(5000)
	go сам.makeTick()
	go сам.run()
	return сам
}

// БлокУст -- установка блокировки выстрела
func (сам *выстрел) БлокУст() {
	сам.еслиБлок.Set()
}

// БлокСброс -- сброс блокировки выстрела
func (сам *выстрел) БлокСброс() {
	сам.еслиБлок.Reset()
}

// ЕслиБлок -- возвращает признак блокировки выстрела
func (сам *выстрел) ЕслиБлок() bool {
	return сам.еслиБлок.Get()
}

// Генерирует тики для атаки каждые 5 секунд
func (сам *выстрел) makeTick() {
	defer func() {
		сам.isEnd.Set()
		close(сам.chTick)
		сам.Отменить()
	}()
	for {
		select {
		case <-сам.Контекст().Done():
			return
		default:
			if сам.isEnd.Get() {
				return
			}
			сам.chTick <- 1
			time.Sleep(time.Millisecond * time.Duration(сам.перезарядка.Get()))
		}
	}
}

// Цикл выстрела (в отдельном потоке)
func (сам *выстрел) run() {
	fmt.Println("выстрел: run() горутина стартовала")
	for {
		select {
		case <-сам.Контекст().Done():
			fmt.Println("выстрел: run() контекст отменён, выходим")
			return
		case <-сам.chTick:
			сам.shot()
			сам.findDamage()
		}
	}
}

// Обновляет возможность выстрела (~)
//
//	Вызывается из отдельного потока
func (сам *выстрел) shot() {
	fmt.Printf("выстрел: shot() вызван, isEnd=%v\n", сам.isEnd.Get())
	сам.Обновить()
	var (
		strOut      string
		lstBattle   = сам.СписПолучить()
		еслиНайдено bool
	)
	fmt.Printf("выстрел: после обновления, строк=%d\n", len(lstBattle))
	// <a href="pve?6-26.ILinkListener-currentControl-attackRegularShellLink" class="simple-but gray"><span><span>ОБЫЧНЫЕ</span></span></a>
	for _, strOut = range lstBattle {
		if strings.Contains(strOut, `-currentControl-attackRegularShellLink`) {
			еслиНайдено = true
			break
		}
	}
	if !еслиНайдено {
		сам.промаховПодряд++
		fmt.Printf("выстрел: кнопка атаки не найдена (%d/%d), строк=%d\n",
			сам.промаховПодряд, 30, len(lstBattle))
		if сам.промаховПодряд >= 30 { // ~3-4 минуты без кнопки — бой точно кончился
			fmt.Println("выстрел: бой окончен — кнопка атаки не появилась")
			сам.isEnd.Set()
			сам.Отменить()
		}
		return // Не убиваем бой — дадим следующему тику попробовать снова
	}
	сам.промаховПодряд = 0 // Сбросить счётчик — кнопка найдена
	fmt.Println("выстрел: кнопка НАЙДЕНА, стреляем!")
	strLink := strings.TrimPrefix(strOut, `<a href="`)
	strLink = strings.TrimSuffix(strLink, `" class="simple-but gray"><span><span>ОБЫЧНЫЕ</span></span></a>`)
	strLink = "https://wartank.ru/" + strLink
	res := сам.Сеть().Get(strLink)
	if res.IsErr() {
		fmt.Println("выстрел: ошибка HTTP при выстреле, пропускаем")
		return // Не убиваем бой — в следующий тик попробуем снова
	}
	lstBattle = res.Unwrap()
	сам.СтрОбновить(lstBattle)
}

// Ищет урон от выстрела (без адаптивного тайминга — фиксированные 5 сек)
func (сам *выстрел) findDamage() {
	var (
		ind         int
		еслиНайдено bool
		lstShot     = сам.СписПолучить()
		strOut      string
	)

	for ind, strOut = range lstShot {
		if strings.Contains(strOut, `<span class="yellow1 td_u">`+сам.логин+`</span>`) {
			ind += 2
			if ind >= len(lstShot) {
				return
			}
			strOut = lstShot[ind]
			еслиНайдено = true
			break
		}
	}
	if !еслиНайдено {
		return
	}
	lstDamage := strings.Split(strOut, ` на  <span class="red1">`)
	if len(lstDamage) < 2 {
		return
	}
	strDamage := lstDamage[1]
	iDamage, err := strconv.Atoi(strDamage)
	if err != nil {
		return
	}
	if iDamage <= 0 {
		iDamage = 0
	}
	сам.уронВсего += alias.АУрон(iDamage)
}

// setDamage -- обновляет время перезарядки в зависимости от произведённого урона
func (сам *выстрел) setDamage(val alias.АУрон) {
	сам.урон.Set(val)
	switch сам.урон.Result() {
	case "none":
		сам.перезарядка.Dec5()
	case "up":
		сам.перезарядка.Dec30()
	case "down":
		сам.перезарядка.Inc210()
	}
}

// IsEnd -- возвращает объект разрешения стрельбы
func (сам *выстрел) IsEnd() ISafeBool {
	return сам.isEnd
}
