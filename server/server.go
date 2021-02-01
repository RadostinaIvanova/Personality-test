package main

import(
	"fmt"
	"bufio"
	"net"
	"log"
	//"io"
	"os"
	 "strconv"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func quiz(questions []string, serverReader bufio.Reader, serverWriter bufio.Writer, indDoc int){
	docName := "answers" + strconv.Itoa(indDoc)
	f,err := os.OpenFile(docName, os.O_WRONLY|os.O_CREATE, 0666)
	check(err)
	for _, question := range questions{
		serverWriter.WriteString(question)
		serverWriter.Flush()
		messageReceived, err2 := serverReader.ReadString('\n')
		check(err2)
		_, errw2 := f.WriteString(messageReceived)
   	    check(errw2)
	}
	f.Close()
}
func handleConnection(conn net.Conn,indDoc int){
	//fmt.Println("Inside handle connection func")
	defer conn.Close()

	serverWriter := bufio.NewWriter(conn)
	serverReader := bufio.NewReader(conn)
	str := []string{ "KOI", "TOI", "AZ"}
	quiz(str, *serverReader, *serverWriter, 2)
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
	var indDoc int = 0
	for{
		fmt.Println("Accept connection on port")
		conn,err := ln.Accept()
		if err != nil{
			log.Fatal(err)
		}
		fmt.Println("Calling go routine - hangle conneciton")
		indDoc++
		go handleConnection(conn,indDoc)
	}
}	