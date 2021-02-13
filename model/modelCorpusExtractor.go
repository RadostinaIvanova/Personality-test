package main 

import (
    "os"
    "fmt"
    "bufio"
    "strings"
    "regexp"
    "log"
)
func FullSentCorpus(sentences []string) [][]string {
	sentences = transform(sentences)
    result := [][]string{}
    for _,value := range sentences{
        sentSplit := []string{}
        sentSplit =  append(sentSplit, startToken)
        value = strings.Trim(value, " ")
		value = strings.ToLower(value) 
        sentSplit =  append(sentSplit, strings.Split(value, " ")...)
        sentSplit = delete_empty(sentSplit)
        sentSplit = append(sentSplit, endToken)
        result = append(result,sentSplit)
    }
    return result
}

func extract(filename string) []string{
    f, err := os.Open(filename)
    if err != nil{
        fmt.Println("could't open file")
    }
    defer f.Close()

    reg, err := regexp.Compile("[^a-zA-Z0-9_\\s]+")
    if err != nil {
        log.Println(err.Error())
    }
    limit := 5
    i := 0
    scanner := bufio.NewScanner(f)
    corpus := []string{}
    for ;i < limit; {
	    for scanner.Scan(){
            text :=  reg.ReplaceAllString(scanner.Text(), "")
	        corpus = append(corpus, text)
        } 
        i++
        }
    return corpus
}

func transform(text []string) []string{
    sentences := []string{}
    for _,paragraph := range text{
         sentences = append(sentences,strings.Split(paragraph, "__eou__")...)
    }
   return sentences
}

func delete_empty (s []string) []string {
    var r []string
    for _, str := range s {
        if str != "" {
            r = append(r, str)
        }
    }
    return r
}

func divideIntoTrainAndTest(percent float64, fullSentCorpus [][]string)([][]string,[][]string){
	portion := int(percent*float64(len(fullSentCorpus)))
	test := fullSentCorpus[:portion]
	train := fullSentCorpus[portion:]
	return train,test
}