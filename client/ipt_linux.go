package main

import (
	"log"
	"os/exec"
	"strconv"
)

func clean_rule_v4() {

	exec.Command("bash", "-c", `iptables-restore --noflush <<-EOF
		*nat
		-D PREROUTING -m addrtype --dst-type LOCAL -j `+get_conf("chain_name").(string)+`
		-F `+get_conf("chain_name").(string)+`
		-X `+get_conf("chain_name").(string)+`
		COMMIT
		EOF`).Run()
}

func clean_rule_v6() {

	exec.Command("bash", "-c", `ip6tables-restore --noflush <<-EOF
		*nat
		-D PREROUTING -m addrtype --dst-type LOCAL -j `+get_conf("chain_name").(string)+`
		-F `+get_conf("chain_name").(string)+`
		-X `+get_conf("chain_name").(string)+`
		COMMIT
		EOF`).Run()
}

func set_rule_v4() {

	local_port := get_conf("local_port").(uint16)
	redir_port := get_conf("redir_port").(uint16)

	if out, err := exec.Command("bash", "-c", `iptables-restore --noflush <<-EOF
		*nat
		:`+get_conf("chain_name").(string)+` -
		-I `+get_conf("chain_name").(string)+` -p tcp -m tcp --dport `+strconv.FormatUint(uint64(local_port), 10)+` -j REDIRECT --to-ports `+strconv.FormatUint(uint64(redir_port), 10)+`
		-I `+get_conf("chain_name").(string)+` -p udp -m udp --dport `+strconv.FormatUint(uint64(local_port), 10)+` -j REDIRECT --to-ports `+strconv.FormatUint(uint64(redir_port), 10)+`
		-A PREROUTING -m addrtype --dst-type LOCAL -j `+get_conf("chain_name").(string)+`	
		COMMIT
		EOF`).CombinedOutput(); err != nil {
		log.Fatalln("iptablesi-restore return a non-zero value while setting ipv4 rules:", string(out))
	}
}

func set_rule_v6() {

	if out, err := exec.Command("bash", "-c", `ip6tables-restore --noflush <<-EOF
		*nat
		:`+get_conf("chain_name").(string)+` -
		-A PREROUTING -m addrtype --dst-type LOCAL -j `+get_conf("chain_name").(string)+`
		COMMIT
		EOF`).CombinedOutput(); err != nil {
		log.Fatalln("ip6tablesi-restore return a non-zero value while setting ipv6 rules:", string(out))
	}
}

func modify_rule_v6(external_port uint16, redir_port uint16) {

	if out, err := exec.Command("bash", "-c", `ip6tables-restore --noflush <<-EOF
		*nat
		-F `+get_conf("chain_name").(string)+`
		-I `+get_conf("chain_name").(string)+` -p tcp -m tcp --dport `+strconv.FormatUint(uint64(external_port), 10)+` -j REDIRECT --to-ports `+strconv.FormatUint(uint64(redir_port), 10)+`
		-I `+get_conf("chain_name").(string)+` -p udp -m udp --dport `+strconv.FormatUint(uint64(external_port), 10)+` -j REDIRECT --to-ports `+strconv.FormatUint(uint64(redir_port), 10)+`
		COMMIT
		EOF`).CombinedOutput(); err != nil {
		log.Fatalln("ip6tablesi-restore return a non-zero value while modifying ipv6 rules:", string(out))
	}
}
