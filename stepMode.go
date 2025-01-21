package main

import (
	"fmt"
	"os"
)

var stepMode = false

// x/4 0x10000
func executeCommand(cpu *CPUState, memory *Memory, commands []string, startAddress uint32, defaultRegisterValue uint32) {
	// if start with 'x/'
	if len(commands) > 0 && commands[0][0:2] == "x/" {
		var count uint32
		var address uint32
		fmt.Sscanf(commands[0], "x/%d", &count)
		fmt.Sscanf(commands[1], "0x%x", &address)
		for i := uint32(0); i < count; i++ {
			fmt.Printf("0x%08x: 0x%08x\n", address+i, readMemory(memory, address+i))
		}
	} else {
		switch commands[0] {
		case "step":
			if cpu.pc/4 >= lenMemory(memory) {
				fmt.Println("PC hors limites mémoire.")
			} else {
				instruction := readMemory(memory, cpu.pc/4)
				opcode, err := GetOpcodeFromInstruction(instruction)

				if err == nil {
					rtnString := opcode.Encoding.Decode(opcode, instruction, cpu, memory)
					logDebug("EXEC", rtnString)
				} else {
					fmt.Println("Erreur lors de l'exécution de l'instruction :", err)
				}
				cpu.pc += 4
			}
		case "continue":
			stepMode = false
			fmt.Println("Sortie du mode pas à pas.")
		case "reset":
			initCPUState(cpu, startAddress, defaultRegisterValue)
			fmt.Println("CPU reset avec PC =", startAddress, "et registre par défaut =", defaultRegisterValue)
		case "exit":
			os.Exit(0)
		default:
			fmt.Println("Commande inconnue.")
		}
	}
}

func handleStepMode(cpu *CPUState, memory *Memory, startAddress uint32, defaultRegisterValue uint32) {
	for stepMode {
		// Affiche l'état des registres
		fmt.Printf("PC: 0x%08x\n", cpu.pc)
		for i := 0; i < 32; i++ {
			fmt.Printf("x%d: 0x%08x\n", i, cpu.x[i])
		}

		// Affiche l'instruction
		if cpu.pc/4 < uint32(lenMemory(memory)) {
			instruction := readMemory(memory, cpu.pc/4)
			opcode, err := GetOpcodeFromInstruction(instruction)
			if err == nil {
				rtnString := opcode.Encoding.Decode(opcode, instruction, cpu, memory)
				logDebug("DISAS", "%s", rtnString)
			} else {
				//logDebug("DISAS", "%08x: %s\n", cpu.pc, err.Error())
			}
		} else {
			fmt.Println("Instruction hors mémoire.")
		}

		// Attend la commande suivante
		fmt.Print("> ")
		var commands = make([]string, 2)
		fmt.Scanln(&commands[0], &commands[1])
		executeCommand(cpu, memory, commands, startAddress, defaultRegisterValue)
	}
}
