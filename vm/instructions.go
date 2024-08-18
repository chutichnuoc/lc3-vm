package vm

import (
	"log"
)

const (
	POS = 1 << iota // P (1 << 0)
	ZRO             // Z (1 << 1)
	NEG             // N (1 << 2)
)

const (
	BR   = iota // branch
	ADD         // add
	LD          // load
	ST          // store
	JSR         // jump register
	AND         // bitwise and
	LDR         // load register
	STR         // store register
	RTI         // unused
	NOT         // bitwise not
	LDI         // load indirect
	STI         // store indirect
	JMP         // jump
	RES         // reserved (unused)
	LEA         // load effective address
	TRAP        // execute trap
)

func (cpu *CPU) executeInstruction() {
	var instr = cpu.memoryRead(cpu.pc)
	var op = instr >> 12
	cpu.pc++

	switch op {
	case BR:
		cpu.branch(instr)
	case ADD:
		cpu.add(instr)
	case LD:
		cpu.load(instr)
	case ST:
		cpu.store(instr)
	case JSR:
		cpu.jumpRegister(instr)
	case AND:
		cpu.bitwiseAnd(instr)
	case LDR:
		cpu.loadRegister(instr)
	case STR:
		cpu.storeRegister(instr)
	case RTI: // unused
	case NOT:
		cpu.bitwiseNot(instr)
	case LDI:
		cpu.loadIndirect(instr)
	case STI:
		cpu.storeIndirect(instr)
	case JMP:
		cpu.jump(instr)
	case RES: // reserved (unused)
	case LEA:
		cpu.loadEffectiveAddress(instr)
	case TRAP:
		cpu.trap(instr)
	default: //  not implemented
		log.Printf("Operation code is not implemented: 0x%04X", instr)
	}
}

func (cpu *CPU) branch(instr uint16) {
	var pcOffset = signExtend(instr&0x1FF, 9)
	var condFlag = (instr >> 9) & 0x7
	if (condFlag & cpu.cond) != 0 {
		cpu.pc += pcOffset
	}
}

func (cpu *CPU) add(instr uint16) {
	var r0 = (instr >> 9) & 0x7
	var r1 = (instr >> 6) & 0x7
	var immFlag = (instr >> 5) & 0x1
	if immFlag != 0 {
		var imm5 = signExtend(instr&0x1F, 5)
		cpu.registers[r0] = cpu.registers[r1] + imm5
	} else {
		var r2 = instr & 0x7
		cpu.registers[r0] = cpu.registers[r1] + cpu.registers[r2]
	}
	cpu.updateFlags(r0)
}

func (cpu *CPU) load(instr uint16) {
	var r0 = (instr >> 9) & 0x7
	var pcOffset = signExtend(instr&0x1FF, 9)
	cpu.registers[r0] = cpu.memoryRead(cpu.pc + pcOffset)
	cpu.updateFlags(r0)
}

func (cpu *CPU) store(instr uint16) {
	var r0 = (instr >> 9) & 0x7
	var pcOffset = signExtend(instr&0x1FF, 9)
	cpu.memoryWrite(cpu.pc+pcOffset, cpu.registers[r0])
}

func (cpu *CPU) jumpRegister(instr uint16) {
	var longFlag = (instr >> 11) & 1
	cpu.registers[7] = cpu.pc
	if longFlag != 0 {
		var longPcOffset = signExtend(instr&0x7FF, 11)
		cpu.pc += longPcOffset
	} else {
		var r1 = (instr >> 6) & 0x7
		cpu.pc = cpu.registers[r1]
	}
}

func (cpu *CPU) bitwiseAnd(instr uint16) {
	var r0 = (instr >> 9) & 0x7
	var r1 = (instr >> 6) & 0x7
	var immFlag = (instr >> 5) & 0x1
	if immFlag != 0 {
		var imm5 = signExtend(instr&0x1F, 5)
		cpu.registers[r0] = cpu.registers[r1] & imm5
	} else {
		var r2 = instr & 0x7
		cpu.registers[r0] = cpu.registers[r1] & cpu.registers[r2]
	}
}

func (cpu *CPU) loadRegister(instr uint16) {
	var r0 = (instr >> 9) & 0x7
	var r1 = (instr >> 6) & 0x7
	var offset = signExtend(instr&0x3F, 6)
	cpu.registers[r0] = cpu.memoryRead(cpu.registers[r1] + offset)
	cpu.updateFlags(r0)
}

func (cpu *CPU) storeRegister(instr uint16) {
	var r0 = (instr >> 9) & 0x7
	var r1 = (instr >> 6) & 0x7
	var offset = signExtend(instr&0x3F, 6)
	cpu.memoryWrite(cpu.registers[r1]+offset, cpu.registers[r0])
}

func (cpu *CPU) bitwiseNot(instr uint16) {
	var r0 = (instr >> 9) & 0x7
	var r1 = (instr >> 6) & 0x7
	cpu.registers[r0] = ^cpu.registers[r1]
	cpu.updateFlags(r0)
}

func (cpu *CPU) loadIndirect(instr uint16) {
	var r0 = (instr >> 9) & 0x7
	var pcOffset = signExtend(instr&0x1FF, 9)
	cpu.registers[r0] = cpu.memoryRead(cpu.memoryRead(cpu.pc + pcOffset))
	cpu.updateFlags(r0)
}

func (cpu *CPU) storeIndirect(instr uint16) {
	var r0 = (instr >> 9) & 0x7
	var pcOffset = signExtend(instr&0x1FF, 9)
	cpu.memoryWrite(cpu.memoryRead(cpu.pc+pcOffset), cpu.registers[r0])
}

func (cpu *CPU) jump(instr uint16) {
	var r1 = (instr >> 6) & 0x7
	cpu.pc = cpu.registers[r1]
}

func (cpu *CPU) loadEffectiveAddress(instr uint16) {
	var r0 = (instr >> 9) & 0x7
	var pcOffset = signExtend(instr&0x1FF, 9)
	cpu.registers[r0] = cpu.pc + pcOffset
	cpu.updateFlags(r0)
}

func signExtend(x uint16, bitCount int) uint16 {
	if ((x >> (bitCount - 1)) & 1) != 0 {
		x |= 0xFFFF << bitCount
	}
	return x
}

func (cpu *CPU) updateFlags(r uint16) {
	if cpu.registers[r] == 0 {
		cpu.cond = ZRO
	} else if cpu.registers[r]>>15 != 0 { // a 1 in the left-most bit indicates negative
		cpu.cond = NEG
	} else {
		cpu.cond = POS
	}
}
