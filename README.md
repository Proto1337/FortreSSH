FortreSSH
===================

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

### TODO / WIP

### Copyright

FortreSSH (C) 2023 Umut Yilmaz.

Golang SSH Tarpit is Copyright (C) 2021 B Tasker. All Rights Reserved. 

Released Under [GNU GPL V3 License](http://www.gnu.org/licenses/gpl-3.0.txt).

