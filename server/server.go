package main

import(
	"bufio"
	"net"
	"log"
	"fmt"
	"os"
	"strconv"
	"strings"
	"io/ioutil"
	"github.com/RadostinaIvanova/Personality-test/classificator"
	"github.com/RadostinaIvanova/Personality-test/corpus"
	"github.com/RadostinaIvanova/Personality-test/model"
)

const questionsDoc string = "C://Users//Radi//Downloads//questions.txt"
const pathToDescriptions string = "D:\\FMI\\Info\\PersonalityTypes4\\" 
const corpusName  string = "D:\\FMI\\golang_workspace\\src\\mbt\\mbt.csv"
const classicatorFileName string = "D:\\FMI\\golang_workspace\\src\\golang-project\\server\\trainedClassificator"
const dialoguesCorpus string = "D:\\FMI\\Info\\dialogues_train.txt"
const modelFileName string = "D:\\FMI\\golang_workspace\\src\\golang-project\\server\\trainedModel"

func extractInfo(option string, path string) string{
	filename := path + option + ".txt"
	buff, err := ioutil.ReadFile(filename)
	if err != nil {
        fmt.Print(err)
    }
	text := string(buff)
	return text
}

func extractClassificator(filename string, corpusName string) classificator.NBclassificator {
	c := classificator.NBclassificator{}
	if !exists(filename){
		trainSet,_ := corpus.MakeClassesFromFile(corpusName)
		c.TrainMultinomialNB(trainSet)
		c.SaveClassificator(filename)
	}else{ 
		c.LoadClassificator(filename)
	}
	return c
}

func extractModel(filename string, dialoguesFile string, limit int) model.MarkovModel{
	m := model.MarkovModel{}
	if !exists(filename){
		corpus := model.Extract(dialoguesFile)
		fullCorpus := model.FullSentCorpus(corpus)
		train, _ := model.DivideIntoTrainAndTest(0.1, fullCorpus)
		numGram := 2
		m.Init(numGram,train,limit)
		fmt.Println("here")
		m.SaveModel(filename)
	}else{ 
		m.LoadModel(filename)
	}
	return m
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

func personalityTypes(pType int) string{
	switch pType{
		case 0:  return "Diplomat"
		case 1:  return "Analyst"
		case 2:  return "Sentinel"
		case 3:  return "Explorer"
	}
	return "Diplomat"
}


func classificate(answers string, c classificator.NBclassificator) string {
	return personalityTypes(c.ApplyMultinomialNB(answers))
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


func sendType(pType string, serverReader bufio.Reader, serverWriter bufio.Writer){
	serverWriter.WriteString(pType + "\n");
	serverWriter.Flush();
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


func optionHandler(documentsPath string, serverReader bufio.Reader, serverWriter bufio.Writer){
	for{
		optionReceived, err2 := serverReader.ReadString('\n')
		if err2!= nil{
			fmt.Println("Unexpected disconnection")
			break;
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

func main(){
	// fmt.Println("Listen on port")
	// ln, err := net.Listen("tcp", ":9000")
	// if err!= nil{
	// 	log.Println(err.Error())
	// }
	
	// questions := extractQuestionsFromFile(questionsDoc)
	// c := extractClassificator(classicatorFileName, corpusName)
	m := extractModel(modelFileName, dialoguesCorpus,400000)
	fmt.Println(m.BestContinuation([]string{"play", "games"}, 0.7, 15))
	// for {
	// 	conn,err := ln.Accept()
	// 	if err!= nil{
	// 		log.Println(err.Error())
	// 	}
	// 	go handleConnection(conn,questions,c)
	// }
}	