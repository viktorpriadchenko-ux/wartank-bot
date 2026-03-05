// package mod_wui -- модуль WUI
package mod_wui

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"

	. "gitp78su.ipnodns.ru/svi/kern/kc/helpers"
	"gitp78su.ipnodns.ru/svi/kern/krn/kctx"
	"gitp78su.ipnodns.ru/svi/kern/krn/kmodule"
	"gitp78su.ipnodns.ru/svi/kern/krn/kserv_http"
	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"
	"gitp78su.ipnodns.ru/svi/kern/mds/mod_serv_http/http_api"
	"gitp78su.ipnodns.ru/svi/kern/mds/mod_serv_http/page_module"
	"gitp78su.ipnodns.ru/svi/kern/mds/mod_serv_http/page_monolit"
	"gitp78su.ipnodns.ru/svi/kern/wui"
	. "gitp78su.ipnodns.ru/svi/kern/wui/wtypes"
)

// ModuleWui -- модуль WUI
type ModuleWui struct {
	IKernelModule
	kCtx      IKernelCtx
	wCtx      IWuiCtx
	kServHttp IKernelServerHttp
	log       ILogBuf
}

var (
	mod   *ModuleWui
	block sync.Mutex
)

// GetModuleWui -- возвращает новый модуль WUI
func GetModuleWui() *ModuleWui {
	block.Lock()
	defer block.Unlock()
	if mod != nil {
		return mod
	}
	sf := &ModuleWui{
		kCtx:          kctx.GetKernelCtx(),
		wCtx:          wui.GetWuiCtx(),
		IKernelModule: kmodule.NewKernelModule("wui"),
		kServHttp:     kserv_http.GetKernelServHttp(),
	}
	sf.log = sf.wCtx.Log()
	_ = page_monolit.GetPageMonolit()
	_ = page_module.GetPageModule()

	_ = http_api.NewHttpApi()
	fibApp := sf.kCtx.Get("fiberApp").Val().(*fiber.App)
	fibApp.Post("/wui/click/:id", adaptor.HTTPHandlerFunc(sf.wuiClick)) // adaptor.HTTPHandlerFunc(greet)
	mod = sf
	return sf
}

// Run -- запускает модуль в работу
func (sf *ModuleWui) Run() {
	sf.log.Info("ModuleWui.Run(): module=%v, is run", sf.Name())
	go sf.kServHttp.Run()
}

// Log -- возвращает буферный лог
func (sf *ModuleWui) Log() ILogBuf {
	return sf.log
}

// IsWork -- признак работы модуля
func (sf *ModuleWui) IsWork() bool {
	return sf.kCtx.Wg().IsWork()
}

// Получает событие из сети
func (sf *ModuleWui) wuiClick(resp http.ResponseWriter, req *http.Request) {
	url := req.RequestURI
	id := strings.TrimPrefix(url, "/wui/click/")
	widget0 := sf.wCtx.Get(id)
	if widget0 == nil {
		strOut := fmt.Sprintf("ModuleWui.wuiClick(): id(%v), widget not exists", id)
		sf.log.Err(strOut)
		fmt.Fprint(resp, strOut)
		return
	}
	widget1, isOk := widget0.Val().(IWuiButton)
	if !isOk {
		strOut := fmt.Sprintf("ModuleWui.wuiClick(): widget(%T) not button", widget0.Val())
		sf.log.Err(strOut)
		fmt.Fprint(resp, strOut)
		return
	}
	dict := map[string]string{}

	// headers := ctx.GetReqHeaders()
	for key, lstVal := range req.Header {
		if len(lstVal) >= 1 {
			dict[key] = lstVal[0]
			continue
		}
	}

	err := req.ParseForm()
	Hassert(err == nil, "ModuleWui.wuiClick(): in parse form, err=\n\t%v", err)
	// Получаем все form-значения
	//values := req.ParseForm()
	for key, lstVal := range req.Form {
		if len(lstVal) >= 1 {
			dict[key] = lstVal[0]
			continue
		}
	}
	strOut := widget1.Click(dict)
	fmt.Fprint(resp, strOut)
}
