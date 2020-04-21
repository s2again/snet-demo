package config

import (
	"encoding/xml"
	"io/ioutil"
	"net"
	"net/http"
)

const ServerConfigURL = "http://51seer.61.com/config/ServerR.xml"

//noinspection SpellCheckingInspection
type ServerConfig struct {
	IpConfig struct {
		HTTP struct {
			URL string `xml:"url,attr"`
		} `xml:"http"`
		Email struct {
			IP   string `xml:"ip,attr"`
			Port int    `xml:"port,attr"`
		} `xml:"Email"`
		Dir struct {
			IP   string `xml:"ip,attr"`
			Port int    `xml:"port,attr"`
		} `xml:"DirSer"`
		Visitor struct {
			IP   string `xml:"ip,attr"`
			Port int    `xml:"port,attr"`
		} `xml:"Visitor"`
		SubServer struct {
			IP   string `xml:"ip,attr"`
			Port int    `xml:"port,attr"`
		} `xml:"SubServer"`
		RegistSer struct {
			IP   string `xml:"ip,attr"`
			Port int    `xml:"port,attr"`
		} `xml:"RegistSer"`
	} `xml:"ipConfig"`
}

func GetServerConfig() (config *ServerConfig, err error) {
	return GetServerConfigFrom(ServerConfigURL)
}

func GetServerConfigFrom(url string) (config *ServerConfig, err error) {
	resp, err := http.Get(url)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return
	}
	config = new(ServerConfig)
	err = xml.NewDecoder(resp.Body).Decode(config)
	return
}

func (c *ServerConfig) GetGuideServerByHTTP() (addr *net.TCPAddr, err error) {
	resp, err := http.Get(c.IpConfig.HTTP.URL)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	addr, err = net.ResolveTCPAddr("tcp", string(respBody))
	return
}
