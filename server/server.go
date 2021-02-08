package main

import(
	"bufio"
	"net"
	"log"
	"fmt"
	"os"
	"strconv"
	 "encoding/gob"
	 "strings"
	//"errors"
	"io/ioutil"
	"github.com/RadostinaIvanova/golang-project/classificator"
	"github.com/RadostinaIvanova/golang-project/corpus"
	
)

func extractInfo(option string, path string) string{
	filename := path + option + ".txt"
	buff, err := ioutil.ReadFile(filename)
	if err != nil {
        fmt.Print(err)
    }
	text := string(buff)
	return text
}

func personalityTypes(pType int) string{
	switch pType{
		case 0:  return "Diplomat"
		case 1:  return "Analyst"
		case 2:  return "Sentinel"
		case 3:  return "Explorer"
	}
	return "Diplomat"
}
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
	var answers string = ""
	serverWriter.WriteString(strconv.Itoa(len(questions)) + "\n")
	serverWriter.Flush()
	for _ , question := range questions{
		serverWriter.WriteString(question)
		serverWriter.Flush()
		messageReceived, err2 := serverReader.ReadString('\n')
		if err2!= nil{
			log.Println(err2.Error())
		}
		messageReceived = strings.TrimSuffix(messageReceived, "\n")
		answers += " "
		answers += messageReceived
	}
	return answers
}
func handleConnection(conn net.Conn, questions []string, c classificator.NBclassificator){
	//fmt.Println("Inside handle connection func")
	defer conn.Close()
	serverWriter := bufio.NewWriterSize(conn,5000)
	serverReader := bufio.NewReaderSize(conn,5000)
	answers := quiz(questions, *serverReader, *serverWriter)
	result := classificate(answers,c)
	pType := personalityTypes(result)
	fmt.Println(pType)
	serverWriter.WriteString(pType + "\n");
	serverWriter.Flush();
	path := "D:\\FMI\\Info\\PersonalityTypes4\\" + pType + "\\"
	var optionReceived string
	for ; optionReceived != "Exit";{
		optionReceived, err2 := serverReader.ReadString('\n')
		if err2!= nil{
			log.Println(err2.Error())
		}
		optionReceived = strings.TrimSuffix(optionReceived, "\n")
		str := extractInfo(optionReceived, path)
		serverWriter.WriteString(str + "\n");
		serverWriter.Flush();
	}
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
	fmt.Println("Listen on port")
	ln, err := net.Listen("tcp", ":9000")
	if err!= nil{
		log.Println(err.Error())
	}
	questionsDoc := "C://Users//Radi//Downloads//questions.txt"
	questions := extractQuestionsFromFile(questionsDoc)

	filename := "trainedClassificator"
	if !exists(filename){
		corpusName := "D:\\FMI\\golang_workspace\\src\\mbt\\mbt.csv"
		trainSet,testSet := corpus.MakeClassesFromFile(corpusName)
		c := classificator.TrainMultinomialNB(trainSet)
		writeEncodedToFile(filename,c )
		classificator.TestClassifier(c,testSet)
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