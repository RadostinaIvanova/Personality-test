package main

import(
	"fmt"
	"bufio"
	"net"
	//"io/ioutil"
)

func handleConnection(conn net.Conn){
	fmt.Println("Inside handle connection func")
	defer conn.Close()

	serverWriter := bufio.NewWriter(conn)
	serverReader := bufio.NewReader(conn)
	for{
	serverWriter.WriteString("Some message from server" + "\n")
	serverWriter.Flush()
	
	messageReceived, err := serverReader.ReadString('\n')
	if err!= nil{
		fmt.Println("Server couldnt read the message")
	}

	fmt.Println(messageReceived)
}
}
func main(){
	fmt.Println("Launching server")
	fmt.Println("Listen on port")
	ln, err := net.Listen("tcp", ":9000")
	if err != nil{
		panic(nil)
	}
	for{
		fmt.Println("Accept connection on port")
		conn,err := ln.Accept()
		if err != nil{
			panic(nil)
		}
		fmt.Println("Calling go routine - hangle conneciton")
		go handleConnection(conn)
	}
}	