// unavailable
package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/s2again/snet"

	"snet_demo/config"
)

const guideAddrString = "129.211.42.134:1863"

var (
	configFile *config.ServerConfig
	guideAddr  *net.TCPAddr
)

func init() {
	var err error
	f, err := os.OpenFile("seer.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
	log.SetOutput(f)
	configFile, err = config.GetServerConfig()
	if err != nil {
		panic(err)
	}
	fmt.Println(configFile)
	guideAddr, err = net.ResolveTCPAddr("tcp4", guideAddrString)
	if err != nil {
		panic(err)
	}
	fmt.Println(guideAddr)
}

const (
	testAccount  = 922610556
	testPassword = "123456zzz"
)

var exit = make(chan interface{})

func main() {
	loginConn, err := snet.ConnectGuideServer(guideAddr)
	if err != nil {
		panic(err)
	}
	p := loginConn.Login(testAccount, testPassword)
	p.OnSuccess(func(v interface{}) {
		info := v.(snet.LoginResponse)
		fmt.Printf("%+v\n", info)
		fmt.Printf("sid: %X%X", loginConn.UserID, info.SessionID)
		exit <- 1
	}).OnFailure(func(v interface{}) {
		fmt.Println(v)
		exit <- 1
	})
	<-exit
}
