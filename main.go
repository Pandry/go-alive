package main

import (
	"fmt"
	"time"

	"github.com/sparrc/go-ping"
)

//"github.com/sparrc/go-ping"

type pingReply struct {
	Source    string
	Reachable bool
	Latency   time.Duration
	Error     error
}

func main() {
	IPs := [...]string{"8.8.8.8", "8.8.4.4", "8.4.8.4", "192.168.1.100", "10.0.0.1", "10.0.1.1", "10.1.1.1"}
	var pingRes pingReply
	for _, ip := range IPs {
		pingRes = pingIP(ip)

		av := "TANGO DOWN!"
		if pingRes.Reachable {
			av = "IT'S ALIVE!"
		}
		fmt.Println("Ping result from " + pingRes.Source + ": " + av)
	}

}

func pingIP(ip string) pingReply {
	reply := pingReply{Reachable: false, Source: ip}
	iping, err := ping.NewPinger(ip)
	if err != nil {
		//panic(err)
		fmt.Println("Error!")
		fmt.Println(err)
		reply.Error = err
		return reply
	}
	iping.SetPrivileged(true)
	iping.Timeout = time.Second * 5
	iping.Count = 1
	iping.Run()                 // blocks until finished
	stats := iping.Statistics() // get send/receive/rtt stats
	if stats.PacketLoss < 100 {
		reply.Reachable = true
	}
	reply.Latency = stats.AvgRtt
	//fmt.Println("Result from IP " + ip + ":")
	//fmt.Println(reply)
	return reply
}
