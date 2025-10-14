package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

type Task struct {
	Id string `json:"id"`;
	Description string `json:"description"`;
	Status string `json:"status"`;
	CreatedAt time.Time `json:"createdAt"`;
	UpdatedAt time.Time `json:"updatedAt"`;
}

func main() {
	// index 0 is path of the binary
	fmt.Println(os.Args[1:])	
	tasks := getTasks("ALL")
	fmt.Println(tasks)
}

func getTasks(status string) []Task {
	fmt.Println(status)
	file, err := os.ReadFile("tasks.json")
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			err = os.WriteFile("tasks.json", []byte("[]"), 0644)
			if err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}
	var tasks []Task
	if file != nil {
		err = json.Unmarshal(file, &tasks)
	}
	if err != nil {
		panic(err)
	}
	return tasks 
}
