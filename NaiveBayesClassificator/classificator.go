package classificator 

import (
	 "strings"
	 "math"
	 "fmt"
)
func convertText(document []string){}

//function trains Multionomal Naive Bayes Classificator and returns vocabulary, priorC and cond probabilty
func TrainMultinomialNB(classes map[int] []string) (map [string] []int,map [string] []float64, []float64 ){
   numOfAllDocs := calNumAllDocs(classes)
   numOfClasses := calNumOfClasses(classes)
   vocabulary := makeVocabulary(classes,numOfClasses)
   makeSliceOfNumDocsInClass := makeSliceOfNumDocsInClass(classes, numOfClasses)
   slicePriorC := makeSlicePriorC(numOfAllDocs, makeSliceOfNumDocsInClass)
   sliceTermCountClass := makeSliceTermCountInClass(numOfClasses, vocabulary)
   sliceCondProb := makeSliceCondProb(vocabulary,sliceTermCountClass)
   return vocabulary, sliceCondProb, slicePriorC
}

//returns the class corresponding to the text given by using formula 	
func applyMultinomialNB(condProb map [string] []float64, priorC []float64, text string ) int{
   
   terms  := extractTerms(text)
   var classificatedAs int 
   var maxScore float64
   for classInd, value := range priorC{
	   score := math.Log(value)
	   for _, term := range terms{
		   if condProbTerm, ok := condProb[term]; ok {
			   score += math.Log(condProbTerm[classInd])
		   }
	   }
	   if classInd == 0 || score > maxScore{
		   maxScore = score
		   classificatedAs = classInd
	   }
   }

   return classificatedAs
}

//testing classifier accuracy
func testClassifier(testClassCorpus map[int] []string, vocabulary map [string] []int, sliceCondProb map [string] []float64, slicePriorC []float64 ){
   numOfClasses := len(testClassCorpus)
   numDocsByClass := numberDocByClass(testClassCorpus)
   confusionMatrix := makeConfMatrix(testClassCorpus, numOfClasses,vocabulary, sliceCondProb,slicePriorC)
   precision, recall, fScore := calcPRF(confusionMatrix, numOfClasses, numDocsByClass)
   precisionOverall, recallOverall,fScoreOverall := calcOverall(precision, recall, fScore, numOfClasses, numDocsByClass)
   fmt.Println("Прецизност: ", precision)
   fmt.Println("Обхват: ", recall)
   fmt.Println("F-score: ", fScore)
   fmt.Println("Обща прецизност: ", precisionOverall)
   fmt.Println("Общ Обхват: ", recallOverall)
   fmt.Println("Обща F-score: ", fScoreOverall)
}

//returns vocabulary of type map[string] []int where the key is a term 
//and the slice contains for each class the term frequency
func makeVocabulary(classes map[int] []string, numOfClasses int) map [string] []int{
   vocabulary := make(map [string] []int)

   for class, docs := range classes{
	   for _, doc := range docs{
		   terms := extractTerms(doc)
		   for _, term := range terms{
			    if  _, ok := vocabulary[term]; !ok {
				   for i := 0; i < numOfClasses; i++{
					vocabulary[term] = append(vocabulary[term], 0)
				   }
				}
			     vocabulary[term][class] +=1
		   }
	   }
   }

   return vocabulary
}


//returns a slice of number of documents in each class
func makeSliceOfNumDocsInClass(classes map[int] []string, numOfClasses int) []int{
	arrOfNumDocs := []int{}
	for classInd := 0; classInd < len(classes);classInd++ {
		arrOfNumDocs = append(arrOfNumDocs, classDocsNum(classes[classInd]))
	}
	return arrOfNumDocs
 }
 
 //returns slice of probabilties of each class with the formula - 
 //count of documents in class divided by all documents in all classes
 func makeSlicePriorC(numOfAllDocs int, sliceOfNumDocsInClass []int) []float64{
	arr := []float64{}
 
	for _ ,docsCount := range sliceOfNumDocsInClass {
		arr = append(arr, float64(docsCount)/ float64(numOfAllDocs))
	}
	return arr
 }
 

 //the function receives as arguments a vocabulary and a slice with number of terms in each class
 //and returns map with key string and value slice of floats
 //the keys represent a term from vocabulary and the slice of floats has the values of the cond probability 
 //inside the innermost cycle is the the making of the slice which we assign to the every term of the vocabulary
 func makeSliceCondProb(vocabulary map [string] []int, sliceNumOfTermClass []int) map [string] []float64{
	sliceCondProb := make(map [string] []float64)
   
   // i := 0
	sizeV := len(vocabulary)
	for term, value := range vocabulary{
		  temp := [] float64{}
		for class,numOfTermsInClass := range sliceNumOfTermClass{
			temp = append(temp, (float64(value[class] + 1)/float64(numOfTermsInClass + sizeV)))		   
 
		}
		sliceCondProb[term] = temp
	 }
	return sliceCondProb
 }

//returns the Precision, F-score and the recall of the classificator for each document of a test set of documents
func calcPRF(confusionMatrix [][]int, numOfClasses int, numAllDocsByClass []int) ([]float64, []float64,[] float64){ 
   
	precision := []float64{}
	recall := []float64{}
	fScore := []float64{}
	for classInd := 0; classInd < numOfClasses; classInd++{
		countDocsExtracted := sumMatrixValues(confusionMatrix, classInd, numOfClasses)
		if confusionMatrix[classInd][classInd] == 0{
				 precision = append(precision, 0.0)
				 recall = append(recall, 0.0)
				 fScore = append(fScore, 0.0)
		}else{
		precision = append(precision, (float64(confusionMatrix[classInd][classInd]) / float64(countDocsExtracted)))
		recall    = append(recall,    (float64(confusionMatrix[classInd][classInd]) / float64(numAllDocsByClass[classInd])))
		fScore    = append(fScore,    ((2.0 * precision[classInd] * recall[classInd]) / (precision[classInd] + recall[classInd])))
		}
	}
	
	return precision, recall, fScore
 }
 
 //returns overall Precision, Recall and F-score for every class 
 func calcOverall(precision []float64,recall  []float64, fScore []float64, 
				 numOfClasses int , countDocsClass[] int) (float64, float64, float64){
 
	var precisionOverall float64 = 0.0
	var recallOverall float64 = 0.0
	var fScoreOverall float64 = 0.0
	allDocs := sum(countDocsClass)
 
	for classInd := 0; classInd < numOfClasses; classInd++{
		precisionOverall += (float64(countDocsClass[classInd]) * precision[classInd])/ float64(allDocs)
		recallOverall    += (float64(countDocsClass[classInd]) * recall[classInd]) / float64(allDocs)
	}
 
	fScoreOverall += (2 * precisionOverall * recallOverall) / (precisionOverall + recallOverall)
	return precisionOverall, recallOverall, fScoreOverall
 }
 

//returns the confusion matrix which shows classification
// accuracy by showing the correct and incorrect predictions on each class.
func makeConfMatrix(testClassCorpus map[int] []string,
				   numOfClasses int,
				   vocabulary map [string] []int,
				   sliceCondProb map [string] []float64, 
				   slicePriorC []float64 ) [][]int{
	confusionMatrix := make([][]int, numOfClasses)
	for i := range confusionMatrix {
		confusionMatrix[i] = make([]int,numOfClasses)
	}
   for classInd := 0; classInd <= numOfClasses;classInd++{
	   for _, doc := range testClassCorpus[classInd]{
		   classified_as_doc := applyMultinomialNB(sliceCondProb, slicePriorC,doc)
		   confusionMatrix[classInd][classified_as_doc] += 1
		   }
	   }
	fmt.Println(confusionMatrix)
   return confusionMatrix
}

 //returns a slice which each index matches the term in vocabulary of the same index 
 // and its value is the number of counts of the term in all documents 
 func makeSliceTermCountInClass(numOfClasses int,vocabulary map [string] []int) []int{
	termCountArr := []int{}
	for _, value := range vocabulary{
		for classInd:=0;classInd < numOfClasses; classInd++ {
			if len(termCountArr) >= numOfClasses{
				 termCountArr[classInd] += value[classInd] 
			}else{
				termCountArr = append(termCountArr,value[classInd])
			}
		}
	}
 
	return termCountArr
 }
 
 
//returns sum of the elements of the matrix
func sumMatrixValues(confussionMatrix [][]int, classInd int, numOfClasses int) int{
   var sum int
	   for i := 0; i < numOfClasses; i++ { 
		   sum += confussionMatrix[i][classInd]
	   }
   return sum
}

//sum the values of slice of ints 
func sum(sl []int) int {  
   result := 0  
   for _, numb := range sl {  
	result += numb  
   }  
   return result  
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
func numberDocByClass(testClassCorpus map[int] []string) []int{
   docsCountByClass := []int{}
   for _, docs := range testClassCorpus{
	   docsCountByClass = append(docsCountByClass, len(docs))
   }
   return docsCountByClass
}


