package main

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

func main() {
	tally := NewTally()
	go InitUI(&tally)

	handle, err := pcap.OpenLive("en1", 1600, true, pcap.BlockForever)
	if err != nil {
		panic(err)
	}

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		p, err := NewPacketData(packet)
		if err != nil {
			continue
		}
		tally.Mux.Lock()
		tally.AddPacket(p)
		tally.Mux.Unlock()
	}
}
