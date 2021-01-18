package main

import ( "fmt"
		 "errors"
		 "time"
)
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
	if len(ops.tasks) == 0 {
		return 0, errors.New("No tasks given")
	}
	for _, task := range ops.tasks{
		addEnd, err = task.Execute(addEnd)
		if err != nil{
			return addEnd, err
		}
	}
	return addEnd, err
}

type operationsConcurrent struct{
	tasks []Task
}


func Fastest(tasks ...Task) operationsConcurrent{
	return operationsConcurrent{tasks}
}

func (ops operationsConcurrent) Execute(addEnd int) (int, error){
	numTasks := len(ops.tasks)
	out := make(chan int, numTasks )
    errs := make(chan error, numTasks)
	if len(ops.tasks) == 0 {
		return 0, errors.New("No tasks given")
	}
	for _, task :=  range ops.tasks{
		go func() {
			taskRes ,err := task.Execute(addEnd)
			errs <- err
		    out <- taskRes
			 }()
	}	
	res := <- out
	err := <- errs
	return res, err
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

type lazyAdder struct {
	adder
	delay time.Duration
}

func (la lazyAdder) Execute(addend int) (int, error) {
	time.Sleep(la.delay * time.Millisecond)
	return la.adder.Execute(addend)
}

func main(){
	if res, err := 	Pipeline(adder{20}, adder{10}, adder{-50}).Execute(100); err != nil {
		fmt.Printf("The pipeline returned an error\n")
	} else {
		fmt.Printf("The pipeline returned %d\n", res)
	}
	f := Fastest()
	fmt.Println(f.Execute(1))

}