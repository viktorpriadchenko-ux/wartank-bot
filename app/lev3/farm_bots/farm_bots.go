// package farm_bots -- бото-ферма
package farm_bots

import (
	"fmt"
	"log"

	. "gitp78su.ipnodns.ru/svi/kern"
	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"

	. "wartank/app/lev0/alias"
	. "wartank/app/lev0/types"
	"wartank/app/lev3/bot"
	"wartank/app/lev3/farm_bots/dict_bot"
)

// БотоФерма -- бото-ферма
type БотоФерма struct {
	конт IKernelCtx
	// прилож  ИПриложение
	хран    IKernelStoreKv
	словБот *dict_bot.СловарьБотов
}

// НовБотоФерма -- возвращает новую бото-ферму
func НовБотоФерма() *БотоФерма {
	конт := GetKernelCtx()
	log.Println("НовБотоФерма()")
	// приложение := конт.Получ("приложение").(ИПриложение)
	сам := &БотоФерма{
		конт: конт,
		// прилож: приложение,
		хран: конт.Get("kernStoreKV").Val().(IKernelStoreKv),
	}
	сам.словБот = dict_bot.НовСловарьБотов(конт)
	_ = ИБотоФерма(сам)
	return сам
}

// Get -- возвращает боевого бота по имени
func (сам *БотоФерма) Get(botNumber АБотНомер) ИБот {
	bot := сам.словБот.Get(botNumber)
	return bot
}

// BotStart -- запускает бота в работу по его имени
func (сам *БотоФерма) BotStart(botNumber АБотНомер) Result[bool] {
	bot := сам.словБот.Get(botNumber)
	if bot == nil {
		err := fmt.Errorf("ServBots.BotStart(): bot(%v) not found", botNumber)
		return NewErr[bool](err)
	}
	bot.Пуск()
	return NewOk(true)
}

// ListBot -- возвращает список существующих ботов
func (сам *БотоФерма) ListBot() []ИБот {
	lst := сам.словБот.ListBot()
	return lst
}

// НовБот -- добавляет нового бота на ферму
func (сам *БотоФерма) НовБот(логин, пароль string, еслиАвто bool) Result[bool] {
	{ // Существует ли такой бот
		for _, бот := range сам.словБот.ListBot() {
			if бот.Имя() == логин {
				err := fmt.Errorf("БотоФерма.НовБот(): логин(%v) уже существует", логин)
				return NewErr[bool](err)
			}
		}
	}
	номер := АБотНомер(len(сам.словБот.ListBot()) + 1)
	фнНомерПровер := func() bool {
		for _, бот := range сам.словБот.ListBot() {
			if бот.Номер() == номер {
				return false
			}
		}
		return true
	}
	for !фнНомерПровер() {
		номер++
	}
	// Нет такого бота, надо его создать
	бот := bot.НовВарБот(сам.конт, номер, логин, пароль, еслиАвто)
	сам.словБот.Add(бот)
	return NewOk(true)
}
