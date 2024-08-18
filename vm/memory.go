package vm

import (
	"encoding/binary"
	"log"
	"os"
)

const (
	KBSR uint16 = 0xFE00 // keyboard status register
	KBDR uint16 = 0xFE02 // keyboard data register
)

func (cpu *CPU) memoryWrite(address uint16, val uint16) {
	cpu.memory[address] = val
}

func (cpu *CPU) memoryRead(address uint16) uint16 {
	// When memory is read from KBSR, check the keyboard and update both memory locations
	// If key is pressed, set the most significant bit of KBSR and set KBDR to key
	if address == KBSR {
		checkKey := cpu.checkKey()

		if checkKey != 0 {
			cpu.memory[KBSR] = 1 << 15
			cpu.memory[KBDR] = checkKey
		} else {
			cpu.memory[KBSR] = 0
		}
	}
	return cpu.memory[address]
}

func (cpu *CPU) LoadImage(path string) error {
	bytes, err := os.ReadFile(path)
	if err != nil {
		log.Printf("Cannot load program %s: %s", path, err)
		return err
	}

	// LC-3 programs are big-endian
	// The first 16 bits of the program file specify the address in memory where the program should start
	origin := binary.BigEndian.Uint16(bytes[:2])
	for i := 2; i < len(bytes); i += 2 {
		cpu.memory[origin] = binary.BigEndian.Uint16(bytes[i : i+2])
		origin++
	}
	return nil
}
