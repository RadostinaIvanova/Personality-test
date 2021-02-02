package main

import(
	"bufio"
	"net"
	"log"
	"fmt"
	"os"
	 "strconv"
	//"errors"
//	"github.com/RadostinaIvanova/golang-course/NaiveBayesClassificator"
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
		if err2 != nil{
			fmt.Println(err2)
		}
		_, err3 := f.WriteString(messageReceived)
		   if err3 != nil{
			   fmt.Println(err3)
			}
		}
	f.Close()
}
func handleConnection(conn net.Conn,indDoc int, questions []string){
	//fmt.Println("Inside handle connection func")
	defer conn.Close()
	serverWriter := bufio.NewWriter(conn)
	serverReader := bufio.NewReader(conn)
	quiz(questions, *serverReader, *serverWriter, indDoc)

}

func extractQuestionsFromFile(filename string) []string{
	questions := []string{}
	f , errf := os.Open(filename)
	check(errf)
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan(){
		questions = append(questions, scanner.Text())
	}
	return questions 
}
func main(){
	// fmt.Println("Launching server")
	fmt.Println("Listen on port")
	ln, err := net.Listen("tcp", ":9000")
	check(err)

	var indDoc int = 0
	questionsDoc := "C://Users//Radi//Downloads//questions.txt"
	questions := extractQuestionsFromFile(questionsDoc)
	for{
		conn,err := ln.Accept()
		if err != nil{
			log.Fatal(err)
		}
		// fmt.Println("Calling go routine - handle conneciton")
		indDoc++
		go handleConnection(conn,indDoc,questions)
	}
}	