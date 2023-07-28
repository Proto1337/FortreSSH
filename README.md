FortreSSH
===================

FortreSSH is a simple SSH tarpit written in Go.
My main inspiration was [endlessh](https://nullprogram.com/blog/2019/03/22/).
After searching around a bit I found [Golang SSH Tarpit](https://github.com/bentasker/Golang-SSH-Tarpit) and decided to fork it.

There was no specific for forking the project or writing another SSH tarpit. I just wanted to practice coding a bit and kinda missed it after not doing it for a long time.

## Idea

I like the idea of keeping automated bots stuck on your obscure setup and keeping them away from bothering others.  
It was important for me to keep the Go binary minimal to ensure security and to keep the setup as easy to understand and basic as possible.  
Also I wanted to create program with a funny name.

Automated bots scan random IPs on port 22 and try to log in with different credentials.  
These attacks are usually [Brute-force attacks](https://en.wikipedia.org/wiki/Brute-force_attack) or [dictionary attacks](https://en.wikipedia.org/wiki/Dictionary_attack).  
If you keep a tarpit on port 22 and move your SSH somewhere else (e.g. change to an obscure port, setup a VPN and a local subnet to allow connects only from there), you can have a safe setup that annoys these automated bots.

It is important to consider that a tarpit does not add any additional security.  
If you keep your SSH port publicly available you should also use tools like [CrowdSec](https://www.crowdsec.net/) or [Fail2ban](https://www.fail2ban.org/wiki/index.php/Main_Page).
The VPN approach and configuring your local firewall to only accept connections from certain IPs is probably the best way to secure your SSH.

## How should it be used

FortreSSH was created with the thought of deploying it in Docker.  
The Dockerfile creates a minimal image with only FortreSSH running on port 2222.

You can tell Docker to bind your port 22/tcp to the containers 2222 but I suggest to keep it on 2222 inside a Docker network and tell nftables (or your preferred solution) to use a dnat to route the traffic.

## Build/Usage

### Docker

You can run with docker, the image uses the default port `2222` so when running the image, just map it across to `22`

    docker run -d -p 22:2222 bentasker12/go_ssh_tarpit

This will fetch it from [Docker Hub](https://hub.docker.com/r/bentasker12/go_ssh_tarpit)

#### Starting On Boot

The easiest way to have the tarpit image start on boot is tell docker to ensure it's always restarted

    docker run -d -p 22:2222 --restart always bentasker12/go_ssh_tarpit


### Manual

If you'd rather not use docker, you just need to build it with `Go`

    go build ssh_tarpit.go

And then run it

    ./ssh_tarpit.go

Or, of course, you can use `go run`

    go run ssh_tarpit.go

However, by default the script binds to port 2222 - this is so that it could easily run as a non-privileged user within the docker container. If you're running directly, you have 2 options

* Edit the constant to bind to `22` and run as root (*very, very, very* bad idea)
* Run as unprivileged user and use IPTables to NAT `22` to `2222`

The latter can be achieved with

    iptables -t nat -A PREROUTING -p tcp --dport 22 -j REDIRECT --to-port 2222


### Example Raspberry Pi Deployment

The following steps can be used to deploy onto a Raspberry Pi running Raspbian

    curl -fsSL https://get.docker.com -o get-docker.sh
    sudo sh get-docker.sh
    sudo usermod -aG docker pi
    logout
    # log back in
    docker run -d -p 2222:2222 --restart always bentasker12/go_ssh_tarpit:armv7
    sudo iptables -t nat -A PREROUTING -p tcp --dport 22 -j REDIRECT --to-port 2222
    sudo apt-get -y install iptables-persistent
    sudo iptables-save > /etc/iptables/rules.v4

You should be able to see the container running with `docker ls` and can use the name/ID from there to view logs with `docker logs`

----

### Copyright

Golang SSH Tarpit is Copyright (C) 2021 B Tasker. All Rights Reserved. 

Released Under [GNU GPL V3 License](http://www.gnu.org/licenses/gpl-3.0.txt).
