package main

import(
	"bufio"
	"net"
	"log"
	"fmt"
	"os"
	"strconv"
	 "encoding/gob"
	//"errors"
	"github.com/RadostinaIvanova/golang-project/classificator"
	"github.com/RadostinaIvanova/golang-project/corpus"
	
)


func loadTrainedClassificator(filename string) classificator.NBclassificator{
	f, err := os.Open(filename)
	if err != nil{
		log.Println(err.Error())
	}
	defer f.Close()

	c := classificator.NBclassificator{}
	decoder := gob.NewDecoder(f)
	decoder.Decode(&c)	
	return c
}

func writeEncodedToFile(filename string, c classificator.NBclassificator){
	f, err := os.Create(filename)
	if err != nil{
		log.Println(err.Error())
	}
	defer f.Close()
	encoder := gob.NewEncoder(f)
	encoder.Encode(c.CondProb)
}

// Exists reports whether the named file or directory exists.
func exists(name string) bool {
    if _, err := os.Stat(name); err != nil {
        if os.IsNotExist(err) {
            return false
        }
    }
    return true
}

func classificate(answers string, c classificator.NBclassificator) int{
	return classificator.ApplyMultinomialNB(c,answers)
}

func quiz(questions []string, serverReader bufio.Reader, serverWriter bufio.Writer) string{
	//docName := "answers" + strconv.Itoa(indDoc)
	//f,err := os.OpenFile(docName, os.O_WRONLY|os.O_CREATE, 0666)
	//check(err)
	var answers string = ""
	for _, question := range questions{
		serverWriter.WriteString(question)
		serverWriter.Flush()

		messageReceived, err2 := serverReader.ReadString('\n')
		if err2!= nil{
			log.Println(err2.Error())
		}
		answers += messageReceived
	//		_, err3 := f.WriteString(messageReceived)
	//		if err3 != nil{
	//			fmt.Println(err3)
	//			}
		}
	//f.Close()
	return answers
}
func handleConnection(conn net.Conn,indDoc int, questions []string, c classificator.NBclassificator){
	//fmt.Println("Inside handle connection func")
	defer conn.Close()
	serverWriter := bufio.NewWriter(conn)
	serverReader := bufio.NewReader(conn)
	answers := quiz(questions, *serverReader, *serverWriter)
	result := classificate(answers,c)
	serverWriter.WriteString(strconv.Itoa(result))
}

func extractQuestionsFromFile(filename string) []string{
	questions := []string{}
	f , errf := os.Open(filename)
	if errf!= nil{
		log.Println(errf.Error())
	}
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
	if err!= nil{
		log.Println(err.Error())
	}
	//var indDoc int = 0
	questionsDoc := "C://Users//Radi//Downloads//questions.txt"
	questions := extractQuestionsFromFile(questionsDoc)
	filename := "classificator.go"
	if !exists(filename){
		trainSet, testSet := 
		writeEncodedToFile(filename)
	}
	c := loadTrainedClassificator(filename)
	
	for{
		conn,err := ln.Accept()
		if err!= nil{
			log.Println(err.Error())
		}
		// fmt.Println("Calling go routine - handle conneciton")
		//indDoc++
		go handleConnection(conn,questions)
	}
}	