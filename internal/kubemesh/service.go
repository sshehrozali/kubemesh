package kubemesh

import (
	"log"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/tcpassembly"
)

type Service struct{}

func (*Service) New() *Service {
	return &Service{}
}

func (*Service) Start() {
	log.Print("starting kubemesh service")

	tsf := &TCPStreamFactory{}
	pool := tcpassembly.NewStreamPool(tsf)
	assembler := tcpassembly.NewAssembler(pool)

	log.Print("tcp assembler successfully started")

	handle, err := pcap.OpenLive("any", 1600, true, pcap.BlockForever)
	if err != nil {
		log.Fatal("error opening handle on 'any' device/slot for attached network interface card")
	}

	handle.SetBPFFilter("tcp port 80")
	packets := gopacket.NewPacketSource(handle, handle.LinkType())

	for packet := range packets.Packets() {
		tcpLayer := packet.Layer(layers.LayerTypeTCP)
		if tcpLayer != nil {
			tcp := tcpLayer.(*layers.TCP)
			assembler.Assemble(packet.NetworkLayer().NetworkFlow(), tcp)
		}
	}
}
