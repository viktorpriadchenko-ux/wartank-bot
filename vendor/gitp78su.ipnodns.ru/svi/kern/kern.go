// package kern -- библиотека гибкого универсального облегчённого ядра для любого микросервиса
package kern

import (
	"context"

	. "gitp78su.ipnodns.ru/svi/kern/kc/helpers"
	"gitp78su.ipnodns.ru/svi/kern/kc/local_ctx"
	"gitp78su.ipnodns.ru/svi/kern/kc/log_buf"
	"gitp78su.ipnodns.ru/svi/kern/kc/safe_bool"
	"gitp78su.ipnodns.ru/svi/kern/kc/safe_bool_react"
	"gitp78su.ipnodns.ru/svi/kern/kc/safe_int"
	"gitp78su.ipnodns.ru/svi/kern/kc/safe_string"
	. "gitp78su.ipnodns.ru/svi/kern/krn/kalias"
	"gitp78su.ipnodns.ru/svi/kern/krn/kbus/kbus_http"
	"gitp78su.ipnodns.ru/svi/kern/krn/kbus/kbus_http/client_bus_http"
	"gitp78su.ipnodns.ru/svi/kern/krn/kbus/kbus_local"
	"gitp78su.ipnodns.ru/svi/kern/krn/kbus/kbus_local/client_bus_local"
	"gitp78su.ipnodns.ru/svi/kern/krn/kctx"
	"gitp78su.ipnodns.ru/svi/kern/krn/kmodule"
	"gitp78su.ipnodns.ru/svi/kern/krn/kmonolit"
	"gitp78su.ipnodns.ru/svi/kern/krn/kserv_http"
	"gitp78su.ipnodns.ru/svi/kern/krn/kstore_kv"
	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"
	"gitp78su.ipnodns.ru/svi/kern/mds/mod_kctx"
	"gitp78su.ipnodns.ru/svi/kern/mds/mod_keeper"
	"gitp78su.ipnodns.ru/svi/kern/mds/mod_serv_http"
	"gitp78su.ipnodns.ru/svi/kern/mds/mod_wui"
)

// GetKernelCtx -- возвращает контекст ядра
func GetKernelCtx() IKernelCtx {
	ctx := kctx.GetKernelCtx()
	return ctx
}

// GetKernelStoreKv -- возвращает быстрое key-value хранилище ядра
func GetKernelStoreKv() IKernelStoreKv {
	store := kstore_kv.GetKernelStore()
	return store
}

// GetKernelServerHttp -- возвращает веб-сервер ядра
func GetKernelServerHttp() IKernelServerHttp {
	kernServHttp := kserv_http.GetKernelServHttp()
	return kernServHttp
}

// NewSafeBool -- возвращает новый потокобезопасный булевый признак
func NewSafeBool() ISafeBool {
	sb := safe_bool.NewSafeBool()
	return sb
}

// GetKernelBusLocal -- возвращает локальную шину данных
func GetKernelBusLocal() IKernelBus {
	ctx := kctx.GetKernelCtx()
	ctx.Set("monolitName", "unknown monolit", "GetKernelBusLocal()")
	bus := kbus_local.GetKernelBusLocal()
	return bus
}

// GetKernelBusHttp -- возвращает HTTP шину данных
func GetKernelBusHttp() IKernelBus {
	bus := kbus_http.GetKernelBusHttp()
	return bus
}

// GetMonolitLocal -- возвращает монолит с локальной шиной
func GetMonolitLocal(name string) IKernelMonolit {
	ctx := kctx.GetKernelCtx()
	ctx.Set("isLocal", true, "bus type")
	for {
		SleepMs()
		if ctx.Get("isLocal") != nil {
			break
		}
	}
	monolit := kmonolit.GetMonolit(name)
	_ = kbus_local.GetKernelBusLocal()
	return monolit
}

// GetMonolitHttp -- возвращает монолит с локальной шиной поверх HTTP
func GetMonolitHttp(name string) IKernelMonolit {
	ctx := kctx.GetKernelCtx()
	_ = kbus_http.GetKernelBusHttp()
	ctx.Set("isLocal", false, "bus type")
	for {
		SleepMs()
		if ctx.Get("isLocal") != nil {
			break
		}
	}
	monolit := kmonolit.GetMonolit(name)
	return monolit
}

// NewKernelModule -- возвращает новый модуль на ядре
func NewKernelModule(name AModuleName) IKernelModule {
	mod := kmodule.NewKernelModule(name)
	return mod
}

// NewClientBusLocal -- возвращает клиент для работы с локальной шиной
func NewClientBusLocal() IBusClient {
	client := client_bus_local.NewClientBusLocal()
	return client
}

// NewClientBusHttp -- возвращает клиент для работы с шиной поверх HTTP
func NewClientBusHttp(url string) IBusClient {
	client := client_bus_http.NewClientBusHttp(url)
	return client
}

// GetModuleServHttp -- возвращает модуль для IKernelServHttp
func GetModuleServHttp() IKernelModule {
	modServHttp := mod_serv_http.GetModuleServHttp()
	return modServHttp
}

// GetModuleKernelCtx -- возвращает модуль для IKernelCtx
func GetModuleKernelCtx() IKernelModule {
	modKernelCtx := mod_kctx.GetModuleKernelCtx()
	return modKernelCtx
}

// GetModuleKernelKeeper -- возвращает модуль для IKernelKeeper
func GetModuleKernelKeeper() IKernelModule {
	modKernelKeeper := mod_keeper.GetModuleKeeper()
	return modKernelKeeper
}

// GetModuleWui -- возвращает модуль для WUI
func GetModuleWui() IKernelModule {
	mod := mod_wui.GetModuleWui()
	return mod
}

// NewLogBuf -- возвращает новый буферизованный лог
func NewLogBuf() ILogBuf {
	log := log_buf.NewLogBuf()
	return log
}

// NewSafeBoolReact -- возвращает новую потокобезопасную реактивную булеву переменную
func NewSafeBoolReact() ISafeBoolReact {
	val := safe_bool_react.NewSafeBoolReact()
	return val
}

// NewSafeInt -- возвращает новую потокобезопасную целочисленную переменную
func NewSafeInt() ISafeInt {
	val := safe_int.NewSafeInt()
	return val
}

// NewLocalCtx -- возвращает новый локальный контекст
func NewLocalCtx(ctx context.Context) ILocalCtx {
	ctx_ := local_ctx.NewLocalCtx(ctx)
	return ctx_
}

// NewSafeString -- возвращает новую потокобезопасную строку
func NewSafeString() ISafeString {
	str := safe_string.NewSafeString()
	return str
}

// MakeOk -- возвращает новый положительный результат операции
func MakeOk[T any](res T) Result[T] {
	return NewOk(res)
}

// MakeErr -- возвращает новую ошибку результат операции
func MakeErr[T any](err error) Result[T] {
	return NewErr[T](err)
}

// MakeSome -- возвращает новый не пустой результат операции
func MakeSome[T any](some T) Option[T] {
	return NewSome(some)
}

// MakeNone -- возвращает новый пустой результат операции
func MakeNone[T any]() Option[T] {
	return NewNone[T]()
}
