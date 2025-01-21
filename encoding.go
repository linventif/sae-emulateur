package main

import "fmt"

func decodeI(opcode Opcode, instruction uint32, cpu *CPUState, memory *Memory) string {
	// [31:20] imm[11:0] [19:15] rs1 [14:12] funct3 [11:7] rd [6:0] opcode
	imm := uint32(instruction) >> 20
	rd := (instruction >> 7) & 0x1F
	rs1 := (instruction >> 15) & 0x1F

	funct3 := (instruction >> 12) & 0x7
	funct7 := uint32(0)
	funct12 := uint32(0)

	if opcode.Type == "OP-IMM" && funct3 == 0b101 {
		funct7 = instruction >> 25
	}
	if opcode.Type == "SYSTEM" {
		funct12 = instruction >> 20
	}

	inst, err := FindInstruction(instruction, funct3, funct7, funct12)
	if err == nil {
		inst.Exec(cpu, memory, rd, rs1, imm)
		if opcode.Type == "OP-IMM" {
			return fmt.Sprintf("%s x%d, x%d, %d\n", inst.Name, rd, rs1, imm)
		} else if opcode.Type == "SYSTEM" {
			return fmt.Sprintf("%s x%d, %d\n", inst.Name, rd, imm)
		}
		return fmt.Sprintf("%s x%d, x%d, %d\n", inst.Name, rd, rs1, imm)
	} else {
		return fmt.Sprintf("%s\n", err.Error())
	}
}

func decodeR(opcode Opcode, instruction uint32, cpu *CPUState, memory *Memory) string {
	// [31:25] funct7 [24:20] rs2 [19:15] rs1 [14:12] funct3 [11:7] rd [6:0] opcode
	rd := (instruction >> 7) & 0x1F
	rs1 := (instruction >> 15) & 0x1F
	rs2 := (instruction >> 20) & 0x1F

	funct3 := (instruction >> 12) & 0x7
	funct7 := instruction >> 25

	inst, err := FindInstruction(instruction, funct3, funct7, 0)
	if err == nil {
		inst.Exec(cpu, memory, rd, rs1, rs2)
		return fmt.Sprintf("%s x%d, x%d, x%d\n", inst.Name, rd, rs1, rs2)
	} else {
		return fmt.Sprintf("%s\n", err.Error())
	}
}

func decodeS(opcode Opcode, instruction uint32, cpu *CPUState, memory *Memory) string {
	// [31:25] imm[11:5] [24:20] rs2 [19:15] rs1 [14:12] funct3 [11:7] imm[4:0] [6:0] opcode
	imm := ((instruction >> 25) << 5) | ((instruction >> 7) & 0x1F)
	rs1 := (instruction >> 15) & 0x1F
	rs2 := (instruction >> 20) & 0x1F

	funct3 := (instruction >> 12) & 0x7

	inst, err := FindInstruction(instruction, funct3, 0, 0)
	if err == nil {
		inst.Exec(cpu, memory, rs1, rs2, imm)
		return fmt.Sprintf("%s x%d, x%d, %d\n", inst.Name, rs1, rs2, imm)
	} else {
		return fmt.Sprintf("%s\n", err.Error())
	}
}

func decodeU(opcode Opcode, instruction uint32, cpu *CPUState, memory *Memory) string {
	// [31:12] imm[31:12] [11:7] rd [6:0] opcode
	imm := instruction >> 12
	rd := (instruction >> 7) & 0x1F

	inst, err := FindInstruction(instruction, 0, 0, 0)
	if err == nil {
		inst.Exec(cpu, memory, rd, imm)
		return fmt.Sprintf("%s x%d, %d\n", inst.Name, rd, imm)
	} else {
		return fmt.Sprintf("%s\n", err.Error())
	}
}

func decodeSB(opcode Opcode, instruction uint32, cpu *CPUState, memory *Memory) string {
	// [31] imm[12] [30:25] imm[10:5] [24:21] funct3 [20:20] imm[11] [19:12] imm[4:1] [11:8] rs2 [7:7] rs1 [6:0] opcode
	imm := ((instruction >> 31) << 12) | ((instruction >> 25) & 0x3F) | ((instruction >> 8) & 0xF) | ((instruction>>7)&0x1)<<11
	rs1 := (instruction >> 15) & 0x1F
	rs2 := (instruction >> 20) & 0x1F

	funct3 := (instruction >> 12) & 0x7

	inst, err := FindInstruction(instruction, funct3, 0, 0)
	if err == nil {
		inst.Exec(cpu, memory, rs1, rs2, imm)
		return fmt.Sprintf("%s x%d, x%d, %d\n", inst.Name, rs1, rs2, imm)
	} else {
		return fmt.Sprintf("%s\n", err.Error())
	}
}

func decodeUJ(opcode Opcode, instruction uint32, cpu *CPUState, memory *Memory) string {
	// [31] imm[20] [30:21] imm[10:1] [20:20] imm[11] [19:12] rd [11:8] imm[19:12] [7:7] imm[10] [6:0] opcode
	imm := ((instruction >> 31) << 20) | ((instruction >> 21) & 0x3FF) | ((instruction>>20)&0x1)<<11 | ((instruction>>12)&0xFF)<<12
	rd := (instruction >> 7) & 0x1F

	inst, err := FindInstruction(instruction, 0, 0, 0)
	if err == nil {
		inst.Exec(cpu, memory, rd, imm)
		return fmt.Sprintf("%s x%d, %d\n", inst.Name, rd, imm)
	} else {
		return fmt.Sprintf("%s\n", err.Error())
	}
}

type Encoding struct {
	Type   string
	Decode func(opcode Opcode, instruction uint32, cpu *CPUState, memory *Memory) string
}

var Encodings = map[string]Encoding{
	"I":  {"I", decodeI},   // Integer-Register-Immediate instructions
	"R":  {"R", decodeR},   // Register-register operations
	"S":  {"S", decodeS},   // Store instructions
	"U":  {"U", decodeU},   // Upper immediate instructions
	"SB": {"SB", decodeSB}, // Branch instructions
	"UJ": {"UJ", decodeUJ}, // Jump instructions
}
