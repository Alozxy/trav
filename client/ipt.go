package main

import (
	"log"
	"os/exec"
	"strconv"
)

func clear_rule_v4() {

	if out, err := exec.Command("bash", "-c", "iptables-save --counters | grep -v TRAVERSAL | iptables-restore --counters --table nat").CombinedOutput(); err != nil {
		log.Fatalln("iptables-save return a non-zero value while clearing ipv4 rules:", string(out))
	}
}

func clear_rule_v6() {

	if out, err := exec.Command("bash", "-c", "ip6tables-save --counters | grep -v TRAVERSAL | ip6tables-restore --counters --table nat").CombinedOutput(); err != nil {
		log.Fatalln("ip6tables-save return a non-zero value while clearing ipv6 rules:", string(out))
	}
}

func set_rule_v4(local_port uint16, redir_port uint16) {

	if out, err := exec.Command("bash", "-c", `iptables-restore --noflush <<-EOF
		*nat
		:TRAVERSAL -
		-A TRAVERSAL -p tcp -m tcp --dport `+strconv.FormatUint(uint64(local_port), 10)+` -j REDIRECT --to-ports `+strconv.FormatUint(uint64(redir_port), 10)+`
		-A PREROUTING -m addrtype --dst-type LOCAL -j TRAVERSAL	
		COMMIT
		EOF`).CombinedOutput(); err != nil {
		log.Fatalln("iptablesi-restore return a non-zero value while setting ipv4 rules:", string(out))
	}
}

func set_rule_v6() {

	if out, err := exec.Command("bash", "-c", `ip6tables-restore --noflush <<-EOF
		*nat
		:TRAVERSAL -
		-A PREROUTING -m addrtype --dst-type LOCAL -j TRAVERSAL
		COMMIT
		EOF`).CombinedOutput(); err != nil {
		log.Fatalln("ip6tablesi-restore return a non-zero value while setting ipv6 rules:", string(out))
	}
}

func modify_rule_v6(external_port uint16, redir_port uint16) {

	if out, err := exec.Command("bash", "-c", `ip6tables-restore --noflush <<-EOF
		*nat
		-F TRAVERSAL
		-A TRAVERSAL -p tcp -m tcp --dport `+strconv.FormatUint(uint64(external_port), 10)+` -j REDIRECT --to-ports `+strconv.FormatUint(uint64(redir_port), 10)+`
		-A TRAVERSAL -p udp -m udp --dport `+strconv.FormatUint(uint64(external_port), 10)+` -j REDIRECT --to-ports `+strconv.FormatUint(uint64(redir_port), 10)+`
		COMMIT
		EOF`).CombinedOutput(); err != nil {
		log.Fatalln("ip6tablesi-restore return a non-zero value while modifying ipv6 rules:", string(out))
	}
}
