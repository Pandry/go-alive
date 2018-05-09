package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api"
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

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		for sig := range c {
			log.Printf("captured %v, stopping profiler and exiting..", sig)
			os.Exit(1)
		}
	}()

	//Init bot
	bot, err := tgbotapi.NewBotAPI("144847736:AAGWBPiFCaejWlFCIW7bnsaNfB3_NUlLqV0")
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = false

	IPs := [...]string{"8.8.8.8", "8.8.4.4", "8.4.8.4", "192.168.1.100", "10.0.0.1", "10.0.1.1", "10.1.1.1"}
	tgAdmins := [...]int64{14092073}

	var pingRes pingReply

	for true {
		for _, ip := range IPs {

			pingRes = pingIP(ip)

			av := "IT'S ALIVE!"
			if !pingRes.Reachable {
				av = "TANGO DOWN!"
				for _, admin := range tgAdmins {
					bot.Send(tgbotapi.NewMessage(admin, ip+" - TANGO DOWN!"))
				}
			}
			fmt.Println("Ping result from " + pingRes.Source + ": " + av)
		}
		fmt.Println("IPs finished, rechecking in 20 seconds")
		time.Sleep(time.Second * 20)
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
	return reply
}
