package main

import (
    "fmt"
    "time"
    "net"
    "strconv"
    "sync"
    "flag"
)

var KnownTcpPorts = map[int]string{
    // Submit a pull request if you want to add more services!
    7:     "Echo",
  	21:    "FTP",
  	22:    "SSH",
  	23:    "telnet",
  	25:    "SMTP",
    53:     "DNS",
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
    443:   "HTTP w/TLS",
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

const DefaultTimeoutSecs = 5
const MaxTcpPort = 65535
const MaxOpenFileDescriptors = 768

func main() {
    start := time.Now()

    // read and parse command line flags
    targetPtr := flag.String("target", "google.com", "target host")
    timeoutSecsPtr := flag.Int("timeout", DefaultTimeoutSecs, "timeout value in seconds")
    sweepPtr := flag.Bool("sweep", false, "sweep through all tcp ports")
    flag.Parse()

    timeoutDuration := time.Duration(time.Second * time.Duration(*timeoutSecsPtr))

    // numRoutines and ports will depend on the value of *sweepPtr
    var numRoutines int
    var ports []int

    if *sweepPtr {
        fmt.Println("Sweep mode active. Set port range to 1-65535.")
        numRoutines = MaxOpenFileDescriptors
        ports = make([]int,MaxTcpPort)
        for i := 0; i < MaxTcpPort; i++ {
            ports[i] = i+1
        }
    } else {
        // get ports as slice from our map[int]string
        ports = make([]int,len(KnownTcpPorts))
        curr := 0
        for k,_ := range KnownTcpPorts {
          ports[curr] = k
          curr++
        }

        // Since the number of known services is much less than MaxOpenFileDescriptors, we can
        // safely allocate this amount of goroutines, making non-sweep runs fast.
        numRoutines = len(ports)
    }

    fmt.Println("Given current execution mode I'm going to use",numRoutines,"goroutines.")

    // set waiting group to wait for all goroutines before finishing
    var wg sync.WaitGroup
    wg.Add(numRoutines)

    // define the number of ports that each goroutine has to scan
    portsPerRoutine := len(ports)/numRoutines

    // the first goroutine will do some extra work if len(Ports)%numRoutines > 0
    currStart := 0
    currEnd := portsPerRoutine + len(ports)%numRoutines

    for currRoutine := 0; currRoutine < numRoutines; currRoutine++ {
        // start new goroutine, pass it everything it needs, including the correct slice of ports
        go scanMultipleTcpPorts(&wg,*targetPtr,ports[currStart:currEnd],timeoutDuration)
        // increment both indexes
        currStart = currEnd
        currEnd += portsPerRoutine
    }

    // every worker needs to call wg.Done in order for .Wait to work
    wg.Wait()

    // finally, print the elapsed time since start
    elapsed := time.Since(start)
    fmt.Println("Finished! Time elapsed:",elapsed)
}

func scanMultipleTcpPorts(wg *sync.WaitGroup, ipAddr string, portRange []int, timeout time.Duration) {
    defer wg.Done()

    for _,port := range portRange {
      // scan individual TCP port
      isOpen := scanTcpPort(ipAddr,port,timeout)

      // get service name from map
      serviceName := KnownTcpPorts[port]

      if isOpen && serviceName != "" {
          fmt.Println("Target is likely running",KnownTcpPorts[port],"on port",port)
      } else if isOpen && serviceName == ""{
          fmt.Println("Port",port,"is open, but no match for known services")
      }

    }

}

func scanTcpPort(ipAddr string, port int, timeout time.Duration) bool {
    // Concat the host and port so that we can use the net package
    addr := ipAddr + ":" + strconv.Itoa(port)

    // Try to resolve TCP address and bounce out if something goes wrong
    tcpAddr, err := net.ResolveTCPAddr("tcp", addr)

    if err != nil {
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
