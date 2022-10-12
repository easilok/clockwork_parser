package models

import (
	"fmt"
	"math"
)

// {"id":"29989","timeSpentSeconds":3480,"started":"2022-10-03T15:19:48Z","author":{"accountId":"6255836b410206006e966f88"},"issue":{"id":"21103"}}

type Worklog struct {
	Id               string `json:"id"`
	TimeSpentSeconds uint64 `json:"timeSpentSeconds"`
	Started          string `json:"started"`
	Author           Author `json:"author"`
	Issue            Issue  `json:"issue"`
}

func (worklog *Worklog) GetIssueCSVLine() (string, error) {
	lineTemplate := "%s,\"%s\",%s,%s,%d,\"%s\","

	// Convert time spent to hours
	timeSpentHours := int32(math.Ceil(float64(worklog.TimeSpentSeconds) / 3600))

	line := fmt.Sprintf(
		lineTemplate,
		worklog.Issue.Fields.Epic.Fields.Summary,
		worklog.Issue.Fields.Summary,
		worklog.Issue.Key,
		worklog.Issue.Fields.IssueType.Name,
		timeSpentHours,
		worklog.Issue.Fields.Assignee.Name,
	)

	return line, nil
}

func (worklog *Worklog) GetEpicCSVLine() (string, error) {
	lineTemplate := "\"%s\",%d,"

	// Convert time spent to hours
	timeSpentHours := int32(math.Ceil(float64(worklog.TimeSpentSeconds) / 3600))

	line := fmt.Sprintf(
		lineTemplate,
		worklog.Issue.Fields.Epic.Fields.Summary,
		timeSpentHours,
	)

	return line, nil
}
