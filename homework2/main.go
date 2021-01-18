package main

import "fmt"
type Task interface {
    Execute(int) (int, error)
}

type operationsConsecutive struct{
	tasks []Task
}

func Pipeline(tasks ...Task) operationsConsecutive{
	return operationsConsecutive{tasks}
}

func (ops operationsConsecutive) Execute(addEnd int) (int, error){
	var err error
	for _, value := range ops.tasks{
		addEnd, err = value.Execute(addEnd)
		if err != nil{
			return addEnd, err
		}
	}
	return addEnd, err
}

type operationsConcurrent struct{
	tasks []Task
}


type adder struct {
	augend int
}

func (a adder) Execute(addend int) (int, error) {
	result := a.augend + addend
	if result > 127 {
		return 0, fmt.Errorf("Result %d exceeds the adder threshold", a)
	}
	return result, nil
}
//func Pipeline(tasks ...Task){

//}
func main(){
	if res, err := 	Pipeline(adder{20}, adder{10}, adder{-50}).Execute(100); err != nil {
		fmt.Printf("The pipeline returned an error\n")
	} else {
		fmt.Printf("The pipeline returned %d\n", res)
	}

}