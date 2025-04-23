package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
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

	switch command {
	case "add":
		addTask(task)
	case "list":
		listTasks()
	case "complete":
		completeTask(task)
	case "delete":
		if len(os.Args) < 3 {
			fmt.Println("Usage: todo delete <task index>")
			return
		}
		index, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Invalid index:", os.Args[2])
			return
		}
		deleteTask(index)
	default:
		fmt.Println("Unknown command:", command)
	}
}

func addTask(task string) {
	todos := postponeTodos(loadTodos())
	today := time.Now().Format("2006-01-02")
	todos = append(todos, Todo{Task: task, Date: today, Completed: false})
	saveTodos(todos)
	fmt.Println("Added:", task, "for", today)
}

func listTasks() {
	todos := postponeTodos(loadTodos())
	sortTodos(todos)
	todayStr := time.Now().Format("2006-01-02")
	for i, todo := range todos {
		var icon string
		if todo.Completed {
			icon = colorGreen + "✔" + colorReset
		} else {
			icon = " " // no icon for unfinished tasks
		}
		var status string
		if todo.Completed {
			status = "completed"
		} else if todo.Date == todayStr {
			status = "ongoing"
		} else {
			status = "pending"
		}
		numCol := fmt.Sprintf("%d.", i+1)
		fmt.Printf("%-4s [%s] %s - %s - %s\n",
			numCol,
			icon,
			todo.Task,
			colorBlue+todo.Date+colorReset,
			status)
	}
}

func completeTask(task string) {
	todos := postponeTodos(loadTodos())
	// Find and update the matching task.
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

func deleteTask(index int) {
	todos := postponeTodos(loadTodos())
	// Sort the todos EXACTLY as they appear in the listTasks function
	sortTodos(todos)

	if index < 1 || index > len(todos) {
		fmt.Println("Invalid index.")
		return
	}
	removedTask := todos[index-1]
	todos = append(todos[:index-1], todos[index:]...)
	saveTodos(todos)
	fmt.Println("Deleted:", removedTask.Task)
}

// sortTodos sorts the todos by date and then by task name.
func sortTodos(todos []Todo) {
	sort.Slice(todos, func(i, j int) bool {
		if todos[i].Date == todos[j].Date {
			return todos[i].Task < todos[j].Task
		}
		return todos[i].Date < todos[j].Date
	})
}

// postponeTodos checks tasks and updates dates as needed.
// If a task's date is empty, it defaults to today;
// if an unfinished task's due date is before today, postpone it to tomorrow.
func postponeTodos(todos []Todo) []Todo {
	todayStr := time.Now().Format("2006-01-02")
	today, _ := time.Parse("2006-01-02", todayStr)
	updated := false
	for i, todo := range todos {
		if todo.Date == "" {
			todos[i].Date = todayStr
			updated = true
			continue
		}
		due, err := time.Parse("2006-01-02", todo.Date)
		if err != nil {
			todos[i].Date = todayStr
			updated = true
			continue
		}
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
