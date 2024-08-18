package vm

import (
	"fmt"
	"log"
)

const (
	GETC  uint16 = 0x20 // get character from keyboard, not echoed onto the terminal
	OUT   uint16 = 0x21 // output a character
	PUTS  uint16 = 0x22 // output a word string
	IN    uint16 = 0x23 // get character from keyboard, echoed onto the terminal
	PUTSP uint16 = 0x24 // output a byte string
	HALT  uint16 = 0x25 // halt the program
)

func (cpu *CPU) trap(instr uint16) {
	switch instr & 0xFF {
	case GETC:
		cpu.getc()
	case OUT:
		cpu.out()
	case PUTS:
		cpu.puts()
	case IN:
		cpu.in()
	case PUTSP:
		cpu.putsp()
	case HALT:
		cpu.halt()
	default:
		log.Printf("Trap code not implemented: 0x%04X", instr)
	}
}

func (cpu *CPU) getc() {
	cpu.registers[0] = cpu.checkKey()
	cpu.updateFlags(0)
}

func (cpu *CPU) out() {
	char := cpu.registers[0]
	fmt.Printf("%c", char)
}

func (cpu *CPU) puts() {
	for address := cpu.registers[0]; cpu.memory[address] != 0x00; address++ {
		fmt.Printf("%c", cpu.memory[address])
	}
}

func (cpu *CPU) in() {
	fmt.Printf("Enter a character: ")
	char := cpu.checkKey()
	fmt.Printf("%c", char)
	cpu.registers[0] = char
	cpu.updateFlags(0)
}

func (cpu *CPU) putsp() {
	for address := cpu.registers[0]; cpu.memory[address] != 0x00; address++ {
		value := cpu.memory[address]
		fmt.Printf("%c", value&0xff)
		char := value & 0xff >> 8
		if char != 0 {
			fmt.Printf("%c", char)
		}
	}
}

func (cpu *CPU) halt() {
	cpu.Stop()
}
