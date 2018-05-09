package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/howeyc/fsnotify"
	"github.com/sparrc/go-ping"
)

type configStruct struct {
	BotToken              string
	IPList                []string
	TelegramNotifiedUsers []int64
	PingAttempts          int
	PingInterval          int
	PingTimeout           int
}

type pingReply struct {
	Source    string
	Reachable bool
	Latency   time.Duration
	Error     error
}

//configuration struct, line 18
var config configStruct

//Configuration file flag
var configFile = flag.String("file", "", "Path to the configuration file. By default ./config.toml")

func main() {
	///
	//	Keyboard interrupts manager
	//		This handles the keyboard interrupts such as CTRL + C
	///
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		for sig := range c {
			log.Printf("captured %v, stopping profiler and exiting..", sig)
			os.Exit(1)
		}
	}()

	///
	//	Config initializer
	//		This handles the keyboard interrupts such as CTRL + C
	///
	//Configuration file flag for shorter flag
	flag.StringVar(configFile, "f", "config.toml", "Path to the configuration file. By default ./config.toml")
	//Parses the flags
	flag.Parse()
	//Reads the config
	config = readConfig()
	//Checks if the config is valid
	checkConfig()

	///
	//	Configuration file watcher
	//		Checks for changes to the configuration file and applies them
	///

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				if ev.IsModify() {
					//Configuration has been changed!
					log.Println("Configuration file has been changed!")
					//Reads the config again
					config = readConfig()
					//Checks again if the config is valid
					checkConfig()
					//if the files gets deleted, kill the process
				} else if ev.IsDelete() {
					log.Panic("The configuration file has been deleted!")
					os.Exit(2)
				}
				break
			case err := <-watcher.Error:
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Watch(*configFile)
	if err != nil {
		log.Fatal(err)
	}
	/////////////////////////

	///
	//	Telegram bot inizialization
	//		Initializes the telegram bot instance
	///
	bot, err := tgbotapi.NewBotAPI(config.BotToken)
	//Checks for errors
	if err != nil {
		log.Panic(err)
	}
	//Silences the debug messages
	bot.Debug = false
	/////////////////////////

	//Sets the IPs to check from the configuration file
	IPs := config.IPList
	//Sets the chats to notify the down to
	tgAdmins := config.TelegramNotifiedUsers

	//Struct referred at line 27, contains the results of the ping
	var pingRes pingReply

	//Cicles forever, until CTRL + C is pressed
	for true {
		//Cicles every IP to check
		for _, ip := range IPs {
			//Pings the IP
			pingRes = pingIP(ip)

			av := "IT'S ALIVE!"
			//If it's not reachable, sends the notification
			if !pingRes.Reachable {
				av = "TANGO DOWN!"
				for _, admin := range tgAdmins {
					bot.Send(tgbotapi.NewMessage(admin, ip+" - TANGO DOWN!"))
				}
			}
			log.Println("Ping result from " + pingRes.Source + ": " + av)
		}
		log.Println("IPs finished, rechecking in " + strconv.Itoa(config.PingInterval) + " seconds")
		time.Sleep(time.Second * time.Duration(config.PingInterval))
	}

}

func pingIP(ip string) pingReply {
	reply := pingReply{Reachable: false, Source: ip}
	iping, err := ping.NewPinger(ip)
	if err != nil {
		//panic(err)
		log.Println("Error!")
		log.Println(err)
		reply.Error = err
		return reply
	}
	iping.SetPrivileged(true)
	iping.Timeout = time.Second * time.Duration(config.PingTimeout)
	iping.Count = config.PingAttempts
	iping.Run()                 // blocks until finished
	stats := iping.Statistics() // get send/receive/rtt stats
	if stats.PacketLoss < 100 {
		reply.Reachable = true
	}
	reply.Latency = stats.AvgRtt
	return reply
}

func readConfig() configStruct {
	_, err := os.Stat(*configFile)
	if err != nil {
		log.Fatal("Config file is missing: ", *configFile)
	}
	log.Println("Reading \"" + *configFile + "\" as toml configuration file")
	var config configStruct
	if _, err := toml.DecodeFile(*configFile, &config); err != nil {
		log.Fatal(err)
	}
	//log.Print(config.Index)
	return config
}

func checkConfig() {
	panic := false
	if len(config.BotToken) == 0 {
		log.Println("Warning, the \"BotToken\" prop is not set in the configuration file!")
		panic = true
	}
	if len(config.IPList) == 0 {
		log.Println("Warning, the \"IPList\" prop is not set in the configuration file!")
		panic = true
	}
	if config.PingInterval == 0 {
		log.Println("Warning, the \"PingInterval\" prop is not set in the configuration file!")
		panic = true
	}
	if config.PingInterval < 0 {
		log.Println("Warning, the \"PingInterval\" prop is invalid!")
		panic = true
	}
	if config.PingTimeout == 0 {
		log.Println("Warning, the \"PingTimeout\" prop is not set in the configuration file!")
		panic = true
	}
	if config.PingTimeout < 0 {
		log.Println("Warning, the \"PingTimeout\" prop is invalid!")
		panic = true
	}
	if config.PingAttempts == 0 {
		log.Println("Warning, the \"PingAttempts\" prop is not set in the configuration file!")
		panic = true
	}
	if config.PingAttempts < 0 {
		log.Println("Warning, the \"PingAttempts\" prop is invalid!")
		panic = true
	}
	if len(config.TelegramNotifiedUsers) == 0 {
		log.Println("Warning, the \"TelegramNotifiedUsers\" prop is not set in the configuration file!")
		panic = true
	}
	if panic {
		os.Exit(1)
	}

}
