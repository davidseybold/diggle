package dns

const (
	maxUDPPacketSize = 512

	udpPacket packetType = "udp"
	tcpPacket packetType = "tcp"
)

type packetType string

type Packet struct {
	header
	Questions   []Question
	Answers     []ResourceRecord
	Authorities []ResourceRecord
	Additional  []ResourceRecord
}
