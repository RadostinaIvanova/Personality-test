package model

import (
    "math"
    "strings"
    "sort"
    "os"
    "log"
    "encoding/gob"
    "fmt"
)
const startToken = "<START>"
const endToken = "<END>"
const unkToken = "<UNK>"


type MarkovModel struct{
    K int
    Kgrams []Kgram
    Tc Tc	
}

type Kgram struct{
    Context string
    WordCount map[string] int
}

type Tc map[string] int

type probWord struct{
	word string
	value float64
}

func (m* MarkovModel) Init(k int, train [][]string,limit int){
    kgrams := make([]Kgram,0, 100000)
    m.Tc = Tc{}
    m.Kgrams = kgrams
    m.K = k
    for i := 0;i <= m.K; i++{
        m.extractKgrams(train,i,limit)
    }
    m.calculateTc()
}

func (m *MarkovModel) SaveModel(filename string){
	f, err := os.Create(filename)
	if err != nil{
		log.Println(err.Error())
	}
	defer f.Close()
	encoder := gob.NewEncoder(f)
    fmt.Println(len(m.Kgrams))
    fmt.Println(m.K)
    fmt.Println(len(m.Tc))
	errEn :=encoder.Encode(m)
    if errEn != nil {
		log.Fatal("encode error:", err)
	}
}

func (m *MarkovModel) LoadModel(filename string){
	f, err := os.Open(filename)
	if err != nil{
		log.Println(err.Error())
	}
	defer f.Close()
	decoder := gob.NewDecoder(f)
	errd := decoder.Decode(&m)	
	if errd != nil {
		log.Fatal("decode error 1:", errd)
	}
}

func (mm *MarkovModel) BestContinuation(sentence []string, alpha float64, l int) []string{
	context := mm.getContext(sentence, mm.K, len(sentence))
	candidates := []string{}
	for k := 0;k < mm.K;k++{
		if ind := mm.existContext(strings.Join(context[k:], " ")); ind > -1 && len(mm.countContext(strings.Join(context[k:], " "))) > l{
            candidates = mm.countContext(strings.Join(context[k:], " "))
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
        if wordProb.word != endToken{
		result = append(result, wordProb.word)
        }
	}
	return result[:l]
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
    monograms := Kgram{"", map[string]int{}}
    for k,v := range dictionary {
        monograms.WordCount[k] = v
    }
    mm.Kgrams = append(mm.Kgrams, monograms)
}



func (mm *MarkovModel) extractKgrams(corpus [][]string, k int,limit int){
    j:=0
    for _, sent := range corpus{
        for i, word := range sent{
            if (word == startToken){
                continue
            }
            context := strings.Join(mm.getContext(sent, mm.K, i), " ")
            
            if ind := mm.existContext(context); ind >= 0{
                 if _, ok := mm.Kgrams[ind].WordCount[word]; !ok{
                    mm.Kgrams[ind].WordCount[word] = 0
                 }
                    mm.Kgrams[ind].WordCount[word] += 1
           }else{
                mm.Kgrams = append(mm.Kgrams, Kgram{context, map[string] int{word : 1,}})
           }
            j++
            if(j == limit){
                return
            }
        }
    } 
}

func (mm *MarkovModel) probMLE(word string ,con string) float64{
    ind := mm.existInKgrams(con,word);
	if ind == -1 {
        return 0.0
    }
	return float64(mm.Kgrams[ind].WordCount[word])/float64(mm.Tc[con])
}


func (mm *MarkovModel) prob(word string, context[] string, alpha float64) float64{
	con := strings.Join(context, " ")
    if con != "" {
        return alpha * mm.probMLE(word, con) + (1-alpha) * mm.prob(word, context[1:], alpha)
    }
    return mm.probMLE(word, con)
}


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


func (mm *MarkovModel) calculateTc(){
    for _, kgram := range(mm.Kgrams){
            mm.Tc[kgram.Context] = 0
        for _,value := range kgram.WordCount{
            mm.Tc[kgram.Context] += value
         }
    }
}

//CHANGED
func (mm *MarkovModel) existInKgrams(context string, word string) int{
    if ind:= mm.existContext(context); ind >= 0{
			if _, ok := mm.Kgrams[ind].WordCount[word];ok{
                return ind
			}
		}
    return -1
}

func (mm *MarkovModel) existContext(context string) int{
    for ind , kgram := range mm.Kgrams{	
        if kgram.Context == context {
            return ind
        }
    }
    return -1
}

func (mm *MarkovModel) countContext(context string)[]string{
	candidates := []string{}
	if ind := mm.existContext(context); ind > -1{
		for key,_ := range mm.Kgrams[ind].WordCount{
		    candidates = append(candidates, key)
		}
    }
	return candidates
}


// func main(){
//     sentences := Extract("D:\\FMI\\Info\\dialogues_train.txt")
//     fullSentCorpus := FullSentCorpus(sentences)
//     train, test := DivideIntoTrainAndTest(0.1, fullSentCorpus)
//     numGram := 2
//     m := MarkovModel{}
//     m.Init(numGram,train,4000000)
//     fmt.Println(m.perplexity(test,0.6))
//     fmt.Println(len(m.Kgrams))
// 	fmt.Println(m.BestContinuation([]string{"<START>", "i", "love", "going"}, 0.6,10))
//     fmt.Println(m.BestContinuation([]string{"<START>", "Watch", "me", "turn", "into"}, 0.6,10))
// }