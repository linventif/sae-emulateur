package main

type CPUState struct {
	x  [32]uint32
	pc uint32
}

func readRegister(state *CPUState, reg uint32) uint32 {
	if reg == 0 {
		return 0
	}
	return state.x[reg]
}

func writeRegister(state *CPUState, reg uint32, value uint32) {
	if reg != 0 {
		state.x[reg] = value
	}
}

func initCPUState(state *CPUState, firstInstruction uint32, defaultMemoryValue uint32) {
	for i := 0; i < 32; i++ {
		writeRegister(state, uint32(i), defaultMemoryValue)
	}
	state.pc = firstInstruction
	logDebug("INIT", "CPU state initialized with default memory value %d\n", defaultMemoryValue)
}
