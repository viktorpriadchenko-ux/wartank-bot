package bot_cookie

import (
	"net/http"
	"sync"

	. "gitp78su.ipnodns.ru/svi/kern"
	. "gitp78su.ipnodns.ru/svi/kern/kc/helpers"
	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"
)

/*
	Предоставляет разделяемый объект кукисов для соединения с игровым сервером
*/

// БотКуки -- кукисы для бота, ничего не требует для своей работы
type БотКуки struct {
	куки []*http.Cookie
	блок sync.RWMutex
	лог  ILogBuf
}

// НовБотКуки -- возвращает новый *БотКуки
func НовБотКуки() БотКуки {
	лог := NewLogBuf()
	лог.Info("БотКуки()\n")
	return БотКуки{
		лог: лог,
	}
}

// Уст -- устанавливает кукисы
func (сам *БотКуки) Уст(cook []*http.Cookie) {
	сам.блок.Lock()
	defer сам.блок.Unlock()
	Hassert(cook != nil, "БотКуки.Уст(): cookie == nil")
	сам.куки = cook
	cookie := http.Cookie{
		Name:  "JSESSIONID",
		Value: сам.куки[0].Value,
		Raw:   "JSESSIONID=" + сам.куки[0].Value + "; _ym_uid=1642083867571238834; _ym_d=1642083867; _ym_isad=2; _ym_visorc=w",
	}
	cookie1 := сам.куки[:0]
	сам.куки = cookie1
	сам.куки = append(сам.куки, &cookie)

	cookie = http.Cookie{
		Name:  "_ym_d",
		Value: "1642083867",
	}
	сам.куки = append(сам.куки, &cookie)

	cookie = http.Cookie{
		Name:  "_ym_isad",
		Value: "2",
	}
	сам.куки = append(сам.куки, &cookie)

	cookie = http.Cookie{
		Name:  "_ym_visorc",
		Value: "w",
	}
	сам.куки = append(сам.куки, &cookie)
}

// Получ -- возвращает хранимые кукисы
func (сам *БотКуки) Получ() []*http.Cookie {
	return сам.куки
}
