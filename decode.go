package main

import (
	"encoding/binary"
	"fmt"
	"os"
)

// OpcodeMap maps opcodes to their type and encoding
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

func printHelp() {
	fmt.Println("Utilisation: decode_riscv [OPTIONS] FICHIER_BIN")
	fmt.Println("")
	fmt.Println("Arguments:")
	fmt.Println("FICHIER_BIN Un fichier au format binaire contenant les instructions à décoder")
	fmt.Println("")
	fmt.Println("Options:")
	fmt.Println("-h Affiche ce message d'aide")
}

func main() {
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

	fmt.Println("offset,valeur,opcode,encoding")

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
		if !ok {
			continue
		}

		fmt.Printf("%08x,%08x,%s,%s\n", offset, instruction, info.Type, info.Encoding)
		offset += 4
	}
}
