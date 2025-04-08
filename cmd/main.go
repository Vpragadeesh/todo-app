package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

const todoFile = "todos.json"

// ANSI color codes for terminal output
const (
	colorGreen = "\033[32m"
	colorRed   = "\033[31m"
	colorBlue  = "\033[34m"
	colorReset = "\033[0m"
)

type Todo struct {
	Task      string `json:"task"`
	Date      string `json:"date"` // Due date in "YYYY-MM-DD" format
	Completed bool   `json:"completed"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: todo <command> [task]")
		return
	}

	command := os.Args[1]
	task := strings.Join(os.Args[2:], " ")

	// Always check and update task dates first.
	switch command {
	case "add":
		addTask(task)
	case "list":
		listTasks()
	case "complete":
		completeTask(task)
	default:
		fmt.Println("Unknown command:", command)
	}
}

func addTask(task string) {
	todos := postponeTodos(loadTodos())
	today := time.Now().Format("2006-01-02")
	// When adding a task, set its date to today.
	todos = append(todos, Todo{Task: task, Date: today, Completed: false})
	saveTodos(todos)
	fmt.Println("Added:", task, "for", today)
}

func listTasks() {
	todos := postponeTodos(loadTodos())
	todayStr := time.Now().Format("2006-01-02")
	for i, todo := range todos {
		// Completed tasks display a green check.
		var icon string
		if todo.Completed {
			icon = colorGreen + "✔" + colorReset
		} else {
			icon = " " // no icon for unfinished tasks
		}

		// Determine the status.
		var status string
		if todo.Completed {
			status = "completed"
		} else {
			if todo.Date == todayStr {
				status = "ongoing"
			} else {
				status = "pending"
			}
		}

		// Print index, icon, task, due date (in blue) and status.
		fmt.Printf("%d. [%s] %s - %s - %s\n",
			i+1,
			icon,
			todo.Task,
			colorBlue+todo.Date+colorReset,
			status)
	}
}

func completeTask(task string) {
	todos := postponeTodos(loadTodos())
	// Search and mark the matching task as completed.
	for i, todo := range todos {
		if todo.Task == task {
			todos[i].Completed = true
			saveTodos(todos)
			fmt.Println("Completed:", task)
			return
		}
	}
	fmt.Println("Task not found:", task)
}

// postponeTodos checks all tasks and applies two rules:
// 1. If a task’s Date is empty, assume today’s date.
// 2. If a task is unfinished and its date is before today, postpone it to tomorrow.
func postponeTodos(todos []Todo) []Todo {
	todayStr := time.Now().Format("2006-01-02")
	today, _ := time.Parse("2006-01-02", todayStr)
	updated := false

	for i, todo := range todos {
		// If Date is empty, assume it's for today.
		if todo.Date == "" {
			todos[i].Date = todayStr
			updated = true
			continue
		}

		due, err := time.Parse("2006-01-02", todo.Date)
		if err != nil {
			// If parsing fails, default the date to today.
			todos[i].Date = todayStr
			updated = true
			continue
		}

		// For unfinished tasks, if the due date is before today, postpone to tomorrow.
		if !todo.Completed && due.Before(today) {
			newDue := today.AddDate(0, 0, 1)
			todos[i].Date = newDue.Format("2006-01-02")
			updated = true
		}
	}

	if updated {
		saveTodos(todos)
	}
	return todos
}

func loadTodos() []Todo {
	file, err := ioutil.ReadFile(todoFile)
	if err != nil {
		return []Todo{}
	}
	var todos []Todo
	json.Unmarshal(file, &todos)
	return todos
}

func saveTodos(todos []Todo) {
	data, _ := json.Marshal(todos)
	ioutil.WriteFile(todoFile, data, 0644)
}
