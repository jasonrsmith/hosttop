package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTally(t *testing.T) {
	tally := NewTally()
	assert.NotNil(t, tally)
}

func TestAddPacketBytes(t *testing.T) {
	tally := NewTally()

	packetData := PacketData{
		SrcHost: "srcHost",
		DstHost: "dstHost",
		SrcPort: "srcPost",
		DstPort: "dstPort",
		Type:    "type",
		Size:    1234,
	}

	tally.AddPacket(packetData)
	assert.Equal(t, packetData.Size, tally.BytesBySrcHost[packetData.SrcHost])
	tally.AddPacket(packetData)
	assert.Equal(t, 2*packetData.Size, tally.BytesBySrcHost[packetData.SrcHost])
	assert.Equal(t, 2*packetData.Size, tally.BytesByDstHost[packetData.DstHost])
}

func TestTopHosts(t *testing.T) {
	tally := NewTally()

	packetData1 := PacketData{
		SrcHost: "srcHost1",
		DstHost: "dstHost",
		SrcPort: "srcPost",
		DstPort: "dstPort",
		Type:    "type",
		Size:    1234,
	}
	packetData2 := PacketData{
		SrcHost: "srcHost2",
		DstHost: "dstHost",
		SrcPort: "srcPost",
		DstPort: "dstPort",
		Type:    "type",
		Size:    4321,
	}

	tally.AddPacket(packetData1)
	tally.AddPacket(packetData2)

	assert.Equal(t, []string{packetData2.SrcHost, packetData1.SrcHost}, tally.TopSrcHosts())
	assert.Equal(t, []string{packetData1.DstHost}, tally.TopDstHosts())
}
