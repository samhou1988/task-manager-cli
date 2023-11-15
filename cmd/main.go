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

type Task struct {
	ID          string
	Description string
}

type TaskController struct {
	tasks []Task
}

func (tc *TaskController) AddTask(task Task) error {
	tc.tasks = append(tc.tasks, task)

	return WriteTaskToFile(TaskFileName, task.ID, task.Description)
}

func (tc *TaskController) ListTasks() []Task {
	file, err := os.Open(TaskFileName)
	if err != nil {
		log.Println("Error opening file:", err)
		return nil
	}
	defer file.Close()

	taskBytes, err := io.ReadAll(file)
	if err != nil {
		log.Println("Error reading file:", err)
		return nil
	}

	tasks := strings.Split(string(taskBytes), "\n")
	return tc.convertToTasks(tasks)
}

func (tc *TaskController) convertToTasks(tasks []string) []Task {
	var convertedTasks []Task
	for _, task := range tasks {
		if task != "" {
			// split the task into id and description
			taskParts := strings.Split(task, TaskSplitString)
			id := strings.TrimSpace(taskParts[0])
			description := strings.TrimSpace(taskParts[1])

			convertedTasks = append(convertedTasks, Task{
				ID:          id,
				Description: description,
			})
		}
	}

	return convertedTasks
}

type TaskView struct {
}

func (tv *TaskView) DisplayTasks(tasks []Task) {
	for _, task := range tasks {
		fmt.Printf("%s %s\n", task.ID, task.Description)
	}
}

func (tv *TaskView) DisplayError(err error) {

}

func main() {
	if len(os.Args) < 2 {
		log.Println("Please provide a command.")
		os.Exit(1)
	}

	command := os.Args[1]
	tc := TaskController{}
	tv := TaskView{}

	switch command {
	case "add":
		description := ReadFromInput()
		task := Task{
			ID:          uuid.New().String(),
			Description: description,
		}
		err := tc.AddTask(task)
		if err != nil {
			log.Println(err)
		}
	case "list":
		tasks := tc.ListTasks()
		tv.DisplayTasks(tasks)
	case "update":
		// update task

	// 	updateTask()
	// case "delete":
	// 	deleteTask()
	// case "deleteAll":
	// 	deleteAllTask()
	default:
		log.Println("Invalid command.")
	}
}

const (
	TaskFileName       = "tasks.txt"
	TaskFilePermission = 0644
	TaskSplitString    = "##"
)

func ReadFromInput() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter task description: ")
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

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
