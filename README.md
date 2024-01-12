<h1 align="center" style="border-bottom: none">
    <a href="https://github.com/Proto1337/FortreSSH" target="_blank"><img alt="FortreSSH" style="width:384px;height:384px;" src="/logo.png"></a><br>FortreSSH
</h1>

---

FortreSSH is a simple SSH tarpit written in Go.
My main inspiration was [endlessh](https://nullprogram.com/blog/2019/03/22/).
After searching around a bit I found [Golang SSH Tarpit](https://github.com/bentasker/Golang-SSH-Tarpit) and decided to fork it.

There was no specific for forking the project or writing another SSH tarpit. I just wanted to practice coding a bit and kinda missed it after not doing it for a long time.

## Idea

I like the idea of keeping automated bots stuck on your obscure setup and keeping them away from bothering others.  
A tarpit sends the bots an endless SSH banner.  
It was important for me to keep the Go binary minimal to ensure security and to keep the setup as easy to understand and basic as possible.  
Also I wanted to create program with a funny name.

Automated bots scan random IPs on port 22 and try to log in with different credentials.  
These attacks are usually [Brute-force attacks](https://en.wikipedia.org/wiki/Brute-force_attack) or [dictionary attacks](https://en.wikipedia.org/wiki/Dictionary_attack).  
If you keep a tarpit on port 22 and move your SSH somewhere else (e.g. change to an obscure port, setup a VPN and a local subnet to allow connects only from there), you can have a safe setup that annoys these automated bots.

It is important to consider that a tarpit does not add any additional security.  
This approach is just made to be annoying, if you want to make research about SSH attacks I recommend honeypots like [Cowrie](https://github.com/cowrie/cowrie).  
If you keep your SSH port publicly available you should also use tools like [CrowdSec](https://www.crowdsec.net/) or [Fail2ban](https://www.fail2ban.org/wiki/index.php/Main_Page).
The VPN approach and configuring your local firewall to only accept connections from certain IPs is probably the best way to secure your SSH.

## How should it be used

FortreSSH was created with the thought of deploying it in Docker.  
I like the ability of containers to isolate processes and [keeping them safe with SELinux](https://opensource.com/article/20/11/selinux-containers).
The Dockerfile creates a minimal image with only FortreSSH running on port 2222.

Of course you could tell Docker to bind your port 22/tcp to the containers 2222 but I suggest to keep it on 2222 inside its own Docker network and tell nftables (or your preferred solution) to use a dnat to route the traffic.  
Especially when using solutions like SELinux you should not mess with policies and should not use labeled ports like 22 for other purposes.

## Build/Usage

### Docker

You can build the Docker image after cloning the repo.

```
git clone https://github.com/Proto1337/FortreSSH
cd FortreSSH
docker build --tag fortressh:latest
```

If you wish to run it directly you could use:

```
docker run -p 2222:2222 fortressh
```

I only recommend this for debugging purposes. Creating the own docker subnet is a better approach.

### Manual deployment (not recommended)

Of course you can compile the binary directly on your machine using Go or gccgo.

```
go build fortressh.go
```

OR

```
gccgo -o fortessh fortressh.go
```

And then run it

```
./fortressh
```

If you use this approach, you should consider that it defaults to port 2222.  
You can change it using "--port" and set it to something else.  
**DO NOT SET IT TO PORT 22 USING ROOT PRIVILEGES!!**
If you really wish to not use Docker please still use the NAT approach!

## nftables configuration

Alot of information about Docker subnets and nftables can be found in [this guide](https://github.com/alexandre-khoury/blog/tree/main/posts/docker-nftables).

### Docker network

It is important to edit `/etc/docker/daemon.json`.

```
{
  "iptables" : false, // Disables iptables management
  "bip": "10.1.1.1/24", // Default bridge network subnet
                        // Must be a valid IP + subnet mask
                        // See https://github.com/docker/for-mac/issues/218#issuecomment-386489719
  "fixed-cidr": "10.1.1.0/25", // Subnet for static IPs in the default bridge network
  "default-address-pools": [
    {
      "base":"10.2.0.0/16", // Available space for custom Docker networks
      "size":24 // Size of a custom network
    }
  ]
}
```
[source](https://github.com/alexandre-khoury/blog/tree/main/posts/docker-nftables#daemonjson)

This way Docker stops messing with your firewall rules and you can start giving IPs out yourself.  
This approach takes more effort but gives you more control.

Now create your Docker network:

```
docker network create tarpit_subnet --subnet 10.2.0.0/24 -o com.docker.network.bridge.name=tarpit0
```

This creates a Docker network named "tarpit_subnet" with possible IPs 10.2.0.1 - 10.2.0.254.  
This custom network can be comfortably used inside of docker-compose.yml files.

For example:

```
version: "3"

services:
    fortressh:
        image: fortressh
        cap_drop:
          - "ALL"
        ports:
            - "10.2.0.1:2222:2222"
        networks:
            tarpit_subnet:
                ipv4_address: 10.2.0.2
        restart: unless-stopped

networks:
    tarpit_subnet:
        external: true
```

Starting the image then with `docker compose up (-d)` sets all the needed network stuff.

### NAT

[Here is more detailed information](https://github.com/alexandre-khoury/blog/tree/main/posts/docker-nftables#nftables-configuration).
Also a lot of inspiration was taken [here](https://wiki.gentoo.org/wiki/Nftables/Examples#Typical_workstation_.28combined_IPv4_and_IPv6.29).

```
#!/sbin/nft -f

define WAN_IFC = {YOUR WAN INTERFACE e.g. eth0, ens3, ...}
define TARPIT_IFC = tarpit0
define SERVER_IP = {YOUR PUBLIC SERVER IP}
define TARPIT_IP = 10.2.0.2

flush ruleset

table inet filter {
	chain prerouting {
		type nat hook prerouting priority dstnat
		policy accept
		ct state invalid drop

		# route to minecraft
		ip daddr $SERVER_IP tcp dport { 22 } dnat to $TARPIT_IP:2222
	}

  	chain input {
		type filter hook input priority 0; policy drop;
		ct state invalid counter drop comment "early drop of invalid packets"
		ct state {established, related} counter accept comment "accept all connections related to connections made by us"
		iif lo accept comment "accept loopback"
		iif != lo ip daddr 127.0.0.1/8 counter drop comment "drop connections to loopback not coming from loopback"
		iif != lo ip6 daddr ::1/128 counter drop comment "drop connections to loopback not coming from loopback"
		ip protocol icmp counter accept comment "accept all ICMP types"
		ip6 nexthdr icmpv6 counter accept comment "accept all ICMP types"

        # Allow outgoing from the Docket network
        iifname $TARPIT_IFC accept

    	counter comment "count dropped packets"
	}

	chain forward {
		type filter hook forward priority 10; policy drop;
		ct state vmap { established : accept, related : accept, invalid : drop }

        # allow Docker Tarpit
        ip daddr $TARPIT_IP tcp dport 2222 accept

        iifname $TARPIT_IFC accept

        counter comment "count dropped packets"
    }

	chain postrouting {
		type nat hook postrouting priority srcnat
		policy accept

		# Modify Docker exiting traffic as coming from server IP
		oifname $WAN_IFC iifname $TARPIT_IFC snat ip to $SERVER_IP
    }

	chain output {
		type filter hook output priority 10; policy accept;
		counter comment "count accepted packets"
	}
}
```

This is a pretty basic nftables config that would allow to route traffic from port 22/tcp to the Docker container 2222/tcp.  
You can load the nftables config with `nft -f /path/to/file`.  
Remember to save the config! Also remember to open a port for a different SSH service or the VPN connection.  
If using SELinux, remember to allow the new port for OpenSSH. Trust me, I locked myself out. :)

-----

### Copyright

FortreSSH (C) 2023 Umut "proto" Yilmaz.

Original code: Golang SSH Tarpit is Copyright (C) 2021 B Tasker. All Rights Reserved. 

Released Under [GNU GPL V3 License](http://www.gnu.org/licenses/gpl-3.0.txt).

