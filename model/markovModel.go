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

//Contains contexts as keys of type string and number occurences as value.
type Tc map[string] int

//Contains words as key of type string and their probability
type probWord struct{
	word string
	value float64
}


//Initializes a MarkovModel object by extracting grams of type <= K from given train set of sentences where 
//every sentence is implemented as slice of strings.
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

//Saves trined model by encoding it and write it to file.
func (m *MarkovModel) SaveModel(filename string){
	f, err := os.Create(filename)
	if err != nil{
		log.Println(err.Error())
	}
	defer f.Close()
	encoder := gob.NewEncoder(f)
	errEn := encoder.Encode(m)
    if errEn != nil {
		log.Fatal("encode error:", err)
	}
}

//Loads trained model by reading (encoded) trained model from file, decodes it and assigns it to m.
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

//Returns l best samples for continuing of given sentence using probability distribution.
func (mm *MarkovModel) BestContinuation(sentence []string, alpha float64, l int) []string{
	context := mm.getContext(sentence, mm.K, len(sentence))
	candidates := []string{}
    check := false
	for k := 0;k < mm.K;k++{
		if ind := mm.existContext(strings.Join(context[k:], " ")); ind > -1 && len(mm.countContext(strings.Join(context[k:], " "))) > l{
            candidates = mm.countContext(strings.Join(context[k:], " "))
            check = true
			break
		}
	}
    if check == true{
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
	return []string{}
}

//Extracts monograms from corpus by making dictionary of unique word and their count 
//and then copy dictionary to type Kgram with context - empty string.
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


//Extracts kgrams from corpus 
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

//Calculates probability of word to be next in a given context.
func (mm *MarkovModel) probMLE(word string ,con string) float64{
    ind := mm.existInKgrams(con,word);
	if ind == -1 {
        return 0.0
    }
	return float64(mm.Kgrams[ind].WordCount[word])/float64(mm.Tc[con])
}


//returning probability distribution by maximizing a likelihood function, 
//so that under the assumed statistical model the observed data is most probable.
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

//Return the the measurement of how well probability model predicts a sample.
func (mm *MarkovModel) Perplexity(corpus [][]string, alpha float64)float64{
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

//Extracts context from sentence for given word at place i with length from i-k to i.
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

//Caluclates for each context all occurences in whole corpus.
func (mm *MarkovModel) calculateTc(){
    for _, kgram := range(mm.Kgrams){
            mm.Tc[kgram.Context] = 0
        for _,value := range kgram.WordCount{
            mm.Tc[kgram.Context] += value
         }
    }
}

//Checks if given context and word are already in Kgrams of the model.
func (mm *MarkovModel) existInKgrams(context string, word string) int{
    if ind:= mm.existContext(context); ind >= 0{
			if _, ok := mm.Kgrams[ind].WordCount[word];ok{
                return ind
			}
		}
    return -1
}

//Checks if context exist in Kgrams of the model.
func (mm *MarkovModel) existContext(context string) int{
    for ind , kgram := range mm.Kgrams{	
        if kgram.Context == context {
            return ind
        }
    }
    return -1
}

//Counts context occurences by summing for every Kgram with given context all its map values from Kgram.WordCount.
func (mm *MarkovModel) countContext(context string)[]string{
	candidates := []string{}
	if ind := mm.existContext(context); ind > -1{
		for key,_ := range mm.Kgrams[ind].WordCount{
		    candidates = append(candidates, key)
		}
    }
	return candidates
}
