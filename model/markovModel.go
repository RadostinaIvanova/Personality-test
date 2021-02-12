package main


import (
    "math"
    "os"
    "fmt"
    "bufio"
    "strings"
    "regexp"
    "log"
	"sort"
)
const startToken = "<START>"
const endToken = "<END>"
const unkToken = "<UNK>"


type MarkovModel struct{
    K int
    kgrams []Kgram
    tc Tc	
}

type Kgram struct{
    context string
    wordCount map[string] int
}

type Tc map[string] int

type probWord struct{
	word string
	value float64
}
   
func (mm *MarkovModel) extractMonograms(corpus [][]string,  limit int){
    dictionary := make(map [string] int)
    for _, sent := range(corpus){
        for _, word := range(sent){
            if word == startToken{
				continue
			}
                if _, ok := dictionary[word]; !ok{
                    dictionary[word] = 0
                }
            	dictionary[word] += 1
		}
    }
    monograms := []Kgram{}
    for k,v := range dictionary {
        monograms = append(monograms, Kgram{"", map[string]int{k:v,}})
    }
    mm.kgrams = append(mm.kgrams, monograms...)
}


//changed
func (mm *MarkovModel) probMLE(word string ,con string) float64{
    ind := mm.existInKgrams(con,word);
	if ind == -1 {
        return 0.0
    }
	return float64(mm.kgrams[ind].wordCount[word])/float64(mm.tc[con])
}


//changed
func (mm *MarkovModel) prob(word string, context[] string, alpha float64) float64{
	con := strings.Join(context, " ")
    if con != "" {
        return alpha * mm.probMLE(word, con) + (1-alpha) * mm.prob(word, context[1:], alpha)
    }
    return mm.probMLE(word, con)
}


//hardcoded 2 because only 1 contex word
func (mm *MarkovModel) sentenceLogProbability(sentence []string, alpha float64) float64{
    sum := 0.0
    for key,value :=  range(sentence){
        if value != startToken {
			if mm.prob(value, mm.getContext(sentence, mm.K, key), alpha) == 0{
				break;
				
			}
            sum += math.Log2(mm.prob(value, mm.getContext(sentence, mm.K, key), alpha))
        }
    }
    return sum
}


func (mm *MarkovModel) extractKgrams(corpus [][]string, k int,limit int){
        j:=0
        for _, sent := range corpus{
            for i, word := range sent{
                if (word == startToken){
                    continue
                }
                context := strings.Join(mm.getContext(sent,  2, i), " ")
				
                if ind := mm.existInKgrams(context, word); ind >= 0{
					if context == ""{
						fmt.Println(mm.kgrams[ind].wordCount[word])
					}
						mm.kgrams[ind].wordCount[word] += 1
                }else{
                    mm.kgrams = append(mm.kgrams, Kgram{context, map[string] int{word : 1,}})
                }
                j++
                if(j == limit){
                    return
                }
            }
        } 
}

func (mm *MarkovModel) perplexity(corpus [][]string, alpha float64)float64{
    sum := 0
    for _, sentence := range(corpus){
        sum += (len(sentence)-1)
    }
    crossEntropy := 0.0
    for _, sentence := range(corpus){
        crossEntropy -= mm.sentenceLogProbability(sentence,alpha)
    }
    crossEntropyRate := crossEntropy / float64(sum)
    return math.Pow(2, crossEntropyRate)
}

func (mm *MarkovModel) getContext(sent []string, k int, i int) []string{
    context := []string{}
    if i-k+1 >= 0{
        context = append(context, sent[i-k+1:i]...)
    }else{
        for j:= 0;j < k-i-1; j++{
            context = append(context, startToken)
        }
        context = append(context, sent[0:i]...)
    }
    return context
}


//changed
func (mm *MarkovModel) calculateTc(){

    for _, kgram := range(mm.kgrams){
        if _, ok := mm.tc[kgram.context]; !ok {
            mm.tc[kgram.context] = 1
        }else{
				for _,value := range kgram.wordCount{
       				mm.tc[kgram.context] += value
			}
        }
    }

}

//CHANGED
func (mm *MarkovModel) existInKgrams(context string, word string) int{
    for ind , kgram := range mm.kgrams{	
        if kgram.context == context {
			if _, ok := mm.kgrams[ind].wordCount[word];ok{
                return ind
			}
		}
	}
    return -1
}

func (mm *MarkovModel) contextInKgrams(context string) int{
    for ind , kgram := range mm.kgrams{	
        if kgram.context == context {
                return ind
			}
		}
    return -1
}

func (mm *MarkovModel) countContext(context string)[]string{
	candidates := []string{}
	for _, kgram := range(mm.kgrams){
		if kgram.context == context{
			for key,_ := range kgram.wordCount{
			candidates = append(candidates, key)
			}
		}
	}
	return candidates
}
func (mm *MarkovModel) bestContinuation(sentence []string, alpha float64, l int) []string{
	context := mm.getContext(sentence, mm.K, len(sentence))
	con := strings.Join(context, " ")
	candidates := []string{}
	for k := 0;k < mm.K;k++{
		can := mm.countContext(con)
		if ind := mm.contextInKgrams(strings.Join(context[k:], " ")); ind > -1 && len(can) > l{
			candidates = can
			break
		}
	}
	wProb := []probWord{}
	result := []string{}
	for _,word := range(candidates){
		wProb = append(wProb, probWord{word, mm.prob(word,context,alpha)})
	}
	sort.SliceStable(wProb, func(i, j int) bool {return wProb[i].value > wProb[j].value})
	for _ ,wordProb  := range wProb{
		result = append(result, wordProb.word)
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

func fullSentCorpus(sentences []string) [][]string {
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

func delete_empty (s []string) []string {
    var r []string
    for _, str := range s {
        if str != "" {
            r = append(r, str)
        }
    }
    return r
}
func main(){
    a := extract("D:\\FMI\\Info\\dialogues_train.txt")
    sentences := (transform(a))
    fullSentCorpus := fullSentCorpus(sentences)
    percent := int(0.1*float64(len(fullSentCorpus)))
    train := fullSentCorpus[percent:]
    kgrams := make([]Kgram,0, 4025000)
    tc := Tc{}
    mm := MarkovModel{2, kgrams, tc}
    mm.extractMonograms(train, 1)
    mm.extractKgrams(train,2,1000000)
    mm.calculateTc()
	fmt.Println(mm.bestContinuation([]string{"<START>", "start", "car","engine"}, 0.6,5))
}