package main

import (
	 "strings"
)
func extractVocabulary(document []string){

}

//returns number of all docs in all classes
func calNumAllDocs(classes map[int] []string) int{
	numOfDocs := 0
	for _,value := range classes{
		numOfDocs += len(value)
	}
	return numOfDocs
}

//returns number of classes 
func calNumOfClasses(classes map[int] []string) int{
	return len(classes)
}

//returns a slice of strings by spliting a given document into terms
func extractTerms(doc string) []string{
	return strings.Fields(doc)
}

//returns number of documents in each class
func classDocsNum(classDocs []string) int{
	return len(classDocs)
}

//returns a slice of number of documents in each class
func makeSliceOfNumDocsInClass(classes map[int] []string) []int{
	arrOfNumDocs := []int{}
	for key,docs := range classes{
		arrOfNumDocs[key] = classDocsNum(docs)
	}
	return arrOfNumDocs
}

//returns slice of probabilties of each class with the formula - 
//count of documents in class divided by all documents in all classes
func makeSlicePriorC(numOfAllDocs int, makeSliceOfNumDocsInClass []int) []float64{
	arr := []float64{}
	for i ,docsCount := range makeSliceOfNumDocsInClass {
		arr[i] = float64(docsCount)/ float64(numOfAllDocs)
	}
	return arr
}

//returns a slice which each index matches the term in vocabulary of the same index 
// and its value is the number of counts of the term in all documents 
func makeSliceTermCountInClass(numOfClasses int,vocabulary map [string] []int) []int{
	termCountArr := []int{}
	for _, value := range vocabulary{
		for i:=0;i <= numOfClasses; i++ {
			termCountArr[i] = value[i]
		}
	}
	return termCountArr
}

//the function receives as arguments a vocabulary and a slice with number of terms in each class
//and returns map with key string and value slice of floats
//the keys represent a term from vocabulary and the slice of floats has the values of the cond probability 
//inside the innermost cycle is the the making of the slice which we assign to the every term of the vocabulary
func makeSliceCondProb(vocabulary map [string] []int, sliceNumOfTermClass []int) map [string] []float64{
	sliceCondProb := map [string] []float64{}
	temp := [] float64{}
	sizeV := len(vocabulary)
	for term, _ := range vocabulary{
		for class,numOfTermInClass := range sliceNumOfTermClass{
			temp[class] = float64(term[class] + 1)/float64(numOfTermInClass + sizeV)
		}
		sliceCondProb[term] = temp
	}
	return sliceCondProb
}

//returns vocabulary of type map[string] []int where the key is a term 
//and the slice contains for each class the term frequency
func makeVocabulary(classes map[int] []string, numOfClasses int) map [string] []int{
	vocabulary := map [string] []int{}
	for class, docs := range classes{
		for _, doc := range docs{
			terms := extractTerms(doc)
			for _, term := range terms{
				if _, ok := vocabulary[term]; !ok {
					for i := 0; i <= numOfClasses; i++{
						vocabulary[term][i] = 0
					}
				}else {
					vocabulary[term][class] +=1
				}
			}
		}
	}
	return vocabulary
}
//function trains Multionomal Naive Bayes Classificator and returns vocabulary, priorC and cond probabilty
func TrainMultinomialNB(classes map[int] []string) (map [string] []int,map [string] []float64, []float64 ){
	numOfAllDocs := calNumAllDocs(classes)
	numOfClasses := calNumOfClasses(classes)
	vocabulary := makeVocabulary(classes,numOfClasses)
	makeSliceOfNumDocsInClass := makeSliceOfNumDocsInClass(classes)
	slicePriorC := makeSlicePriorC(numOfAllDocs, makeSliceOfNumDocsInClass)
	sliceTermCountClass := makeSliceTermCountInClass(numOfClasses, vocabulary)
	sliceCondProb := makeSliceCondProb(vocabulary,sliceTermCountClass)
	return vocabulary, sliceCondProb, slicePriorC
}
func main(){
	
}