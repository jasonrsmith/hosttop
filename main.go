package main

import (
	"fmt"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

type PacketData struct {
	SrcHost string
	DstHost string
	SrcPort string
	DstPort string
	Type    string
	Size    int
}

func main() {
	handle, err := pcap.OpenLive("en1", 1024, true, pcap.BlockForever)
	if err != nil {
		panic(err)
	}

	tally := NewTally()
	go InitUI(&tally)

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		p, err := parsePacket(packet)
		if err != nil {
			continue
		}
		tally.Mux.Lock()
		tally.AddPacket(p)
		tally.Mux.Unlock()
	}
}

func parsePacket(packet gopacket.Packet) (PacketData, error) {
	if packet.NetworkLayer() == nil || packet.TransportLayer() == nil {
		return PacketData{}, fmt.Errorf("Invalid packet")
	}
	size := len(packet.Data())
	srcHost := packet.NetworkLayer().NetworkFlow().Src().String()
	dstHost := packet.NetworkLayer().NetworkFlow().Dst().String()
	srcPort := packet.TransportLayer().TransportFlow().Src().String()
	dstPort := packet.TransportLayer().TransportFlow().Dst().String()
	tlType := packet.TransportLayer().TransportFlow().Src().EndpointType().String()

	return PacketData{
		srcHost,
		dstHost,
		srcPort,
		dstPort,
		tlType,
		size,
	}, nil
}
