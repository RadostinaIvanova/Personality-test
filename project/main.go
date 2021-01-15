package main

import (
	//"encoding/csv"
	"fmt"
	"strings"
	"strconv"
	"io/ioutil"
	// "path/filepath"
    "log"
)

import 
	"math"


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
   fmt.Println(len(priorC))
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

//the function receives as arguments a vocabulary and a slice with number of terms in each class
//and returns map with key string and value slice of floats
//the keys represent a term from vocabulary and the slice of floats has the values of the cond probability 
//inside the innermost cycle is the the making of the slice which we assign to the every term of the vocabulary
func makeSliceCondProb(vocabulary map [string] []int, sliceNumOfTermClass []int) map [string] []float64{
   sliceCondProb := make(map [string] []float64)
   temp := [] float64{}
  // i := 0
   sizeV := len(vocabulary)
   for term, value := range vocabulary{
	   for class,numOfTermsInClass := range sliceNumOfTermClass{
		   temp = append(temp, (float64(value[class] + 1)/float64(numOfTermsInClass + sizeV)))		   
	   }
	   sliceCondProb[term] = temp
	}
   return sliceCondProb
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

//returns sum of the elements of the matrix
func sumMatrixValues(confusionMatrix [][]int) int{
   var sum int
   for _, col := range confusionMatrix{
	   for _, value := range col{
		   sum+= value
	   }
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

func numberDocByClass(testClassCorpus map[int] []string) []int{
   docsCountByClass := []int{}
   for _, docs := range testClassCorpus{
	   docsCountByClass = append(docsCountByClass, len(docs))
   }
   return docsCountByClass
}

//returns the Precision, F-score and the recall of the classificator for each document of a test set of documents
func calcPRF(confusionMatrix [][]int, numOfClasses int, numAllDocsByClass []int) ([]float64, []float64,[] float64){ 
   
   countDocsExtracted := sumMatrixValues(confusionMatrix)
   precision := []float64{}
   recall := []float64{}
   fScore := []float64{}

   for classInd := 0; classInd < numOfClasses; classInd++{
	   if confusionMatrix[classInd][classInd] == 0{
			precision = append(precision, 0.0)
			recall = append(precision, 0.0)
			fScore = append(precision, 0.0)
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
	   precisionOverall += (float64(precision[classInd] * precision[classInd])/ float64(allDocs))
	   recallOverall    += (float64(countDocsClass[classInd]) * recall[classInd]) / float64(allDocs)
   }

   fScoreOverall += (2 * precisionOverall * recallOverall) / (precisionOverall + recallOverall)
   return precisionOverall, recallOverall, fScoreOverall
}


func notCorruptedRecord(record []string) bool{
	check := false
		for _, s := range record {
			isLetter := func(c rune) bool { return (c > 'a' && c < 'z')}
			b := strings.IndexFunc(s, isLetter) != -1
			if b {
				check = true
			}
		}
		return check
}

const EX_IND uint = 6
const N_IND  uint = 7
const O_IND  uint = 8
const A_IND  uint = 9

func convertFloatToStrRecord(record []string) []string{
	if len(record) > 9{
		if j,_ := strconv.ParseFloat(record[EX_IND], 64);j <= 2.5{
			record[EX_IND] = "i"
		}else{
			record[EX_IND] = "e"
		}
		if j,_ := strconv.ParseFloat(record[N_IND], 64); j <= 2.5{
			record[N_IND] = "s"
		}else{
			record[N_IND] = "n"
		}
		if j,_ := strconv.ParseFloat(record[O_IND], 64); j <= 2.5{
			record[O_IND] = "t"
		}else{
			record[O_IND] = "f"
		}
		if  j,_ := strconv.ParseFloat(record[A_IND], 64); j<= 2.5{
			record[A_IND] = "j"
		}else{
			record[A_IND] = "p"
		}
	}
	classType := record[EX_IND] + record[N_IND] + record[O_IND] + record[A_IND]
	record = record[0:len(record) - 4]
	record = append(record, classType)
	return record
}

// }
func change(classType string) int{
	switch classType{
	case "intj": return 0
	case "intp" : return 1
	case "entp" : return 2
	case "infj" : return 3
	case "infp" : return 4
	case "enfj" : return 5
	case "enfp" : return 6
	case "esfj" : return 7
	case "istp" : return 8
	case "isfp" : return 9
	case "estp" : return 10
	case "esfp" : return 11
	}
	return 0
}
func putInClass(records [][]string)map[int] []string{

	classes :=  make(map[int] []string)
	for _, record := range records{
		classTypeStr := record[len(record) - 1]
		classType := change(classTypeStr)
		record = record[0 : len(record) - 1]
		newRec := strings.Join(record[:]," ")
		classes[classType] = append(classes[classType], newRec)
		
	}
	return classes
}

//func main() {
// 	csvFile, err := os.Open("D:\\FMI\\golang_workspace\\src\\golang_course\\project\\data1.csv")
// 	defer csvFile.Close()

// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	fmt.Println("Successfully Opened CSV file")
//     csvLines := csv.NewReader(csvFile)
// 	check := false
// 	var str1 [][]string
// 	for {
		
// 		record, err1 := csvLines.Read()
// 		if(err1 == io.EOF){
// 			break
// 		}
		
// 		if notCorruptedRecord(record){
// 			if check == true{
// 			record = record[2:len(record)-1]
// 			record = convertFloatToStrRecord(record) 
// 			str1 = append(str1, record)
// 			}
// 		}
// 		check = true
// 	 }
// 	//fmt.Println(str1)
// 	classes := putInClass(str1)
// 	testSet := make(map[int] []string)
// 	trainSet := make(map[int] []string)
// 	for class, docs:= range classes{
// 		fmt.Println(len(docs))
// 		testCount := float64(len(docs)) * 0.1
// 		testSet[class] = docs[:int(testCount)]
// 		trainSet[class] = docs[int(testCount):]
// 	}
// 	vocabulary, condProb, prior := TrainMultinomialNB(trainSet)
// 	fmt.Println(applyMultinomialNB(condProb, prior,"i didn't know humanity back love and family"))
// 	testClassifier(testSet, vocabulary,condProb, prior)
// }

func convert(record string) string{
	return strings.Replace(record, "***", "", -1)
}
func main() {
	classes := make(map [int] []string)
	filesAll, err := ioutil.ReadDir("src\\golang_course\\project\\all")
	if err != nil {
        log.Fatal(err)
	}
	for i := 0; i < len(filesAll) - 1; i++{
		for _, classF := range filesAll {
			classNames := "src\\golang_course\\project\\all\\" + classF.Name()
			files, err := ioutil.ReadDir(classNames)
			if err != nil {
				log.Fatal(err)
			}
			records := []string{}
			for _, f := range files {
				fileName := "src\\golang_course\\project\\all\\C-Culture\\" + f.Name()
				record,_ := ioutil.ReadFile(fileName)
				recordStr := convert(string(record))
			//	fmt.Println(recordStr)
				records = append(records, recordStr)
			}
			classes[i] = records
		}
	}
	testSet := make(map[int] []string)
	trainSet := make(map[int] []string)
	for class, docs:= range classes{
		testCount := float64(len(docs)) * 0.1
		testSet[class] = docs[:int(testCount)]
		trainSet[class] = docs[int(testCount):]
	}
	_, condProb, prior := TrainMultinomialNB(trainSet)
	text1 := "Ние виждаме една серия от управленски провали. Укрепено ли е правителството - не то не съществува. равителството е несъществуващото - в този тежък момент, когато достигаме близо 1000 заразени на ден, ако говорим за здравния проблем "
	fmt.Println(applyMultinomialNB(condProb, prior, text1 ))
	//testClassifier(testSet, vocabulary,condProb, prior)
}
	
