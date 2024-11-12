package main

import (
	"encoding/binary"
	"fmt"
	"os"
)

func printHelp() {
	fmt.Println("Utilisation: sae-emulateur [OPTIONS] FICHIER_BIN")
	fmt.Println("")
	fmt.Println("Arguments:")
	fmt.Println("  FICHIER_BIN Un fichier au format binaire contenant les instructions à décoder")
	fmt.Println("")
	fmt.Println("Options:")
	fmt.Println("  -h \t\t\t Affiche ce message d'aide")
	fmt.Println("  -m <uint32> \t Définir la taille de la mémoire en octets (par défaut 512 Ko)")
	fmt.Println("  -d <uint32> \t Définir la valeur par défaut de la mémoire (par défaut 0)")
}

func main() {
	// default memory size & default memory value
	var memorySize uint32 = 512 * 1024
	var registerDefault uint32 = 0
	var cpu CPUState
	var memory []uint32
	var startAddress uint32 = 0
	var endAddress uint32 = 0

	// extract options
	for i, arg := range os.Args {
		if arg == "-h" {
			printHelp()
			os.Exit(0)
		}

		if arg == "-m" {
			if i+1 < len(os.Args) {
				if _, err := fmt.Sscanf(os.Args[i+1], "%d", &memorySize); err != nil {
					printHelp()
					os.Exit(1)
				}
			}
		}

		if arg == "-d" {
			if i+1 < len(os.Args) {
				if _, err := fmt.Sscanf(os.Args[i+1], "%d", &registerDefault); err != nil {
					printHelp()
					os.Exit(1)
				}
			}
		}
	}

	// extract filename from last argument and open file
	filename := os.Args[len(os.Args)-1]
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		os.Exit(1)
	}

	// init memory and cpu state
	initMemory(&memory, memorySize)
	initCPUState(&cpu, startAddress, registerDefault)

	// read binary file and load instructions into memory
	offset := startAddress / 4
	for {
		var instruction uint32
		err = binary.Read(file, binary.LittleEndian, &instruction)
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			fmt.Printf("Error reading file: %v\n", err)
			os.Exit(1)
		}
		if offset < uint32(len(memory)) {
			memory[offset] = instruction
			offset++
		} else {
			fmt.Printf("Binary file too large for memory, actual size: %d, requested size: %d\n", len(memory), offset)
			os.Exit(1)
		}
	}

	// set end address
	endAddress = offset * 4

	// loop through memory and decode instructions
	for {
		// check if pc is out of memory bounds
		if cpu.pc/4 >= uint32(len(memory)) {
			break
		}

		// check if program is finished
		if cpu.pc == endAddress {
			break
		}

		// decode instruction
		instruction := memory[cpu.pc/4]
		opcode, err := GetOpcodeFromInstruction(instruction)

		if err == nil {
			rtnString := opcode.Encoding.Decode(instruction, &cpu, &memory)
			logDebug("DISAS", "%s", rtnString)
		} else {
			logDebug("DISAS", "%08x: %s\n", cpu.pc, err.Error())
		}
		cpu.pc += 4
	}
}
