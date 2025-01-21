package main

import (
	"fmt"
)

type Instruction struct {
	Name string
	Exec func(cpu *CPUState, memory *Memory, args ...uint32)
}

var Instructions = map[[4]uint32]Instruction{
	// BRANCH
	// BEQ : Branch if Equal
	{0b1100011, 0b000, 0, 0}: {
		"BEQ",
		func(cpu *CPUState, memory *Memory, args ...uint32) {
			rs1, rs2, imm := args[0], args[1], args[2]
			if readRegister(cpu, rs1) == readRegister(cpu, rs2) {
				cpu.pc += imm
			}
		},
	},
	// BNE : Branch if Not Equal
	{0b1100011, 0b001, 0, 0}: {
		"BNE",
		func(cpu *CPUState, memory *Memory, args ...uint32) {
			rs1, rs2, imm := args[0], args[1], args[2]
			if readRegister(cpu, rs1) != readRegister(cpu, rs2) {
				cpu.pc += imm
			}
		},
	},
	// BLT : Branch if Less Than
	{0b1100011, 0b100, 0, 0}: {
		"BLT",
		func(cpu *CPUState, memory *Memory, args ...uint32) {
			rs1, rs2, imm := args[0], args[1], args[2]
			if int32(readRegister(cpu, rs1)) < int32(readRegister(cpu, rs2)) {
				cpu.pc += imm
			}
		},
	},
	// BGE : Branch if Greater or Equal
	{0b1100011, 0b101, 0, 0}: {
		"BGE",
		func(cpu *CPUState, memory *Memory, args ...uint32) {
			rs1, rs2, imm := args[0], args[1], args[2]
			if int32(readRegister(cpu, rs1)) >= int32(readRegister(cpu, rs2)) {
				cpu.pc += imm
			}
		},
	},
	// BLTU : Branch if Less Than Unsigned
	{0b1100011, 0b110, 0, 0}: {
		"BLTU",
		func(cpu *CPUState, memory *Memory, args ...uint32) {
			rs1, rs2, imm := args[0], args[1], args[2]
			if readRegister(cpu, rs1) < readRegister(cpu, rs2) {
				cpu.pc += imm
			}
		},
	},
	// BGEU : Branch if Greater or Equal Unsigned
	{0b1100011, 0b111, 0, 0}: {
		"BGEU",
		func(cpu *CPUState, memory *Memory, args ...uint32) {
			rs1, rs2, imm := args[0], args[1], args[2]
			if readRegister(cpu, rs1) >= readRegister(cpu, rs2) {
				cpu.pc += imm
			}
		},
	},
	// LOAD
	// LB : Load Byte
	{0b0000011, 0b000, 0, 0}: {
		"LB",
		func(cpu *CPUState, memory *Memory, args ...uint32) {
			rd, rs1, imm := args[0], args[1], args[2]
			address := readRegister(cpu, rs1) + imm
			writeRegister(cpu, rd, readMemory(memory, address))
		},
	},
	// LH : Load Halfword
	{0b0000011, 0b001, 0, 0}: {
		"LH",
		func(cpu *CPUState, memory *Memory, args ...uint32) {
			rd, rs1, imm := args[0], args[1], args[2]
			address := readRegister(cpu, rs1) + imm
			writeRegister(cpu, rd, readMemory(memory, address))
		},
	},
	// LW : Load Word
	{0b0000011, 0b010, 0, 0}: {
		"LW",
		func(cpu *CPUState, memory *Memory, args ...uint32) {
			rd, rs1, imm := args[0], args[1], args[2]
			address := readRegister(cpu, rs1) + imm
			writeRegister(cpu, rd, readMemory(memory, address))
		},
	},
	// LBU : Load Byte Unsigned
	{0b0000011, 0b100, 0, 0}: {
		"LBU",
		func(cpu *CPUState, memory *Memory, args ...uint32) {
			rd, rs1, imm := args[0], args[1], args[2]
			address := readRegister(cpu, rs1) + imm
			writeRegister(cpu, rd, readMemory(memory, address))
		},
	},
	// LHU : Load Halfword Unsigned
	{0b0000011, 0b101, 0, 0}: {
		"LHU",
		func(cpu *CPUState, memory *Memory, args ...uint32) {
			rd, rs1, imm := args[0], args[1], args[2]
			address := readRegister(cpu, rs1) + imm
			writeRegister(cpu, rd, readMemory(memory, address))
		},
	},
	// MISC-MEM
	// FENCE : Fence Instruction
	{0b0001111, 0, 0, 0}: {
		"FENCE",
		func(cpu *CPUState, memory *Memory, args ...uint32) {
			// Do nothing
		},
	},
	// OP-IMM
	// ADDI : Add Immediate
	{0b0010011, 0b000, 0, 0}: {
		"ADDI",
		func(cpu *CPUState, memory *Memory, args ...uint32) {
			rd, rs1, imm := args[0], args[1], args[2]
			writeRegister(cpu, rd, readRegister(cpu, rs1)+imm)
		},
	},
	// SLTI : Set Less Than Immediate
	{0b0010011, 0b010, 0, 0}: {
		"SLTI",
		func(cpu *CPUState, memory *Memory, args ...uint32) {
			rd, rs1, imm := args[0], args[1], args[2]
			if int32(readRegister(cpu, rs1)) < int32(imm) {
				writeRegister(cpu, rd, 1)
			} else {
				writeRegister(cpu, rd, 0)
			}
		},
	},
	// SLTIU : Set Less Than Immediate Unsigned
	{0b0010011, 0b011, 0, 0}: {
		"SLTIU",
		func(cpu *CPUState, memory *Memory, args ...uint32) {
			rd, rs1, imm := args[0], args[1], args[2]
			signExtendedImm := uint32(int32(imm<<20) >> 20)
			if readRegister(cpu, rs1) < signExtendedImm {
				writeRegister(cpu, rd, 1)
			} else {
				writeRegister(cpu, rd, 0)
			}
		},
	},
	// XORI : XOR Immediate
	{0b0010011, 0b100, 0, 0}: {
		"XORI",
		func(cpu *CPUState, memory *Memory, args ...uint32) {
			rd, rs1, imm := args[0], args[1], args[2]
			writeRegister(cpu, rd, readRegister(cpu, rs1)^imm)
		},
	},
	// ORI : OR Immediate
	{0b0010011, 0b110, 0, 0}: {
		"ORI",
		func(cpu *CPUState, memory *Memory, args ...uint32) {
			rd, rs1, imm := args[0], args[1], args[2]
			writeRegister(cpu, rd, readRegister(cpu, rs1)|imm)
		},
	},
	// ANDI : AND Immediate
	{0b0010011, 0b111, 0, 0}: {
		"ANDI",
		func(cpu *CPUState, memory *Memory, args ...uint32) {
			rd, rs1, imm := args[0], args[1], args[2]
			writeRegister(cpu, rd, readRegister(cpu, rs1)&imm)
		},
	},
	// OP-IMM-32 : 32-bit Immediate
	{0b0010011, 0b001, 0, 0}: {
		"SLLI",
		func(cpu *CPUState, memory *Memory, args ...uint32) {
			rd, rs1, imm := args[0], args[1], args[2]
			writeRegister(cpu, rd, readRegister(cpu, rs1)<<imm)
		},
	},
	// SRLI : Shift Right Logical Immediate
	{0b0010011, 0b101, 0b0000000, 0}: {
		"SRLI",
		func(cpu *CPUState, memory *Memory, args ...uint32) {
			rd, rs1, imm := args[0], args[1], args[2]
			writeRegister(cpu, rd, readRegister(cpu, rs1)>>imm)
		},
	},
	// SRAI : Shift Right Arithmetic Immediate
	{0b0010011, 0b101, 0b0100000, 0}: {
		"SRAI",
		func(cpu *CPUState, memory *Memory, args ...uint32) {
			rd, rs1, imm := args[0], args[1], args[2]
			writeRegister(cpu, rd, uint32(int32(readRegister(cpu, rs1))>>imm))
		},
	},
	// JALR
	// JALR : Jump and Link Register
	{0b1100111, 0, 0, 0}: {
		"JALR",
		func(cpu *CPUState, memory *Memory, args ...uint32) {
			rd, rs1, imm := args[0], args[1], args[2]
			writeRegister(cpu, rd, cpu.pc+4)
			var targetAddress = readRegister(cpu, rs1) + imm
			cpu.pc = targetAddress
		},
	},
	// SYSTEM
	// ECALL : Environment Call
	{0b1110011, 0, 0, 0}: {
		"ECALL",
		func(cpu *CPUState, memory *Memory, args ...uint32) {
			// Do nothing
		},
	},
	// EBREAK : Environment Break
	{0b1110011, 0, 0, 1}: {
		"EBREAK",
		func(cpu *CPUState, memory *Memory, args ...uint32) {
			stepMode = true
			fmt.Println("Step by step mode enabled")
		},
	},
	// JAL
	// JAL : Jump and Link
	{0b1101111, 0, 0, 0}: {
		"JAL",
		func(cpu *CPUState, memory *Memory, args ...uint32) {
			rd, imm := args[0], args[1]
			writeRegister(cpu, rd, cpu.pc+4)
			cpu.pc = cpu.pc + imm
		},
	},
	// OP
	// ADD : Add
	{0b0110011, 0b000, 0b0000000, 0}: {
		"ADD",
		func(cpu *CPUState, memory *Memory, args ...uint32) {
			rd, rs1, rs2 := args[0], args[1], args[2]
			writeRegister(cpu, rd, readRegister(cpu, rs1)+readRegister(cpu, rs2))
		},
	},
	// SUB : Subtract
	{0b0110011, 0b000, 0b0100000, 0}: {
		"SUB",
		func(cpu *CPUState, memory *Memory, args ...uint32) {
			rd, rs1, rs2 := args[0], args[1], args[2]
			writeRegister(cpu, rd, readRegister(cpu, rs1)-readRegister(cpu, rs2))
		},
	},
	// SLL : Shift Left Logical
	{0b0110011, 0b001, 0, 0}: {
		"SLL",
		func(cpu *CPUState, memory *Memory, args ...uint32) {
			rd, rs1, rs2 := args[0], args[1], args[2]
			writeRegister(cpu, rd, readRegister(cpu, rs1)<<readRegister(cpu, rs2))
		},
	},
	// SLT : Set Less Than
	{0b0110011, 0b010, 0, 0}: {
		"SLT",
		func(cpu *CPUState, memory *Memory, args ...uint32) {
			rd, rs1, rs2 := args[0], args[1], args[2]
			if readRegister(cpu, rs1) < readRegister(cpu, rs2) {
				writeRegister(cpu, rd, 1)
			} else {
				writeRegister(cpu, rd, 0)
			}
		},
	},
	// SLTU : Set Less Than Unsigned
	{0b0110011, 0b011, 0, 0}: {
		"SLTU",
		func(cpu *CPUState, memory *Memory, args ...uint32) {
			rd, rs1, rs2 := args[0], args[1], args[2]
			if readRegister(cpu, rs1) < readRegister(cpu, rs2) {
				writeRegister(cpu, rd, 1)
			} else {
				writeRegister(cpu, rd, 0)
			}
		},
	},
	// XOR : XOR
	{0b0110011, 0b100, 0, 0}: {
		"XOR",
		func(cpu *CPUState, memory *Memory, args ...uint32) {
			rd, rs1, rs2 := args[0], args[1], args[2]
			writeRegister(cpu, rd, readRegister(cpu, rs1)^readRegister(cpu, rs2))
		},
	},
	// SRL : Shift Right Logical
	{0b0110011, 0b101, 0b0000000, 0}: {
		"SRL",
		func(cpu *CPUState, memory *Memory, args ...uint32) {
			rd, rs1, rs2 := args[0], args[1], args[2]
			writeRegister(cpu, rd, readRegister(cpu, rs1)>>readRegister(cpu, rs2))
		},
	},
	// SRA : Shift Right Arithmetic
	{0b0110011, 0b101, 0b0100000, 0}: {
		"SRA",
		func(cpu *CPUState, memory *Memory, args ...uint32) {
			rd, rs1, rs2 := args[0], args[1], args[2]
			writeRegister(cpu, rd, uint32(int32(readRegister(cpu, rs1))>>readRegister(cpu, rs2)))
		},
	},
	// OR : OR
	{0b0110011, 0b110, 0, 0}: {
		"OR",
		func(cpu *CPUState, memory *Memory, args ...uint32) {
			rd, rs1, rs2 := args[0], args[1], args[2]
			writeRegister(cpu, rd, readRegister(cpu, rs1)|readRegister(cpu, rs2))
		},
	},
	// AND : AND
	{0b0110011, 0b111, 0, 0}: {
		"AND",
		func(cpu *CPUState, memory *Memory, args ...uint32) {
			rd, rs1, rs2 := args[0], args[1], args[2]
			writeRegister(cpu, rd, readRegister(cpu, rs1)&readRegister(cpu, rs2))
		},
	},
	// STORE
	// SB : Store Byte
	{0b0100011, 0b000, 0, 0}: {
		"SB",
		func(cpu *CPUState, memory *Memory, args ...uint32) {
			rs1, rs2, imm := args[0], args[1], args[2]
			address := readRegister(cpu, rs1) + imm
			value := readRegister(cpu, rs2) & 0xFF // Mask to keep only the lower 8 bits
			writeByte(memory, address, value)      // Write only a byte to memory
		},
	},
	// SH : Store Halfword
	{0b0100011, 0b001, 0, 0}: {
		"SH",
		func(cpu *CPUState, memory *Memory, args ...uint32) {
			rs1, rs2, imm := args[0], args[1], args[2]
			address := readRegister(cpu, rs1) + imm
			value := readRegister(cpu, rs2) & 0xFFFF // Mask lower 16 bits
			//writeMemory(memory, address, value)      // Write the entire 32-bit value
			writeHalfword(memory, address, value) // Write only a halfword to memory
			fmt.Printf("SH: address: %d, value: %d\n", address, value)
		},
	},
	// SW : Store Word
	{0b0100011, 0b010, 0, 0}: {
		"SW",
		func(cpu *CPUState, memory *Memory, args ...uint32) {
			rs1, rs2, imm := args[0], args[1], args[2]
			address := readRegister(cpu, rs1) + imm
			value := readRegister(cpu, rs2) // Use the entire 32-bit value
			writeWord(memory, address, value)
		},
	},
	// AU-IPC
	// AUIPC : Add Upper Immediate to PC
	{0b0010111, 0, 0, 0}: {
		"AUIPC",
		func(cpu *CPUState, memory *Memory, args ...uint32) {
			rd, imm := args[0], args[1]
			imm = imm << 12
			writeRegister(cpu, rd, cpu.pc+imm)
		},
	},
	// LUI
	// LUI : Load Upper Immediate
	{0b0110111, 0, 0, 0}: {
		"LUI",
		func(cpu *CPUState, memory *Memory, args ...uint32) {
			rd, imm := args[0], args[1]
			imm = imm << 12
			writeRegister(cpu, rd, imm)
		},
	},
}

func FindInstruction(instruction uint32, funct3 uint32, funct7 uint32, funct12 uint32) (Instruction, error) {
	opcode := instruction & 0x7F
	if instr, ok := Instructions[[4]uint32{opcode, funct3, funct7, funct12}]; ok {
		return instr, nil
	}
	return Instruction{}, fmt.Errorf("instruction {opcode: %b, funct3: %b, funct7: %b, funct12: %b} not found", opcode, funct3, funct7, funct12)
}
