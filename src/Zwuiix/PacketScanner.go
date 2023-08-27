package Zwuiix

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

type PacketScanner struct{}

func (s PacketScanner) Start() {
	interfaces, err := pcap.FindAllDevs()
	if err != nil {
		log.Fatal(err)
	}

	var handles []*pcap.Handle
	for _, iface := range interfaces {
		handle, err := pcap.OpenLive(iface.Name, 1600, true, pcap.BlockForever)
		if err != nil {
			log.Printf("Error opening adapter %s: %v", iface.Name, err)
			continue
		}
		defer handle.Close()
		handles = append(handles, handle)
	}

	fmt.Println("Capturing outgoing traffic on all interfaces...")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	packetChans := make([]chan gopacket.Packet, len(handles))
	for i := range handles {
		packetChans[i] = make(chan gopacket.Packet)
		go capturePackets(handles[i], packetChans[i])
	}

	for {
		select {
		case packet := <-mergeChannels(packetChans...):
			printPacketInfo(packet)
		case <-sigChan:
			fmt.Println("Exiting...")
			return
		}
	}
}

func capturePackets(handle *pcap.Handle, packetChan chan gopacket.Packet) {
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		packetChan <- packet
	}
}

func mergeChannels(channels ...chan gopacket.Packet) chan gopacket.Packet {
	merged := make(chan gopacket.Packet)
	for _, ch := range channels {
		go func(c chan gopacket.Packet) {
			for {
				merged <- <-c
			}
		}(ch)
	}
	return merged
}

func printPacketInfo(packet gopacket.Packet) {
	ipv4Layer := packet.Layer(layers.LayerTypeIPv4)
	if ipv4Layer != nil {
		ipv4, _ := ipv4Layer.(*layers.IPv4)
		srcIP := ipv4.SrcIP.String()
		dstIP := ipv4.DstIP.String()

		tcpLayer := packet.Layer(layers.LayerTypeTCP)
		if tcpLayer != nil {
			tcp, _ := tcpLayer.(*layers.TCP)
			dstPort := tcp.DstPort.String()

			ip, _ := GetLocalIP()
			if ip == srcIP && strings.HasPrefix(dstIP, "162.159.") && strings.HasSuffix(dstIP, ".232") {
				fmt.Println("[PACKET SCANNER] You -> " + dstIP + ":" + dstPort + " (Non-normal action if no browser is open or discord is launched (small exception, if there are a lot of logs that do not redirect to the same ip address)!)")
			}
		}
	}
}

func GetLocalIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	ip := localAddr.IP.String()

	return ip, nil
}
