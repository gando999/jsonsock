package jsonsock

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

type SimpleClient struct {
	tcpConn *net.TCPConn
}

func CreateClient(connection string) *SimpleClient {
	c := new(SimpleClient)
	service := connection
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)
	c.tcpConn = conn
	return c
}

func (simpleClient *SimpleClient) Send(targetFunc string, params []interface{}) string {
	req, err := MarshalRequest(targetFunc, params)
	checkError(err)
	fmt.Fprintf(simpleClient.tcpConn, req+"\n")
	checkError(err)
	message, _ := bufio.NewReader(simpleClient.tcpConn).ReadString('\n')
	return message
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
