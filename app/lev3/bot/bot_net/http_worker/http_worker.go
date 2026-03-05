// package http_worker -- веб-воркер бота
package http_worker

import (
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	. "wartank/app/lev0/types"
	"wartank/app/lev3/bot/bot_net/http_worker/net_stat"

	. "gitp78su.ipnodns.ru/svi/kern"
	. "gitp78su.ipnodns.ru/svi/kern/kc/helpers"
	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"
)

// ХттпВоркер -- объект сетевого исполнителя бота
type ХттпВоркер struct {
	sync.Mutex
	конт       ILocalCtx
	ботСеть    ИБотСеть
	хттпКлиент *http.Client
	статистика *net_stat.СетьСтата
	лог        ILogBuf
	канТик     chan int
}

// НовХттпВоркер -- возвращает веб-воркера бота
func НовХттпВоркер(конт ILocalCtx) *ХттпВоркер {
	лог := NewLogBuf()
	лог.Info("НовХттпВоркер()\n")

	ботСеть := конт.Get("ботСеть").Val().(ИБотСеть)
	хттпКлиент := конт.Get("хттпКлиент").Val().(*http.Client)

	сам := &ХттпВоркер{
		конт:       конт,
		ботСеть:    ботСеть,
		хттпКлиент: хттпКлиент,
		статистика: net_stat.НовСетьСтата(конт),
		лог:        лог,
		канТик:     make(chan int, 3),
	}
	_ = ИХттпВоркер(сам)
	конт.Set("хттпВоркер", сам, "HTTP-воркер бота")
	go сам.тик()
	return сам
}

// Делает тики в отдельном канале, не более 3 шт за раз (тротлинг для обмана сервера)
//
//	Если тиками никто не пользуются -- они накапливаются.
//	А если пользуются больше 3х за раз -- включается задержка по 1000 мсек
func (сам *ХттпВоркер) тик() {
	defer close(сам.канТик) // Закрываем канал при выходе, чтобы Получ() не зависал навечно
	for {
		select {
		case <-сам.конт.Ctx().Done():
			return
		default:
			time.Sleep(time.Millisecond * 1000)
			сам.канТик <- 1
		}
	}
}

// Получ -- потокобезопасно возвращает список строк по ссылке
func (сам *ХттпВоркер) Получ(strLink string) []string {
	сам.Lock()
	defer сам.Unlock()
	// Ждём тик или отмену контекста
	select {
	case _, ok := <-сам.канТик:
		if !ok { // Канал закрыт — контекст отменён
			return nil
		}
	case <-сам.конт.Ctx().Done():
		return nil
	}
	var ответ *http.Response
	запрос, ош := http.NewRequest("GET", strLink, nil)
	Hassert(ош == nil, "ХттпВоркер.Получ(): при создании запроса, err=\n\t%v\n", ош)
	запрос.Header.Set("User-Agent", "Mozilla Firefox 120.1")
	ответ, ош = сам.хттпКлиент.Do(запрос)
	// ответ, ош := http.Get(strLink)
	Hassert(ош == nil, "ХттпВоркер.Получ(): при ответе на запрос, err=\n\t%v\n", ош)
	defer ответ.Body.Close()
	Hassert(ответ.StatusCode == http.StatusOK, "ХттпВоркер.Получ(): code=%v, status=%v", ответ.StatusCode, ответ.Status)
	binData, ош := io.ReadAll(ответ.Body)
	Hassert(ош == nil, "ХттпВоркер.Получ(): при чтении тела ответа, err=\n\t%v\n", ош)

	Hassert(len(binData) != 0, "ХттпВоркер.Получ(): пустое тело ответа")
	lenData := len(binData) + len(strLink)
	сам.статистика.AddByte(lenData)

	lstString := strings.Split(string(binData), "\n")
	Hassert(len(lstString) != 0, "ХттпВоркер.Получ(): пустая строка ответа")
	return lstString
}
