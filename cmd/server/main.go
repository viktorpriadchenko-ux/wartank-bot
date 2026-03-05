// package main -- пускач для сервера на файбере
//
//	Команда запуска обновлятора noip.com
//	noip-duc -g p78su.ddns.net,p78git.ddns.net --daemonize -u prospero78su -p Lera_07091978
//
// Профилирование:
//
//	go tool pprof http://localhost:29080/debug/pprof/profile?seconds=30
package main

import (
	"gitp78su.ipnodns.ru/svi/kern"
	. "gitp78su.ipnodns.ru/svi/kern"

	"wartank/app/lev0/bot_log"
	"wartank/app/lev4/build"
)

func main() {
	// Перехват stdout в кольцевой буфер (лог в веб-UI)
	bot_log.Init()

	монолит := GetMonolitLocal("WarTank")

	модКонт := GetModuleKernelCtx()
	монолит.Add(модКонт)

	модВеб := GetModuleServHttp()
	монолит.Add(модВеб)

	модВуи := kern.GetModuleWui()
	монолит.Add(модВуи)

	сервер := build.НовМодСервер()
	монолит.Add(сервер)

	монолит.Run()
	монолит.Wait()
}
