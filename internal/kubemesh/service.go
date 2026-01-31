package kubemesh

import (
	"log"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/tcpassembly"
)

type Service struct{}

func New() *Service {
	return &Service{}
}

func (*Service) Start() *pcap.Handle {
	log.Print("Starting kubemesh service")

	handle, err := pcap.OpenLive("any", 1600, true, pcap.BlockForever)
	if err != nil {
		log.Fatal("Error opening handle on 'any' device/slot for attached network interface card")
	}
	log.Print("Handle opened on 'any'")

	handle.SetBPFFilter("tcp port 80")
	log.Print("BPF filter set successfully")

	return handle
}

func (*Service) Assemble(deviceHandle *pcap.Handle, factory *TCPStreamFactory) *tcpassembly.Assembler {
	pool := tcpassembly.NewStreamPool(factory)
	assembler := tcpassembly.NewAssembler(pool)

	log.Print("TCP assembler successfully initialised")
	return assembler
}

func (*Service) Stream(deviceHandle *pcap.Handle, assembler *tcpassembly.Assembler) {
	packets := gopacket.NewPacketSource(deviceHandle, deviceHandle.LinkType())

	log.Print("Capturing TCP streams...")
	for packet := range packets.Packets() {
		tcpLayer := packet.Layer(layers.LayerTypeTCP)
		if tcpLayer != nil {
			tcp := tcpLayer.(*layers.TCP)
			assembler.Assemble(packet.NetworkLayer().NetworkFlow(), tcp)
		}
	}
}
