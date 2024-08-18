package vm

import (
	"github.com/eiannone/keyboard"
	"log"
)

func (cpu *CPU) checkKey() uint16 {
	char, key, err := keyboard.GetSingleKey()
	if err != nil {
		log.Printf("Error getting key, %s", err)
	}
	if key == keyboard.KeyEsc || key == keyboard.KeyCtrlC {
		cpu.Stop()
	}
	return uint16(char)
}
