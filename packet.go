package main

import (
	"fmt"

	"github.com/google/gopacket"
)

type PacketData struct {
	SrcHost string
	DstHost string
	SrcPort string
	DstPort string
	Type    string
	Size    int
}

func NewPacketData(packet gopacket.Packet) (PacketData, error) {
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
