package main

import (
	"log"
	"os"
	"os/exec"
	"strconv"
)

func set_rule_v4(local_port uint16, redir_port uint16) {

	if out, err := exec.Command("iptables", "-t", "nat", "-N", "TRAVERSAL").CombinedOutput(); err != nil {
		log.Println("iptables return a non-zero value:", string(out))
		log.Println(err)
	}
	if out, err := exec.Command("iptables", "-t", "nat", "-F", "TRAVERSAL").CombinedOutput(); err != nil {
		log.Println("iptables return a non-zero value:", string(out))
		log.Println(err)
		os.Exit(1)
	}
	if out, err := exec.Command("iptables", "-t", "nat", "-D", "PREROUTING", "-m", "addrtype", "--dst-type", "LOCAL", "-j", "TRAVERSAL").CombinedOutput(); err != nil {
		log.Println("iptables return a non-zero value:", string(out))
		log.Println(err)
	}
	if out, err := exec.Command("iptables", "-t", "nat", "-A", "TRAVERSAL", "-p", "tcp", "--dport", strconv.FormatUint(uint64(local_port), 10), "-j", "REDIRECT", "--to-ports", strconv.FormatUint(uint64(redir_port), 10)).CombinedOutput(); err != nil {
		log.Println("iptables return a non-zero value:", string(out))
		log.Println(err)
		os.Exit(1)
	}
	if out, err := exec.Command("iptables", "-t", "nat", "-A", "PREROUTING", "-m", "addrtype", "--dst-type", "LOCAL", "-j", "TRAVERSAL").CombinedOutput(); err != nil {
		log.Println("iptables return a non-zero value:", string(out))
		log.Println(err)
		os.Exit(1)
	}
}

func set_rule_v6() {

	if out, err := exec.Command("ip6tables", "-t", "nat", "-N", "TRAVERSAL").CombinedOutput(); err != nil {
		log.Println("ip6tables return a non-zero value:", string(out))
		log.Println(err)
	}
	if out, err := exec.Command("ip6tables", "-t", "nat", "-F", "TRAVERSAL").CombinedOutput(); err != nil {
		log.Println("ip6tables return a non-zero value:", string(out))
		log.Println(err)
		os.Exit(1)
	}
	if out, err := exec.Command("ip6tables", "-t", "nat", "-D", "PREROUTING", "-m", "addrtype", "--dst-type", "LOCAL", "-j", "TRAVERSAL").CombinedOutput(); err != nil {
		log.Println("ip6tables return a non-zero value:", string(out))
		log.Println(err)
	}
	if out, err := exec.Command("ip6tables", "-t", "nat", "-A", "PREROUTING", "-m", "addrtype", "--dst-type", "LOCAL", "-j", "TRAVERSAL").CombinedOutput(); err != nil {
		log.Println("ip6tables return a non-zero value:", string(out))
		log.Println(err)
		os.Exit(1)
	}

}

func modify_rule_v6(external_port uint16, redir_port uint16) {
	
	if out, err := exec.Command("ip6tables", "-t", "nat", "-F", "TRAVERSAL").CombinedOutput(); err != nil {
		log.Println("ip6tables return a non-zero value:", string(out))
	}
	if out, err := exec.Command("ip6tables", "-t", "nat", "-A", "TRAVERSAL", "-p", "tcp", "--dport", strconv.FormatUint(uint64(external_port), 10), "-j", "REDIRECT", "--to-ports", strconv.FormatUint(uint64(redir_port), 10)).CombinedOutput(); err != nil {
		log.Println("ip6tables return a non-zero value:", string(out))
	}
	if out, err := exec.Command("ip6tables", "-t", "nat", "-A", "TRAVERSAL", "-p", "udp", "--dport", strconv.FormatUint(uint64(external_port), 10), "-j", "REDIRECT", "--to-ports", strconv.FormatUint(uint64(redir_port), 10)).CombinedOutput(); err != nil {
		log.Println("ip6tables return a non-zero value:", string(out))
	}
}
