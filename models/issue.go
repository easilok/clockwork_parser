package models

type Assignee struct {
	AccountId string `json:"accountId"`
	Name      string `json:"displayName"`
}

type EpicFields struct {
	Summary string `json:"summary"`
}

type Epic struct {
	Id     string     `json:"id"`
	Key    string     `json:"key"`
	Fields EpicFields `json:"fields"`
}

type IssueType struct {
	Name string `json:"name"`
}

type IssueFields struct {
	IssueType   IssueType `json:"issuetype"`
	Epic        Epic      `json:"parent"`
	Description string    `json:"description"`
	Summary     string    `json:"summary"`
	Assignee    Assignee  `json:"assignee"`
}

type Issue struct {
	Id     string      `json:"id"`
	Key    string      `json:"key"`
	Fields IssueFields `json:"fields"`
}
