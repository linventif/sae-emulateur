package main

import "fmt"

type Opcode struct {
	Type     string
	Encoding Encoding
}

var OpcodeMap = map[uint32]Opcode{
	0b1100011: {"BRANCH", Encodings["SB"]},  // Conditional branches
	0b1100111: {"JALR", Encodings["I"]},     // Jump and link register
	0b0000011: {"LOAD", Encodings["I"]},     // Load instructions
	0b0001111: {"MISC-MEM", Encodings["I"]}, // Misc memory instructions
	0b0010011: {"OP-IMM", Encodings["I"]},   // Integer-Register-Immediate instructions
	0b1110011: {"SYSTEM", Encodings["I"]},   // System instructions
	0b1101111: {"JAL", Encodings["UJ"]},     // Jump and link
	0b0110011: {"OP", Encodings["R"]},       // Register-register operations
	0b0100011: {"STORE", Encodings["S"]},    // Store instructions
	0b0010111: {"AUIPC", Encodings["U"]},    // Add upper immediate to PC
	0b0110111: {"LUI", Encodings["U"]},      // Load upper immediate
}

func GetOpcode(opcode uint32) (Opcode, error) {
	if op, found := OpcodeMap[opcode]; found {
		return op, nil
	}
	return Opcode{}, fmt.Errorf("opcode '%d' not found", opcode)
}

func GetOpcodeFromInstruction(instruction uint32) (Opcode, error) {
	opcode := instruction & 0x7F
	return GetOpcode(opcode)
}
