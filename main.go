package main

import (
	"fmt"
	"sync"
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

	/*
		pinger, err := ping.NewPinger("www.google.com")
		if err != nil {
			panic(err)
		}
		pinger.SetPrivileged(true)
		pinger.Timeout = time.Second * 30

		pinger.OnRecv = func(pkt *ping.Packet) {
			fmt.Printf("%d bytes from %s: icmp_seq=%d time=%v\n",
				pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt)
		}
		pinger.OnFinish = func(stats *ping.Statistics) {
			fmt.Printf("\n--- %s ping statistics ---\n", stats.Addr)
			fmt.Printf("%d packets transmitted, %d packets received, %v%% packet loss\n",
				stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss)
			fmt.Printf("round-trip min/avg/max/stddev = %v/%v/%v/%v\n",
				stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt)
		}

		fmt.Printf("PING %s (%s):\n", pinger.Addr(), pinger.IPAddr())
		pinger.Run()
	*/

	messages := make(chan pingReply)
	var wg sync.WaitGroup
	wg.Add(1)
	go pingIP("8.8.8.0", &wg, messages)
	wg.Add(1)
	go pingIP("8.8.8.9", &wg, messages)
	wg.Add(1)
	go pingIP("8.8.8.8", &wg, messages)
	wg.Wait()
	for i := range messages {
		av := "TANGO DOWN!"
		if i.Reachable {
			av = "IT'S ALIVE!"
		}
		fmt.Println("Ping result from " + i.Source + ": " + av)
	}

}

func pingIP(ip string, wg *sync.WaitGroup, ch chan pingReply) pingReply {
	defer wg.Done()
	reply := pingReply{Reachable: false, Source: ip}
	iping, err := ping.NewPinger(ip)
	if err != nil {
		//panic(err)
		fmt.Println("Error!")
		fmt.Println(err)
		reply.Error = err
		wg.Done()
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
	ch <- reply
	return reply
}
