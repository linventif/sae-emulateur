package main

import (
	"encoding/binary"
	"fmt"
	"os"
)

// Remove the redeclaration of OpcodeMap if it exists elsewhere in the package
var OpcodeMap = map[uint32]struct {
	Type     string
	Encoding string
}{
	0b1100011: {"BRANCH", "SB"},  // Conditional branches
	0b1100111: {"JALR", "I"},     // Jump and link register
	0b0000011: {"LOAD", "I"},     // Load instructions
	0b0001111: {"MISC-MEM", "I"}, // Misc memory instructions
	0b0010011: {"OP-IMM", "I"},   // Integer-Register-Immediate instructions
	0b1110011: {"SYSTEM", "I"},   // System instructions
	0b1101111: {"JAL", "UJ"},     // Jump and link
	0b0110011: {"OP", "R"},       // Register-register operations
	0b0100011: {"STORE", "S"},    // Store instructions
	0b0010111: {"AUIPC", "U"},    // Add upper immediate to PC
	0b0110111: {"LUI", "U"},      // Load upper immediate
}

/*
struct Instruction {
	uint32_t opcode: {
		uint32_t funct3: {
			string ||
			uint32_t funct7: {
				string ||
				uint32_t funct12: {
					string
				}
			}
		}
	}
}
*/

var instructionSet = map[uint32]interface{}{
	0b1100011: map[uint32]interface{}{
		0b000: "BEQ",
		0b001: "BNE",
		0b100: "BLT",
		0b101: "BGE",
		0b110: "BLTU",
		0b111: "BGEU",
	},
	0b0000011: map[uint32]interface{}{
		0b000: "LB",
		0b001: "LH",
		0b010: "LW",
		0b100: "LBU",
		0b101: "LHU",
	},
	0b0001111: map[uint32]interface{}{
		0b000: "FENCE",
	},
	0b0010011: map[uint32]interface{}{
		0b000: "ADDI",
		0b010: "SLTI",
		0b011: "SLTIU",
		0b100: "XORI",
		0b110: "ORI",
		0b111: "ANDI",
		0b001: "SLLI",
		0b101: map[uint32]string{
			0b0000000: "SRLI",
			0b0100000: "SRAI",
		},
	},
	0b1100111: "JALR",
	0b1110011: map[uint32]interface{}{
		0: map[uint32]interface{}{
			0: map[uint32]interface{}{
				0b000000000000: "ECALL",
				0b000000000001: "EBREAK",
			},
		},
	},
	0b1101111: "JAL",
	0b0110011: map[uint32]interface{}{
		0b000: map[uint32]string{
			0b0000000: "ADD",
			0b0100000: "SUB",
		},
		0b001: "SLL",
		0b010: "SLT",
		0b011: "SLTU",
		0b100: "XOR",
		0b101: map[uint32]string{
			0b0000000: "SRL",
			0b0100000: "SRA",
		},
		0b110: "OR",
		0b111: "AND",
	},
	0b0100011: map[uint32]interface{}{
		0b000: "SB",
		0b001: "SH",
		0b010: "SW",
	},
	0b0010111: "AUIPC",
	0b0110111: "LUI",
}

func findInstruction(opcode uint32, funct3 *uint32, funct7 *uint32, funct12 *uint32) string {
	// Step 1: Find opcode
	if level1, ok := instructionSet[opcode]; ok {
		if str, ok := level1.(string); ok {
			return str // Instruction directly at opcode level
		}

		// Step 2: Check funct3 if provided
		if funct3 != nil {
			if level2, ok := level1.(map[uint32]interface{})[*funct3]; ok {
				if str, ok := level2.(string); ok {
					return str // Instruction found at funct3 level
				}

				// Step 3: Check funct7 if provided
				if funct7 != nil {
					switch level3 := level2.(type) {
					case map[uint32]interface{}:
						if level3Value, ok := level3[*funct7]; ok {
							if str, ok := level3Value.(string); ok {
								return str // Instruction found at funct7 level
							}
						}
					case map[uint32]string:
						if str, ok := level3[*funct7]; ok {
							return str // Instruction found at funct7 level (final string map)
						}
					default:
						return "Instruction not found"
					}

					// Step 4: Check funct12 if provided
					if funct12 != nil {
						// Safely retrieve level3 from level2 using funct7
						if level3Map, ok := level2.(map[uint32]interface{}); ok {
							if level3, ok := level3Map[*funct7]; ok {
								// Level3 is a map[uint32]interface{}, drill down again
								if level3Map, ok := level3.(map[uint32]interface{}); ok {
									if level4, ok := level3Map[*funct12]; ok {
										if str, ok := level4.(string); ok {
											return str // Instruction found at funct12 level
										}
									}
								}
							}
						}
					}

				}
			}
		}
	}

	return "Instruction not found"
}

// Remove the redeclaration of printHelp if it exists elsewhere in the package
func printHelp() {
	fmt.Println("Utilisation: disas [OPTIONS] FICHIER_BIN")
	fmt.Println("")
	fmt.Println("Arguments:")
	fmt.Println("FICHIER_BIN Un fichier au format binaire contenant les instructions à décoder")
	fmt.Println("")
	fmt.Println("Options:")
	fmt.Println("-h \t Affiche ce message d'aide")
}

func decodeEncodingI(instruction uint32) string {
	imm := (instruction >> 20) & 0xFFF
	rd := (instruction >> 7) & 0x1F
	rs1 := (instruction >> 15) & 0x1F

	funct3 := (instruction >> 12) & 0x7
	funct7 := (instruction >> 25) & 0x7F
	instructionCode := findInstruction(instruction&0x7F, &funct3, &funct7, &imm)

	if instructionCode != "Unknown" {
		return fmt.Sprintf("%s, x%d, x%d, %d // 0x%02x", instructionCode, rd, rs1, imm, imm)
	}
	return fmt.Sprintf("Unknown")
}

func decodeEncodingR(instruction uint32) string {
	rd := (instruction >> 7) & 0x1F
	rs1 := (instruction >> 15) & 0x1F
	rs2 := (instruction >> 20) & 0x1F
	funct3 := (instruction >> 12) & 0x7
	funct7 := (instruction >> 25) & 0x7F

	instructionCode := findInstruction(instruction&0x7F, &funct3, &funct7, nil)
	if instructionCode != "Unknown" {
		return fmt.Sprintf("%s, x%d, x%d, x%d // 0x%02x", instructionCode, rd, rs1, rs2, rs2)
	}
	return fmt.Sprintf("Unknown")
}

func decodeEncodingS(instruction uint32) string {
	imm := ((instruction >> 25) << 5) | ((instruction >> 7) & 0x1F)
	rs1 := (instruction >> 15) & 0x1F
	rs2 := (instruction >> 20) & 0x1F
	funct3 := (instruction >> 12) & 0x7
	instructionCode := findInstruction(instruction&0x7F, &funct3, nil, &imm)
	if instructionCode != "Unknown" {
		return fmt.Sprintf("%s, x%d, x%d, %d // 0x%02x", instructionCode, rs1, rs2, imm, imm)
	}
	return fmt.Sprintf("Unknown")
}

func decodeEncodingU(instruction uint32) string {
	imm := instruction & 0xFFFFF000
	rd := (instruction >> 7) & 0x1F
	instructionCode := findInstruction(instruction&0x7F, nil, nil, &imm)
	if instructionCode != "Unknown" {
		return fmt.Sprintf("%s, x%d, %d // 0x%02x", instructionCode, rd, imm, imm)
	}
	return fmt.Sprintf("Unknown")
}

func decodeEncodingSB(instruction uint32) string {
	imm := ((instruction >> 31) << 12) | ((instruction >> 7) & 0x1) | ((instruction >> 25) & 0x3F) | ((instruction>>8)&0xF)<<5
	rs1 := (instruction >> 15) & 0x1F
	rs2 := (instruction >> 20) & 0x1F
	funct3 := (instruction >> 12) & 0x7
	instructionCode := findInstruction(instruction&0x7F, &funct3, nil, &imm)
	if instructionCode != "Unknown" {
		return fmt.Sprintf("%s, x%d, x%d, %d // 0x%02x", instructionCode, rs1, rs2, imm, imm)
	}
	return fmt.Sprintf("Unknown")
}

func decodeEncodingUJ(instruction uint32) string {
	imm := ((instruction >> 31) << 20) | ((instruction >> 12) & 0xFF) | ((instruction >> 20) & 0x1) | ((instruction>>21)&0x3FF)<<1
	rd := (instruction >> 7) & 0x1F
	instructionCode := findInstruction(instruction&0x7F, nil, nil, &imm)
	if instructionCode != "Unknown" {
		return fmt.Sprintf("%s, x%d, %d // 0x%02x", instructionCode, rd, imm, imm)
	}
	return fmt.Sprintf("Unknown")
}

func main() {
	// - /home/linventif/sae-emulateur/riscv-samples/assembly/test1.bin
	os.Args = []string{"disas", "/home/linventif/sae-emulateur/riscv-samples/assembly/test2.bin"}

	if len(os.Args) < 2 || os.Args[1] == "-h" {
		printHelp()
		os.Exit(0)
	}

	filename := os.Args[1]
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	fmt.Println("offset,instruction, ...")

	offset := uint32(0)
	var instruction uint32
	for {
		err := binary.Read(file, binary.LittleEndian, &instruction)
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			fmt.Printf("Error reading file: %v\n", err)
			os.Exit(1)
		}

		opcode := instruction & 0x7F
		info, ok := OpcodeMap[opcode]
		if ok {
			switch info.Encoding {
			case "I":
				info.Encoding = decodeEncodingI(instruction)
			case "R":
				info.Encoding = decodeEncodingR(instruction)
			case "S":
				info.Encoding = decodeEncodingS(instruction)
			case "SB":
				info.Encoding = decodeEncodingSB(instruction)
			case "U":
				info.Encoding = decodeEncodingU(instruction)
			case "UJ":
				info.Encoding = decodeEncodingUJ(instruction)
			}

			fmt.Printf("%08x,%s\n", offset, info.Encoding)
		} else {
			fmt.Printf("%08x,Unknown\n", offset)
			continue
		}

		offset += 4
	}
}
