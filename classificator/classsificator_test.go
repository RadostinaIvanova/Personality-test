  package classificator
  
  import (
	"testing"
	"strconv"
  )

func equal(a []int, b []int) bool {
	if len(a) != len(b) {
			return false
	}
	for i, v := range a {
			if v != b[i] {
					return false
			}
		}
	return true
}

func TestCalNumAllDocs(t *testing.T){
	c := NBclassificator{}
	m := map [int] []string{
		1 : {"Hello", "Goodbye", "GoodAfternoon"},
		2 : {"Something", "Anything", "Other"},
	}
	numAllDocs := 6
	result := c.calNumAllDocs(m)
	if  result != numAllDocs{
		t.Error("calNumAllDocs returned" +  strconv.Itoa(result) + ", we expected " + strconv.Itoa(numAllDocs))
	}
}


func TestCalNumOfClassess(t *testing.T){
	c := NBclassificator{}
	m := map [int] []string{
		1 : {"Hello", "Goodbye", "GoodAfternoon"},
		2 : {"Something", "Anything", "Other"},
	}
	numOfClasses := 2
	result := c.calNumOfClasses(m)
	if  result != numOfClasses{
		t.Error("calNumAllDocs returned" +  strconv.Itoa(result) + ", we expected " + strconv.Itoa(numOfClasses))
	}
}

func TestExtractTerms(t *testing.T){
	c := NBclassificator{}
	str := "Hello my friend"
	splitted := []string{"Hello", "my", "friend"}
	result := c.extractTerms(str) 
	for i := 0; i < len(splitted);i++{
		if splitted[i] != result[i] {
			t.Error("extractTerms returned" + result[i] + ", we expected " + result[i])
		}
	}
 }

 func TestClassDocsNum(t *testing.T){
	c := NBclassificator{}
	classDocs := []string{"Hello", "my", "friend"}
	expected := 3
	result := c.classDocsNum(classDocs)
	if  result != expected{
		t.Error("classDocsNum returned" +  strconv.Itoa(result) + ", we expected " + strconv.Itoa(expected))
	}
 }

 func TestMakeSliceOfNumDocsInClass(t *testing.T){
	c := NBclassificator{}
	m := map [int] []string{
		0 : {"Hello", "Goodbye", "GoodAfternoon"},
		1 : {"Something", "Anything"},
	}
	expected := []int{3,2}
	numOfClasses := 2
	result := c.makeSliceOfNumDocsInClass(m, numOfClasses)
	for i := 0;i < len(expected); i++{
		if  result[i] != expected[i]{
			t.Error("makeSliceOfNumDocsInClass returned" +  strconv.Itoa(result[i]) + ", we expected " + strconv.Itoa(expected[i]))
		}
	}
 }

 
 func TestMakeSliceTermCountInClass(t *testing.T){
	v := map [string] []int{
		"Hello" : {3, 4, 5},
		"Goodbye" : {1, 2, 3},
	}
	c := NBclassificator{v,map [string] []float64{},  []float64{}}
	expected := []int{4,6,8}
	numOfClasses := 3
	result := c.makeSliceTermCountInClass(numOfClasses)
	for i := 0;i < len(expected); i++{
		if  result[i] != expected[i]{
			t.Error("makeSliceTermCountInClass returned" +  strconv.Itoa(result[i]) + ", we expected " + strconv.Itoa(expected[i]))
		}
	}
 }

 func TestMakeVocabulary(t *testing.T){
	c := NBclassificator{}
	m := map [int] []string{
		0 : {"Hello my friend", "Goodbye my friend", "Anything fellow"},
		1 : {"Hello something went wrong", "Anything my friend"},
	}
	expected := map [string] []int{ 
		"Hello" : {1, 1},
		"my"    : {2, 1},
		"friend" : {2, 1},
		"Goodbye" : {1, 0},
		"Anything":{1,1},
		"fellow" :{1,0},
		"something" : {0,1},
		"went" : {0,1},
		"wrong" :{0,1},
	}
	numOfClasses := 2
	result := c.makeVocabulary(m,numOfClasses)
	for key, value := range expected{
		if val, ok := result[key]; !ok || !equal(val,value){
			//t.Error("Expected with " + key + " value: " + strconv.Itoa(value[0]) + " and " + "value " + strconv.Itoa(value[1]) + ". Received " + key + " value: " + strconv.Itoa(val[0]) + " and " + "value " + strconv.Itoa(val[1]))
			t.Error("makeVocabulary wrong vocabulary")
		}
	}
 }


func TestSumMatrixValues(t *testing.T){
	c := NBclassificator{}
	confussionMatrix := [][]int{{1,2},
								{4,5}}
	classInd :=  0
	numOfClasses := 2 
	expected := 5
	result := c.sumMatrixValues(confussionMatrix, classInd,numOfClasses)
	if  result != expected{
		t.Error("sumMatrixValues returned" +  strconv.Itoa(result) + ", we expected " + strconv.Itoa(expected))
	}
 }

