package health

import (
	"fmt"
	"strings"
	"time"

	. "wartank/app/lev0/types"
)

/*
	Контролирует состояние здоровья танка
*/

// Здоровье -- контроль здоровья танка
type Здоровье struct {
	ИСражениеПроцесс
	канТик     chan int // Канал для ровной отправки тиков
	счётЛечить int      // Счётчик вызовов лечения
}

// НовЗдоровье -- возвращает новый *Здоровье
func НовЗдоровье(проц ИСражениеПроцесс) (*Здоровье, error) {
	{ // Предусловия
		if проц == nil {
			return nil, fmt.Errorf("НовЗдоровье(): действие==nil")
		}
	}
	сам := &Здоровье{
		ИСражениеПроцесс: проц,
		канТик:           make(chan int, 2),
	}
	go сам.makeTik()
	go сам.пуск()
	return сам, nil
}

// Отправляет тики с заданным равным интервалом
func (сам *Здоровье) makeTik() {
	defer func() {
		close(сам.канТик)
		сам.Отменить()
		// log._rintf("Здоровье.makeTick(): сражение завершёно\n")
	}()
	for {
		select {
		case <-сам.Контекст().Done():
			return
		default:
			сам.канТик <- 1
			time.Sleep(time.Second * 10)
		}
	}
}

// Главный цикл обработки здоровья в сражении
func (сам *Здоровье) пуск() {
	defer func() {
		сам.Отменить()
	}()
	for range сам.канТик {
		сам.лечить()
	}
}

// Полное -- возвращает объект полного здоровья танка
func (сам *Здоровье) Полное() int {
	return 0
}

// ЕслиМёртвый -- возвращает признак мертвичины танка
func (сам *Здоровье) ЕслиМёртвый() bool {
	lstBattle := сам.СписПолучить()
	for _, strOut := range lstBattle {
		if strings.Contains(strOut, `>Ваш танк подбит.`) {
			// log._rintf("INFO Здоровье.repair(): танк подбит\n")
			сам.Отменить()
			return true
		}
	}
	return false
}

// Восстанавливает здоровье (~)
func (сам *Здоровье) лечить() {

	// <span>Ремкомплект</span>
	// <a href="pve?19-14.ILinkListener-currentControl-repairLink" class="simple-but blue"><span><span>Ремкомплект</span></span></a>

	strLink := "https://wartank.ru/pve?19-{count}.ILinkListener-currentControl-repairLink"
	// <a href="pve?6-26.ILinkListener-currentControl-repairLink" class="simple-but blue"><span><span>Ремкомплект</span></span></a>
	strLink = strings.ReplaceAll(strLink, "{count}", fmt.Sprint(сам.счётЛечить))
	res := сам.Сеть().Get(strLink)
	if res.IsErr() {
		// log._rintf("ERRO Здоровье.repair(): при выполнении GET-команды ремонта, err=\n\t%v\n", err)
		return
	}
	сам.СтрОбновить(res.Unwrap())
}
