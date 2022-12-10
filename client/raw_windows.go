package main

import (
	"log"
	"net"
	"time"
)

func send_syn() {

	local_port := get_conf("local_port").(uint16)
	lAddr := &net.TCPAddr{
		Port: int(local_port),
	}
	d := net.Dialer{
		Timeout:   1 * time.Millisecond,
		LocalAddr: lAddr,
	}
	conn, err := d.Dial("tcp", "111.111.111.111:513")
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("connect successfully")
	conn.Close()
}
