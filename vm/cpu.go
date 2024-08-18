package vm

type CPU struct {
	memory    [1 << 16]uint16 // 65536 locations
	isRunning bool

	registers [8]uint16 // general purpose
	pc        uint16    // program counter
	cond      uint16    // condition flags
}

func NewCPU() *CPU {
	cpu := &CPU{}
	// set the pc to starting position, 0x3000 is the default
	cpu.pc = 0x3000
	return cpu
}

func (cpu *CPU) Run() {
	cpu.isRunning = true
	for cpu.isRunning {
		cpu.executeInstruction()
		// increment memory cycle counter (MMC)
		// to keep track of how many instructions have been executed
		cpu.memory[0xFFFF]++
	}
}

func (cpu *CPU) Stop() {
	cpu.isRunning = false
}
