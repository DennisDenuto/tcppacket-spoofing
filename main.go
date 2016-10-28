package main

import (
	"fmt"
	"github.com/google/gopacket/pcap"
	"time"
	"github.com/google/gopacket"
	"os"
	"github.com/google/gopacket/pcapgo"
	"github.com/google/gopacket/layers"
	"github.com/DennisDenuto/wifi-redirector/listener/http"
	sender_http "github.com/DennisDenuto/wifi-redirector/sender/http"
)

func main() {
	httpListener := http.NewHttpListener("en1")
	httpPackets, err := httpListener.Listen(http.HttpPacketReader{})

	if err != nil {
		fmt.Errorf("error listening", err)
		return
	}

	handle, _ := pcap.OpenLive("en1",
		int32(65535),
		true,
		-1 * time.Second)

	httpInterceptor := sender_http.HttpInterceptor{Sender: sender_http.PacketSender{Handler: handle}}

	for packet := range httpPackets {
		httpInterceptor.Intercept("en1", packet, sender_http.PacketSender{})
	}

}

func main1() {
	fmt.Println("hi")

	ifs, _ := pcap.FindAllDevs()

	for _, value := range ifs {
		fmt.Println(value)
	}

	handle, err := pcap.OpenLive("en1",
		int32(65535),
		true,
		-1 * time.Second)
	//handle.SetBPFFilter("tcp and port 80")
	defer handle.Close()

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	//packet, _ := packetSource.NextPacket()
	writer := WritePacketsToFile()
	for packet := range packetSource.Packets() {
		writer.WritePacket(packet.Metadata().CaptureInfo, packet.Data())

		tcpLayer := packet.Layer(layers.LayerTypeTCP)
		if tcpLayer != nil {
			tcp, _ := tcpLayer.(*layers.TCP)
			fmt.Println(tcp.SrcPort)
			fmt.Println(tcp.DstPort)
		}

		if packet.ApplicationLayer() != nil {
			fmt.Println(string(packet.ApplicationLayer().Payload()))
		}

	}

	fmt.Println(err)
	fmt.Println(handle)
	fmt.Println(pcap.Version())
}

func SendPacket(handler *pcap.Handle) {
	buffer := gopacket.NewSerializeBuffer()
	options := gopacket.SerializeOptions{}

	gopacket.SerializeLayers(buffer, options,
		&layers.Ethernet{},
		&layers.IPv4{},
		&layers.TCP{},
		gopacket.Payload([]byte{65, 66, 67}),
	)
	handler.WritePacketData(buffer.Bytes())
}

func WritePacketsToFile() *pcapgo.Writer {
	dumpFile, _ := os.Create("dump.pcap")
	defer dumpFile.Close()

	packetWrite := pcapgo.NewWriter(dumpFile)
	packetWrite.WriteFileHeader(
		65535,
		layers.LinkTypeEthernet,
	)

	return packetWrite
}