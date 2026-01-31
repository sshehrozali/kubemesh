package main

import (
	"log"
	"kubemesh/internal/kubemesh"
)

func main() {
	log.Print("Initialising kubemesh")

	service := kubemesh.New()
	handle := service.Start()

	tsf := &kubemesh.TCPStreamFactory{}
	assembler := service.Assemble(handle, tsf)

	service.Stream(handle, assembler)
}
