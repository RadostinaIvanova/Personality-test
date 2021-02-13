package classificator 

import (
	 "strings"
	 "math"
	 "fmt"
)

type NBclassificator struct{
	Vocabulary map[string] []int
	CondProb map [string] []float64
	PriorC []float64
}

//function trains Multionomal Naive Bayes Classificator and returns vocabulary, priorC and cond probabilty
func (c *NBclassificator) TrainMultinomialNB(classes map[int] []string){
	numOfAllDocs := c.calNumAllDocs(classes)
	numOfClasses := c.calNumOfClasses(classes)
	c.Vocabulary = c.makeVocabulary(classes,numOfClasses)
	makeSliceOfNumDocsInClass := c.makeSliceOfNumDocsInClass(classes, numOfClasses)
	c.PriorC = c.makeSlicePriorC(numOfAllDocs, makeSliceOfNumDocsInClass)
	sliceTermCountClass := c.makeSliceTermCountInClass(numOfClasses)
	c.CondProb = c.makeSliceCondProb(sliceTermCountClass)
 }
 
 //returns the class corresponding to the text given by using formula 	
 func (c *NBclassificator) ApplyMultinomialNB(text string ) int{
	terms  := c.extractTerms(text)
	var classificatedAs int 
	var maxScore float64

	for classInd, value := range c.PriorC{
		score := math.Log(value)
		for _, term := range terms{
			if condProbTerm, ok := c.CondProb[term]; ok {
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
 func (c *NBclassificator)TestClassifier(testClassCorpus map[int] []string){
	numOfClasses := len(testClassCorpus)
	numDocsByClass := c.numberDocByClass(testClassCorpus)
	confusionMatrix := c.makeConfMatrix(testClassCorpus, numOfClasses)
	precision, recall, fScore := c.calcPRF(confusionMatrix, numOfClasses, numDocsByClass)
	precisionOverall, recallOverall,fScoreOverall := c.calcOverall(precision, recall, fScore, numOfClasses, numDocsByClass)
	fmt.Println("Прецизност: ", precision)
	fmt.Println("Обхват: ", recall)
	fmt.Println("F-score: ", fScore)
	fmt.Println("Обща прецизност: ", precisionOverall)
	fmt.Println("Общ Обхват: ", recallOverall)
	fmt.Println("Обща F-score: ", fScoreOverall)
 }

 //returns number of all docs in all classes
 func (c *NBclassificator) calNumAllDocs(classes map[int] []string) int{
	numOfDocs := 0
	for _,value := range classes{
		numOfDocs += len(value)
	}
	return numOfDocs
 }
 
 //returns number of classes 
 func (c *NBclassificator) calNumOfClasses(classes map[int] []string) int{
	return len(classes)
 }
 
 //returns a slice of strings by spliting a given document into terms
 func (c *NBclassificator) extractTerms(doc string) []string{
	return strings.Fields(doc)
 }
 
 //returns number of documents in each class
 func (c *NBclassificator) classDocsNum(classDocs []string) int{
	return len(classDocs)
 }
 
 //returns a slice of number of documents in each class
 func (c *NBclassificator) makeSliceOfNumDocsInClass(classes map[int] []string, numOfClasses int) []int{
	arrOfNumDocs := []int{}
	for classInd := 0; classInd < len(classes);classInd++ {
		arrOfNumDocs = append(arrOfNumDocs, c.classDocsNum(classes[classInd]))
	}
	return arrOfNumDocs
 }
 
 //returns slice of probabilties of each class with the formula - 
 //count of documents in class divided by all documents in all classes
 func (c *NBclassificator) makeSlicePriorC(numOfAllDocs int, sliceOfNumDocsInClass []int) []float64{
	arr := []float64{}
 
	for _ ,docsCount := range sliceOfNumDocsInClass {
		arr = append(arr, float64(docsCount)/ float64(numOfAllDocs))
	}
	return arr
 }
 
 //returns a slice which each index matches the term in vocabulary of the same index 
 // and its value is the number of counts of the term in all documents 
 func (c *NBclassificator) makeSliceTermCountInClass(numOfClasses int) []int{
	termCountArr := []int{}
	for _, value := range c.Vocabulary{
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
 
 //the function receives as arguments a vocabulary and a slice with number of terms in each class
 //and returns map with key string and value slice of floats
 //the keys represent a term from vocabulary and the slice of floats has the values of the cond probability 
 //inside the innermost cycle is the the making of the slice which we assign to the every term of the vocabulary
 func (c *NBclassificator) makeSliceCondProb(sliceNumOfTermClass []int) map [string] []float64{
	sliceCondProb := make(map [string] []float64)
   
   // i := 0
	sizeV := len(c.Vocabulary)
	for term, value := range c.Vocabulary{
		  temp := [] float64{}
		for class,numOfTermsInClass := range sliceNumOfTermClass{
			temp = append(temp, (float64(value[class] + 1)/float64(numOfTermsInClass + sizeV)))		   
 
		}
		sliceCondProb[term] = temp
	 }
	return sliceCondProb
 }
 
 //returns vocabulary of type map[string] []int where the key is a term 
 //and the slice contains for each class the term frequency
 func (c *NBclassificator) makeVocabulary(classes map[int] []string, numOfClasses int) map [string] []int{
	vocabulary := make(map [string] []int)
 
	for class, docs := range classes{
		for _, doc := range docs{
			terms := c.extractTerms(doc)
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
 
 
 //returns the confusion matrix which shows classification
 // accuracy by showing the correct and incorrect predictions on each class.
 func (c *NBclassificator) makeConfMatrix(testClassCorpus map[int] []string,
					numOfClasses int) [][]int{
	 confusionMatrix := make([][]int, numOfClasses)
	 for i := range confusionMatrix {
		 confusionMatrix[i] = make([]int,numOfClasses)
	 }
	for classInd := 0; classInd <= numOfClasses;classInd++{
		for _, doc := range testClassCorpus[classInd]{
			classified_as_doc := c.ApplyMultinomialNB(doc)
			confusionMatrix[classInd][classified_as_doc] += 1
			}
		}
	 for _, value := range confusionMatrix{
		 fmt.Println(value)
	 }
	 //fmt.Println(confusionMatrix)
	return confusionMatrix
 }
 
 //returns sum of the elements of the matrix by given column(classInd)
 func (c *NBclassificator) sumMatrixValues(confussionMatrix [][]int, classInd int, numOfClasses int) int{
	var sum int
		for i := 0; i < numOfClasses; i++ { 
			sum += confussionMatrix[i][classInd]
		}
	return sum
 }
 
 //sum the values of slice of ints 
 func (c *NBclassificator) sum(sl []int) int {  
	result := 0  
	for _, numb := range sl {  
	 result += numb  
	}  
	return result  
 }  
 
 func (c *NBclassificator) numberDocByClass(testClassCorpus map[int] []string) []int{
	docsCountByClass := []int{}
	for _, docs := range testClassCorpus{
		docsCountByClass = append(docsCountByClass, len(docs))
	}
	return docsCountByClass
 }
 
 //returns the Precision, F-score and the recall of the classificator for each document of a test set of documents
 func (c *NBclassificator) calcPRF(confusionMatrix [][]int, numOfClasses int, numAllDocsByClass []int) ([]float64, []float64,[] float64){ 

	precision := []float64{}
	recall := []float64{}
	fScore := []float64{}
	for classInd := 0; classInd < numOfClasses; classInd++{
		countDocsExtracted := c.sumMatrixValues(confusionMatrix, classInd, numOfClasses)
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
 func (c *NBclassificator) calcOverall(precision []float64,recall  []float64, fScore []float64, 
				 numOfClasses int , countDocsClass[] int) (float64, float64, float64){
 
	var precisionOverall float64 = 0.0
	var recallOverall float64 = 0.0
	var fScoreOverall float64 = 0.0
	allDocs := c.sum(countDocsClass)
 
	for classInd := 0; classInd < numOfClasses; classInd++{
		precisionOverall += (float64(countDocsClass[classInd]) * precision[classInd])/ float64(allDocs)
		recallOverall    += (float64(countDocsClass[classInd]) * recall[classInd]) / float64(allDocs)
	}
 
	fScoreOverall += (2 * precisionOverall * recallOverall) / (precisionOverall + recallOverall)
	return precisionOverall, recallOverall, fScoreOverall
 }