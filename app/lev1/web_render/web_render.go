// package web_render -- возвращает веб-рендер
package web_render

import (
	_ "embed"
	"fmt"
	"strings"

	. "wartank/app/lev0/types"

	. "gitp78su.ipnodns.ru/svi/kern"
	. "gitp78su.ipnodns.ru/svi/kern/kc/helpers"
	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"
)

// ВебРендер -- простой веб-рендер
type ВебРендер struct {
	лог     ILogBuf
	слвБлок map[string]interface{}
	шаблон  string
}

// НовВебРендер -- возвращает новый веб-рендер
func НовВебРендер(шаблон string) *ВебРендер {
	лог := NewLogBuf()
	сам := &ВебРендер{
		лог:     лог,
		слвБлок: map[string]interface{}{},
		шаблон:  шаблон,
	}
	_ = ИВебРендер(сам)
	return сам
}

// Получ -- возвращает результат рендеринга
func (сам *ВебРендер) Получ() string {
	рез := сам.шаблон
	for ключ, блок := range сам.слвБлок {
		знач := fmt.Sprint(блок)
		рез = strings.ReplaceAll(рез, ключ, знач)
	}
	return рез
}

// Доб -- добавляет блок замещения
func (сам *ВебРендер) Доб(ключ string, блок any) {
	Hassert(ключ != "", "ВебРендер.Доб(): пустой ключ")
	сам.слвБлок[ключ] = блок
}
