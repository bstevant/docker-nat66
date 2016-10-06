package main

import (
	"log"
	"os/exec"
	"strings"
)

func set_nat66_postrouting(dev, prefix string) {
	cmd := new(exec.Cmd)
	cmd.Path = "/sbin/ip6tables"
	cmd.Args = []string{
		"/sbin/ip6tables", "-t", "nat", "-A", "POSTROUTING",
		"-o", dev, "-s", prefix,
		"-j", "MASQUERADE",
	}
	log.Println("NAT66 set postrouting: " + strings.Join(cmd.Args, " "))
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func unset_nat66_postrouting(dev, prefix string) {
	cmd := new(exec.Cmd)
	cmd.Path = "/sbin/ip6tables"
	cmd.Args = []string{
		"/sbin/ip6tables", "-t", "nat", "-D", "POSTROUTING",
		"-o", dev, "-s", prefix,
		"-j", "MASQUERADE",
	}
	log.Println("NAT66 unset postrouting: " + strings.Join(cmd.Args, " "))
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func add_nat66_portfwd(dev, proto, dport, target_ipv6, target_port string) {
	cmd := new(exec.Cmd)
	cmd.Path = "/sbin/ip6tables"
	cmd.Args = []string{
		"/sbin/ip6tables", "-t", "nat", "-A", "PREROUTING",
		"-i", dev, "-p", proto, "--dport", dport,
		"-j", "DNAT",
		"--to-destination", "[" + target_ipv6 + "]:" + target_port,
	}
	log.Println("NAT66 add portfwd: " + strings.Join(cmd.Args, " "))
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func del_nat66_portfwd(dev, proto, dport, target_ipv6, target_port string) {
	cmd := new(exec.Cmd)
	cmd.Path = "/sbin/ip6tables"
	cmd.Args = []string{
		"/sbin/ip6tables", "-t", "nat", "-D", "PREROUTING",
		"-i", dev, "-p", proto, "--dport", dport,
		"-j", "DNAT",
		"--to-destination", "[" + target_ipv6 + "]:" + target_port,
	}
	log.Println("NAT66 del portfwd: " + strings.Join(cmd.Args, " "))
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
