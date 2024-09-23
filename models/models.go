package models

type Task struct {
	ID      string `json:"id, omitempty"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment, omitempty"`
	Repeat  string `json:"repeat"`
}

type AddTaskResponse struct {
	ID    int    `json:"id,omitempty"`
	Error string `json:"error,omitempty"`
}

type GetTaskResponse struct {
	Tasks []Task `json:"tasks,omitempty"`
	Error string `json:"error,omitempty"`
}
