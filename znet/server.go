package znet

import (
	"fmt"
	"net"
	"zinx/ziface"
	"zinx/zutils"
)

type Server struct {
	Name      string
	IPVersion string
	IP        string
	Port      int16

	// 消息处理器
	MessageHandler ziface.IMessageHandler

	// 连接管理器
	ConnManager ziface.IConnManager

	StartHock func(ziface.IConnection)
	StopHock  func(ziface.IConnection)
}

func NewServer(name string) ziface.IServer {
	zutils.ZINX.Load()

	s := Server{
		Name:           zutils.ZINX.Name,
		IPVersion:      "tcp4",
		IP:             zutils.ZINX.Host,
		Port:           zutils.ZINX.TcpPort,
		MessageHandler: NewMessageHandler(),
		ConnManager:    NewConnManager(),
	}
	return &s
}

func (s *Server) Start() {
	fmt.Printf("[START] Server listener at IP: %s, Port %d, is starting\n", s.IP, s.Port)
	fmt.Printf("[Zinx] Version: %s, MaxConn: %d, MaxPacketSize: %d\n",
		zutils.ZINX.Version,
		zutils.ZINX.MaxConn,
		zutils.ZINX.MaxPacketSize)

	go func() {
		// 0. 开启Worker池
		s.MessageHandler.StartWorkerPool()

		// 1. 解析IP、端口以及协议族，返回一个Addr对象
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("reslove addr err: ", err)
			return
		}

		// 2. 开始监听
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen tcp err ", err)
			return
		}

		// 3. 处理连接
		// 连接对象的ID，初始值为1000
		var cid int32 = 1000
		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("accept err ", err)
				continue
			}

			dealConn := NewConntion(s, conn, cid, s.MessageHandler)
			cid++

			go dealConn.Start()
		}
	}()
}

func (s *Server) Stop() {
	fmt.Println("[STOP] Zinx server , name ", s.Name)
	s.ConnManager.ClearAll()
}

func (s *Server) Serve() {
	s.Start()

	select {}
}

func (s *Server) AddRouter(messageId uint32, router ziface.IRouter) {
	s.MessageHandler.AddRouter(messageId, router)
}

func (s *Server) GetConnManager() ziface.IConnManager {
	return s.ConnManager
}

func (s *Server) SetOnConnStart(hock func(ziface.IConnection)) {
	s.StartHock = hock
}

func (s *Server) SetOnConnStop(hock func(ziface.IConnection)) {
	s.StopHock = hock
}

func (s *Server) CallOnConnStart(conn ziface.IConnection) {
	if s.StartHock != nil {
		s.StartHock(conn)
	}
}

func (s *Server) CallOnConnStop(conn ziface.IConnection) {
	if s.StopHock != nil {
		s.StopHock(conn)
	}
}
