package main

import (
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/pion/stun"
)

func request(external_port *uint16) {

	local_port := get_conf("local_port").(uint16)

	lAddr := &net.TCPAddr{
		Port: int(local_port),
	}
	d := &net.Dialer{
		Timeout:   3 * time.Second,
		LocalAddr: lAddr,
	}

	stun_dial(d, external_port)
}

func stun_dial(d *net.Dialer, external_port *uint16) {

	server_ip := get_conf("server_ip").(string)
	server_port := get_conf("server_port").(uint16)
	redir_port := get_conf("redir_port").(uint16)
	enable_ipv6 := get_conf("enable_ipv6").(bool)
	output := get_conf("output").(string)

	log.Println("connecting to stun server...")
	conn, err := d.Dial("tcp4", server_ip+":"+strconv.FormatUint(uint64(server_port), 10))
	if err != nil {
		log.Println(err)
		return
	}

	c, err := stun.NewClient(conn)
	if err != nil {
		log.Println(err)
		return
	}
	if err = c.Do(stun.MustBuild(stun.TransactionID, stun.BindingRequest), func(res stun.Event) {

		if res.Error != nil {
			log.Println(res.Error)
			return
		}
		var xorAddr stun.XORMappedAddress
		if getErr := xorAddr.GetFrom(res.Message); getErr != nil {
			log.Println(getErr)
			if err := c.Close(); err != nil {
				log.Println(err)
				return
			}
			return
		}

		if int(*external_port) == xorAddr.Port {
			log.Println("stun: external port:", xorAddr.Port, "no change")
			return
		} else {
			log.Println("stun: external port:", xorAddr.Port, ", updating file...")
			err = os.WriteFile(output, []byte(strconv.Itoa(xorAddr.Port)), 0644)
			if err != nil {
				log.Fatalln(err)
			}

			if !get_conf("enable_redirect").(bool) {

				if enable_ipv6 {
					log.Println("updating ipv6 firewall rules...")
					modify_rule_v6(uint16(xorAddr.Port), redir_port)
				}
			} else {

				log.Println("port forwarding disabled, skipping...")
			}

			*external_port = uint16(xorAddr.Port)

			return
		}

	}); err != nil {
		log.Println("do:", err)
		return
	}
	if err := c.Close(); err != nil {
		log.Println(err)
		return
	}
}
