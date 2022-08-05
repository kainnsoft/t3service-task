package entity

type Task struct {
	Finished  bool   `json:"finished,omitempty"`
	ID        int    `json:"id,omitempty"`
	Descr     string `json:"descr,omitempty"`
	Body      string `json:"body,omitempty"`
	Author    User   `json:"author,omitempty"`
	Approvers []User `json:"approvers,omitempty"`
}
