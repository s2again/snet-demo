package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/fanliao/go-promise"
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
	fmt.Print("输入新手宠物任务参数:")
	var noviceParam1 uint32
	n, err := fmt.Scanf("%d\n", &noviceParam1)
	if n < 1 {
		fmt.Println("too few input")
		os.Exit(-1)
	}
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}

	for {
		var sid string
		fmt.Print("输入SID:")
		n, err := fmt.Scanf("%40s\n", &sid)
		if n < 1 {
			fmt.Println("too few input")
			continue
		}
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		v, err := createFreshmenRole(sid, noviceParam1).Get()
		if err != nil {
			log.Println(err)
			fmt.Println("Error: ", err)
			continue
		}
		petinfo := v.(snet.PetInfo)
		fmt.Printf("精灵信息：\n%+v\n", petinfo)
		fmt.Printf("个体： %d 性格：%v\n", petinfo.Dv, petinfo.Nature)
	}
}

func createFreshmenRole(sid string, noviceParam1 uint32) (task *promise.Promise) {
	task = promise.NewPromise()
	go func() {
		// connect guide server
		conn, err := utils.ConnectGuideServer(guideAddr, sid)
		if err != nil {
			task.Reject(errors.New("ConnectGuideServer failed: " + err.Error()))
			return
		}
		// create
		err = createRole(conn)
		if err != nil {
			// reconnect guide server
			conn, err = utils.ConnectGuideServer(guideAddr, sid)
			if err != nil {
				task.Reject(errors.New("ConnectGuideServer failed: " + err.Error()))
				return
			}
		}
		// get online server list
		list, err := utils.GetOnlineServerList(conn)
		if err != nil || len(list) < 1 {
			task.Reject(errors.New("getServerList failed: " + err.Error()))
			return
		}
		// connect first online server, and close connection to sub server
		for _, server := range list {
			fmt.Printf("登录%d号服务器%s:%d，当前在线人数%d \n",
				server.OnlineID, server.IP, server.Port, server.UserCnt)
			var (
				onlineConn *snet.OnlineServerConnection
				resp       interface{}
				pet        snet.PetInfo
				err        error
			)
			onlineConn, err = utils.Guide2Online(conn, server)
			if err != nil {
				goto failOnce
			}
			// login online server
			resp, err = onlineConn.LoginOnline().Get()
			if err != nil {
				goto failOnce
			}
			goto connected
		connected:
			pet, err = finishNoviceAfterLogin(onlineConn, resp.(snet.LoginResponseFromOnline), noviceParam1)
			if err != nil {
				task.Reject(errors.New("connectOnline promise rejected: " + err.Error()))
				return
			}
			onlineConn.Close()
			task.Resolve(pet)
			return
		failOnce:
			fmt.Println("连接服务器失败 ", server.OnlineID, server.IP, server.Port, err.Error())
			log.Println(err.Error())
			continue
		}
		task.Reject(errors.New("connectOnline promise rejected: all servers inaccessible"))
	}()
	return task
}

func finishNoviceAfterLogin(conn *snet.OnlineServerConnection, info snet.LoginResponseFromOnline, noviceParam1 uint32) (pet snet.PetInfo, err error) {
	defer func() {
		x := recover()
		if x != nil {
			err = errors.New("Error: " + fmt.Sprint(x))
			log.Println(x)
			fmt.Println("发生错误：", x)
		}
	}()
	fmt.Printf("%+v\n", info)
	fmt.Println("登录成功")
	fmt.Printf("userID: %v sessionID: %X\n", conn.UserID, conn.SessionID)
	// Command_SYSTEM_TIME
	// Command_SYSTEM_TIME
	// Command_MAIL_GET_UNREAD
	// Command_NONO_INFO
	utils.MustAcceptAndCompleteTask(conn, 0x55, 1)                          // freshman suit
	task1stPet := utils.MustAcceptAndCompleteTask(conn, 0x56, noviceParam1) // freshman pet
	petTm := task1stPet.CaptureTm
	utils.MustAcceptAndCompleteTask(conn, 0x57, 1) // freshman pet ball
	utils.MustAcceptAndCompleteTask(conn, 0x58, 1) // freshman money
	utils.MustResolvePromise(conn.ReleasePet(petTm, 1))
	utils.MustResolvePromise(conn.LeaveMap())
	utils.MustResolvePromise(conn.EnterMap(0, 8, 0x1c8, 0x8f))
	utils.MustResolvePromise(conn.ListMapPlayer())

	petinfo := utils.MustResolvePromise(conn.GetPetInfo(petTm))
	return petinfo.(snet.PetInfo), nil
}

// noinspection GoUnusedFunction
func createRole(conn *snet.GuideServerConnection) error {
	var nickname [16]byte
	copy(nickname[:], "小沙雕")
	_, err := conn.CreateRole(nickname, snet.RoleGreen).Get()
	if err != nil {
		return err
	}
	return nil
}
