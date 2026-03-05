package div_war_sound

import (
	"time"
	"wartank/app/lev1/sound"
	"wartank/app/lev1/sound/is_sound_play"
)

/*
	Выполняет контроль за запуском одной озвучки битвы
*/

// DivWarSound -- контроль одного раза запуска звука битвы
type DivWarSound struct {
	isPlay *is_sound_play.IsPlay
}

// NewDivWarSound -- возвращает новый  *DivWarSound
func NewDivWarSound() *DivWarSound {
	return &DivWarSound{
		isPlay: is_sound_play.NewIsPlay(),
	}
}

// Play -- играет музончик, если можно
func (сам *DivWarSound) Play() {
	if сам.isPlay.Get() {
		return
	}
	go сам.play()
}

// Проигрывает экслюзивно в отдельном потоке звук
func (сам *DivWarSound) play() {
	сам.isPlay.Set()
	val := 7
	for val > 0 {
		sound.DivWar()
		val--
		time.Sleep(time.Second * 1)
	}
	val = 600 // Пауза для блокировки повторного включения начатой битвы
	for val >= 0 {
		val--
		time.Sleep(time.Second * 1)
	}
	сам.isPlay.Reset()
}
