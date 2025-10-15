package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"
)

type Task struct {
	Id          int       `json:"id"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type TaskUpdate struct {
	Description *string
	Status      *string
}

const filePath = "tasks.json"

func main() {
	if len(os.Args) < 2 {
		fmt.Println("please insert command")
		return
	}

	command := os.Args[1]
	if command == "list" {
		status := "ALL"
		if len(os.Args) > 2 {
			status = os.Args[2]
			if status != "done" && status != "todo" && status != "in-progress" {
				fmt.Println("please fill valid status")
				return
			}
		}
		tasks, err := getTasks(status)
		if err != nil {
			fmt.Println(err)
			return
		}
		for _, task := range *tasks {
			fmt.Println(task)
		}
		return
	} else if command == "add" {
		if len(os.Args) < 3 {
			fmt.Println("please fill task description")
			return
		}
		description := os.Args[2]
		task, err := insertTask(description)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("success add task with id: %d\n", task.Id)
		return
	} else if command == "update" {
		if len(os.Args) < 3 {
			fmt.Println("please fill task id")
			return
		}
		if len(os.Args) < 4 {
			fmt.Println("please fill task description")
			return
		}

		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Printf("failed cast to integer: %s\n", err)
			return
		}
		description := os.Args[3]
		taskUpdate := &TaskUpdate{
			Description: &description,
		}
		task, err := updateTask(id, *taskUpdate)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(task)
		return
	}
}

func getTasks(status string) (*[]Task, error) {
	allTasks, err := safeReadFile()
	if err != nil {
		return nil, err
	}
	if status == "ALL" {
		return allTasks, nil
	}

	var filteredTasks []Task
	for _, task := range *allTasks {
		if task.Status == status {
			filteredTasks = append(filteredTasks, task)
		}
	}

	return &filteredTasks, nil
}

func insertTask(description string) (*Task, error) {
	allTasks, err := safeReadFile()
	if err != nil {
		return nil, err
	}
	lastId := 0
	if len(*allTasks) > 0 {
		lastTask := (*allTasks)[len(*allTasks)-1]
		lastId = lastTask.Id
	}
	task := Task{
		Id:          lastId + 1,
		Description: description,
	}
	*allTasks = append(*allTasks, task)

	// write file
	jsonBytes, err := json.Marshal(allTasks)
	if err != nil {
		return nil, err
	}
	if err := os.WriteFile(filePath, jsonBytes, 0644); err != nil {
		return nil, err
	}

	return &task, nil
}

func updateTask(id int, update TaskUpdate) (*Task, error) {
	allTasks, err := safeReadFile()
	if err != nil {
		return nil, err
	}

	for i := range *allTasks {
		task := &(*allTasks)[i]
		if task.Id == id {
			fmt.Println("masa ga masuk sini ya")
			if update.Description != nil {
				task.Description = *update.Description
			}
			if update.Status != nil {
				task.Status = *update.Status
			}

			jsonBytes, err := json.MarshalIndent(allTasks, "", "  ")
			if err != nil {
				return nil, err
			}
			if err := os.WriteFile(filePath, jsonBytes, 0644); err != nil {
				return nil, err
			}

			return task, nil
		}
	}
	return nil, fmt.Errorf("task with id %d not found", id)
}

func safeReadFile() (*[]Task, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if err := os.WriteFile(filePath, []byte("[]"), 0644); err != nil {
				return nil, fmt.Errorf("failed to create %s: %w", filePath, err)
			}
			return &[]Task{}, nil
		}
		return nil, fmt.Errorf("failed to read %s: %w", filePath, err)
	}

	var tasks []Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, fmt.Errorf("failed to parse %s: %w", filePath, err)
	}

	return &tasks, err
}
