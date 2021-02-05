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
	errd := decoder.Decode(&c)	
	if errd != nil {
		log.Fatal("decode error 1:", errd)
	}
	return c
}

func writeEncodedToFile(filename string, c classificator.NBclassificator){
	f, err := os.Create(filename)
	if err != nil{
		log.Println(err.Error())
	}
	defer f.Close()
	encoder := gob.NewEncoder(f)
	encoder.Encode(c)
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
	fmt.Println(questions)
	return answers
}
func handleConnection(conn net.Conn, questions []string, c classificator.NBclassificator){
	//fmt.Println("Inside handle connection func")
	defer conn.Close()
	serverWriter := bufio.NewWriter(conn)
	serverReader := bufio.NewReader(conn)
	answers := quiz(questions, *serverReader, *serverWriter)
	result := classificate(answers,c)
	fmt.Println("classificate")
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
	
	filename := "trainedClassificator"
	if !exists(filename){
		corpusName := "D:\\FMI\\golang_workspace\\src\\mbt\\mbt.csv"
		trainSet,_ := corpus.MakeClassesFromFile(corpusName)
		c := classificator.TrainMultinomialNB(trainSet)
		writeEncodedToFile(filename,c )
	}
	c := loadTrainedClassificator(filename)
	for{
		conn,err := ln.Accept()
		if err!= nil{
			log.Println(err.Error())
		}
		go handleConnection(conn,questions,c)
	}
}	