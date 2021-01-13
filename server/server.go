package main

import(
	"fmt"
	"bufio"
	"net"
	"log"
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
		fmt.Println(err)
	}

	fmt.Println(messageReceived)
	}
}
func main(){
	fmt.Println("Launching server")
	fmt.Println("Listen on port")
	ln, err := net.Listen("tcp", ":9000")
	if err != nil{
		log.Fatal(err)
	}
	for{
		fmt.Println("Accept connection on port")
		conn,err := ln.Accept()
		if err != nil{
			log.Fatal(err)
		}
		fmt.Println("Calling go routine - hangle conneciton")
		go handleConnection(conn)
	}
}	