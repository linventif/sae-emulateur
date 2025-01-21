package main

import (
	"testing"
)

func TestGetOpcode(t *testing.T) {
	tests := []struct {
		opcode     uint32
		expected   string
		shouldFail bool
	}{
		{0b1100011, "BRANCH", false},
		{0b1100111, "JALR", false},
		{0b0000011, "LOAD", false},
		{0b0001111, "MISC-MEM", false},
		{0b0010011, "OP-IMM", false},
		{0b1110011, "SYSTEM", false},
		{0b1101111, "JAL", false},
		{0b0110011, "OP", false},
		{0b0100011, "STORE", false},
		{0b0010111, "AUIPC", false},
		{0b0110111, "LUI", false},
		{0b0000000, "", true}, // Invalid opcode for testing
	}

	for _, test := range tests {
		t.Run(test.expected, func(t *testing.T) {
			op, err := GetOpcode(test.opcode)

			if test.shouldFail {
				if err == nil {
					t.Errorf("expected failure for opcode '%b', but got success", test.opcode)
				}
			} else {
				if err != nil {
					t.Errorf("expected success for opcode '%b', but got error: %v", test.opcode, err)
				} else if op.Type != test.expected {
					t.Errorf("expected opcode type '%s', but got '%s'", test.expected, op.Type)
				}
			}
		})
	}
}

func TestGetOpcodeFromInstruction(t *testing.T) {
	tests := []struct {
		instruction uint32
		expected    string
		shouldFail  bool
	}{
		{0x63, "BRANCH", false},   // BRANCH instruction with opcode 0b1100011
		{0x67, "JALR", false},     // JALR instruction with opcode 0b1100111
		{0x03, "LOAD", false},     // LOAD instruction with opcode 0b0000011
		{0x0F, "MISC-MEM", false}, // MISC-MEM instruction with opcode 0b0001111
		{0x13, "OP-IMM", false},   // OP-IMM instruction with opcode 0b0010011
		{0x73, "SYSTEM", false},   // SYSTEM instruction with opcode 0b1110011
		{0x6F, "JAL", false},      // JAL instruction with opcode 0b1101111
		{0x33, "OP", false},       // OP instruction with opcode 0b0110011
		{0x23, "STORE", false},    // STORE instruction with opcode 0b0100011
		{0x17, "AUIPC", false},    // AUIPC instruction with opcode 0b0010111
		{0x37, "LUI", false},      // LUI instruction with opcode 0b0110111
		{0x00, "", true},          // Invalid instruction for testing
	}

	for _, test := range tests {
		t.Run(test.expected, func(t *testing.T) {
			op, err := GetOpcodeFromInstruction(test.instruction)

			if test.shouldFail {
				if err == nil {
					t.Errorf("expected failure for instruction '%b', but got success", test.instruction)
				}
			} else {
				if err != nil {
					t.Errorf("expected success for instruction '%b', but got error: %v", test.instruction, err)
				} else if op.Type != test.expected {
					t.Errorf("expected opcode type '%s', but got '%s'", test.expected, op.Type)
				}
			}
		})
	}
}
