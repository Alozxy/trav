package main

import (
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func syn_loop(local_port uint16, server_ip string, server_port uint16) {
	for true {
		send_syn(local_port, server_ip, server_port)
		time.Sleep(1 * time.Second)
	}
}

func main() {

	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		for s := range c {
			switch s {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				clear_rule_v4()
				clear_rule_v6()
				log.Fatalln(s)
			}
		}
	}()

	var stun_server string
	var local_port_64 uint64
	var redir_port_64 uint64
	var interval int
	var enable_ipv6 bool

	flag.StringVar(&stun_server, "s", "stun.mixvoip.com:3478", "stun server address in [addr:port] format, must support stun over tcp.")
	flag.Uint64Var(&local_port_64, "l", 12345, "local port")
	flag.Uint64Var(&redir_port_64, "r", 14885, "redir port")
	flag.IntVar(&interval, "i", 120, "interval between two stun request in second")
	flag.BoolVar(&enable_ipv6, "6", false, "whether to enable ipv6 forwarding. Note that the forwarding port for ipv6 is the external port rather than local port, and will be modified when nat mapping change")
	flag.Parse()

	var local_port uint16 = uint16(local_port_64)
	var redir_port uint16 = uint16(redir_port_64)

	server_ip_list, err := net.LookupIP(strings.Split(stun_server, ":")[0])
	if err != nil {
		log.Fatalln("can't resolve stun server's hostname", err)
	}
	server_ip := server_ip_list[0].String()
	server_port_64, err := strconv.ParseUint(strings.Split(stun_server, ":")[1], 10, 16)
	server_port := uint16(server_port_64)

	clear_rule_v4()
	clear_rule_v6()
	log.Println("creating firewall rules...")
	set_rule_v4(local_port, redir_port)
	if enable_ipv6 {
		log.Println("ipv6 firewall rules is enabled")
		set_rule_v6()
	}

	go syn_loop(local_port, server_ip, server_port)

	var external_port uint16 = 0
	for true {

		request(local_port, &external_port, server_ip, server_port, redir_port, enable_ipv6)
		time.Sleep(time.Duration(interval) * time.Second)
	}

}
