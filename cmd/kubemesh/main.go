package main

import (
	"kubemesh/internal/kubemesh"
	"log"
)

const (
	TRAFFIC_PORT           = "TRAFFIC_PORT"
	NODE_NETWORK_INTERFACE = "NODE_NETWORK_INTERFACE"
)

func main() {
	log.Print("Initialising kubemesh")

	log.Print("Retriving env secrets...")
	tp := kubemesh.GetEnv(TRAFFIC_PORT, "80")
	if !kubemesh.IsValidPort(tp) {
		log.Fatal("Traffic port is invalid")
	}

	nic := kubemesh.GetEnv(NODE_NETWORK_INTERFACE, "any")
	if !kubemesh.IsValidNodeNic(nic) {
		log.Fatal("Invalid node network interface")
	}

	service := kubemesh.New(tp, nic)
	handle := service.Start()

	tsf := &kubemesh.TCPStreamFactory{}
	assembler := service.Assemble(handle, tsf)

	service.Stream(handle, assembler)
}
