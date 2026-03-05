package types

/*
	Интерфейс для приложения
*/

import (
	. "gitp78su.ipnodns.ru/svi/kern/krn/ktypes"
)

// ИСервер -- интерфейс для приложения
type ИПриложение interface {
	IKernelModule
	// ServBots -- словарь имеющихся ботов
	ServBots() ИБотоФерма
	// Стат -- возвращает статистику сервера
	Стат() ИСерверСтат
}
