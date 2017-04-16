# Gringo
The ~~dumbest~~ simplest Golang Port Scanner.

## What it does
So far, it does a simple sweep of all TCP ports to see which ports are open.

## Does it work if the target system has a firewall?
No. You're welcome to submit a pull request.

## Does it scan UDP ports?
No. UDP ports are particularly difficult to scan when you consider custom designed protocols and/or firewalls. Again, you're welcome to submit a pull request.

## Does it paint a big bullseye on my back?
Yes, it does a simple sweep on TCP ports. Use it only on systems you control or if you have permission.  
This can be extended to work with [gopacket][1] and send SYN/FIN packets instead of completely opening TCP connections. This would greatly reduce the chance of detection.

## Why?
~~H4xx0r Sk1llz~~ It seems like a good exercise to learn how to use go routines, and I've had my mind set on building a thing like this for a long time. I can always improve later if I can find the time for it.

## Further reading
[Universal Source of Truth - Port Scanner Page][2]

[1]: https://github.com/google/gopacket
[2]: https://en.wikipedia.org/wiki/Port_scanner
