package main

import (
	"log"
	"os/exec"
	"strconv"
)

func clear_rule_v4() {

	if out, err := exec.Command("netsh", "interface", "portproxy", "reset").CombinedOutput(); err != nil {
		log.Fatalln("netsh return a non-zero value while clearing ipv4 rules:", string(out))
	}

}

func clear_rule_v6() {

}

func set_rule_v4() {

	local_port := get_conf("local_port").(uint16)
	redir_port := get_conf("redir_port").(uint16)
	server_ip := get_conf("server_ip").(string)
	server_port := get_conf("server_port").(uint16)
	src_ip, err := local_ip(server_ip + ":" + strconv.Itoa(int(server_port)))
	if err != nil {
		log.Println("failed to get local ip, won't set ipv4 redirect rules")
		log.Println(err)
		return
	}

	if out, err := exec.Command("netsh", "interface", "portproxy", "add", "v4tov4", "listenport="+strconv.FormatUint(uint64(local_port), 10), "listenaddress="+src_ip.String(), "connectport="+strconv.FormatUint(uint64(redir_port), 10), "connectaddress=127.0.0.1").CombinedOutput(); err != nil {
		log.Fatalln("netsh return a non-zero value while setting ipv4 rules:", string(out))
	}

}

func set_rule_v6() {

}

func modify_rule_v6(external_port uint16, redir_port uint16) {

}
