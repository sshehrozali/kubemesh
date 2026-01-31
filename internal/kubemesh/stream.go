package kubemesh

import (
	"fmt"

	"github.com/google/gopacket"
	"github.com/google/gopacket/tcpassembly"
	"github.com/google/gopacket/tcpassembly/tcpreader"
)

type TCPStream struct {
	reader tcpreader.ReaderStream
	data   []byte
	net,
	transport gopacket.Flow
}

type TCPStreamFactory struct{}

func (ts *TCPStream) Capture() {
	buf := make([]byte, 4096)

	for {
		noOfBytes, err := ts.reader.Read(buf)
		ts.data = append(ts.data, buf[:noOfBytes]...)

		if err != nil {
			fmt.Printf("\nData: \n%s", string(ts.data))
			ts.data = nil
			break
		}
	}
}

func (tsf *TCPStreamFactory) New(net, transport gopacket.Flow) tcpassembly.Stream {
	ts := &TCPStream{
		reader:    tcpreader.NewReaderStream(),
		net:       net,
		transport: transport,
	}

	go ts.Capture()
	return &ts.reader
}
