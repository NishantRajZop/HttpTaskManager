package main

import "fmt"

/*  a Struct task consist of Three Attributes id , taskName and its Status */
type task struct {
	id        int
	taskName  string
	completed bool
}

/*
this generateID returns You a Counter which will initialize Your getNextID() to 0
and then after on every getNextID()  call it will increment the id and return you a new ID.
*/
func generateID() func() int {
	id := 0
	returningFunction := func() int {
		id++
		return id
	}

	return returningFunction
}

// addTask appends the current task to the tasks slice and returns the new slice
// addTask takes Two Params.  CurrTask and TaskList , append the CurrTask and Returns the Updated Task

func ListTasks(tasks []task, completed map[int]bool) string {
	var s string

	s = fmt.Sprintf("total tasks: %d\n", len(tasks))

	for _, tasksOneByOne := range tasks {
		status := "Pending"
		if completed[tasksOneByOne.id] {
			status = "Completed"
		}

		s += fmt.Sprintf("ID: %d | %-15s | %s\n", tasksOneByOne.id, tasksOneByOne.taskName, status)
	}

	fmt.Println(s)

	return s
}

func markTaskAsCompleted(id int, tasks []task, completed map[int]bool) {
	for _, tasksOnebyOne := range tasks {
		if tasksOnebyOne.id == id {
			completed[tasksOnebyOne.id] = true
			return
		}
	}

	fmt.Println("Sorry This id Doesnt Exist ->", id)
}

// func main() {
//	var tasks []task // task will hold an array of TASK [{id,name,completed} , {id,name,completed} , {}]
//	var completed map[int]bool
//	completed = make(map[int]bool)
//	getNextId := generateID()
//
//	task1 := task{getNextId(), "Have BreakFast", false}
//	tasks = addTask(task1, tasks)
//
//	ListTasks(tasks, completed)
//
//	task2 := task{getNextId(), "Have Lunch", false}
//	tasks = addTask(task2, tasks)
//
//	ListTasks(tasks, completed)
//
//	markTaskAsCompleted(1, tasks, completed)
//	markTaskAsCompleted(2, tasks, completed)
//	markTaskAsCompleted(3, tasks, completed)
//	//markTaskAsCompleted(4, tasks, completed)
//
//	ListTasks(tasks, completed)
//
// }
