package main

func initMemory(memory *[]uint32, size uint32) {
	*memory = make([]uint32, size)
	logDebug("INIT", "Memory initialized with size %d\n", size)
}
