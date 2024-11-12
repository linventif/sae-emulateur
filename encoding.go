package main

import "fmt"

func decodeI(instruction uint32, cpu *CPUState, memory *[]uint32) string {
	imm := (instruction >> 20) & 0xFFF
	rd := (instruction >> 7) & 0x1F
	rs1 := (instruction >> 15) & 0x1F

	funct3 := (instruction >> 12) & 0x7
	funct7 := (instruction >> 25) & 0x7F

	inst, err := FindInstruction(instruction, funct3, funct7, 0)
	if err == nil {
		inst.Exec(cpu, memory, rd, rs1, imm)
		return fmt.Sprintf("%s x%d, x%d, %d\n", inst.Name, rd, rs1, imm)
	} else {
		return fmt.Sprintf("Instruction not found: %d\n", instruction)
	}
}

func decodeR(instruction uint32, cpu *CPUState, memory *[]uint32) string {
	rd := (instruction >> 7) & 0x1F
	rs1 := (instruction >> 15) & 0x1F
	rs2 := (instruction >> 20) & 0x1F

	funct3 := (instruction >> 12) & 0x7
	funct7 := (instruction >> 25) & 0x7F

	inst, err := FindInstruction(instruction, funct3, funct7, 0)
	if err == nil {
		inst.Exec(cpu, memory, rd, rs1, rs2)
		return fmt.Sprintf("%s x%d, x%d, x%d\n", inst.Name, rd, rs1, rs2)
	} else {
		return fmt.Sprintf("Instruction not found: %d\n", instruction)
	}
}

func decodeS(instruction uint32, cpu *CPUState, memory *[]uint32) string {
	imm := ((instruction >> 25) << 5) | ((instruction >> 7) & 0x1F)
	rs1 := (instruction >> 15) & 0x1F
	rs2 := (instruction >> 20) & 0x1F

	funct3 := (instruction >> 12) & 0x7

	inst, err := FindInstruction(instruction, funct3, 0, 0)
	if err == nil {
		inst.Exec(cpu, memory, rs1, rs2, imm)
		return fmt.Sprintf("%s x%d, x%d, %d\n", inst.Name, rs1, rs2, imm)
	} else {
		return fmt.Sprintf("Instruction not found: %d\n", instruction)
	}
}

func decodeU(instruction uint32, cpu *CPUState, memory *[]uint32) string {
	imm := instruction & 0xFFFFF000
	rd := (instruction >> 7) & 0x1F

	inst, err := FindInstruction(instruction, 0, 0, 0)
	if err == nil {
		inst.Exec(cpu, memory, rd, imm)
		return fmt.Sprintf("%s x%d, %d\n", inst.Name, rd, imm)
	} else {
		return fmt.Sprintf("Instruction not found: %d\n", instruction)
	}
}

func decodeSB(instruction uint32, cpu *CPUState, memory *[]uint32) string {
	imm := ((instruction >> 31) << 12) | ((instruction >> 7) & 0x1) | ((instruction >> 25) & 0x3F) | ((instruction>>8)&0xF)<<5
	rs1 := (instruction >> 15) & 0x1F
	rs2 := (instruction >> 20) & 0x1F

	funct3 := (instruction >> 12) & 0x7
	inst, err := FindInstruction(instruction, funct3, 0, 0)
	if err == nil {
		inst.Exec(cpu, memory, rs1, rs2, imm)
		return fmt.Sprintf("%s x%d, x%d, %d\n", inst.Name, rs1, rs2, imm)
	} else {
		return fmt.Sprintf("Instruction not found: %d\n", instruction)
	}
}

func decodeUJ(instruction uint32, cpu *CPUState, memory *[]uint32) string {
	imm := ((instruction >> 31) << 20) | ((instruction >> 12) & 0xFF) | ((instruction >> 20) & 0x1) | ((instruction>>21)&0x3FF)<<1
	rd := (instruction >> 7) & 0x1F

	inst, err := FindInstruction(instruction, 0, 0, 0)
	if err == nil {
		inst.Exec(cpu, memory, rd, imm)
		return fmt.Sprintf("%s x%d, %d\n", inst.Name, rd, imm)
	} else {
		return fmt.Sprintf("Instruction not found: %d\n", instruction)
	}
}

type Encoding struct {
	Type   string
	Decode func(opCode uint32, cpu *CPUState, memory *[]uint32) string
}

var Encodings = map[string]Encoding{
	"I":  {"I", decodeI},   // Integer-Register-Immediate instructions
	"R":  {"R", decodeR},   // Register-register operations
	"S":  {"S", decodeS},   // Store instructions
	"U":  {"U", decodeU},   // Upper immediate instructions
	"SB": {"SB", decodeSB}, // Branch instructions
	"UJ": {"UJ", decodeUJ}, // Jump instructions
}
