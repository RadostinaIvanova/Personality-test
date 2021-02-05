package main
import(
	"fmt"
	"bufio"
	"net"
	"os"
	"log"
	
)

func main(){
	conn,err := net.Dial("tcp", "localhost:9000")
	if err != nil{
		log.Fatal(err)
	}
	defer conn.Close()

	clientReader := bufio.NewReader(conn)
	clientWriter := bufio.NewWriter(conn)
	for{
		_, err:= clientReader.ReadString('?')
		if err != nil{
			fmt.Println(err)
		}
		readStd := bufio.NewReader(os.Stdin)
		text, err2 := readStd.ReadString('\n')
		if err2 != nil{
			fmt.Println(err2)
		}
		clientWriter.WriteString(text + "\n")
		clientWriter.Flush()
	}
}

