package main

import (
    "fmt"
    "runtime"
    "time"
    "net"
    "strconv"
    "sync"
)

/*
 * A single threaded execution would take too much time, so I chose to use a worker pool.
 * I imagine that one can increase performance by using multiple go routines, but I don't
 * imagine performance increasing when numWorkers >> number of operating system threads.
 * Shout at me via twitter @goncaloptomas if you think I'm wrong.
 */
var numWorkers = runtime.NumCPU()

const DefaultTimeout = time.Duration(time.Second * 2)

var Ports = []int {
    7,
    21,
    22,
    23,
    25,
    66,
    69,
    80,
    88,
    110,
    123,
    137,
    139,
    194,
    118,
    150,
    445,
    554,
    631,
    1433,
    1434,
    3306,
    3535,
    5800,
    8080,
    9160,
    27017,
    28017,
}

var KnownTcpPorts = map[int]string{
    7:     "Echo",
  	21:    "FTP",
  	22:    "SSH",
  	23:    "telnet",
  	25:    "SMTP",
  	66:    "Oracle SQL*NET?",
  	69:    "tftp",
    80:    "HTTP",
  	88:    "kerberos",
  	110:   "POP3",
  	123:   "NTP",
  	137:   "netbios",
  	139:   "netbios",
  	194:   "IRC",
  	118:   "SQL service?",
  	150:   "SQL-net?",
    445:   "Samba",
    554:   "RTSP",
    631:   "CUPS",
    1433:  "Microsoft SQL server",
    1434:  "Microsoft SQL monitor",
    3306:  "MySQL/MariaDB ",
    3535:  "SMTP (alternate)",
    5800:  "VNC remote desktop",
    8080:  "HTTP",
  	9160:  "Cassandra [ http://cassandra.apache.org/ ]",
    27017: "mongodb [ http://www.mongodb.org/ ]",
  	28017: "mongodb web admin [ http://www.mongodb.org/ ]",
}

func main() {
    fmt.Println("Going to scan target for some services I know mm'Kay? Hold on to your hat.")

    // we are going to use several numWorkers (set to 1 to measure single threaded performance)
    var wg sync.WaitGroup
    wg.Add(numWorkers)

    // define the number of ports that each worker has to scan
    portsPerRoutine := len(Ports)/numWorkers

    currStart := 0
    // the first worker will do some extra work if len(Ports)%numWorkers > 0
    // I'm sure there are other ways to do this but I thought it was a good idea.
    currEnd := portsPerRoutine + len(Ports)%numWorkers

    for currentWorker := 0; currentWorker < numWorkers; currentWorker++ {
        // start new go routine, pass it everything it needs, including the correct slice of ports
        go scanMultipleTcpPorts(&wg,"goncalotomas.com",Ports[currStart:currEnd],DefaultTimeout)
        // increment both indexes
        currStart = currEnd
        currEnd += portsPerRoutine
    }

    // every worker needs to call wg.Done in order for .Wait to work
    wg.Wait()
    fmt.Println("Finished!")

}

func scanMultipleTcpPorts(wg *sync.WaitGroup, ipAddr string, portRange []int, timeout time.Duration) {
    defer wg.Done() // can't forget that...

    for _,port := range portRange {

      isOpen := scanTcpPort(ipAddr,port,timeout)

      serviceName := KnownTcpPorts[port]

      if isOpen && serviceName != "" {
          fmt.Println("Target is likely running",KnownTcpPorts[port],"on port",port)
      } else if isOpen && serviceName == ""{
          fmt.Println("Port",port,"is open, no match for known services")
      }

    }

}

func scanTcpPort(ipAddr string, port int, timeout time.Duration) bool {
    // Concat the host and port so that we can use the net package
    addr := ipAddr + ":" + strconv.Itoa(port)
    // Try to resolve TCP address and bounce out if something goes wrong
    tcpAddr, err := net.ResolveTCPAddr("tcp", addr)

    if err != nil {
        fmt.Println(err)
  		  return false
  	}

    // Tries to open a TCP connection to the TCP address, with a specific timeout value.
    // Bounce out if something goes wrong.
    conn, err := net.DialTimeout("tcp", tcpAddr.String(), timeout)

    if err != nil {
  		  return false
  	}

    // if we made it this far then we were able to open a connection
    defer conn.Close()

  	return true
}
