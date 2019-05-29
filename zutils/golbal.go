package zutils

import (
	"encoding/json"
	"io/ioutil"
	"zinx/ziface"
)

type Global struct {
	//项目目录
	AppDir string

	// 全局Server对象
	TcpServer ziface.IServer

	// 当前服务器主机
	Host string

	// 当前服务器监听端口
	TcpPort int16

	// 服务器名字
	Name string

	// 当前框架版本
	Version string

	// 数据报最大值
	MaxPacketSize uint32

	// 允许的最大连接数
	MaxConn int32

	// Worker池的Worker最大数量
	WorkerPoolSize uint32

	// Response数据缓冲通道的大小
	DataBufSize int32
}

var ZINX *Global

func (g *Global) loadConfig() {
	data, err := ioutil.ReadFile("D:\\go\\src\\zinx\\conf\\config.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(data, &ZINX)
	if err != nil {
		panic(err)
	}
}

func init() {
	ZINX = &Global{
		Name:           "ZINX_APP",
		Version:        "v0.4",
		TcpPort:        7777,
		Host:           "0.0.0.0",
		MaxConn:        12000,
		MaxPacketSize:  4096,
		WorkerPoolSize: 0,
		DataBufSize:    5,
	}
	ZINX.loadConfig()
}

// 空方法，主要用于初始化全局对象
func (g *Global) Load() {}
