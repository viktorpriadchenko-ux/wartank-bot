// package dict_bot -- потокобезопасный словарь ботов
package dict_bot

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	. "gitp78su.ipnodns.ru/svi/kern"
	. "gitp78su.ipnodns.ru/svi/kern/kc/helpers"
	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"

	. "wartank/app/lev0/alias"
	. "wartank/app/lev0/types"
	"wartank/app/lev3/bot"
)

const (
	strBotList = "/bots/list" // Константа путь к списку ботов в базе
)

// СловарьБотов -- потокобезопасный словарь ботов
type СловарьБотов struct {
	конт    IKernelCtx
	прилож  ИПриложение
	хран    IKernelStoreKv
	словарь map[АБотНомер]ИБот
	блок    sync.RWMutex
	лог     ILogBuf
}

// НовСловарьБотов -- возвращает новый потокобезопасный словарь ботов
func НовСловарьБотов(конт IKernelCtx) *СловарьБотов {
	лог := NewLogBuf()
	лог.Info("НовСловарьБотов()\n")
	сам := &СловарьБотов{
		конт:    конт,
		прилож:  конт.Get("мод_сервер").Val().(ИПриложение),
		хран:    конт.Get("kernStoreKV").Val().(IKernelStoreKv),
		словарь: map[АБотНомер]ИБот{},
		лог:     лог,
	}
	сам.load()
	return сам
}

// ListBot -- возвращает список существующих ботов
func (сам *СловарьБотов) ListBot() []ИБот {
	сам.блок.RLock()
	defer сам.блок.RUnlock()
	lst := make([]ИБот, 0)
	for _, bot := range сам.словарь {
		lst = append(lst, bot)
	}
	return lst
}

// Get -- возвращает бота по имени
func (сам *СловарьБотов) Get(botNumber АБотНомер) ИБот {
	сам.блок.RLock()
	defer сам.блок.RUnlock()
	bot := сам.словарь[botNumber]
	return bot
}

// Add -- добавляет нового бота в словарь
func (сам *СловарьБотов) Add(bot ИБот) {
	сам.блок.Lock()
	defer сам.блок.Unlock()
	if bot == nil {
		return
	}
	сам.словарь[bot.Номер()] = bot
	сам.save()
}

// Сохраняет словарь ботов в базе
func (сам *СловарьБотов) save() {
	strNumber := ""
	for botNumber := range сам.словарь {
		strNumber += fmt.Sprint(botNumber) + ";"
	}
	strNumber = strNumber[:len(strNumber)-1]
	res := сам.хран.Set(strBotList, []byte(strNumber))
	if res.IsErr() {
		сам.конт.Cancel()
	}
}

// Загружает всех ботов с базы
func (сам *СловарьБотов) load() {
	res := сам.хран.Get(strBotList)
	if res.IsErr() {
		if !strings.Contains(res.Error().Error(), "not found") {
			res.Hassert("СловарьБотов.load(): при загрузке списка ботов")
		}
		return
	}
	strNumbers := string(res.Unwrap())
	if strNumbers == "" {
		return
	}
	lstNumbers := strings.Split(strNumbers, ";")
	for _, стрНомер := range lstNumbers {
		if стрНомер == "" {
			continue
		}
		иНомер, err := strconv.Atoi(стрНомер)
		Hassert(err == nil, "СловарьБотов.load(): при получении номера бота, ош=\n\t%v\n", err)
		ботНомер := АБотНомер(иНомер)
		_, isOk := сам.словарь[ботНомер]
		if isOk {
			continue
		}
		bot := bot.ЗагрузитьВарБот(ботНомер)
		if bot.Автозапуск().Get() {
			go bot.Пуск()
		}
		сам.словарь[ботНомер] = bot
	}
}
