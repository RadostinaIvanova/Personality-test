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
	"io/ioutil"
	"github.com/RadostinaIvanova/golang-project/classificator"
	"github.com/RadostinaIvanova/golang-project/corpus"
	
)

const questionsDoc string = "C://Users//Radi//Downloads//questions.txt"
const pathToDescriptions string = "D:\\FMI\\Info\\PersonalityTypes4\\" 
const corpusName  string = "D:\\FMI\\golang_workspace\\src\\mbt\\mbt.csv"
const classicatorFileName string = "trainedClassificator"

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

func saveTrainedClassificator(filename string){
	trainSet,_ := corpus.MakeClassesFromFile(corpusName)
	c := classificator.TrainMultinomialNB(trainSet)
	writeEncodedClassificatorToFile(filename,c )
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

func extractClassificator(filename string) classificator.NBclassificator{
	if !exists(filename){
		saveTrainedClassificator(filename)
	}
	c := loadTrainedClassificator(filename)
	return c
}
func writeEncodedClassificatorToFile(filename string, c classificator.NBclassificator){
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

func classificate(answers string, c classificator.NBclassificator) string {
	return personalityTypes(classificator.ApplyMultinomialNB(c,answers))
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

func optionHandler(documentsPath string, serverReader bufio.Reader, serverWriter bufio.Writer){
	for{
		optionReceived, err2 := serverReader.ReadString('\n')
		if err2!= nil{
			log.Println(err2.Error())
		}
		optionReceived = strings.TrimSuffix(optionReceived, "\n")
		if optionReceived == "Quit"{
			break
		}
		str := extractInfo(optionReceived, documentsPath)
		serverWriter.WriteString(str + "\n");
		serverWriter.Flush();
	}
}

func sendType(pType string, serverReader bufio.Reader, serverWriter bufio.Writer){
	serverWriter.WriteString(pType + "\n");
	serverWriter.Flush();
}
func handleConnection(conn net.Conn, questions []string, c classificator.NBclassificator){
	//fmt.Println("Inside handle connection func")
	defer conn.Close()
	serverWriter := bufio.NewWriterSize(conn,5000)
	serverReader := bufio.NewReaderSize(conn,5000)

	answers := quiz(questions, *serverReader, *serverWriter)
	pType := classificate(answers,c)
	sendType(pType, *serverReader, *serverWriter)

	documentsPath := pathToDescriptions + pType + "\\"
	optionHandler(documentsPath,*serverReader, *serverWriter)
	
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
	
	questions := extractQuestionsFromFile(questionsDoc)
	c := extractClassificator(classicatorFileName)

	for {
		conn,err := ln.Accept()
		if err!= nil{
			log.Println(err.Error())
		}
		go handleConnection(conn,questions,c)
	}
}	