package main

import (
	// "encoding/hex"
	"fmt"

	"github.com/google/gopacket/pcap"
)

func main() {
	nics, _ := pcap.FindAllDevs()

	for _, nic := range nics {
		if nic.Name == "en0" {
			handle, err := pcap.OpenLive("en0", 1600, true, pcap.BlockForever) // keep capturing until packet arrives

			if err != nil {
				fmt.Printf("error opening device %d\n\n", err)
			}

			handle.SetBPFFilter("tcp port 80")

			for {
				data, _, _ := handle.ReadPacketData()
				// fmt.Printf("packets: %s", hex.Dump(data))

				ipHeaderSize := int(data[14]&0x0F) * 4
				fmt.Printf("\nIP header size: %d", ipHeaderSize)

				tcpHeaderSize := int(data[46]>>4) * 4
				fmt.Printf("\nTCP header size: %d", tcpHeaderSize)

				srcIP := data[26:30]
				dstIP := data[30:34]

				srcPort := uint16(data[34])<<8 | uint16(data[35])
				dstPort := uint16(data[36])<<8 | uint16(data[37])

				fmt.Printf("\nSource IP: %d.%d.%d.%d:%d", srcIP[0], srcIP[1], srcIP[2], srcIP[3], srcPort)
				fmt.Printf("\nDestination IP: %d.%d.%d.%d:%d", dstIP[0], dstIP[1], dstIP[2], dstIP[3], dstPort)

				packetType := ""
				
				if (dstPort == 80) {
					packetType = "[REQUEST]"
				} else if (srcPort == 80) {
					packetType = "[RESPONSE]"
				}

				flags := data[47]
				info := ""

				if (flags & 0x02) != 0 {
					info += "[SYN] "
				}
				if (flags & 0x10) != 0 {
					info += "[ACK] "
				}
				if (flags & 0x08) != 0 {
					info += "[PSH] "
				}
				if (flags & 0x01) != 0 {
					info += "[FIN] "
				}
				if (flags & 0x04) != 0 {
					info += "[RST] "
				}

				fmt.Printf("\nFlags: %s", info)
				fmt.Printf("\nType: %s\n", packetType)

				payloadIndx := 14 + ipHeaderSize + tcpHeaderSize

				// if a packet size is greater than payload index, then payload exists
				if len(data) > payloadIndx {
					fmt.Printf("\n%s", string(data[payloadIndx:]))
				}

				// fmt.Printf("Payload: %s", string(data[66:]))
				fmt.Println("\n----")
			}
		}
	}
}
