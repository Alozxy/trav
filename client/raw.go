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

func send_udp() {

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
	src_port := layers.UDPPort(local_port)
	dst_port := layers.UDPPort(513)

	ip_header := &layers.IPv4{
		SrcIP:    src_ip,
		DstIP:    dst_ip,
		Version:  4,
		TTL:      64,
		Protocol: layers.IPProtocolUDP,
	}

	rand.Seed(time.Now().UnixNano())
	udp_header := &layers.UDP{
		SrcPort: src_port,
		DstPort: dst_port,
	}
	udp_header.SetNetworkLayerForChecksum(ip_header)

	payload := gopacket.Payload("SYN")
	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{
		ComputeChecksums: true,
		FixLengths:       true,
	}
	if err := gopacket.SerializeLayers(buf, opts, udp_header, payload); err != nil {
		log.Fatal(err)
	}

	conn, err := net.ListenPacket("ip4:udp", "0.0.0.0")
	if err != nil {
		log.Fatalln(err)
	}

	if _, err := conn.WriteTo(buf.Bytes(), &net.IPAddr{IP: dst_ip}); err != nil {
		log.Fatalln(err)
	}
	log.Println("send udp packet to " + dst_ip.String() + ":" + strconv.FormatUint(uint64(dst_port), 10))

	conn.Close()
}

func udp_loop() {

	for {
		send_udp()
		time.Sleep(1 * time.Second)
	}
}
