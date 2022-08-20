# trav

## Intro

Trav is a simple tool to expose a local port behind full cone nat to public network.

## Compile

```bash
git clone https://github.com/Alozxy/trav
cd trav/client
go build -o trav
```

## Usage

Currently only supports Linux system, and need iptables to create redirect rules.

```
trav -i 1500 -l 12345 -r 14885 -s stun.mixvoip.com:3478 -o /tmp/external.port
```

Then you can find the port exposed on the public network in /tmp/external.port. When someone connect to this port of nat firewall, the traffic will be redirected to port 14885. The file will be updated with next stun request when the nat mapping changes.


For a detailed and complete list of supported arguments, you can use -h flag to view a full version of help.
```
trav -h                                                                     
Usage of trav:                                                                          
  -6    
      enable ipv6 forwarding. Note that the forwarding port for ipv6 is the external port rather than local port,
      and will be modified when nat mapping change.
  -D
      disable iptables or netsh's port forwarding      
  -i int                                                                                  
      interval between two stun request in second (default 120)                         
  -l uint                                                                                 
      local port (default 12345)  
  -o string
      Write output to <file-path> (default "./external.port")
  -r uint                                                                                 
      redir port (default 14885)                                                        
  -s string                                                                               
      stun server address in [addr:port] format, must support stun over tcp. (default "stun.mixvoip.com:3478") 
```
