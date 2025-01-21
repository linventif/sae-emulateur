package main

import (
	"testing"
)

func TestInstructions(t *testing.T) {
	var memorySize uint32 = 512 * 1024
	var registerDefault uint32 = 0
	var cpu CPUState
	var memory Memory
	var startAddress uint32 = 0

	tests := []struct {
		name         string
		instruction  Instruction
		args         []uint32
		defaultRegs  map[uint32]uint32
		defaultMem   map[uint32]uint32
		expectedPC   uint32
		expectedRegs map[uint32]uint32
		expectedMem  map[uint32]uint32
	}{
		// OP (ADD, SUB, ...)
		{
			name:         "ADD",
			instruction:  Instructions[[4]uint32{0b0110011, 0b000, 0, 0}],
			args:         []uint32{2, 1, 3}, // ADD x2, x1, x3 => x2 = x1 + x3
			defaultRegs:  map[uint32]uint32{1: 1, 3: 2},
			defaultMem:   map[uint32]uint32{},
			expectedPC:   0,
			expectedRegs: map[uint32]uint32{2: 3},
			expectedMem:  map[uint32]uint32{},
		},
		{
			name:         "SUB",
			instruction:  Instructions[[4]uint32{0b0110011, 0b000, 0b0100000}],
			args:         []uint32{2, 1, 3}, // SUB x2, x1, x3 => x2 = x1 - x3
			defaultRegs:  map[uint32]uint32{1: 4, 3: 2},
			defaultMem:   map[uint32]uint32{},
			expectedPC:   0,
			expectedRegs: map[uint32]uint32{2: 2},
			expectedMem:  map[uint32]uint32{},
		},
		{
			name:         "SLL",
			instruction:  Instructions[[4]uint32{0b0110011, 0b001, 0, 0}],
			args:         []uint32{2, 1, 3}, // SLL x2, x1, x3 => x2 = x1 << x3
			defaultRegs:  map[uint32]uint32{1: 1, 3: 2},
			defaultMem:   map[uint32]uint32{},
			expectedPC:   0,
			expectedRegs: map[uint32]uint32{2: 4},
			expectedMem:  map[uint32]uint32{},
		},
		{
			name:         "SLT",
			instruction:  Instructions[[4]uint32{0b0110011, 0b010, 0, 0}],
			args:         []uint32{2, 1, 3}, // SLT x2, x1, x3 => x2 = (x1 < x3) ? 1 : 0
			defaultRegs:  map[uint32]uint32{1: 1, 3: 2},
			defaultMem:   map[uint32]uint32{},
			expectedPC:   0,
			expectedRegs: map[uint32]uint32{2: 1},
			expectedMem:  map[uint32]uint32{},
		},
		{
			name:         "SLTU",
			instruction:  Instructions[[4]uint32{0b0110011, 0b011, 0, 0}],
			args:         []uint32{2, 1, 3}, // SLTU x2, x1, x3 => x2 = (x1 < x3) ? 1 : 0
			defaultRegs:  map[uint32]uint32{1: 1, 3: 2},
			defaultMem:   map[uint32]uint32{},
			expectedPC:   0,
			expectedRegs: map[uint32]uint32{2: 1},
			expectedMem:  map[uint32]uint32{},
		},
		{
			name:         "XOR",
			instruction:  Instructions[[4]uint32{0b0110011, 0b100, 0, 0}],
			args:         []uint32{2, 1, 3}, // XOR x2, x1, x3 => x2 = x1 ^ x3
			defaultRegs:  map[uint32]uint32{1: 1, 3: 2},
			defaultMem:   map[uint32]uint32{},
			expectedPC:   0,
			expectedRegs: map[uint32]uint32{2: 3},
			expectedMem:  map[uint32]uint32{},
		},
		{
			name:         "SRL",
			instruction:  Instructions[[4]uint32{0b0110011, 0b101, 0b0000000, 0}],
			args:         []uint32{2, 1, 3}, // SRL x2, x1, x3 => x2 = x1 >> x3
			defaultRegs:  map[uint32]uint32{1: 8, 3: 1},
			defaultMem:   map[uint32]uint32{},
			expectedPC:   0,
			expectedRegs: map[uint32]uint32{2: 4},
			expectedMem:  map[uint32]uint32{},
		},
		{
			name:         "SRA",
			instruction:  Instructions[[4]uint32{0b0110011, 0b101, 0b0100000, 0}],
			args:         []uint32{2, 1, 3}, // SRA x2, x1, x3 => x2 = x1 >> x3 (arithmetic)
			defaultRegs:  map[uint32]uint32{1: 8, 3: 1},
			defaultMem:   map[uint32]uint32{},
			expectedPC:   0,
			expectedRegs: map[uint32]uint32{2: 4},
			expectedMem:  map[uint32]uint32{},
		},
		{
			name:         "OR",
			instruction:  Instructions[[4]uint32{0b0110011, 0b110, 0, 0}],
			args:         []uint32{2, 1, 3}, // OR x2, x1, x3 => x2 = x1 | x3
			defaultRegs:  map[uint32]uint32{1: 1, 3: 2},
			defaultMem:   map[uint32]uint32{},
			expectedPC:   0,
			expectedRegs: map[uint32]uint32{2: 3},
			expectedMem:  map[uint32]uint32{},
		},
		{
			name:         "AND",
			instruction:  Instructions[[4]uint32{0b0110011, 0b111, 0, 0}],
			args:         []uint32{2, 1, 3}, // AND x2, x1, x3 => x2 = x1 & x3
			defaultRegs:  map[uint32]uint32{1: 1, 3: 2},
			defaultMem:   map[uint32]uint32{},
			expectedPC:   0,
			expectedRegs: map[uint32]uint32{2: 0},
			expectedMem:  map[uint32]uint32{},
		},
		// OP-IMM (ADDI, SLTI, ...)
		{
			name:         "ADDI",
			instruction:  Instructions[[4]uint32{0b0010011, 0b000, 0, 0}],
			args:         []uint32{2, 1, 5}, // ADDI x2, x1, 5 => x2 = x1 + 5
			defaultRegs:  map[uint32]uint32{1: 0},
			defaultMem:   map[uint32]uint32{},
			expectedPC:   0,
			expectedRegs: map[uint32]uint32{2: 5},
			expectedMem:  map[uint32]uint32{},
		},
		{
			name:         "SLTI",
			instruction:  Instructions[[4]uint32{0b0010011, 0b010, 0, 0}],
			args:         []uint32{2, 1, 5}, // SLTI x2, x1, 5 => x2 = (x1 < 5) ? 1 : 0
			defaultRegs:  map[uint32]uint32{1: 4},
			defaultMem:   map[uint32]uint32{},
			expectedPC:   0,
			expectedRegs: map[uint32]uint32{2: 1},
			expectedMem:  map[uint32]uint32{},
		},
		{
			name:         "SLTIU",
			instruction:  Instructions[[4]uint32{0b0010011, 0b011, 0, 0}],
			args:         []uint32{2, 1, 5},       // SLTIU x2, x1, 5 => x2 = (x1 < 5) ? 1 : 0
			defaultRegs:  map[uint32]uint32{1: 4}, // x1 = 4
			defaultMem:   map[uint32]uint32{},
			expectedPC:   0,
			expectedRegs: map[uint32]uint32{2: 1}, // x2 = 1 car 4 < 5
			expectedMem:  map[uint32]uint32{},
		},
		{
			name:         "XORI",
			instruction:  Instructions[[4]uint32{0b0010011, 0b100, 0, 0}],
			args:         []uint32{2, 1, 5}, // XORI x2, x1, 5 => x2 = x1 ^ 5
			defaultRegs:  map[uint32]uint32{1: 1},
			defaultMem:   map[uint32]uint32{},
			expectedPC:   0,
			expectedRegs: map[uint32]uint32{2: 4},
			expectedMem:  map[uint32]uint32{},
		},
		{
			name:         "ORI",
			instruction:  Instructions[[4]uint32{0b0010011, 0b110, 0, 0}],
			args:         []uint32{2, 1, 5}, // ORI x2, x1, 5 => x2 = x1 | 5
			defaultRegs:  map[uint32]uint32{1: 1},
			defaultMem:   map[uint32]uint32{},
			expectedPC:   0,
			expectedRegs: map[uint32]uint32{2: 5},
			expectedMem:  map[uint32]uint32{},
		},
		{
			name:         "ANDI",
			instruction:  Instructions[[4]uint32{0b0010011, 0b111, 0, 0}],
			args:         []uint32{2, 1, 5}, // ANDI x2, x1, 5 => x2 = x1 & 5
			defaultRegs:  map[uint32]uint32{1: 7},
			defaultMem:   map[uint32]uint32{},
			expectedPC:   0,
			expectedRegs: map[uint32]uint32{2: 5},
			expectedMem:  map[uint32]uint32{},
		},
		{
			name:         "SLLI",
			instruction:  Instructions[[4]uint32{0b0010011, 0b001, 0, 0}],
			args:         []uint32{2, 1, 5}, // SLLI x2, x1, 5 => x2 = x1 << 5
			defaultRegs:  map[uint32]uint32{1: 1},
			defaultMem:   map[uint32]uint32{},
			expectedPC:   0,
			expectedRegs: map[uint32]uint32{2: 32},
			expectedMem:  map[uint32]uint32{},
		},
		{
			name:         "SRLI",
			instruction:  Instructions[[4]uint32{0b0010011, 0b101, 0b0000000, 0}],
			args:         []uint32{2, 1, 5}, // SRLI x2, x1, 5 => x2 = x1 >> 5
			defaultRegs:  map[uint32]uint32{1: 32},
			defaultMem:   map[uint32]uint32{},
			expectedPC:   0,
			expectedRegs: map[uint32]uint32{2: 1},
			expectedMem:  map[uint32]uint32{},
		},
		{
			name:         "SRAI",
			instruction:  Instructions[[4]uint32{0b0010011, 0b101, 0b0100000, 0}],
			args:         []uint32{2, 1, 5}, // SRAI x2, x1, 5 => x2 = x1 >> 5 (arithmetic)
			defaultRegs:  map[uint32]uint32{1: 32},
			defaultMem:   map[uint32]uint32{},
			expectedPC:   0,
			expectedRegs: map[uint32]uint32{2: 1},
			expectedMem:  map[uint32]uint32{},
		},
		// Store
		{
			name:         "SB",
			instruction:  Instructions[[4]uint32{0b0100011, 0b000, 0, 0}],
			args:         []uint32{1, 2, 0}, // SB x2, 0(x1)
			defaultRegs:  map[uint32]uint32{1: 4, 2: 0x12345678},
			defaultMem:   map[uint32]uint32{1: 0xFFFFFFFF}, // Memory initialized with all bits set
			expectedPC:   0,
			expectedRegs: map[uint32]uint32{},
			expectedMem:  map[uint32]uint32{1: 0xFFFFFF78}, // Only the last byte is stored
		},
		{
			name:         "SH", // Store halfword, it
			instruction:  Instructions[[4]uint32{0b0100011, 0b001, 0, 0}],
			args:         []uint32{1, 2, 2},                      // SH x2, 2(x1)
			defaultRegs:  map[uint32]uint32{1: 0, 2: 0x12345678}, // x1 = base address, x2 = value to store
			defaultMem:   map[uint32]uint32{1: 0xFFFFFFFF},       // Memory initialized with all bits set
			expectedPC:   0,
			expectedRegs: map[uint32]uint32{},
			expectedMem:  map[uint32]uint32{0: 0x12345678}, // Only the last two bytes are stored
		},
		{
			name:         "SW", // Store Word, it will store 4 bytes in memory starting from the address x1 + 4
			instruction:  Instructions[[4]uint32{0b0100011, 0b010, 0, 0}],
			args:         []uint32{1, 2, 4},                      // SW x2, 4(x1)
			defaultRegs:  map[uint32]uint32{1: 0, 2: 0x12345678}, // x1 = base address, x2 = value to store
			defaultMem:   map[uint32]uint32{1: 0xFFFFFFFF},       // Memory initialized with all bits set
			expectedPC:   0,
			expectedRegs: map[uint32]uint32{},
			expectedMem:  map[uint32]uint32{1: 0x12345678}, // Value is stored at the base address
		},
		// Load
		{
			name:         "LB",
			instruction:  Instructions[[4]uint32{0b0000011, 0b000, 0, 0}],
			args:         []uint32{1, 2, 4}, // LB x1, 4(x2) => x1 = memory[x2+4]
			defaultRegs:  map[uint32]uint32{2: 0},
			defaultMem:   map[uint32]uint32{},
			expectedPC:   0,
			expectedRegs: map[uint32]uint32{1: 0},
			expectedMem:  map[uint32]uint32{4: 0},
		},
		{
			name:         "LH",
			instruction:  Instructions[[4]uint32{0b0000011, 0b001, 0, 0}],
			args:         []uint32{1, 2, 4}, // LH x1, 4(x2) => x1 = memory[x2+4]
			defaultRegs:  map[uint32]uint32{2: 0},
			defaultMem:   map[uint32]uint32{},
			expectedPC:   0,
			expectedRegs: map[uint32]uint32{1: 0},
			expectedMem:  map[uint32]uint32{4: 0},
		},
		{
			name:         "LW",
			instruction:  Instructions[[4]uint32{0b0000011, 0b010, 0, 0}],
			args:         []uint32{1, 2, 4}, // LW x1, 4(x2) => x1 = memory[x2+4]
			defaultRegs:  map[uint32]uint32{2: 0},
			defaultMem:   map[uint32]uint32{},
			expectedPC:   0,
			expectedRegs: map[uint32]uint32{1: 0},
			expectedMem:  map[uint32]uint32{4: 0},
		},
		{
			name:         "LBU",
			instruction:  Instructions[[4]uint32{0b0000011, 0b100, 0, 0}],
			args:         []uint32{1, 2, 4}, // LBU x1, 4(x2) => x1 = memory[x2+4]
			defaultRegs:  map[uint32]uint32{2: 0},
			defaultMem:   map[uint32]uint32{},
			expectedPC:   0,
			expectedRegs: map[uint32]uint32{1: 0},
			expectedMem:  map[uint32]uint32{4: 0},
		},
		{
			name:         "LHU",
			instruction:  Instructions[[4]uint32{0b0000011, 0b101, 0, 0}],
			args:         []uint32{1, 2, 4}, // LHU x1, 4(x2) => x1 = memory[x2+4]
			defaultRegs:  map[uint32]uint32{2: 0},
			defaultMem:   map[uint32]uint32{},
			expectedPC:   0,
			expectedRegs: map[uint32]uint32{1: 0},
			expectedMem:  map[uint32]uint32{4: 0},
		},
		// Branch
		{
			name:         "BEQ",
			instruction:  Instructions[[4]uint32{0b1100011, 0b000, 0, 0}],
			args:         []uint32{1, 2, 4}, // BEQ x1, x2, 4 => if x1 == x2 { pc += 4 }
			defaultRegs:  map[uint32]uint32{1: 1, 2: 1},
			defaultMem:   map[uint32]uint32{},
			expectedPC:   4,
			expectedRegs: map[uint32]uint32{},
			expectedMem:  map[uint32]uint32{},
		},
		{
			name:         "BNE",
			instruction:  Instructions[[4]uint32{0b1100011, 0b001, 0, 0}],
			args:         []uint32{1, 2, 4}, // BNE x1, x2, 4 => if x1 != x2 { pc += 4 }
			defaultRegs:  map[uint32]uint32{1: 1, 2: 2},
			defaultMem:   map[uint32]uint32{},
			expectedPC:   4,
			expectedRegs: map[uint32]uint32{},
			expectedMem:  map[uint32]uint32{},
		},
		{
			name:         "BLT",
			instruction:  Instructions[[4]uint32{0b1100011, 0b100, 0, 0}],
			args:         []uint32{1, 2, 4}, // BLT x1, x2, 4 => if x1 < x2 { pc += 4 }
			defaultRegs:  map[uint32]uint32{1: 1, 2: 2},
			defaultMem:   map[uint32]uint32{},
			expectedPC:   4,
			expectedRegs: map[uint32]uint32{},
			expectedMem:  map[uint32]uint32{},
		},
		{
			name:         "BGE",
			instruction:  Instructions[[4]uint32{0b1100011, 0b101, 0, 0}],
			args:         []uint32{1, 2, 4}, // BGE x1, x2, 4 => if x1 >= x2 { pc += 4 }
			defaultRegs:  map[uint32]uint32{1: 2, 2: 1},
			defaultMem:   map[uint32]uint32{},
			expectedPC:   4,
			expectedRegs: map[uint32]uint32{},
			expectedMem:  map[uint32]uint32{},
		},
		{
			name:         "BLTU",
			instruction:  Instructions[[4]uint32{0b1100011, 0b110, 0, 0}],
			args:         []uint32{1, 2, 4}, // BLTU x1, x2, 4 => if x1 < x2 { pc += 4 }
			defaultRegs:  map[uint32]uint32{1: 1, 2: 2},
			defaultMem:   map[uint32]uint32{},
			expectedPC:   4,
			expectedRegs: map[uint32]uint32{},
			expectedMem:  map[uint32]uint32{},
		},
		{
			name:         "BGEU",
			instruction:  Instructions[[4]uint32{0b1100011, 0b111, 0, 0}],
			args:         []uint32{1, 2, 4}, // BGEU x1, x2, 4 => if x1 >= x2 { pc += 4 }
			defaultRegs:  map[uint32]uint32{1: 2, 2: 1},
			defaultMem:   map[uint32]uint32{},
			expectedPC:   4,
			expectedRegs: map[uint32]uint32{},
			expectedMem:  map[uint32]uint32{},
		},
		// LUI
		{
			name:         "LUI",
			instruction:  Instructions[[4]uint32{0b0110111, 0, 0, 0}],
			args:         []uint32{1, 0x12345}, // LUI x1, 0x12345 => x1 = 0x12345000
			defaultRegs:  map[uint32]uint32{},
			defaultMem:   map[uint32]uint32{},
			expectedPC:   0,
			expectedRegs: map[uint32]uint32{1: 0x12345000},
			expectedMem:  map[uint32]uint32{},
		},
		// AUIPC
		{
			name:         "AUIPC",
			instruction:  Instructions[[4]uint32{0b0010111, 0, 0, 0}],
			args:         []uint32{1, 0x12345}, // AUIPC x1, 0x12345 => x1 = pc + 0x12345000
			defaultRegs:  map[uint32]uint32{},
			defaultMem:   map[uint32]uint32{},
			expectedPC:   0,
			expectedRegs: map[uint32]uint32{1: 0x12345000},
			expectedMem:  map[uint32]uint32{},
		},
		// JAL
		{
			name:         "JAL",
			instruction:  Instructions[[4]uint32{0b1101111, 0, 0, 0}],
			args:         []uint32{1, 4}, // JAL x1, 4 => x1 = pc + 4; pc += 4
			defaultRegs:  map[uint32]uint32{},
			defaultMem:   map[uint32]uint32{},
			expectedPC:   4,
			expectedRegs: map[uint32]uint32{1: 4},
			expectedMem:  map[uint32]uint32{},
		},
		// JALR
		{
			name:         "JALR",
			instruction:  Instructions[[4]uint32{0b1100111, 0, 0, 0}],
			args:         []uint32{1, 2, 4}, // JALR x1, x2, 4 => x1 = pc + 4; pc = x2 + 4
			defaultRegs:  map[uint32]uint32{2: 0},
			defaultMem:   map[uint32]uint32{},
			expectedPC:   4,
			expectedRegs: map[uint32]uint32{1: 4},
			expectedMem:  map[uint32]uint32{},
		},
		// SYSTEM
		{
			name:         "ECALL",
			instruction:  Instructions[[4]uint32{0b1110011, 0, 0, 0}],
			args:         []uint32{},
			defaultRegs:  map[uint32]uint32{},
			defaultMem:   map[uint32]uint32{},
			expectedPC:   0,
			expectedRegs: map[uint32]uint32{},
			expectedMem:  map[uint32]uint32{},
		},
		{
			name:         "EBREAK",
			instruction:  Instructions[[4]uint32{0b1110011, 0, 0, 1}],
			args:         []uint32{},
			defaultRegs:  map[uint32]uint32{},
			defaultMem:   map[uint32]uint32{},
			expectedPC:   0,
			expectedRegs: map[uint32]uint32{},
			expectedMem:  map[uint32]uint32{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			initMemory(&memory, memorySize, 0)
			initCPUState(&cpu, startAddress, registerDefault)

			for reg, value := range test.defaultRegs {
				writeRegister(&cpu, reg, value)
			}

			for addr, value := range test.defaultMem {
				if addr >= memorySize {
					t.Errorf("memory address %d out of bounds", addr)
					continue
				}
				writeMemory(&memory, addr, value)
			}

			test.instruction.Exec(&cpu, &memory, test.args...)

			if cpu.pc != test.expectedPC {
				t.Errorf("expected pc=%d, got pc=%d", test.expectedPC, cpu.pc)
			}

			for reg, expected := range test.expectedRegs {
				if cpu.x[reg] != expected {
					t.Errorf("expected x%d=%d, got x%d=%d", reg, expected, reg, cpu.x[reg])
				}
			}

			for addr, expected := range test.expectedMem {
				if addr >= memorySize {
					t.Errorf("memory address %d out of bounds", addr)
					continue
				}
				if readMemory(&memory, addr) != expected {
					t.Errorf("expected memory[%d]=%d, got memory[%d]=%d", addr, expected, addr, readMemory(&memory, addr))
				}
			}
		})
	}
}
