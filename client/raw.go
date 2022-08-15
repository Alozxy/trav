package main

import (
	"log"
	"math/rand"
	"net"
	"strconv"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

func local_ip(server_addr string) net.IP {

	conn, err := net.Dial("udp4", server_addr)
	if err != nil {
		log.Println("failed to get local ip")
		log.Fatalln(err)
	}
	defer conn.Close()

	local_addr := conn.LocalAddr().(*net.UDPAddr)
	return local_addr.IP
}

func send_syn(local_port uint16, server_ip string, server_port uint16) {

	src_ip := local_ip(server_ip + ":" + strconv.Itoa(int(server_port)))
	dst_ip := net.ParseIP("1.0.0.1").To4()
	src_port := layers.TCPPort(local_port)
	dst_port := layers.TCPPort(443)

	ip_header := &layers.IPv4{
		SrcIP:    src_ip,
		DstIP:    dst_ip,
		Version:  4,
		TTL:      64,
		Protocol: layers.IPProtocolTCP,
	}

	rand.Seed(time.Now().UnixNano())
	tcp_header := &layers.TCP{
		SrcPort: src_port,
		DstPort: dst_port,
		Seq:     rand.Uint32(),
		SYN:     true,
		Window:  65535,
	}
	tcp_header.SetNetworkLayerForChecksum(ip_header)

	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{
		ComputeChecksums: true,
		FixLengths:       true,
	}
	if err := gopacket.SerializeLayers(buf, opts, tcp_header); err != nil {
		log.Fatal(err)
	}

	conn, err := net.ListenPacket("ip4:tcp", "0.0.0.0")
	if err != nil {
		log.Fatalln(err)
	}

	if _, err := conn.WriteTo(buf.Bytes(), &net.IPAddr{IP: dst_ip}); err != nil {
		log.Fatalln(err)
	}
	log.Println("send syn packet to " + dst_ip.String() + ":" + strconv.FormatUint(uint64(dst_port), 10))

	conn.Close()
}
