// package build -- различные строители
package build

import (
	"net/http"
	_ "net/http/pprof"
	"runtime/debug"
	"time"

	. "gitp78su.ipnodns.ru/svi/kern"
	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"

	"wartank/app/lev4/mod_serv"
)

func НовМодСервер() IKernelModule {
	конт := GetKernelCtx()
	go ЗапуститьПрофиль()
	bi, _ := debug.ReadBuildInfo()
	лог := NewLogBuf()
	лог.Info("ИНФО \tgo = %v\n\tvers = %v\n", bi.GoVersion, bi.Main.Version)
	прил := mod_serv.НовМодСервер()
	go func() {
		time.Sleep(time.Minute * 20)
		конт.Cancel()
	}()
	return прил
}

func ЗапуститьПрофиль() {
	лог := NewLogBuf()
	port := "29081"
	for {
		err := http.ListenAndServe("0.0.0.0:"+port, nil)
		if err != nil {
			лог.Err("profile(): ошибка при запуске профилировщика, err=\n\t%v\n", err)
		}
		time.Sleep(time.Second * 1)
	}
}
