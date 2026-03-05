package battle_sound

import (
	"time"
	"wartank/app/lev1/sound"
	"wartank/app/lev1/sound/is_sound_play"
)

/*
	Выполняет контроль за запуском одной озвучки битвы
*/
// BattleSound -- контроль одного раза запуска звука битвы
type BattleSound struct {
	isPlay *is_sound_play.IsPlay
}

// NewBattleSound -- возвращает новый  *BattleSound
func NewBattleSound() *BattleSound {
	return &BattleSound{
		isPlay: is_sound_play.NewIsPlay(),
	}
}

// Play -- играет музончик, если можно
func (сам *BattleSound) Play() {
	if сам.isPlay.Get() {
		return
	}
	go сам.play()
}

// Проигрывает экслюзивно в отдельном потоке звук
func (сам *BattleSound) play() {
	сам.isPlay.Set()
	val := 7
	for val > 0 {
		sound.Battle()
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
