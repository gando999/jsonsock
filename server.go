package jsonsock

import (
	"fmt"
	"net"
	"os"
)

type SimpleServer struct {
	listener   net.Listener
	dispatcher Dispatcher
}

func CreateServer(connection string) *SimpleServer {
	s := new(SimpleServer)
	l, err := net.Listen("tcp", connection)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	fmt.Println("Listening on " + connection)
	s.listener = l
	s.dispatcher = CreateDispatcher()
	return s
}

func (simpleServer *SimpleServer) Start() {
	defer simpleServer.listener.Close()
	for {
		conn, err := simpleServer.listener.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		go handleRequest(conn, simpleServer.dispatcher)
	}
}

func (simpleServer *SimpleServer) Register(namespace string, targetImpl interface{}) {
	simpleServer.dispatcher.RegisterImpl(namespace, targetImpl)
}

func handleRequest(conn net.Conn, dispatcher Dispatcher) {
	buf := make([]byte, 1024)
	readLen, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	if readLen != 0 {
		request := buf[:readLen]
		requestJson, err := UnmarshalRequest(string(request))
		resolved, err := dispatcher.CallFunc(requestJson.Method, requestJson.Params)
		if err == nil {
			resp, err := MarshalResponse(resolved[0].Interface())
			if err == nil {
				conn.Write([]byte(resp + "\n"))
			}
		}
	}
	conn.Close()
}
