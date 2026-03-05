package sound

import (
	"fmt"
	"os/exec"
)

/*
	Играет простые мелодии в кончоле при наступлении каких-либо событий
*/

const ( // Ноты для проигрывания, Гц
	do  = 261.63
	re  = 293.66
	mi  = 329.63
	fa  = 349.23
	sol = 392.00
	la  = 440.00
	si  = 493.88
)

// Battle -- играет звук начала битвы
func Battle() {
	play(do)
	play(re)
	play(re)
}

// DivWar -- играет звук начала сражения дивизий
func DivWar() {
	play(mi)
	play(fa)
	play(fa)
}

// MineForce -- играет звук ускорения апгрейда шахты
func MineForce() {
	play(sol)
	play(sol)
	play(sol)
}

// ArsenalForce -- играет звук ускорения апгрейда арсенала
func ArsenalForce() {
	play(re)
	play(re)
	play(re)
}

// BankTake -- играет звук забрать серебро
func BankTake() {
	play(si)
	play(la)
	play(sol)
}

// BankForce -- играет звук ускорения апгрейда банка
func BankForce() {
	play(fa)
	play(fa)
	play(fa)
}

// Polygon -- играет звук работы полигона
func Polygon() {
	play(la)
	play(re)
	play(mi)
}

// Shot -- звук выстрела
func Shot() {
	cmd := exec.Command("./beep", "-f=1200.0", "-t=50", "-v=25")
	_ = cmd.Start()
	_ = cmd.Wait()
	cmd = exec.Command("./beep", "-f=1300.0", "-t=50", "-v=25")
	_ = cmd.Start()
	_ = cmd.Wait()
	cmd = exec.Command("./beep", "-f=1400.0", "-t=50", "-v=25")
	_ = cmd.Start()
	_ = cmd.Wait()
}

// Repair -- звук восстановления здоровья
func Repair() {
	cmd := exec.Command("./beep", "-f=1500.0", "-t=50", "-v=25")
	_ = cmd.Start()
	_ = cmd.Wait()
	cmd = exec.Command("./beep", "-f=1750.0", "-t=50", "-v=25")
	_ = cmd.Start()
	_ = cmd.Wait()
	cmd = exec.Command("./beep", "-f=2000.0", "-t=50", "-v=25")
	_ = cmd.Start()
	_ = cmd.Wait()
}

// HelthDown -- звук потери здоровья
func HelthDown() {
	cmd := exec.Command("./beep", "-f=800.0", "-t=50", "-v=25")
	_ = cmd.Start()
	_ = cmd.Wait()
	cmd = exec.Command("./beep", "-f=700.0", "-t=50", "-v=25")
	_ = cmd.Start()
	_ = cmd.Wait()
	cmd = exec.Command("./beep", "-f=6000.0", "-t=50", "-v=25")
	_ = cmd.Start()
	_ = cmd.Wait()
}

// Manver -- звук выполнения манёвра
func Manevr() {
	cmd := exec.Command("./beep", "-f=150.0", "-t=25", "-v=25")
	_ = cmd.Start()
	_ = cmd.Wait()
	cmd = exec.Command("./beep", "-f=150.0", "-t=25", "-v=25")
	_ = cmd.Start()
	_ = cmd.Wait()
	cmd = exec.Command("./beep", "-f=150.0", "-t=25", "-v=25")
	_ = cmd.Start()
	_ = cmd.Wait()
}

func play(note float32) {
	cmd := exec.Command("./beep", "-f="+fmt.Sprintf("%0.2f", note), "-t=100", "-v=25")
	_ = cmd.Start()
	_ = cmd.Wait()
	cmd = exec.Command("./beep", "-f=0", "-t=20", "-v=25")
	_ = cmd.Start()
	_ = cmd.Wait()
}
