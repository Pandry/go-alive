# go-alive
A simple program that checks every n seconds a given array of IP addresses and notifies via Telegram if one or more of the don't reply to the Echo request
(Where n is a user-chosen variable)

# Usage
To install the software:

```
go get github.com/Pandry/go-alive
go build main.go
# Edit your configuration file
./goalive
```
# Options
Unitl now there's only a `-f` option to set a different configuration file
```
-f, -file, --file       sets the location of the configuration file

./goalive -f conf.alternative.toml
```

# Configuration settings
## BotToken
The `BotToken` variable is your telegram API key.
If you don't have one you need to get it from [https://t.me/BotFather](Botfather)  
```
BotToken="123456:YouTGBOTAPIKEY"
```

## IPList
The `IPList` variable is an array that contains the hosts you would like to check  
```
IPList=["1.2.3.4", "5.6.7.8", "1.1.1.1", "1.0.0.1"]
```

## TelegramNotifiedUsers
The `TelegramNotifiedUsers` variable is an array that contains the **IDs** of the chats you want to text to if there's a down IP  
```
TelegramNotifiedUsers=[14123456, 14123457]
```

## PingAttempts
The `PingAttempts` variable is an integer that indicates the number of ping requests the bot sould doevery time it checks an IP  
```
PingAttempts=2
```

## PingInterval
The `PingInterval` variable is an integer that indicates the seconds of pause the program should take before staring again to check the IPs  
```
PingInterval=2
```

## PingTimeout
The `PingTimeout` variable is an integer that indicates the seconds of max timeout allowed for every ping request  
```
PingTimeout=2
```

# TODOs:
- Generate configuration file if not found nor inserted
- Custom interval for every IP or class of IP
- Async 
- Custom DNS server
- Disable auto-reload of the configuration
