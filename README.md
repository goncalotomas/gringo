# Gringo
The ~~dumbest~~ simplest Golang Port Scanner.

## What it does
By default it searches for known TCP services, but it can be set to sweep the whole TCP port range with the `sweep` flag.
### Wi-Fi Performance
Regular mode: under 5 seconds  
Sweep mode: usually around 90 seconds

## Usage
```bash
# Assuming you've cloned the repo and you ran "go build" in the root
./gringo -target={YOUR_TARGET_URL} [-sweep] [-timeout={SECONDS}]
```

## Does it scan UDP ports?
No. UDP ports are particularly difficult to scan when you consider custom designed protocols and/or firewalls. You're welcome to submit a pull request.

## Does it paint a big bullseye on my forehead?
Yes, since you can set it to sweep the whole TCP port range and that _will_ raise alarms on certain systems. Use it only on systems you control or if you have permission.  
This can be extended to work with [gopacket][1] and scan using SYN/FIN packets instead of opening and closing TCP connections. This significantly reduces the chance of detection.

## Why?
~~H4xx0r Sk1llz~~ It seems like a good exercise to learn how to use go routines, and I've had my mind set on building a thing like this for a long time. I can always improve later if I can find the time for it.

## Further reading
[Universal Source of Truth - Port Scanner Page][2]

[1]: https://github.com/google/gopacket
[2]: https://en.wikipedia.org/wiki/Port_scanner
