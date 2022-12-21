package main

import (
	"fmt"
	"time"
)

type Ttype struct {
	id         int
	cT         string // время создания
	fT         string // время выполнения
	taskString string
}

type resultStruct struct {
	taskRESULT []byte
}

func taskCreturer(superChan chan Ttype, a Ttype) {
	for {
		ct := time.Now().Format(time.RFC3339)
		if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков.
			ct = "Some error occured"
		}

		time.Sleep(time.Millisecond * 1000)                                              // это для того, чтобы id не повторялись
		superChan <- Ttype{cT: ct, id: int(time.Now().Unix()), taskString: a.taskString} // передаем таск на выполнение
	}
}

func task_worker(a Ttype, r resultStruct) (Ttype, resultStruct) {
	for {
		tt, _ := time.Parse(time.RFC3339, a.cT)
		if tt.After(time.Now().Add(-20 * time.Second)) {
			r.taskRESULT = []byte("task has been successed")
		} else {
			r.taskRESULT = []byte("something went wrong")
		}
		a.fT = time.Now().Format(time.RFC3339Nano)

		time.Sleep(time.Millisecond * 150)

		return a, r
	}
}

func taskSorter(superChan chan Ttype, doneTasks chan Ttype, undoneTasks chan error, t Ttype, r resultStruct) {
	for t := range superChan {

		l, r := task_worker(t, r)

		if string(r.taskRESULT[14:]) == "successed" {

			t.taskString = string(r.taskRESULT)
			doneTasks <- t
		} else {
			undoneTasks <- fmt.Errorf("Task id: %d \nCreation time: %s, \nExecution time: %s, \nError: %s", t.id, t.cT, l.fT, r.taskRESULT)
		}
	}
}

func Result(doneTasks chan Ttype, undoneTasks chan error) {
	go func() {
		for d := range doneTasks {
			fmt.Println("\nDone tasks:")
			fmt.Println(d)
		}
	}()

	for u := range undoneTasks {
		fmt.Println("\nErrors!")
		fmt.Println(u)
	}
}

func main() {
	superChan := make(chan Ttype, 1)
	doneTasks := make(chan Ttype, 1)
	undoneTasks := make(chan error, 1)

	go taskCreturer(superChan, Ttype{})
	go task_worker(Ttype{}, resultStruct{})
	go taskSorter(superChan, doneTasks, undoneTasks, Ttype{}, resultStruct{})
	Result(doneTasks, undoneTasks)
}
