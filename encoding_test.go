package main

import (
	"testing"
)

func encodeIType(rd, rs1, imm, funct3 uint32, funct7 uint32) uint32 {
	// [31:20] imm[11:0] [19:15] rs1 [14:12] funct3 [11:7] rd [6:0] opcode
	var instruction = imm<<20 | rs1<<15 | funct3<<12 | rd<<7 | 0b0010011
	if funct7 != 0 {
		instruction |= funct7 << 25
	}
	return instruction
}

func encodeRType(rd, rs1, rs2, funct3, funct7 uint32) uint32 {
	// [31:25] funct7 [24:20] rs2 [19:15] rs1 [14:12] funct3 [11:7] rd [6:0] opcode
	return funct7<<25 | rs2<<20 | rs1<<15 | funct3<<12 | rd<<7 | 0b0110011
}

func encodeSType(rs1, rs2, imm, funct3 uint32) uint32 {
	// [31:25] imm[11:5] [24:20] rs2 [19:15] rs1 [14:12] funct3 [11:7] imm[4:0] [6:0] opcode
	return imm<<25 | rs2<<20 | rs1<<15 | funct3<<12 | 0b0100011
}

func encodeUType(rd, imm uint32) uint32 {
	// [31:12] imm[31:12] [11:7] rd [6:0] opcode
	return imm<<12 | rd<<7 | 0b0110111
}

func encodeJType(rd, imm uint32) uint32 {
	// [31:12] imm[20:1] [11] imm[11] [10:7] rd [6:0] opcode
	return imm<<12 | rd<<7 | 0b1101111
}

func encodeBType(rs1, rs2, imm, funct3 uint32) uint32 {
	// [31] imm[12] [30:25] imm[10:5] [24:21] funct3 [20] imm[11] [19:12] imm[4:1] [11:8] rs2 [7:0] rs1 [6:0] opcode
	return imm<<31 | (imm>>1&0x3F)<<25 | funct3<<12 | (imm>>11&0x1)<<7 | rs2<<7 | rs1<<15 | 0b1100011
}

func TestEncodings(t *testing.T) {
	var memorySize uint32 = 512 * 1024
	var registerDefault uint32 = 0
	var cpu CPUState
	var memory Memory
	var startAddress uint32 = 0

	initMemory(&memory, memorySize, 0)
	initCPUState(&cpu, startAddress, registerDefault)

	tests := []struct {
		name        string
		encoding    string
		instruction uint32
		expected    string
	}{
		// OP (ADD, SUB, ...)
		{
			name:        "R-type ADD",
			encoding:    "R",
			instruction: encodeRType(2, 1, 2, 0b000, 0b0000000), // ADD x2, x1, x2
			expected:    "ADD x2, x1, x2\n",
		},
		{
			name:        "R-type SUB",
			encoding:    "R",
			instruction: encodeRType(2, 1, 2, 0b000, 0b0100000), // SUB x2, x1, x2
			expected:    "SUB x2, x1, x2\n",
		},
		{
			name:        "R-type SLL",
			encoding:    "R",
			instruction: encodeRType(2, 1, 2, 0b001, 0), // SLL x2, x1, x2
			expected:    "SLL x2, x1, x2\n",
		},
		{
			name:        "R-type SLT",
			encoding:    "R",
			instruction: encodeRType(2, 1, 2, 0b010, 0), // SLT x2, x1, x2
			expected:    "SLT x2, x1, x2\n",
		},
		{
			name:        "R-type SLTU",
			encoding:    "R",
			instruction: encodeRType(2, 1, 2, 0b011, 0), // SLTU x2, x1, x2
			expected:    "SLTU x2, x1, x2\n",
		},
		{
			name:        "R-type XOR",
			encoding:    "R",
			instruction: encodeRType(2, 1, 2, 0b100, 0), // XOR x2, x1, x2
			expected:    "XOR x2, x1, x2\n",
		},
		// OP-IMM (ADDI, SLTI, ...)
		{
			name:        "I-type ADDI",
			encoding:    "I",
			instruction: encodeIType(2, 1, 5, 0b000, 0), // ADDI x2, x1, 5
			expected:    "ADDI x2, x1, 5\n",
		},
		{
			name:        "I-type SLTI",
			encoding:    "I",
			instruction: encodeIType(2, 1, 5, 0b010, 0), // SLTI x2, x1, 5
			expected:    "SLTI x2, x1, 5\n",
		},
		{
			name:        "I-type SLTIU",
			encoding:    "I",
			instruction: encodeIType(2, 1, 5, 0b011, 0), // SLTIU x2, x1, 5
			expected:    "SLTIU x2, x1, 5\n",
		},
		{
			name:        "I-type XORI",
			encoding:    "I",
			instruction: encodeIType(2, 1, 5, 0b100, 0), // XORI x2, x1, 5
			expected:    "XORI x2, x1, 5\n",
		},
		{
			name:        "I-type SLLI",
			encoding:    "I",
			instruction: encodeIType(2, 1, 5, 0b001, 0), // SLLI x2, x1, 5
			expected:    "SLLI x2, x1, 5\n",
		},
		{
			name:        "I-type SRLI",
			encoding:    "I",
			instruction: encodeIType(2, 1, 5, 0b101, 0b0000000), // SRLI x2, x1, 5
			expected:    "SRLI x2, x1, 5\n",
		},
		{
			name:        "I-type SRAI",
			encoding:    "I",
			instruction: encodeIType(2, 1, 5, 0b101, 0b0100000), // SRAI x2, x1, 5
			expected:    "SRAI x2, x1, 5\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Récupérer l'encodage correspondant
			encoding, ok := Encodings[test.encoding]
			if !ok {
				t.Fatalf("encoding '%s' not found", test.encoding)
			}

			// Exécuter la fonction de décodage
			opcode, _ := GetOpcodeFromInstruction(test.instruction)
			result := encoding.Decode(opcode, test.instruction, &cpu, &memory)

			// Vérifier le résultat
			if result != test.expected {
				t.Errorf("expected '%s', got '%s'", test.expected, result)
			}
		})
	}
}
