package main

type Memory struct {
	data []uint32
}

func initMemory(memory *Memory, size uint32, defaultValue uint32) {
	memory.data = make([]uint32, size)
	for i := 0; i < len(memory.data); i++ {
		memory.data[i] = defaultValue
	}
	logDebug("INIT", "Memory initialized with default value %d\n", defaultValue)
}

func readMemory(memory *Memory, address uint32) uint32 {
	if address < uint32(len(memory.data)) {
		return memory.data[address]
	}
	return 0
}

func writeMemory(memory *Memory, address uint32, value uint32) {
	wordIndex := address / 4 // Convert byte address to word index
	if wordIndex < uint32(len(memory.data)) {
		memory.data[wordIndex] = value
	}
}

func writeByte(memory *Memory, address uint32, value uint32) {
	wordIndex := address / 4                                                                   // Find the 32-bit word index
	byteOffset := (address % 4) * 8                                                            // Calculate the byte's position (0, 8, 16, or 24 bits)
	mask := uint32(0xFF << byteOffset)                                                         // Create a mask to isolate the byte
	memory.data[wordIndex] = (memory.data[wordIndex] & ^mask) | ((value & 0xFF) << byteOffset) // Clear the byte and write the new value
}

func writeHalfword(memory *Memory, address uint32, value uint32) {
	wordIndex := address / 4                                                                         // Find the 32-bit word index
	halfwordOffset := (address % 4) * 16                                                             // Calculate the halfword's position (0 or 16 bits)
	mask := uint32(0xFFFF << halfwordOffset)                                                         // Create a mask to isolate the halfword
	memory.data[wordIndex] = (memory.data[wordIndex] & ^mask) | ((value & 0xFFFF) << halfwordOffset) // Clear the halfword and write the new value
}

func writeWord(memory *Memory, address uint32, value uint32) {
	wordIndex := address / 4       // Find the 32-bit word index
	memory.data[wordIndex] = value // Write the value to memory
}

func lenMemory(memory *Memory) uint32 {
	return uint32(len(memory.data))
}
