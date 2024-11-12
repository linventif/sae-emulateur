package main

import "fmt"

type Instruction struct {
	Name string
	Exec func(cpu *CPUState, memory *[]uint32, args ...uint32)
}

var Instructions = map[[4]uint32]Instruction{
	// BRANCH
	{0b1100011, 0b000, 0, 0}: {
		"BEQ",
		func(cpu *CPUState, memory *[]uint32, args ...uint32) {
			rs1, rs2, imm := args[0], args[1], args[2]
			if readRegister(cpu, rs1) == readRegister(cpu, rs2) {
				cpu.pc += imm
			}
		},
	},
	{0b1100011, 0b001, 0, 0}: {
		"BNE",
		func(cpu *CPUState, memory *[]uint32, args ...uint32) {
			rs1, rs2, imm := args[0], args[1], args[2]
			if readRegister(cpu, rs1) != readRegister(cpu, rs2) {
				cpu.pc += imm
			}
		},
	},
	{0b1100011, 0b100, 0, 0}: {
		"BLT",
		func(cpu *CPUState, memory *[]uint32, args ...uint32) {
			rs1, rs2, imm := args[0], args[1], args[2]
			if int32(readRegister(cpu, rs1)) < int32(readRegister(cpu, rs2)) {
				cpu.pc += imm
			}
		},
	},
	{0b1100011, 0b101, 0, 0}: {
		"BGE",
		func(cpu *CPUState, memory *[]uint32, args ...uint32) {
			rs1, rs2, imm := args[0], args[1], args[2]
			if int32(readRegister(cpu, rs1)) >= int32(readRegister(cpu, rs2)) {
				cpu.pc += imm
			}
		},
	},
	{0b1100011, 0b110, 0, 0}: {
		"BLTU",
		func(cpu *CPUState, memory *[]uint32, args ...uint32) {
			rs1, rs2, imm := args[0], args[1], args[2]
			if readRegister(cpu, rs1) < readRegister(cpu, rs2) {
				cpu.pc += imm
			}
		},
	},
	{0b1100011, 0b111, 0, 0}: {
		"BGEU",
		func(cpu *CPUState, memory *[]uint32, args ...uint32) {
			rs1, rs2, imm := args[0], args[1], args[2]
			if readRegister(cpu, rs1) >= readRegister(cpu, rs2) {
				cpu.pc += imm
			}
		},
	},
	// LOAD
	{0b0000011, 0b000, 0, 0}: {
		"LB",
		func(cpu *CPUState, memory *[]uint32, args ...uint32) {
			rd, rs1, imm := args[0], args[1], args[2]
			address := readRegister(cpu, rs1) + imm
			writeRegister(cpu, rd, uint32(int8((*memory)[address])))
		},
	},
	{0b0000011, 0b001, 0, 0}: {
		"LH",
		func(cpu *CPUState, memory *[]uint32, args ...uint32) {
			rd, rs1, imm := args[0], args[1], args[2]
			address := readRegister(cpu, rs1) + imm
			writeRegister(cpu, rd, uint32(int16((*memory)[address])))
		},
	},
	{0b0000011, 0b010, 0, 0}: {
		"LW",
		func(cpu *CPUState, memory *[]uint32, args ...uint32) {
			rd, rs1, imm := args[0], args[1], args[2]
			address := readRegister(cpu, rs1) + imm
			writeRegister(cpu, rd, (*memory)[address])
		},
	},
	{0b0000011, 0b100, 0, 0}: {
		"LBU",
		func(cpu *CPUState, memory *[]uint32, args ...uint32) {
			rd, rs1, imm := args[0], args[1], args[2]
			address := readRegister(cpu, rs1) + imm
			writeRegister(cpu, rd, (*memory)[address]&0xFF)
		},
	},
	{0b0000011, 0b101, 0, 0}: {
		"LHU",
		func(cpu *CPUState, memory *[]uint32, args ...uint32) {
			rd, rs1, imm := args[0], args[1], args[2]
			address := readRegister(cpu, rs1) + imm
			writeRegister(cpu, rd, (*memory)[address]&0xFFFF)
		},
	},
	// MISC-MEM
	{0b0001111, 0, 0, 0}: {
		"FENCE",
		func(cpu *CPUState, memory *[]uint32, args ...uint32) {
			// Do nothing
		},
	},
	// OP-IMM
	{0b0010011, 0b000, 0, 0}: {
		"ADDI",
		func(cpu *CPUState, memory *[]uint32, args ...uint32) {
			rd, rs1, imm := args[0], args[1], args[2]
			writeRegister(cpu, rd, readRegister(cpu, rs1)+imm)
		},
	},
	{0b0010011, 0b010, 0, 0}: {
		"SLTI",
		func(cpu *CPUState, memory *[]uint32, args ...uint32) {
			rd, rs1, imm := args[0], args[1], args[2]
			if int32(readRegister(cpu, rs1)) < int32(imm) {
				writeRegister(cpu, rd, 1)
			} else {
				writeRegister(cpu, rd, 0)
			}
		},
	},
	{0b0010011, 0b011, 0, 0}: {
		"SLTIU",
		func(cpu *CPUState, memory *[]uint32, args ...uint32) {
			rd, rs1, imm := args[0], args[1], args[2]
			if readRegister(cpu, rs1) < imm {
				writeRegister(cpu, rd, 1)
			} else {
				writeRegister(cpu, rd, 0)
			}
		},
	},
	{0b0010011, 0b100, 0, 0}: {
		"XORI",
		func(cpu *CPUState, memory *[]uint32, args ...uint32) {
			rd, rs1, imm := args[0], args[1], args[2]
			writeRegister(cpu, rd, readRegister(cpu, rs1)^imm)
		},
	},
	{0b0010011, 0b110, 0, 0}: {
		"ORI",
		func(cpu *CPUState, memory *[]uint32, args ...uint32) {
			rd, rs1, imm := args[0], args[1], args[2]
			writeRegister(cpu, rd, readRegister(cpu, rs1)|imm)
		},
	},
	{0b0010011, 0b111, 0, 0}: {
		"ANDI",
		func(cpu *CPUState, memory *[]uint32, args ...uint32) {
			rd, rs1, imm := args[0], args[1], args[2]
			writeRegister(cpu, rd, readRegister(cpu, rs1)&imm)
		},
	},
	{0b0010011, 0b001, 0, 0}: {
		"SLLI",
		func(cpu *CPUState, memory *[]uint32, args ...uint32) {
			rd, rs1, imm := args[0], args[1], args[2]
			writeRegister(cpu, rd, readRegister(cpu, rs1)<<imm)
		},
	},
	{0b0010011, 0b101, 0b0000000, 0}: {
		"SRLI",
		func(cpu *CPUState, memory *[]uint32, args ...uint32) {
			rd, rs1, imm := args[0], args[1], args[2]
			writeRegister(cpu, rd, readRegister(cpu, rs1)>>imm)
		},
	},
	{0b0010011, 0b101, 0b0100000, 0}: {
		"SRAI",
		func(cpu *CPUState, memory *[]uint32, args ...uint32) {
			rd, rs1, imm := args[0], args[1], args[2]
			writeRegister(cpu, rd, uint32(int32(readRegister(cpu, rs1))>>imm))
		},
	},
	// JALR
	{0b1100111, 0, 0, 0}: {
		"JALR",
		func(cpu *CPUState, memory *[]uint32, args ...uint32) {
			rd, rs1, imm := args[0], args[1], args[2]
			writeRegister(cpu, rd, cpu.pc+4)
			cpu.pc = (readRegister(cpu, rs1) + imm) & 0xFFFFFFFE
		},
	},
	// SYSTEM
	{0b1110011, 0, 0, 0}: {
		"ECALL",
		func(cpu *CPUState, memory *[]uint32, args ...uint32) {
			switch readRegister(cpu, 10) {
			case 1:
				fmt.Printf("%d", readRegister(cpu, 11))
			case 10:
				cpu.pc = 0
			}
		},
	},
	{0b1110011, 0, 0, 1}: {
		"EBREAK",
		func(cpu *CPUState, memory *[]uint32, args ...uint32) {
			cpu.pc = 0
		},
	},
	// JAL
	{0b1101111, 0, 0, 0}: {
		"JAL",
		func(cpu *CPUState, memory *[]uint32, args ...uint32) {
			rd, imm := args[0], args[1]
			writeRegister(cpu, rd, cpu.pc+4)
			cpu.pc = cpu.pc + imm
		},
	},
	// OP
	{0b0110011, 0b000, 0b0000000, 0}: {
		"ADD",
		func(cpu *CPUState, memory *[]uint32, args ...uint32) {
			rd, rs1, rs2 := args[0], args[1], args[2]
			writeRegister(cpu, rd, readRegister(cpu, rs1)+readRegister(cpu, rs2))
		},
	},
	{0b0110011, 0b000, 0b0100000, 0}: {
		"SUB",
		func(cpu *CPUState, memory *[]uint32, args ...uint32) {
			rd, rs1, rs2 := args[0], args[1], args[2]
			writeRegister(cpu, rd, readRegister(cpu, rs1)-readRegister(cpu, rs2))
		},
	},
	{0b0110011, 0b001, 0, 0}: {
		"SLL",
		func(cpu *CPUState, memory *[]uint32, args ...uint32) {
			rd, rs1, rs2 := args[0], args[1], args[2]
			writeRegister(cpu, rd, readRegister(cpu, rs1)<<readRegister(cpu, rs2))
		},
	},
	{0b0110011, 0b010, 0, 0}: {
		"SLT",
		func(cpu *CPUState, memory *[]uint32, args ...uint32) {
			rd, rs1, rs2 := args[0], args[1], args[2]
			if readRegister(cpu, rs1) < readRegister(cpu, rs2) {
				writeRegister(cpu, rd, 1)
			} else {
				writeRegister(cpu, rd, 0)
			}
		},
	},
	{0b0110011, 0b011, 0, 0}: {
		"SLTU",
		func(cpu *CPUState, memory *[]uint32, args ...uint32) {
			rd, rs1, rs2 := args[0], args[1], args[2]
			if readRegister(cpu, rs1) < readRegister(cpu, rs2) {
				writeRegister(cpu, rd, 1)
			} else {
				writeRegister(cpu, rd, 0)
			}
		},
	},
	{0b0110011, 0b100, 0, 0}: {
		"XOR",
		func(cpu *CPUState, memory *[]uint32, args ...uint32) {
			rd, rs1, rs2 := args[0], args[1], args[2]
			writeRegister(cpu, rd, readRegister(cpu, rs1)^readRegister(cpu, rs2))
		},
	},
	{0b0110011, 0b101, 0b0000000, 0}: {
		"SRL",
		func(cpu *CPUState, memory *[]uint32, args ...uint32) {
			rd, rs1, rs2 := args[0], args[1], args[2]
			writeRegister(cpu, rd, readRegister(cpu, rs1)>>readRegister(cpu, rs2))
		},
	},
	{0b0110011, 0b101, 0b0100000, 0}: {
		"SRA",
		func(cpu *CPUState, memory *[]uint32, args ...uint32) {
			rd, rs1, rs2 := args[0], args[1], args[2]
			writeRegister(cpu, rd, uint32(int32(readRegister(cpu, rs1))>>readRegister(cpu, rs2)))
		},
	},
	{0b0110011, 0b110, 0, 0}: {
		"OR",
		func(cpu *CPUState, memory *[]uint32, args ...uint32) {
			rd, rs1, rs2 := args[0], args[1], args[2]
			writeRegister(cpu, rd, readRegister(cpu, rs1)|readRegister(cpu, rs2))
		},
	},
	{0b0110011, 0b111, 0, 0}: {
		"AND",
		func(cpu *CPUState, memory *[]uint32, args ...uint32) {
			rd, rs1, rs2 := args[0], args[1], args[2]
			writeRegister(cpu, rd, readRegister(cpu, rs1)&readRegister(cpu, rs2))
		},
	},
	// STORE
	{0b0100011, 0b000, 0, 0}: {
		"SB",
		func(cpu *CPUState, memory *[]uint32, args ...uint32) {
			rs1, rs2, imm := args[0], args[1], args[2]
			address := readRegister(cpu, rs1) + imm
			(*memory)[address] = uint32(readRegister(cpu, rs2) & 0xFF)
		},
	},
	{0b0100011, 0b001, 0, 0}: {
		"SH",
		func(cpu *CPUState, memory *[]uint32, args ...uint32) {
			rs1, rs2, imm := args[0], args[1], args[2]
			address := readRegister(cpu, rs1) + imm
			(*memory)[address] = uint32(readRegister(cpu, rs2) & 0xFFFF)
		},
	},
	{0b0100011, 0b010, 0, 0}: {
		"SW",
		func(cpu *CPUState, memory *[]uint32, args ...uint32) {
			rs1, rs2, imm := args[0], args[1], args[2]
			address := readRegister(cpu, rs1) + imm
			(*memory)[address] = readRegister(cpu, rs2)
		},
	},
	// AU-IPC
	{0b0010111, 0, 0, 0}: {
		"AUIPC",
		func(cpu *CPUState, memory *[]uint32, args ...uint32) {
			rd, imm := args[0], args[1]
			writeRegister(cpu, rd, cpu.pc+imm)
		},
	},
	// LUI
	{0b0110111, 0, 0, 0}: {
		"LUI",
		func(cpu *CPUState, memory *[]uint32, args ...uint32) {
			rd, imm := args[0], args[1]
			writeRegister(cpu, rd, imm)
		},
	},
}

func FindInstruction(instruction uint32, funct3 uint32, funct7 uint32, funct12 uint32) (Instruction, error) {
	opcode := instruction & 0x7F
	if instr, ok := Instructions[[4]uint32{opcode, funct3, funct7, funct12}]; ok {
		return instr, nil
	}
	return Instruction{}, fmt.Errorf("instruction not found")
}

//var instructionSet = map[uint32]interface{}{
//	0b1100011: map[uint32]interface{}{
//		0b000: "BEQ",
//		0b001: "BNE",
//		0b100: "BLT",
//		0b101: "BGE",
//		0b110: "BLTU",
//		0b111: "BGEU",
//	},
//	0b0000011: map[uint32]interface{}{
//		0b000: "LB",
//		0b001: "LH",
//		0b010: "LW",
//		0b100: "LBU",
//		0b101: "LHU",
//	},
//	0b0001111: map[uint32]interface{}{
//		0b000: "FENCE",
//	},
//	0b0010011: map[uint32]interface{}{
//		0b000: "ADDI",
//		0b010: "SLTI",
//		0b011: "SLTIU",
//		0b100: "XORI",
//		0b110: "ORI",
//		0b111: "ANDI",
//		0b001: "SLLI",
//		0b101: map[uint32]string{
//			0b0000000: "SRLI",
//			0b0100000: "SRAI",
//		},
//	},
//	0b1100111: "JALR",
//	0b1110011: map[uint32]interface{}{
//		0: map[uint32]interface{}{
//			0: map[uint32]interface{}{
//				0b000000000000: "ECALL",
//				0b000000000001: "EBREAK",
//			},
//		},
//	},
//	0b1101111: "JAL",
//	0b0110011: map[uint32]interface{}{
//		0b000: map[uint32]string{
//			0b0000000: "ADD",
//			0b0100000: "SUB",
//		},
//		0b001: "SLL",
//		0b010: "SLT",
//		0b011: "SLTU",
//		0b100: "XOR",
//		0b101: map[uint32]string{
//			0b0000000: "SRL",
//			0b0100000: "SRA",
//		},
//		0b110: "OR",
//		0b111: "AND",
//	},
//	0b0100011: map[uint32]interface{}{
//		0b000: "SB",
//		0b001: "SH",
//		0b010: "SW",
//	},
//	0b0010111: "AUIPC",
//	0b0110111: "LUI",
//}
//
//func findInstruction(opcode uint32, funct3 *uint32, funct7 *uint32, funct12 *uint32) string {
//	// Step 1: Find opcode
//	if level1, ok := instructionSet[opcode]; ok {
//		if str, ok := level1.(string); ok {
//			return str // Instruction directly at opcode level
//		}
//
//		// Step 2: Check funct3 if provided
//		if funct3 != nil {
//			if level2, ok := level1.(map[uint32]interface{})[*funct3]; ok {
//				if str, ok := level2.(string); ok {
//					return str // Instruction found at funct3 level
//				}
//
//				// Step 3: Check funct7 if provided
//				if funct7 != nil {
//					switch level3 := level2.(type) {
//					case map[uint32]interface{}:
//						if level3Value, ok := level3[*funct7]; ok {
//							if str, ok := level3Value.(string); ok {
//								return str // Instruction found at funct7 level
//							}
//						}
//					case map[uint32]string:
//						if str, ok := level3[*funct7]; ok {
//							return str // Instruction found at funct7 level (final string map)
//						}
//					default:
//						return "Instruction not found"
//					}
//
//					// Step 4: Check funct12 if provided
//					if funct12 != nil {
//						// Safely retrieve level3 from level2 using funct7
//						if level3Map, ok := level2.(map[uint32]interface{}); ok {
//							if level3, ok := level3Map[*funct7]; ok {
//								// Level3 is a map[uint32]interface{}, drill down again
//								if level3Map, ok := level3.(map[uint32]interface{}); ok {
//									if level4, ok := level3Map[*funct12]; ok {
//										if str, ok := level4.(string); ok {
//											return str // Instruction found at funct12 level
//										}
//									}
//								}
//							}
//						}
//					}
//
//				}
//			}
//		}
//	}
//
//	return "Instruction not found"
//}
