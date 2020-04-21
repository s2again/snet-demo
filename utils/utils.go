package utils

import (
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"strconv"

	"github.com/fanliao/go-promise"
	"github.com/s2again/snet"
	"github.com/s2again/snet/core"
)

// noinspection GoUnusedFunction
func ConnectGuideServer(loginAddr *net.TCPAddr, sid string) (conn *snet.GuideServerConnection, err error) {
	defer func() {
		x := recover()
		if x != nil {
			err = fmt.Errorf("unknown error: %v", x)
		}
	}()
	uid, sessionID, err := core.ParseSIDString(sid)
	if err != nil {
		return nil, fmt.Errorf("ParseSIDString Error: %v", err)
	}

	conn, err = snet.ConnectGuideServer(loginAddr)
	if err != nil {
		return nil, fmt.Errorf("ConnectGuideServer Error: %v", err)
	}
	conn.SetSession(uid, sessionID)
	return
}

// noinspection GoUnusedFunction
func Guide2Online(conn *snet.GuideServerConnection, server snet.OnlineServerInfo) (onlineConn *snet.OnlineServerConnection, err error) {
	// login online
	onlineConn, err = snet.ConnectOnlineServer(server, conn.UserID, conn.SessionID)
	if err != nil {
		return nil, err
	}
	conn.Close()
	return
}

func GetOnlineServerList(conn *snet.GuideServerConnection) ([]snet.OnlineServerInfo, error) {
	v, err := conn.ListOnlineServers().Get()
	if err != nil {
		return nil, err
	}
	info := v.(snet.CommendSvrInfo)
	fmt.Printf("CommendSvrInfo %+v\n", info)
	return info.SvrList, nil
}

func ParseSID(sid string) (userID uint32, session [16]byte, err error) {
	if len(sid) != 40 {
		err = errors.New("illegal parameter")
		return
	}
	userIDtmp, err := strconv.ParseUint(sid[:8], 16, 32)
	userID = uint32(userIDtmp)

	sessiontmp, err := hex.DecodeString(sid[8:40])
	if err != nil {
		return
	}
	copy(session[:], sessiontmp[:32])
	return
}

func MustResolvePromise(p *promise.Promise) interface{} {
	v, err := p.Get()
	if err != nil {
		panic(err)
	}
	return v
}

func MustAcceptAndCompleteTask(conn *snet.OnlineServerConnection, taskID uint32, param uint32) snet.NoviceFinishInfo {
	_, err := conn.AcceptTask(taskID).Get()
	if err != nil {
		panic(err)
	}
	result := MustResolvePromise(conn.CompleteTask(taskID, param))
	fmt.Println("finish novice", result)
	return result.(snet.NoviceFinishInfo)
}
