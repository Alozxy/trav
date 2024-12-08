package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"
)

func ifconfig_me_request(external_port *uint16) {

	local_port := get_conf("local_port").(uint16)

	lAddr := &net.TCPAddr{
		Port: int(local_port),
	}
	d := &net.Dialer{
		Timeout:   5 * time.Second,
		LocalAddr: lAddr,
	}

	ipv4DialContext := func(ctx context.Context, network, addr string) (net.Conn, error) {
		return d.DialContext(ctx, "tcp4", addr)
	}

	log.Println("connecting to ifconfig.me")
	client := &http.Client{
		Transport: &http.Transport{
			DialContext:       ipv4DialContext,
			ForceAttemptHTTP2: false,
		},
	}

	resp, err := client.Get("https://ifconfig.me/all.json")
	if err != nil {
		log.Println(err)
		return
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error:", err)
		return
	}

	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Println("Error parsing JSON:", err)
		return
	}

	if value, ok := data["ip_addr"]; ok {
		log.Printf("IP: %v\n", value)
	} else {
		log.Println("Key 'ip_addr' not found")
		return
	}

	if value, ok := data["port"]; ok {
		log.Printf("Port: %v\n", value)
		port, err := strconv.ParseUint(value.(string), 10, 16)
		if err != nil {
			log.Println("Error converting port to uint16:", err)
			return
		}
		*external_port = uint16(port)
	} else {
		log.Println("Key 'port' not found")
		return
	}

}
