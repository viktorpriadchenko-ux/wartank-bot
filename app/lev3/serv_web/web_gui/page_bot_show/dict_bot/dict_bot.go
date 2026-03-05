// package dict_bot -- костыль для кеширования ботов
package dict_bot

import (
	"strconv"
	"sync"

	. "gitp78su.ipnodns.ru/svi/kern"
	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"

	"wartank/app/lev0/alias"
	. "wartank/app/lev0/types"
)

type DictBot struct {
	sync.Mutex
	ctx  IKernelCtx
	dict map[string]ИБот
	app  ИПриложение
}

var (
	sf    *DictBot
	block sync.Mutex
)

func GetDictBot() *DictBot {
	block.Lock()
	defer block.Unlock()
	if sf != nil {
		return sf
	}
	ctx := GetKernelCtx()
	sf := &DictBot{
		ctx:  ctx,
		dict: map[string]ИБот{},
		app:  ctx.Get("мод_сервер").Val().(ИПриложение),
	}
	return sf
}

func (sf *DictBot) Get(dict map[string]string) ИБот {
	strId, isOk := dict["id"]
	if !isOk {
		return nil
	}
	bot, isOk := sf.dict[strId]
	if isOk {
		return bot
	}
	iId, err := strconv.Atoi(strId)
	if err != nil {
		return nil
	}
	бот := sf.app.ServBots().Get(alias.АБотНомер(iId))
	if бот == nil {
		return nil
	}
	sf.dict[strId] = бот
	return бот
}
