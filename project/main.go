package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"io"
	"strings"
	"strconv"

)

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
func putInClass(records [][]string){

	classes :=  make(map[int] []string)
	for _, record := range records{
		classTypeStr := record[len(record) - 1]
		classType := change(classTypeStr)
		record = record[0 : len(record) - 1]
		newRec := strings.Join(record[:]," ")
		classes[classType] = append(classes[classType], newRec)
		
	}
	for key, _ := range classes{
		fmt.Print(key)
	}
}

func main() {
	csvFile, err := os.Open("D://FMI//golang_workspace//src//project//data1.csv")
	defer csvFile.Close()

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened CSV file")
    csvLines := csv.NewReader(csvFile)
	check := false
	var str1 [][]string
	for {
		
		record, err1 := csvLines.Read()
		if(err1 == io.EOF){
			break
		}
		
		if notCorruptedRecord(record){
			if check == true{
			record = record[2:len(record)-1]
			record = convertFloatToStrRecord(record) 
			str1 = append(str1, record)
			}
		}
		check = true
	 }
	putInClass(str1)

}
