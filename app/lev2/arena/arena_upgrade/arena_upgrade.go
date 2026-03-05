// package arena_upgrade -- арена улучшения параметров
package arena_upgrade

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	. "wartank/app/lev0/types"

	. "gitp78su.ipnodns.ru/svi/kern"
	. "gitp78su.ipnodns.ru/svi/kern/kc/helpers"
	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"
)

// ТанкПараметры -- параметры танка повышение
type АренаАпгрейд struct {
	конт   ILocalCtx
	прилож ИПриложение
	номер  string // Номер танка в игре
	лог    ILogBuf
}

// НовТанкПараметры -- возвращает новые параметры танка
func НовТанкПараметры(конт ILocalCtx) *АренаАпгрейд {
	лог := NewLogBuf()
	лог.Info("НовТанкПараметры()\n")

	сам := &АренаАпгрейд{
		конт:   конт,
		прилож: конт.Get("прилож").(ИПриложение),
		лог:    лог,
	}
	return сам
}

// Пуск -- запуск в работу
func (сам *АренаАпгрейд) Пуск() {
	go сам.пуск()
}

// Запускает в работу в отдельном потоке
func (сам *АренаАпгрейд) пуск() {
	time.Sleep(time.Second * 4)
	ош := сам.номерПолуч()
	if ош != nil {
		log.Printf("ТанкПараметры.пуск(): при получении параметров танка, ош=\n\t%v\n", ош)
		сам.конт.Cancel()
		return
	}
	for {
		select {
		case <-сам.конт.Ctx().Done():
			return
		default:
			сам.работать()
		}
	}
}

// Основной метод работы
func (сам *АренаАпгрейд) работать() {
	defer time.Sleep(time.Second * 300)
	сам.улучшить()
}

// Улучшает параметры танка
func (сам *АренаАпгрейд) улучшить() {
	// https://wartank.ru/pimp/34479487
	клиент := сам.конт.Get("хттпВоркер").(ИХттпВоркер)
	фнУлучшить := func() bool {
		лстСтр := клиент.Получ("https://wartank.ru/pimp/" + сам.номер)
		var (
			стрВых    string
			еслиНашли bool
		)
		// <a class="simple-but border mb5" href="34479487?22-1.ILinkListener-modules-slots-0-slot-root-pimpLink-link">
		for _, стрВых = range лстСтр {
			if strings.Contains(стрВых, `<a class="simple-but border mb5" href="`+сам.номер) {
				еслиНашли = true
				break
			}
		}
		Hassert(еслиНашли, "ТанкПараметры.улучшить(): не нашёл кнопку улучшения")
		стрВых = strings.TrimPrefix(стрВых, `<a class="simple-but border mb5" href="`)
		стрВых = strings.TrimSuffix(стрВых, `">`)
		// https://wartank.ru/pimp/34479487?21-1.ILinkListener-modules-slots-0-slot-root-pimpLink-link
		стрСсылка := "https://wartank.ru/pimp/" + стрВых
		_ = клиент.Получ(стрСсылка)
		return true
	}
	счётОш := 5
	for счётОш > 0 {
		if фнУлучшить() {
			break
		}
		счётОш--
		time.Sleep(time.Millisecond * 350)
	}
}

// Получает собственный номер танка с сервера
func (сам *АренаАпгрейд) номерПолуч() error {
	клиент := сам.конт.Get("хттпВоркер").(ИХттпВоркер)
	лстСтр := клиент.Получ("https://wartank.ru/angar")
	var (
		стрНомер  string
		еслиНашёл bool
	)
	// https://wartank.ru/power/34479487
	for _, стрНомер = range лстСтр {
		if strings.Contains(стрНомер, `href="power/`) {
			еслиНашёл = true
			break
		}
	}
	if !еслиНашёл {
		return fmt.Errorf("ТанкПараметры.номерПолуч(): не нашёл собственный номер на сервере")
	}
	стрВыход := strings.TrimPrefix(стрНомер, `<a class="simple-but border" href="power/`)
	стрВыход = strings.TrimSuffix(стрВыход, `"><span><span>Повысить параметры</span></span></a>`)
	_, ош := strconv.Atoi(стрВыход)
	if ош != nil {
		return fmt.Errorf("ТанкПараметры.номерПолуч(): ошибка преобразования собственного номера(%s), ош=\n\t%w", стрВыход, ош)
	}
	сам.номер = стрВыход
	return nil
}
