package main

import (
	"time"
)

type TodoItem struct {
	Id      string `gorethink:"id,omitempty"`
	Text    string
	Status  string
	Created time.Time
}

func (t *TodoItem) Completed() bool {
	return t.Status == "complete"
}

func NewTodoItem(text string) *TodoItem {
	return &TodoItem{
		Text:   text,
		Status: "active",
	}
}
