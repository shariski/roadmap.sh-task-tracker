package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

type Task struct {
	Id          string    `json:"id"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

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
		tasks := getTasks(status)
		for _, task := range tasks {
			fmt.Println(task)
		}
		return
	}
}

func getTasks(status string) []Task {
	allTasks, err := safeReadFile()
	if err != nil {
		panic(err)
	}
	if status == "ALL" {
		return allTasks
	}

	var filteredTasks []Task
	for _, task := range allTasks {
		if task.Status == status {
			filteredTasks = append(filteredTasks, task)
		}
	}

	return filteredTasks
}

func safeReadFile() ([]Task, error) {
	filePath := "tasks.json"

	data, err := os.ReadFile(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if err := os.WriteFile(filePath, []byte("[]"), 0644); err != nil {
				return nil, fmt.Errorf("failed to create %s: %w", filePath, err)
			}
			return []Task{}, nil
		}
		return nil, fmt.Errorf("failed to read %s: %w", filePath, err)
	}

	var tasks []Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, fmt.Errorf("failed to parse %s: %w", filePath, err)
	}

	return tasks, err
}
