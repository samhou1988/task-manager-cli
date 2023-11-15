package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/google/uuid"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide a command.")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "add":
		addTask()
	case "list":
		listTasks()
	case "update":
		updateTask()
	default:
		fmt.Println("Invalid command.")
	}
}

// create a constant for the file name
const fileName = "tasks.txt"
const splitString = "##"

func addTask() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter task description: ")
	description, _ := reader.ReadString('\n')
	description = strings.TrimSpace(description)
	taskId := uuid.New().String()

	fmt.Println("TaskId is:", taskId)

	// Write the task to a file
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	_, err = fmt.Fprintf(file, "%s##%s\n", taskId, description)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	fmt.Println("Task added successfully.")
}

// listTasks reads tasks from a file and prints them to the console
func listTasks() {
	// Read tasks from a file
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	taskBytes, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	tasks := strings.Split(string(taskBytes), "\n")

	// Print the tasks
	fmt.Println("Tasks:")
	for _, task := range tasks {
		if task != "" {
			// split the task into id and description
			taskParts := strings.Split(task, splitString)
			id := strings.TrimSpace(taskParts[0])
			description := strings.TrimSpace(taskParts[1])

			fmt.Printf("%s %s\n", id, description)
		}
	}
}

func updateTask() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter task ID: ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "" {
		fmt.Println("Invalid task ID.")
		return
	}

	fmt.Print("Enter new task description: ")
	newDescription, _ := reader.ReadString('\n')
	newDescription = strings.TrimSpace(newDescription)

	// Read tasks from a file
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	taskBytes, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	tasks := strings.Split(string(taskBytes), "\n")

	// Update the task by ID
	found := false
	updatedTasks := []string{}
	for _, task := range tasks {
		if task != "" {
			// split the task into id and description
			taskParts := strings.Split(task, splitString)
			id := strings.TrimSpace(taskParts[0])

			if id == input {
				// update the description not work
				taskParts[1] = newDescription
				found = true
			}

			updateTask := strings.Join(taskParts, splitString)
			updatedTasks = append(updatedTasks, updateTask)
		}
	}

	if !found {
		fmt.Println("Task not found.")
		return
	}

	tasks = updatedTasks

	// Write the updated tasks to the file
	output := strings.Join(tasks, "\n")
	err = os.WriteFile(fileName, []byte(output), 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	fmt.Println("Task updated successfully.")
}
