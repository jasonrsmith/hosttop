package main

import (
	"net"
	"sort"
	"sync"
)

type Tally struct {
	BytesBySrcHost     map[string]int
	BytesByDstHost     map[string]int
	TotalBytesReceived int
	TotalBytesSent     int
	Mux                sync.Mutex

	lookupTable map[string]string
	hostIP      string
}

func NewTally() Tally {
	return Tally{
		BytesBySrcHost: make(map[string]int),
		BytesByDstHost: make(map[string]int),
		lookupTable:    make(map[string]string),
		hostIP:         getHostIP(),
	}
}

func (t *Tally) AddPacket(p PacketData) {
	if p.SrcHost == t.hostIP {
		t.TotalBytesReceived += p.Size
	}
	if p.DstHost == t.hostIP {
		t.TotalBytesSent += p.Size
	}
	t.BytesBySrcHost[t.reverseLookupHost(p.SrcHost)] += p.Size
	t.BytesByDstHost[t.reverseLookupHost(p.DstHost)] += p.Size
}

func (t *Tally) TopSrcHosts() []string {
	hosts := make([]string, len(t.BytesBySrcHost))
	i := 0
	for host := range t.BytesBySrcHost {
		hosts[i] = host
		i++
	}
	sort.Slice(hosts, func(i, j int) bool {
		return t.BytesBySrcHost[hosts[i]] > t.BytesBySrcHost[hosts[j]]
	})
	return hosts
}

func (t *Tally) TopDstHosts() []string {
	hosts := make([]string, len(t.BytesByDstHost))
	i := 0
	for host := range t.BytesByDstHost {
		hosts[i] = host
		i++
	}
	sort.Slice(hosts, func(i, j int) bool {
		return t.BytesByDstHost[hosts[i]] > t.BytesByDstHost[hosts[j]]
	})
	return hosts
}

func (t *Tally) reverseLookupHost(host string) string {
	resolvedHost, exists := t.lookupTable[host]
	if !exists {
		resolvedHostList, err := net.LookupAddr(host)
		if err != nil {
			resolvedHost = host
		} else {
			t.lookupTable[host] = resolvedHostList[0]
			resolvedHost = resolvedHostList[0]
		}
	}
	return resolvedHost
}

func getHostIP() string {
	ifaces, err := net.Interfaces()
	if err != nil {
		return ""
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return ""
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String()
		}
	}
	return ""
}
