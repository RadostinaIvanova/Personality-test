package corpus

import(
	"encoding/csv"
	"fmt"
	"strings"
	"os"
)
func transform(filename string) (map [int] []string, map[int] []string){
	all := extractAndChange(filename)
	classes := putInClass(all)
	testSet, trainSet := makeSets(classes)
	return testSet, trainSet
}

func extractAndChange(fileName string) [][]string{
	csvFile, err := os.Open("D:\\FMI\\golang_workspace\\src\\mbt\\mbt.csv")
	defer csvFile.Close()
	if err != nil {
		fmt.Println(err)
	}
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
			doc := transformToLowerAndEraseSymbols(record[1])
			classAndDoc := append(classAndDoc, record[0])
			classAndDoc = append(classAndDoc, doc)
			all = append(all, classAndDoc)
			}
		check = true
	 }
	return all
}
func makeSets(classes map[int] []string) (map[int] []string, map[int] []string){
	testSet := make(map[int] []string)
	trainSet := make(map[int] []string)
	for i := 200; i < 256;i+=100{
		for class, docs:= range classes{
			testCount := 50
			testSet[class] = docs[ : testCount]
			trainSet[class] = docs[int(testCount):i] //700
		}
	}
	return testSet, trainSet
}
func putInClass(records [][]string)map[int] []string{
	classes :=  make(map[int] []string)
	for _, record := range records{
		classTypeStr := record[0]
		classType := changeClassToInt(classTypeStr)
		classes[classType] = append(classes[classType], record[len(record) - 1])
		
	}
	return classes
}
func transformToLowerAndEraseSymbols(str string) string{

		newValue := strings.ToLower(str)
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


func changeClassToInt(classType string) int{
	switch classType{
	case "INTJ": return 2
	case "INTP" : return 0
	case "ENTJ" : return 1
	case "INFJ" : return 2
	case "INFP" : return 0
	case "ENFJ" : return 1
	case "ENFP" : return 0 
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