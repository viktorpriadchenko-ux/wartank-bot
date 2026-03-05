// package death_register -- регистрирует танк в схватке
package death_register

import (
	"log"
	"strings"
	"time"

	. "wartank/app/lev0/types"
	"wartank/app/lev2/arena"

	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"
)

// СхваткаРегистрация -- регистрирует танк к началу схватки
type СхваткаРегистрация struct {
	ИАрена
	конт         ILocalCtx
	счётРегистер int // Счётчик регистраций на сражение
}

// НовСхваткаРегистрация -- возвращает новый ожидатель битвы
func НовСхваткаРегистрация(конт IKernelCtx) *СхваткаРегистрация {
	сам := &СхваткаРегистрация{
		конт:         конт,
		счётРегистер: 10_000,
	}
	аренаКонфиг := arena.АренаКонфиг{
		Конт_:        конт,
		АренаИмя_:    "Сражение",
		СтрКонтроль_: `<title>Сражения</title>`,
		СтрУрл_:      "https://wartank.ru/dm",
	}
	сам.ИАрена = arena.НовАрена(конт, аренаКонфиг)
	return сам
}

// Зарегистрироваться -- регистрирует танк на сражение
func (сам *СхваткаРегистрация) Зарегистрироваться() {
	// Найдено приглашение на участие
	// https://wartank.ru/dm?{count}-1.ILinkListener-currentOverview-apply
	фнРегис := func() []string {
		стрСсылка := "https://wartank.ru/dm?0-1.ILinkListener-currentOverview-apply"
		стрКонтроль := "" // "https://wartank.ru/dm?0-1.ILinkListener-currentOverview-apply"
		for {
			time.Sleep(time.Second * 1)
			res := сам.Сеть().Get(стрСсылка)
			if res.IsErr() {
				log.Printf("ERRO СхваткаРегистрация.Зарегистрироваться(): при выполнении GET-команды на подъём в атаку, err=\n\t%v\n", res.Error())
			}
			лстСражение := res.Unwrap()
			if len(лстСражение) < 113 {
				continue
			}
			стрКонтроль = лстСражение[113]
			if !strings.Contains(стрКонтроль, "ILinkListener-currentOverview-apply") {
				return лстСражение
			}
			log.Printf("СхваткаРегистрация.Зарегистрироваться(): регистрация не прошла\n")
			стрСсылка = strings.TrimPrefix(стрКонтроль, `<a class="simple-but border" href="`)
			стрСсылка = strings.TrimSuffix(стрСсылка, `.ILinkListener-currentOverview-apply"><span><span>Взвод, подъем! В атаку!</span></span></a>`)
			стрСсылка = "https://wartank.ru/" + стрСсылка + ".ILinkListener-currentOverview-apply"
		}
	}

	сам.СтрОбновить(фнРегис())
}
