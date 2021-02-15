package main

import(
	"fmt"
	"bufio"
	"net"
	"log"
	"github.com/rivo/tview"
	"strings"
	"strconv"	
)

//receives questions from server and send their answers to server
func answerQuiz(clientReader bufio.Reader, clientWriter bufio.Writer, app *tview.Application) {
	received, err:= clientReader.ReadString('\n')
	if err != nil{
		fmt.Println(err)
	}
	received = strings.TrimSuffix(received, "\n")
	numOfQuestions,_ := strconv.Atoi(received)
	for i := 0; i < numOfQuestions; i++ {
		received, err:= clientReader.ReadString('?')
		if err != nil{
			fmt.Println(err)
		}
		text := getInputFieldText(app,received)
		clientWriter.WriteString(text + "\n")
		clientWriter.Flush()
	}
}

//receives the result from server of the classification 
func receiveClassification(clientReader bufio.Reader, app *tview.Application){
	classificationResult, err:= clientReader.ReadString('\n')
	if err != nil{
		fmt.Println(err)
	}
	classificationResult = strings.TrimSuffix(classificationResult, "\n")
	modalClassification(app, classificationResult)
}

func handleConnection(clientReader bufio.Reader, clientWriter bufio.Writer){
	app := tview.NewApplication()
	setStartTestButton(startTest,app)
	answerQuiz(clientReader, clientWriter, app)	
	receiveClassification(clientReader, app)
	options(clientReader,clientWriter)
}


func main(){
	conn,err := net.Dial("tcp", "localhost:9000")
	if err != nil{
		log.Fatal(err)
	}
	defer conn.Close()
	clientReader := bufio.NewReaderSize(conn,5000)
	clientWriter := bufio.NewWriterSize(conn,5000)
	handleConnection(*clientReader, *clientWriter)	
}

