package main

import (
	 "fmt"
	 "strings"
)
func extractVocabulary(document []string){

}

//func countDocument()

//Each class has its own file 
// class [0] [[document1],[document2]] vector<STRING> 
// class [0] [document2]
// vector<string> = class["key"]


// etc
//class[1][document3]
//var map with classes and documents map[string]int
// classes := map[int][]string{
// 	0 : {"lajdjadiadiw jbbj", "jajidad"},
// 	1 : {"laino", "zadnik"},
// }
//fmt.Println(calNumDocsOfClasses(classes))

func calNumAllDocs(classes map[int] []string) int{
	numOfDocs := 0
	for _,value := range classes{
		numOfDocs += len(value)
	}
	return numOfDocs
}

func calNumOfClasses(classes map[int] []string) int{
	return len(classes)
}

func extractTerms(doc string) []string{
	return strings.Fields(doc)
}

func classDocsNum(classDocs []string) int{
	return len(classDocs)
}

func makeArrayOfNumDocsInClass(classes map[int] []string) []int{
	arrOfNumDocs := []int{}
	for key,docs := range classes{
		arrOfNumDocs[key] = classDocsNum(docs)
	}
	return arrOfNumDocs
}

func makeArrPriorC(numOfAllDocs int, arrNumDocsInClass []int) []float64{
	arr := []float64{}
	for i ,docsCount := range arrNumDocsInClass {
		arr[i] = float64(docsCount)/ float64(numOfAllDocs)
	}
	return arr
}

func makeArrTermCountInClass(numOfClasses int,vocabulary map [string] []int) []int{
	termCountArr := []int{}
	for _, value := range vocabulary{
		for i:=0;i <= numOfClasses; i++ {
			termCountArr[i] = value[i]
		}
	}
	return termCountArr
}
func TrainMultinomialNB(classes map[int] []string, document []string, vocabulary map [string] []int){
	numOfAllDocs := calNumAllDocs(classes)
	numOfClasses := calNumOfClasses(classes)
	for class, docs := range classes{
		for _, doc := range docs{
			terms := extractTerms(doc)
			for _, term := range terms{
				if val, ok := vocabulary[term]; !ok {
					for i := 0; i <= numOfClasses; i++{
						vocabulary[term][i] = 0
					}
				}else {
					vocabulary[term][class] +=1
				}
			}
		}
	}
	arrNumDocsInClass := makeArrayOfNumDocsInClass(classes)
	arrPriorC := makeArrPriorC(numOfAllDocs, arrNumDocsInClass)
	arrTermCountClass := makeArrTermCountInClass(numOfClasses, vocabulary)
	ararCondProb := 
}
func main(){
	
}