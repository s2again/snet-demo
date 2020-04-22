// 客户端模拟Demo
// http://51seer.61.com/?sid=
package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/s2again/snet"

	"snet_demo/config"
	"snet_demo/utils"
)

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
	guideAddr, err = configFile.GetGuideServerByHTTP()
	if err != nil {
		panic(err)
	}
	fmt.Println(guideAddr)
}

func main() {
	// Login
	var sid string
	fmt.Println("Input SID:")
	n, err := fmt.Scanf("%s", &sid)
	if n < 1 {
		return
	}
	if err != nil {
		panic(err)
	}

	loginConn, err := snet.ConnectGuideServer(guideAddr)
	if err != nil {
		panic(err)
	}
	conn, err := login(loginConn, sid)
	if err != nil {
		panic(err)
	}
	loginConn.Close()
	fmt.Printf("userID: %v sessionID: %v\n", loginConn.UserID, loginConn.SessionID)

	petlist := utils.MustResolvePromise(conn.GetPetList())
	fmt.Printf("精灵列表： \n")
	for _, pet := range petlist.([]snet.PetListInfo) {
		fmt.Printf("%+v\n", pet)
	}
	for _, pet := range petlist.([]snet.PetListInfo) {
		petinfo := utils.MustResolvePromise(conn.GetPetInfo(pet.CatchTime))
		fmt.Printf("%+v\n", petinfo)
	}
	select {}
}

func login(loginConn *snet.GuideServerConnection, sid string) (conn *snet.OnlineServerConnection, err error) {
	userID, session, err := utils.ParseSIDString(sid)
	if err != nil {
		panic(err)
	}
	loginConn.SetSession(userID, session)
	v, err := loginConn.ListOnlineServers().Get()
	if err != nil {
		panic(err)
	}
	info := v.(snet.CommendSvrInfo)
	firstOnline := info.SvrList[0]
	fmt.Println("Login into Online", firstOnline.OnlineID, firstOnline.IP, firstOnline.Port)
	return snet.ConnectOnlineServer(firstOnline, userID, session)
}
