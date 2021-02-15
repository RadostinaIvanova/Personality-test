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
	"github.com/RadostinaIvanova/Personality-test/model"
)

const questionsDoc string = "D:\\FMI\\golang_workspace\\src\\golang-project\\resources\\questions.txt"
const pathToDescriptions string = "D:\\FMI\\golang_workspace\\src\\golang-project\\resources\\PersonalityTypes4\\" 
const corpusName  string = "D:\\FMI\\golang_workspace\\src\\golang-project\\resources\\mbt.csv"
const dialoguesCorpus string = "D:\\FMI\\golang_workspace\\src\\golang-project\\resources\\dialogues_train.txt"
const classicatorFileName string = "D:\\FMI\\golang_workspace\\src\\golang-project\\server\\trainedClassificator"
const modelFileName string = "D:\\FMI\\golang_workspace\\src\\golang-project\\server\\trainedModel"
const expanding int = 2



//return trained classificator and if file with already trained model - loads it
//if not trains one and save it to file and returns it
func extractClassificator(filename string, corpusName string) classificator.NBclassificator {
	c := classificator.NBclassificator{}
	if !exists(filename){
		trainSet,_ := classificator.MakeClassesFromFile(corpusName)
		c.TrainMultinomialNB(trainSet)
		c.SaveClassificator(filename)
	}else{ 
		c.LoadClassificator(filename)
	}
	return c
}

//returns trained language Markov Model and if file with already trained model - loads it
//if not trains one and save it to file and returns it
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

//converts types from encoded ints to their real names as strings 
func personalityTypes(pType int) string{
	switch pType{
		case 0:  return "Diplomat"
		case 1:  return "Analyst"
		case 2:  return "Sentinel"
		case 3:  return "Explorer"
	}
	return "Analyst"
}

//returns result from classification after applying trained classificator
func classificate(answers string, c classificator.NBclassificator) string {
	return personalityTypes(c.ApplyMultinomialNB(answers))
}

//sends questions and receives answers from client
func quiz(questions []string, serverReader bufio.Reader, serverWriter bufio.Writer, m model.MarkovModel) string{
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
		if len(m.BestContinuation(strings.Split(messageReceived, " "),0.6, expanding)) <= expanding{
			addContext := strings.Join(m.BestContinuation(strings.Split(messageReceived, " "),0.6, 3), " ")
			answers += " "
			messageReceived += addContext
		}
		
		answers += " "
		answers += messageReceived
	}
	return answers
}

//send the result from classification
func sendType(pType string, serverReader bufio.Reader, serverWriter bufio.Writer){
	serverWriter.WriteString(pType + "\n");
	serverWriter.Flush();
}

//extract questions which will ask the client
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
// Exists reports whether the named file or directory exists.
func exists(name string) bool {
    if _, err := os.Stat(name); err != nil {
        if os.IsNotExist(err) {
            return false
        }
    }
    return true
}

//extract info which is demanded from the client
func extractInfo(option string, path string) string{
	filename := path + option + ".txt"
	buff, err := ioutil.ReadFile(filename)
	if err != nil {
        fmt.Print(err)
    }

	text := string(buff)
	return text
}

//reveiving  which information for the personality type is wanted from the client
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

func handleConnection(conn net.Conn, questions []string, c classificator.NBclassificator,m model.MarkovModel){
	//fmt.Println("Inside handle connection func")
	defer conn.Close()
	serverWriter := bufio.NewWriterSize(conn,5000)
	serverReader := bufio.NewReaderSize(conn,5000)

	answers := quiz(questions, *serverReader, *serverWriter, m)
	pType := classificate(answers,c)
	sendType(pType, *serverReader, *serverWriter)

	documentsPath := pathToDescriptions + pType + "\\"
	fmt.Println(documentsPath)
	optionHandler(documentsPath,*serverReader, *serverWriter)
	
}

func main(){
	fmt.Println("Listen on port")
	ln, err := net.Listen("tcp", ":9000")
	if err!= nil{
		log.Println(err.Error())
	}
	
	questions := extractQuestionsFromFile(questionsDoc)
	c := extractClassificator(classicatorFileName, corpusName)
	m := extractModel(modelFileName, dialoguesCorpus,400000)
	
	for {
		conn,err := ln.Accept()
		if err!= nil{
			log.Println(err.Error())
		}
		go handleConnection(conn,questions,c,m)
	}
}	