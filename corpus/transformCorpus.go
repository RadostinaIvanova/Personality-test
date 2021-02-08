package corpus

import(
	"encoding/csv"
	"fmt"
	"strings"
	"os"
	"regexp"
	"log"
)

func MakeClassesFromFile(filename string) (map[int][] string,map [int][]string){
	all := readCsvFile(filename)	
	classes := divideIntoClasses(all)
	trainSet, testSet := divideIntoTrainTestSets(classes)
	return trainSet,testSet
}

func readCsvFile(filename string) [][]string{
	csvFile, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
	}
	defer csvFile.Close()
    csvLines := csv.NewReader(csvFile)
	
	check := false
	var all [][]string
	var classAndDoc[] string

	for {
		record, err1 := csvLines.Read()
		if(err1 != nil){
			break
		}
		if check == true{
			doc := transformToLowerAndEraseDots(record[1])
			classAndDoc := append(classAndDoc, record[0])
			classAndDoc = append(classAndDoc, doc)
			all = append(all, classAndDoc)
			}
		check = true
	 }
	 return all	
}

func divideIntoTrainTestSets(classes map [int] []string)(map [int] []string,map [int] []string){
	testSet := make(map[int] []string)
	trainSet := make(map[int] []string)
	max := 250
		for class, docs:= range classes{
			testCount := 10
			testSet[class] = docs[ :testCount]
			trainSet[class] = docs[int(testCount):max] //700
			}
	return trainSet, testSet
}	

func divideIntoClasses(records [][]string)map[int] []string{
	classes :=  make(map[int] []string)
	for _, record := range records{
		classTypeStr := record[0]
		classType := encodeClassToInt(classTypeStr)
		classes[classType] = append(classes[classType], record[len(record) - 1])
	}
	return classes
}

func encodeClassToInt(classType string) int{
	switch classType{
		case "INTJ": return 1
		case "ISTJ" : return 2
		case "ISFJ" : return 2
		case "INFJ" : return 0
		case "INTP" : return 1
		case "INFP" : return 0
		case "ISTP" : return 3
		case "ISFP" : return 3
		case "ENTJ" : return 1
		case "ENFJ" : return 0
		case "ENFP" : return 0
		case "ENTP" : return 1	
		case "ESFJ" : return 2
		case "ESTP" : return 3
		case "ESFP" : return 3	
		case "ESTJ" : return 2
		}
	return 0
}


func transformToLowerAndEraseDots(str string) string{
	reg, err := regexp.Compile("[^a-zA-Z0-9\\s]+")
    if err != nil {
        log.Println(err.Error())
    }
	newValue := reg.ReplaceAllString(str, "")
	newValue = strings.ToLower(newValue)
	newValue = strings.Replace(newValue, "and", "", -1)
	newValue = strings.Replace(newValue, "or", "", -1)
 return newValue
}