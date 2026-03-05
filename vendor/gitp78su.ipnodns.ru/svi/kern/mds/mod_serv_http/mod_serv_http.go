// package mod_serv_http -- модуль HTTP-сервера
package mod_serv_http

import (
	"sync"

	"gitp78su.ipnodns.ru/svi/kern/krn/kmodule"
	"gitp78su.ipnodns.ru/svi/kern/krn/kserv_http"
	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"
	"gitp78su.ipnodns.ru/svi/kern/mds/mod_serv_http/http_api"
	"gitp78su.ipnodns.ru/svi/kern/mds/mod_serv_http/page_module"
	"gitp78su.ipnodns.ru/svi/kern/mds/mod_serv_http/page_monolit"
)

// ModuleServHttp -- модуль HTTP-сервера
type ModuleServHttp struct {
	IKernelModule
	kServHttp IKernelServerHttp
	log       ILogBuf
}

var (
	mod   *ModuleServHttp
	block sync.RWMutex
)

// GetModuleServHttp -- возвращает новый модуль HTTP-сервера
func GetModuleServHttp() *ModuleServHttp {
	block.Lock()
	defer block.Unlock()
	if mod != nil {
		return mod
	}
	sf := &ModuleServHttp{
		IKernelModule: kmodule.NewKernelModule("kServHttp"),
		kServHttp:     kserv_http.GetKernelServHttp(),
	}
	sf.log = sf.Ctx().Log()
	_ = page_monolit.GetPageMonolit()
	_ = page_module.GetPageModule()

	_ = http_api.NewHttpApi()
	mod = sf
	return sf
}

// Run -- запускает модуль в работу
func (sf *ModuleServHttp) Run() {
	sf.log.Info("ModuleServHttp.Run(): module=%v, is run", sf.Name())
	go sf.kServHttp.Run()
}

// IsWork -- признак работы модуля
func (sf *ModuleServHttp) IsWork() bool {
	return sf.kServHttp.IsWork()
}
