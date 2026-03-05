// package arena_laborator -- лаборатория на базе
package arena_laborator

import (
	"fmt"
	"log"
	"strings"
	"time"
	. "wartank/app/lev0/types"
)

// АренаЛаборатория на базе
type АренаЛаборатория struct {
	бот ИБот
}

// НоваяЛаборатория -- возвращает новую лабораторию
func НоваяЛаборатория(бот ИБот) (*АренаЛаборатория, error) {
	if бот == nil {
		return nil, fmt.Errorf("НоваяЛаборатория(): ИБот == nil")
	}
	сам := &АренаЛаборатория{
		бот: бот,
	}
	return сам, nil
}

// Пуск -- запуск в работу
func (сам *АренаЛаборатория) Пуск() {
	go сам.пуск()
}

// Запускает в работу в отдельном потоке
func (сам *АренаЛаборатория) пуск() {
	time.Sleep(time.Millisecond * 4500)
	for {
		select {
		case <-сам.бот.КонтБот().Ctx().Done():
			return
		default:
			сам.работать()
		}
	}
}

// Основной метод работы
func (сам *АренаЛаборатория) работать() {
	defer time.Sleep(time.Second * 300)
	if ош := сам.улучшить(); ош != nil {
		log.Printf("Лаборатория.работать(): ош=\n\t%v\n", ош)
		return
	}
}

// Улучшает параметры лаборатории
func (сам *АренаЛаборатория) улучшить() error {
	// https://wartank.ru/buildings
	клиент := сам.бот.Сеть().ВебВоркер()
	фнПостроить := func() error {
		лстСтр := клиент.Получ("https://wartank.ru/buildings")
		еслиНашли := false
		// <td style="width:50%;padding-left:1px;"><a class="simple-but border mb5" href="building-upgrade/Laboratory"><span><span>Построить</span></span></a></td>
		for _, стр := range лстСтр {
			if strings.Contains(стр, `<td style="width:50%;padding-left:1px;"><a class="simple-but border mb5" href="building-upgrade/Laboratory"><span><span>Построить</span></span></a></td>`) {
				еслиНашли = true
				break
			}
		}
		if !еслиНашли {
			return nil
		}
		// https://wartank.ru/building-upgrade/Laboratory
		_ = клиент.Получ("https://wartank.ru/building-upgrade/Laboratory")
		return nil
	}
	фнКупить := func() error {
		лстСтр := клиент.Получ("https://wartank.ru/building-upgrade/Laboratory")
		стрВых := ""
		// <a class="simple-but border mb5" href="Laboratory?118-1.ILinkListener-upgradeLink-link">
		for _, стрВых = range лстСтр {
			if strings.Contains(стрВых, `<a class="simple-but border mb5" href="Laboratory?`) {
				break
			}
		}
		if стрВых == "" {
			return nil
		}
		стрВых = strings.TrimPrefix(стрВых, `<a class="simple-but border mb5" href="`)
		стрВых = strings.TrimSuffix(стрВых, `">`)
		// https://wartank.ru/building-upgrade/Laboratory?117-1.ILinkListener-upgradeLink-link
		стрВых = "https://wartank.ru/building-upgrade/" + стрВых
		_ = клиент.Получ(стрВых)
		return nil
	}
	счётОш := 5
	for счётОш > 0 {
		time.Sleep(time.Millisecond * 350)
		счётОш--
		ош := фнПостроить()
		if ош != nil {
			log.Printf("Лаборатория.улучшить(): получить, ош=\n\t%v\n", ош)
			continue
		}
		ош = фнКупить()
		if ош != nil {
			log.Printf("Лаборатория.улучшить(): оплатить, ош=\n\t%v\n", ош)
			continue
		}
		// FIXME: надо сделать подтверждение оплаты
		break
	}
	return nil
}
