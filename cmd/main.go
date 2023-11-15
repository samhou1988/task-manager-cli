package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/google/uuid"
)

func main() {
	if len(os.Args) < 2 {
		log.Println("Please provide a command.")
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
	case "delete":
		deleteTask()
	case "deleteAll":
		deleteAllTask()
	default:
		log.Println("Invalid command.")
	}
}

const (
	TaskFileName       = "tasks.txt"
	TaskFilePermission = 0644
	TaskSplitString    = "##"
)

func WriteTaskToFile(fileName, taskId, description string) error {
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, TaskFilePermission)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = fmt.Fprintf(file, "%s##%s\n", taskId, description)
	if err != nil {
		return err
	}

	return nil
}

func addTask() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter task description: ")
	description, _ := reader.ReadString('\n')
	description = strings.TrimSpace(description)
	taskId := uuid.New().String()

	log.Println("TaskId is:", taskId)

	// Write the task to a file
	err := WriteTaskToFile(TaskFileName, taskId, description)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Task added successfully.")
}

// listTasks reads tasks from a file and prints them to the console
func listTasks() {
	// Read tasks from a file
	file, err := os.Open(TaskFileName)
	if err != nil {
		log.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	taskBytes, err := io.ReadAll(file)
	if err != nil {
		log.Println("Error reading file:", err)
		return
	}

	tasks := strings.Split(string(taskBytes), "\n")

	// Print the tasks
	log.Println("Tasks:")
	for _, task := range tasks {
		if task != "" {
			// split the task into id and description
			taskParts := strings.Split(task, TaskSplitString)
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
		log.Println("Invalid task ID.")
		return
	}

	fmt.Print("Enter new task description: ")
	newDescription, _ := reader.ReadString('\n')
	newDescription = strings.TrimSpace(newDescription)

	// Read tasks from a file
	file, err := os.Open(TaskFileName)
	if err != nil {
		log.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	taskBytes, err := io.ReadAll(file)
	if err != nil {
		log.Println("Error reading file:", err)
		return
	}

	tasks := strings.Split(string(taskBytes), "\n")

	// Update the task by ID
	found := false
	updatedTasks := []string{}
	for _, task := range tasks {
		if task != "" {
			// split the task into id and description
			taskParts := strings.Split(task, TaskSplitString)
			id := strings.TrimSpace(taskParts[0])

			if id == input {
				// update the description not work
				taskParts[1] = newDescription
				found = true
			}

			updateTask := strings.Join(taskParts, TaskSplitString)
			updatedTasks = append(updatedTasks, updateTask)
		}
	}

	if !found {
		log.Println("Task not found.")
		return
	}

	tasks = updatedTasks

	// Write the updated tasks to the file
	output := strings.Join(tasks, "\n")
	err = os.WriteFile(TaskFileName, []byte(output), 0644)
	if err != nil {
		log.Println("Error writing to file:", err)
		return
	}

	log.Println("Task updated successfully.")
}

// Delete a task by taskId
func deleteTask() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter task ID: ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "" {
		log.Println("Invalid task ID.")
		return
	}

	// Read tasks from a file
	file, err := os.Open(TaskFileName)
	if err != nil {
		log.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	taskBytes, err := io.ReadAll(file)
	if err != nil {
		log.Println("Error reading file:", err)
		return
	}

	tasks := strings.Split(string(taskBytes), "\n")

	// Update the task by ID
	found := false
	updatedTasks := []string{}
	for _, task := range tasks {
		if task != "" {
			// split the task into id and description
			taskParts := strings.Split(task, TaskSplitString)
			id := strings.TrimSpace(taskParts[0])

			if id == input {
				found = true
			} else {
				updateTask := strings.Join(taskParts, TaskSplitString)
				updatedTasks = append(updatedTasks, updateTask)
			}
		}
	}

	if !found {
		log.Println("Task not found.")
		return
	}

	tasks = updatedTasks

	// Write the updated tasks to the file
	output := strings.Join(tasks, "\n")
	err = os.WriteFile(TaskFileName, []byte(output), 0644)
	if err != nil {
		log.Println("Error writing to file:", err)
		return
	}

	log.Println("Task Remove successfully.")
}

// Delete all Task.
func deleteAllTask() {
	// Read tasks from a file
	file, err := os.Open(TaskFileName)
	if err != nil {
		log.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// remove file
	err = os.Remove(TaskFileName)
	if err != nil {
		log.Println("Error remove file:", err)
		return
	}

	log.Println("All Task Remove successfully.")
}
