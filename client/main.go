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
)

func main() {

	var stun_server string
	var local_port_64 uint64
	var redir_port_64 uint64
	var interval int
	var enable_ipv6 bool
	var enable_redirect bool
	var print_version bool
	var output string
	flag.StringVar(&stun_server, "s", "stun.mixvoip.com:3478", "stun server address in [addr:port] format, must support stun over tcp.")
	flag.Uint64Var(&local_port_64, "l", 12345, "local port")
	flag.Uint64Var(&redir_port_64, "r", 14885, "redir port")
	flag.IntVar(&interval, "i", 300, "interval between two stun request in second")
	flag.BoolVar(&enable_ipv6, "6", false, "enable ipv6 forwarding. Note that the forwarding port for ipv6 is the external port rather than local port, and will be modified when nat mapping change")
	flag.BoolVar(&enable_redirect, "D", false, "disable iptables or netsh's port forwarding")
	flag.BoolVar(&print_version, "v", false, "show current version")
	flag.StringVar(&output, "o", "./external.port", "Write output to <file-path>")
	flag.Parse()
	if print_version {
		println("trav", version)
		os.Exit(0)
	}
	var local_port uint16 = uint16(local_port_64)
	var redir_port uint16 = uint16(redir_port_64)
	server_ip_list, err := net.LookupIP(strings.Split(stun_server, ":")[0])
	if err != nil {
		log.Fatalln("can't resolve stun server's hostname", err)
	}
	server_ip := server_ip_list[0].String()
	server_port_uint64, err := strconv.ParseUint(strings.Split(stun_server, ":")[1], 10, 16)
	if err != nil {
		log.Fatalln("stun server format error", err)
	}
	server_port := uint16(server_port_uint64)

	set_conf("stun_server", stun_server)
	set_conf("local_port", local_port)
	set_conf("redir_port", redir_port)
	set_conf("interval", interval)
	set_conf("enable_ipv6", enable_ipv6)
	set_conf("enable_redirect", enable_redirect)
	set_conf("output", output)
	set_conf("server_ip", server_ip)
	set_conf("server_port", server_port)

	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		for s := range c {
			switch s {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				if !get_conf("enable_redirect").(bool) {
					clear_rule_v4()
					clear_rule_v6()
				}
				log.Fatalln(s)
			}
		}
	}()

	if !get_conf("enable_redirect").(bool) {

		clear_rule_v4()
		clear_rule_v6()
		log.Println("creating firewall rules...")
		set_rule_v4()
		if get_conf("enable_ipv6").(bool) {

			log.Println("ipv6 firewall rules is enabled")
			set_rule_v6()
		}
	} else {

		log.Println("port forwarding disabled, skipping...")
	}

	start()

}

func local_ip(server_addr string) (net.IP, error) {

	conn, err := net.Dial("udp4", server_addr)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	local_addr := conn.LocalAddr().(*net.UDPAddr)
	return local_addr.IP, nil
}
