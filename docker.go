package main

import (
	"github.com/fsouza/go-dockerclient"
)

// bindings store network settings for each known container
var bindings map[string]*docker.NetworkSettings

// clear bindings by deleting corresponding port forwarding
func clear_bindings() {
	for id, net := range bindings {
		ipv6addr := net.GlobalIPv6Address
		for port, binds := range net.Ports {
			for b := range binds {
				del_nat66_portfwd(*dev, port.Proto(), binds[b].HostPort, ipv6addr, port.Port())
			}
		}
		delete(bindings, id)
	}
}

func init_docker() {
	endpoint := "unix:///var/run/docker.sock"
	bindings := make(map[string]*docker.NetworkSettings)

	// Init client to Docker socket
	client, _ := docker.NewClient(endpoint)
	events := make(chan *docker.APIEvents)
	client.AddEventListener(events)

	// Iterate over docker event
	for msg := range events {

		// When a container start ...
		if msg.Action == "start" {
			c, err := client.InspectContainer(msg.ID)
			if err == nil {
				net := c.NetworkSettings
				bindings[msg.ID] = net
				ipv6addr := net.GlobalIPv6Address
				for port, binds := range net.Ports {
					for b := range binds {
						add_nat66_portfwd(*dev, port.Proto(), binds[b].HostPort, ipv6addr, port.Port())
					}
				}
			}
		}

		// When a container die
		if msg.Action == "die" {
			_, err := client.InspectContainer(msg.ID)
			if err == nil {
				if net, ok := bindings[msg.ID]; ok {
					ipv6addr := net.GlobalIPv6Address
					for port, binds := range net.Ports {
						for b := range binds {
							del_nat66_portfwd(*dev, port.Proto(), binds[b].HostPort, ipv6addr, port.Port())
						}
					}
					delete(bindings, msg.ID)
				}
			}
		}
	}

	unset_nat66_postrouting(*dev, *prefix)
}
