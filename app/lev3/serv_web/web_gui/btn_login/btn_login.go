// package btn_login -- кнопка логина
package btn_login

import (
	"crypto/rand"
	_ "embed"
	"strings"

	. "gitp78su.ipnodns.ru/svi/kern"
	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"
	"gitp78su.ipnodns.ru/svi/kern/wui"
	. "gitp78su.ipnodns.ru/svi/kern/wui/wtypes"
)

// КнпЛогин -- кнопка логина на сайте
type КнпЛогин struct {
	Кнп  IWuiButton
	конт IKernelCtx
}

// НовКнпЛогин -- возвращает новую кнопку логина
func НовКнпЛогин() *КнпЛогин {
	sf := &КнпЛогин{
		конт: GetKernelCtx(),
	}
	sf.Кнп = wui.NewWuiButton("Логин", sf.КликЛогин)
	sf.Кнп.Hx().Target().Set("#main")
	return sf
}

// Html -- возвращает HTML-представление кнопки
func (sf *КнпЛогин) Html() string {
	return sf.Кнп.Html()
}

//go:embed block_login.html
var стрБлокЛогин string

type ИВебКнопка interface {
	Html() string
	КликСписБот(map[string]string) string
}

// Событие клика по кнопке
//
//	FIXME: здесь не прокидываются кукисы. Надо в form вставлять скрытое поле 'token'
func (сам *КнпЛогин) КликЛогин(слв map[string]string) string {
	стрРез := strings.ReplaceAll(стрБлокЛогин, "{.id}", сам.Кнп.Hx().Url().String())
	кнпСписБот := сам.конт.Get("кнп_спис_бот").Val().(ИВебКнопка)
	еслиЛогин := сам.конт.Get("еслиЛогин")
	if еслиЛогин != nil {
		токенЗапр := еслиЛогин.Val().(string)
		токенСвой_ := сам.Кнп.Hx().Vals().Get("токен")
		if токенСвой_ == nil {
			return стрРез
		}
		токенСвой := токенСвой_.(string)
		if токенСвой != токенЗапр {
			return стрРез
		}
		return кнпСписБот.КликСписБот(слв)
	}
	логин, еслиОк := слв["login"]
	if !еслиОк {
		return стрРез
	}
	if логин != "svi" {
		return стрРез
	}
	пароль, еслиОк := слв["password"]
	if !еслиОк {
		return стрРез
	}
	if пароль != "Lera_07091978" {
		return стрРез
	}
	токен := rand.Text()
	сам.Кнп.Hx().Vals().Set("токен", токен)
	сам.конт.Set("еслиЛогин", токен, "С какого IP зашёл царь")
	return кнпСписБот.КликСписБот(слв)
}
