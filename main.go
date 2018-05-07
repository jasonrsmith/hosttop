package main

import (
	"flag"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

type Config struct {
	iface string
}

func readConfig() Config {
	config := Config{}
	flag.StringVar(&config.iface, "i", "en1", "interface to listen on")
	flag.Parse()
	return config
}

func main() {
	config := readConfig()
	tally := NewTally()
	go InitUI(&tally)

	handle, err := pcap.OpenLive(config.iface, 1600, true, pcap.BlockForever)
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
