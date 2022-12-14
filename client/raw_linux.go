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

func syn_loop() {

	for {
		send_syn()
		time.Sleep(1 * time.Second)
	}
}

func send_syn() {

	local_port := get_conf("local_port").(uint16)
	server_ip := get_conf("server_ip").(string)
	server_port := get_conf("server_port").(uint16)

	src_ip, err := local_ip(server_ip + ":" + strconv.Itoa(int(server_port)))
	if err != nil {
		log.Println("failed to get local ip")
		log.Println(err)
		return
	}
	dst_ip := net.ParseIP("111.111.111.111").To4()
	src_port := layers.TCPPort(local_port)
	dst_port := layers.TCPPort(513)

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
