package shot

import (
	"fmt"
	// "log"
	"strings"
	"time"

	. "wartank/app/lev0/types"
)

/*
	Исходник предоставляет выстрел со свойствами:
		- время до выстрела
		- длительность перезарядки

	Первый параметр постоянно изменяется (после выстрела восстанавливается)
	Второй параметр меняется медленно (в зависимости от количества очков после выстрела)
*/

// Выстрел -- объект выстрела
type Выстрел struct {
	ИСражениеПроцесс          // FIXME!!!
	канТик           chan int // Тик для выстрела
	выстрелСчёт      int      // Счётчик выстрелов для сервера
}

// НовВыстрел -- возвращает новый *Shot
func НовВыстрел(проц ИСражениеПроцесс) (*Выстрел, error) {
	{ // Предусловия
		if проц == nil {
			return nil, fmt.Errorf("НовВыстрел(): действие==nil")
		}
	}
	сам := &Выстрел{
		ИСражениеПроцесс: проц,
		выстрелСчёт:      1,
		канТик:           make(chan int, 2),
	}
	go сам.делайТик()
	go сам.пуск()
	return сам, nil
}

// Генерирует тики, когда можно стрелять
func (сам *Выстрел) делайТик() {
	defer func() {
		close(сам.канТик)
		// log._rintf("Shot.makeTick(): сражение завершёно\n")
	}()
	for {
		select {
		case <-сам.Контекст().Done():
			return
		default:
			сам.канТик <- 1 // Первый выстрел, потом спать по таймингу
			// log._rintf("INFO Shot.run() перезарядка=%v msec\n", recharge)
			// Если идёт перезарядка -- постепенно обнуляем время ожидания
			time.Sleep(time.Millisecond * 6800)
		}
	}
}

// Цикл выстрела (в отдельном потоке)
func (сам *Выстрел) пуск() {
	for range сам.канТик {
		// Стрелять можно, стандартное ожидание
		сам.выстрел()
	}
}

// Обновляет возможность выстрела (~)
//
//	Вызывается из отдельного потока
func (сам *Выстрел) выстрел() {
	сам.Сеть().Обновить()

	// <a href="pve?6-26.ILinkListener-currentControl-attackRegularShellLink" class="simple-but gray"><span><span>ОБЫЧНЫЕ</span></span></a>

	strLink := "https://wartank.ru/pve?6-{count}.ILinkListener-currentControl-attackRegularShellLink"
	strLink = strings.ReplaceAll(strLink, "{count}", fmt.Sprint(сам.выстрелСчёт))
	сам.выстрелСчёт++
	res := сам.Сеть().Get(strLink)
	if res.IsErr() {
		// log._rintf("ERRO Shot.shot(): при исполнении GET-команды выстрела обычным снарядом, err=\n\t%v\n", err)
		return
	}
	сам.Манёвр().УстНадо()
	// sound.Shot()
}
