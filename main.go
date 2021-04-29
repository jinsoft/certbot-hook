package main

import (
	"certbot-hook/providers"
	"flag"
	"fmt"
	"log"
	"os"
)

var (
	validationKey string
	domain        string
)

var (
	action = flag.String("action", "", "add or del")
	key    = flag.String("accessKey_id", "", "AccessKey ID")
	secret = flag.String("accessKey_secret", "", "AccessKey Secret")
)

const DNSTYPE = "TXT"
const RR = "_acme-challenge"

func SetupLogger() {
	logFileLocation, _ := os.OpenFile("./log.txt", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0744)
	log.SetOutput(logFileLocation)
}

func main() {
	SetupLogger()
	flag.Parse()
	parseCertbotEnvValue()
	var p *providers.Aliyun
	p = providers.NewAliyun(domain, *key, *secret)
	run(p)
}

func run(p *providers.Aliyun) {
	log.Println("Current action is", *action)
	switch *action {
	case "add":
		log.Println("will add record TXT _acme-challenge ", validationKey)
		recordId, err := p.ResolveDomainName(DNSTYPE, RR, validationKey)
		if err != nil {
			log.Fatalf("Auto resolve record TXT _acme-challenge %s comes to error %s \n", validationKey, err)
		}
		if *recordId != "" {
			log.Println("Auto resolve txt record success", *recordId)
		} else {
			log.Fatalln("Auto resolve txt record failed")
		}
		fmt.Println("record id:", *recordId)
		log.Println("Hook finish")
		//time.Sleep(time.Second * 20)
	case "del":
		log.Println("Will delete record TXT _acme-challenge")
		recordId, err := p.DeleteResolveDomainName(DNSTYPE, RR)
		if err != nil {
			log.Printf("Auto delete record TXT _acme-challenge %s comes to error %s \n", validationKey, err)
		}
		if *recordId != "" {
			log.Println("Auto delete txt record success")
		} else {
			log.Println("Auto delete txt record failed")
		}
	default:
		panic("action only support add or del ")
	}
	// 配置全栈加速
}

func parseCertbotEnvValue() {
	validationKey = os.Getenv("CERTBOT_VALIDATION")
	if validationKey == "" {
		panic("没有获取到 CERTBOT_VALIDATION")
	}

	domain = os.Getenv("CERTBOT_DOMAIN")
	if domain == "" {
		panic("没有获取到 CERTBOT_DOMAIN")
	}
}
