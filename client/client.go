package main
import(
	"fmt"
	"bufio"
	"net"
	"os"
	
)

func main(){
	conn,err := net.Dial("tcp", "localhost:9000")
	if err != nil{
		panic(err)
	}
	defer conn.Close()
	fmt.Println("Server accessed")

	clientReader := bufio.NewReader(conn)
	clientWriter := bufio.NewWriter(conn)
	for{
	messageReceived, err:= clientReader.ReadString('\n')
	if err != nil{
		fmt.Println("Client couldn't receive a message")
	}
	fmt.Println(messageReceived)

	readStd := bufio.NewReader(os.Stdin)
    fmt.Print("Enter text: ")
	text, _ := readStd.ReadString('\n')
	
	
	clientWriter.WriteString(text + "\n")
	clientWriter.Flush()
}
}

