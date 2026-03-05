// package bot_net_login -- сетевой вход на сервер бота
package bot_net_login

import (
	"io"
	"net/http"
	"net/url"
	"strings"

	. "gitp78su.ipnodns.ru/svi/kern"
	. "gitp78su.ipnodns.ru/svi/kern/kc/helpers"
	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"

	. "wartank/app/lev0/types"
)

// БотСетьЛогин -- объект сетевого входа на сервер
type БотСетьЛогин struct {
	сеть       ИБотСеть
	хттпКлиент *http.Client // Сырой клиент HTTP
	куки       ИБотКуки
	лог        ILogBuf
}

// НовБотСетьЛогин -- возвращает новый *БотСетьЛогин
func НовБотСетьЛогин(конт ILocalCtx) *БотСетьЛогин {
	лог := NewLogBuf()
	лог.Info("НовБотСетьЛогин()\n")

	ботСеть := конт.Get("ботСеть").Val().(ИБотСеть)
	хттпКлиент := конт.Get("хттпКлиент").Val().(*http.Client)
	сам := &БотСетьЛогин{
		сеть:       ботСеть,
		хттпКлиент: хттпКлиент,
		куки:       ботСеть.Куки(),
		лог:        лог,
	}
	сам.подключить()
	return сам
}

// Вызывается один раз
func (сам *БотСетьЛогин) подключить() {
	ссылкаНачать := сам.получСтрГлав()
	if ссылкаНачать != "" { // Это вход на базу
		стрСтрЛогин := сам.получСтрЛогин(ссылкаНачать)
		сам.выполнитьЛогин(стрСтрЛогин)
		return
	}
	// Логин уже был
	запр, ош1 := http.NewRequest("GET", "https://wartank.ru/angar", nil)
	Hassert(ош1 == nil, "БотСетьЛогин.подключить(): при получении страницы ангара, err=\n\t%v\n", ош1)
	resp, ош2 := сам.хттпКлиент.Do(запр)
	Hassert(ош2 == nil, "БотСетьЛогин.подключить(): при выполнении запроса, err=\n\t%v\n", ош2)
	defer resp.Body.Close()
	_, ош := io.ReadAll(resp.Body)
	Hassert(ош == nil, "БотСетьЛогин.подключить(): при получении тела страницы ангара, err=\n\t%v\n", ош)
}

// Прочитать главную страницу для получения кукисов
func (сам *БотСетьЛогин) получСтрГлав() string {
	ответ, ош := сам.хттпКлиент.Get("https://wartank.ru/")
	Hassert(ош == nil, "БотСетьЛогин.получСтрГлав(): err=\n\t%v\n", ош)
	defer ответ.Body.Close()
	// Получить куки из ответа
	куки := ответ.Cookies()
	if len(куки) > 0 {
		сам.куки.Уст(куки)
	}
	бинДанные, ош := io.ReadAll(ответ.Body)
	Hassert(ош == nil, "БотСетьЛогин.получСтрГлав(): при чтении тела ответа, err=\n\t%v\n", ош)
	// Вырезать из тела страницы ссылку на вход
	списСтр := strings.Split(string(бинДанные), "\n")
	var (
		стрВых    string
		ссылкаНач string
	)
	for _, стрСсылка := range списСтр {
		if strings.Contains(стрСсылка, `w:id="showSigninLink"`) {
			стрВых = стрСсылка
			break
		}
	}
	if стрВых == "" { // Уже был логин
		return ""
	}
	{ // Получить ссылку на вход
		списСсылка := strings.Split(стрВых, `href="`)
		списСсылка = strings.Split(списСсылка[1], `"><span><span>`)
		ссылкаНач = "https://wartank.ru/" + списСсылка[0]
	}
	return ссылкаНач
}

// Получает страницу логина
func (сам *БотСетьЛогин) получСтрЛогин(linkBegin string) string {
	// _mt.Println("\БотСетьЛогин.getPageLogin()")
	ответ, ош := сам.хттпКлиент.Get(linkBegin)
	Hassert(ош == nil, "БотСетьЛогин.получСтрЛогин(): in make GET-request for send login page, err=\n\t%v\n", ош)
	defer ответ.Body.Close()
	данные, ош := io.ReadAll(ответ.Body)
	Hassert(ош == nil, "БотСетьЛогин.получСтрЛогин(): in read body response, err=\n\t%v\n", ош)
	стрТелоЛогин := string(данные)
	return стрТелоЛогин
}

// Выполняет логин через POST-запрос
func (сам *БотСетьЛогин) выполнитьЛогин(strBody string) {
	// _mt.Println("\БотСетьЛогин.makePostLogin()")
	var strBodyMain string
	{ // Получить форму логина
		lstLogin := strings.Split(strBody, "\n")
		for _, strForm := range lstLogin {
			if strings.Contains(strForm, `<form w:id="loginForm" id=`) {
				strBodyMain = strForm
				break
			}
		}
	}
	var postLink string
	{ // Получить ссылку на POST-запрос
		lstLink := strings.Split(strBodyMain, ` action="`)
		strLink := lstLink[1]
		lstLink = strings.Split(strLink, `"><div style=`)
		postLink = "https://wartank.ru/" + lstLink[0]
	}
	// Конструируем ПОСТ-форму логина
	form := url.Values{}
	form.Add("id1_hf_0", "")
	form.Add("login", сам.сеть.Бот().Имя())
	form.Add("password", сам.сеть.Бот().Пароль())

	resp, err := сам.хттпКлиент.PostForm(postLink, form)
	Hassert(resp.StatusCode == 200, "БотСетьЛогин.makePostLogin():  in get POST-login response, err=\n\t%v\n", err)
	defer resp.Body.Close()
	Hassert(err == nil, "БотСетьЛогин.makePostLogin(): in read body POST-login response, err=\n\t%v\n", err)
	urlObj, _ := url.Parse("https://wartank.ru/")
	сам.хттпКлиент.Jar.SetCookies(urlObj, сам.куки.Получ())

}
