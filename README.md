# docker-nat66
Docker-nat66 is a daemon managing IPv6 port forwarding for Docker containers. It listens to the docker daemon socket for new containers and adjusts Linux netfilter IPv6 tables according to exposed ports.

With these IPv6 port-forwarding rules place, the container is accessible using IPv6 the same way as it is using the IPv4 port-forwarding managed by the Docker daemon.

## Getting docker-nat66
Docker-nat66 code is written in Go and available on GitHub under Apache license. To download and build the code (you need the Golang framework):
````
$ git clone https://github.com/bstevant/docker-nat66.git
$ make
````
An initial release [0.1](https://github.com/bstevant/docker-nat66/releases/tag/0.1) is available if you want a pre-build binary.

**DISCLAIMER: This is an early-released code, intended for debugging and getting feedbacks. This code not ready for production. Use at your own risk !**

## Using docker-nat66

### Docker daemon
Your local Docker should allow IPv6 networking for hosted containers. The Docker daemon should run with the following 2 options:
````
dockerd --ipv6 --fixed-cidr-v6=<Docker IPv6 prefix>
````
The Docker IPv6 prefix is a /64 prefix used to assign IPv6 addresses to the containers. See the [Docker IPv6 Networking Guide](https://docs.docker.com/engine/userguide/networking/default_network/ipv6/) for more informations.

**We strongly recommend to use as Docker IPv6 prefix an Unicast-Local-Address (RFC 4193) prefix, instead of `2001:db8::/16` or any fancy hexspeak IPv6 prefix. To get your own /48 prefix, just go to this [IPv6 ULA Prefix generator](http://unique-local-ipv6.com/).**

### `iptables` tools
Docker-nat66 requires the `ip6tables` tool to interact with the netfilter IPv6 tables. This tool should be installed on your Linux system.

### Start the daemon
Docker-nat66 is supposed to run on the host, besides the Docker daemon.

Docker-nat66 requires 2 arguments to be run:
````
docker-nat66 -dev <egress interface> -prefix <Docker IPv6 prefix>
````
The **egress interface** is the network interface receiving incoming requests to your containers. This interface should be connected to an IPv6-enabled network (e.g. it should have a global IPv6 address).

The **Docker IPv6 prefix** is the prefix used for container addressing, as explained above. 

The daemon requires to be started as root:
````
# ./docker-nat66 -dev eth0 -prefix fdfd:5898:4917:1:/64
````
Once the docker-nat66 daemon is started, it initializes NAT66 in the netfilter tables and listens to docker events to add/remove port-forwarding rules as containers are started or terminated.

### Example
````
# ./docker-nat66 -dev=eth0 -prefix=fdfd:5898:4917:1:/64
2016/10/06 14:23:42 NAT66 set postrouting: /sbin/ip6tables -t nat -A POSTROUTING -o eth0 -s fdfd:5898:4917:1:/64 -j MASQUERADE
````
This first rule is mandatory for translating outgoing packets from your containers.

In another terminal, we start a `redis` container. Notice we first indicate the port to be exposed, but not the port to be mapped with.
````
$ docker run -ti redis -p 6379
````
The docker-nat66 terminal automatically adjusts the NAT66 rules to forward incoming packets to the mapped port to the container port and address:
````
2016/10/06 14:23:48 NAT66 add portfwd: /sbin/ip6tables -t nat -A PREROUTING -i eth0 -p tcp --dport 32770 -j DNAT --to-destination [fdfd:5898:4917:1:0:242:ac11:2]:6379
````
Your `redis` container is now reachable using the IPv6 address configured on `eth0` and port `32770`

When the `redis` container is terminated, the port-forwarding rule is removed from the table:
````
2016/10/06 14:23:54 NAT66 del portfwd: /sbin/ip6tables -t nat -D PREROUTING -i eth0 -p tcp --dport 32770 -j DNAT --to-destination [fdfd:5898:4917:1:0:242:ac11:2]:6379
````
When you explicitely give a port to be mapped for the container, it is taken into account in the port-forwarding rule:
````
$ docker run -ti redis -p 2356:6379
````
Docker-nat66 outputs: (notice the change for the `--dport` argument)
````
2016/10/06 14:24:14 NAT66 add portfwd: /sbin/ip6tables -t nat -A PREROUTING -i eth0 -p tcp --dport 2356 -j DNAT --to-destination [fdfd:5898:4917:1:0:242:ac11:2]:6379
````
When the docker-nat66 daemon is terminated, it should remove any residing NAT66 rules it inserted and the global NAT66 configuration.

## Feedbacks and Bug Reports
As this code is in its very early stage, any feedback is warmly welcome. I am very interested in discussing how this tool can be adapted for Docker production environments where IPv6 is deployed and used.

I am committed to promoting IPv6 usage for more than 10 years and I will really appreciate that the Docker community also embraces IPv6. I hope this tool will help!
