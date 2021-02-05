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

	clientReader := bufio.NewReaderSize(conn,5000)
	clientWriter := bufio.NewWriterSize(conn,5000)
	for{
		received, err:= clientReader.ReadString('?')
		if err != nil{
			fmt.Println("tuka e ggreshkata")
			fmt.Println(err)
		}
		fmt.Println(received)
		readStd := bufio.NewReader(os.Stdin)
		text, err2 := readStd.ReadString('.')
		if err2 != nil{
			fmt.Println(err2)
		}
		clientWriter.WriteString(text)
		clientWriter.Flush()
	}
}

