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
	max := 200
		for class, docs:= range classes{
			testCount := 50
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
	case "INTJ": return 2
	case "INTP" : return 0
	case "ENTJ" : return 1
	case "INFJ" : return 2
	case "INFP" : return 0
	case "ENFJ" : return 1
	case "ENFP" : return 1 
	case "ESFJ" : return 3
	case "ISTP" : return 0
	case "ISFP" : return 0
	case "ESTP" : return 3
	case "ESFP" : return 3
	case "ENTP" : return 1
	case "ISTJ" : return 2
	case "ISFJ" : return 2
	case "ESTJ" : return 3
	}
	return 0
}


func transformToLowerAndEraseDots(str string) string{
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
    if err != nil {
        log.Println(err.Error())
    } 
	newValue := reg.ReplaceAllString(str, "")
	newValue = strings.ToLower(str)
	newValue = strings.Replace(newValue, ".", "", -1)
	newValue = strings.Replace(newValue, "|||", "", -1)
	newValue = strings.Replace(newValue, "[", "", -1)
	newValue = strings.Replace(newValue, "]", "", -1)
	newValue = strings.Replace(newValue, "!", "", -1)
	newValue = strings.Replace(newValue, "?", "", -1)
	newValue = strings.Replace(newValue, ",", "", -1)
	newValue = strings.Replace(newValue, "and", "", -1)
	newValue = strings.Replace(newValue, "or", "", -1)
 return newValue
}