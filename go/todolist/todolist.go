package main

import (
	"fmt"
	"sort"
)

type Task struct {
	TaskID      int
	Description string
	DueDate     int
	Tags        map[string]bool
	Completed   bool
}

type TodoList struct {
	nextTaskID int
	userTasks  map[int][]*Task // userId → slice of pointers to Task
}

func Constructor() *TodoList {
	return &TodoList{
		nextTaskID: 1,
		userTasks:  make(map[int][]*Task),
	}
}

func (tl *TodoList) AddTask(userId int, description string, dueDate int, tags []string) int {
	task := &Task{
		TaskID:      tl.nextTaskID,
		Description: description,
		DueDate:     dueDate,
		Tags:        make(map[string]bool),
		Completed:   false,
	}
	for _, t := range tags {
		task.Tags[t] = true
	}
	tl.userTasks[userId] = append(tl.userTasks[userId], task)
	tl.nextTaskID++
	return task.TaskID
}

func (tl *TodoList) GetAllTasks(userId int) []string {
	tasks, exists := tl.userTasks[userId]
	if !exists {
		return []string{}
	}
	var filtered []*Task
	for _, t := range tasks {
		if !t.Completed {
			filtered = append(filtered, t)
		}
	}
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].DueDate < filtered[j].DueDate
	})
	res := make([]string, len(filtered))
	for i, t := range filtered {
		res[i] = t.Description
	}
	return res
}

func (tl *TodoList) GetTasksForTag(userId int, tag string) []string {
	tasks, exists := tl.userTasks[userId]
	if !exists {
		return []string{}
	}
	var filtered []*Task
	for _, t := range tasks {
		if !t.Completed && t.Tags[tag] {
			filtered = append(filtered, t)
		}
	}
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].DueDate < filtered[j].DueDate
	})
	res := make([]string, len(filtered))
	for i, t := range filtered {
		res[i] = t.Description
	}
	return res
}

func (tl *TodoList) CompleteTask(userId int, taskId int) {
	tasks, exists := tl.userTasks[userId]
	if !exists {
		return
	}
	for _, t := range tasks {
		if t.TaskID == taskId && !t.Completed {
			t.Completed = true
			break
		}
	}
}

// (Optional) Example usage
func main() {
	todo := Constructor()
	fmt.Println(todo.AddTask(1, "Task1", 50, []string{}))      // 1
	fmt.Println(todo.AddTask(1, "Task2", 100, []string{"P1"})) // 2
	fmt.Println(todo.GetAllTasks(1))                           // ["Task1", "Task2"]
	fmt.Println(todo.GetAllTasks(5))                           // []
	fmt.Println(todo.AddTask(1, "Task3", 30, []string{"P1"}))  // 3
	fmt.Println(todo.GetTasksForTag(1, "P1"))                  // ["Task3", "Task2"]
	todo.CompleteTask(5, 1)                                    // nothing
	todo.CompleteTask(1, 2)                                    // complete Task2
	fmt.Println(todo.GetTasksForTag(1, "P1"))                  // ["Task3"]
	fmt.Println(todo.GetAllTasks(1))                           // ["Task3", "Task1"]
}
