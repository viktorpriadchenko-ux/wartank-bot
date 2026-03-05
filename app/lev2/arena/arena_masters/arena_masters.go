package arena_masters

import (
	"strings"

	. "gitp78su.ipnodns.ru/svi/kern"
	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"

	. "wartank/app/lev0/alias"
	. "wartank/app/lev0/types"
	"wartank/app/lev2/arena"
	"wartank/app/lev2/arena/arena_build"
	"wartank/app/lev2/arena/arena_masters/bf_masters_fight"
	"wartank/app/lev2/arena/arena_masters/bf_masters_register"
	"wartank/app/lev2/arena/arena_masters/bf_masters_wait"
)

/*
	Битва мастеров. Работает примерно раз в сутки.
	Требуется три победы, потом нужно загрести золотишко.
	Между битвами надо удерживать рейтинг, чтобы не кидало к монстрам.
*/

// БитваМастеров -- объект битвы мастеров
type БитваМастеров struct {
	ИАренаСтроение
	конт ILocalCtx
	лог  ILogBuf
}

// НовБитваМастеров -- возвращает новый *BatMas
func НовБитваМастеров(конт ILocalCtx) *БитваМастеров {
	лог := NewLogBuf()
	лог.Info("НовБитваМастеров()\n")
	сам := &БитваМастеров{
		конт: конт,
		лог:  лог,
	}
	аренаКонфиг := arena.АренаКонфиг{
		Конт_:        конт,
		АренаИмя_:    "Битва мастеров",
		СтрКонтроль_: `/> Битва мастеров <`,
		СтрУрл_:      "https://wartank.ru/pvp",
	}
	сам.ИАренаСтроение = arena_build.НовАренаСтроение(конт, аренаКонфиг)
	конт.Set("pvp", сам, "Арена битвы мастеров")
	return сам
}

func (сам *БитваМастеров) Пуск() {
	сам.ИАренаСтроение.Обновить()
	bf_masters_register.СражениеРегистрация(сам.конт)
	bf_masters_wait.МастераОжидать(сам.конт)
	bf_masters_fight.МастераВыполнить(сам.конт)
}

// Вычисляет нужно ли идти в битву мастеров
//
//	если нужно, то время проверять уже не надо
func (сам *БитваМастеров) goBatMas() bool {
	сам.findTimeCount()
	if !сам.upBattle() {
		return false
	}
	countTime := сам.ВремяОстат().String()
	if countTime > "00:25:00" {
		сам.ОбратВремяУст(АВремя(countTime))
	}

	// Время меньше 25 сек, надо уточнять (тут возможна ошибка с экраном ожидания)
	сам.Обновить()
	// Время ожидания вышло, надо начать атаку
	сам.ОбратВремяУст("00")
	return false
}

// Ищет время до начала битвы мастеров
func (сам *БитваМастеров) findTimeCount() {
	var (
		strOut      string
		lstBattle   = сам.СписПолучить()
		еслиНайдено bool
	)
	// Обновление через: 12:02:22
	for _, strOut = range lstBattle {
		if strings.Contains(strOut, `Обновление через: `) {
			еслиНайдено = true
			break
		}
	}
	if еслиНайдено { // Ждём начала битвы мастеров
		lstTime := strings.Split(strOut, `Обновление через: `)
		strTime := lstTime[1]
		lstTime = strings.Split(strTime, ` (`)
		strTime = lstTime[0]

		сам.ОбратВремяУст(АВремя(strTime))
	}
}

// При необходимости даёт команду на участие в битве мастеров,
//
//	вызывается только если есть награда
func (сам *БитваМастеров) upBattle() bool {
	сам.Обновить()
	// log.Error("BatMas.upBattle(): доделать")
	// var (
	// 	strOut    string
	// 	lstBattle = сам.GetLst()
	// 	еслиНайдено    bool
	// )
	// for _, strOut = range lstBattle {
	// 	if strings.Contains(strOut, `>Взвод, подъем! В атаку!<`) {
	// 		еслиНайдено = true
	// 		break
	// 	}
	// }
	// if еслиНайдено {
	// 	lstUp := strings.Split(strOut, `<a class="simple-but border" href="`)
	// 	linkUp := lstUp[1]
	// 	lstUp = strings.Split(linkUp, `"><span><span>Взвод, подъем! В атаку!</span></span></a>`)
	// 	linkUp = "https://wartank.ru/" + lstUp[0]
	// 	lstBattle, err := сам.net.Get(linkUp)
	// 	if err != nil {
	// 		log.WithError(err).Error("Battle.upBattle(): при выполнении GET-команды на подъём в атаку")
	// 		return false
	// 	}
	// 	if err = сам.Update(lstBattle); err != nil {
	// 		log.WithError(err).Error("Battle.upBattle(): при обновлении lstBattle")
	// 	}
	// }
	return false
}
