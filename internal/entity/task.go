package entity

type Task struct {
	ID        int    `json:"id,omitempty"`
	Author    User   `json:"author,omitempty"`
	Descr     string `json:"descr,omitempty"`
	Body      string `json:"body,omitempty"`
	Approvers []User `json:"approvers,omitempty"`
	Finished  bool   `json:"finished,omitempty"`
}
