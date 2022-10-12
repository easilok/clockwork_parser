package parsers

import "github.com/easilok/clockwork_parser/models"

func GroupWorklogsByIssue(worklogs []models.Worklog) ([]models.Worklog, error) {

	// Receive a list of worklogs
	// create a map[string]worklog with the key being the issue ID/KEY
	worklogGroup := make(map[string]models.Worklog)

	for _, worklog := range worklogs {
		if val, ok := worklogGroup[worklog.Issue.Id]; ok {
			// If a issue ID/KEY does not exist as map key, create with worklog data
			val.TimeSpentSeconds = val.TimeSpentSeconds + worklog.TimeSpentSeconds
			worklogGroup[worklog.Issue.Id] = val
		} else {
			// else sum the TimeSpentSeconds field with existing map and new work
			worklogGroup[worklog.Issue.Id] = models.Worklog(worklog)
		}

	}

	// Flatten the created map into an array
	result := []models.Worklog{}

	for _, worklog := range worklogGroup {
		result = append(result, worklog)
	}

	return result, nil
}

func GroupWorklogsByEpics(worklogs []models.Worklog) ([]models.Worklog, error) {

	// Receive a list of worklogs
	// create a map[string]worklog with the key being the epic name
	worklogGroup := make(map[string]models.Worklog)

	for _, worklog := range worklogs {

		worklogKey := worklog.Issue.Fields.Epic.Id
		if worklog.Author.AccountId != worklog.Issue.Fields.Assignee.AccountId {
			// Check if worklog is a peer review
			worklogKey = "Review"
			worklog.Issue.Fields.Epic.Fields.Summary = "Review"
		} else if worklog.Issue.Fields.IssueType.Name == "Bug" {
			// Check if worklog is a bug
			worklogKey = "Bug"
			worklog.Issue.Fields.Epic.Fields.Summary = "Bug"
		} else if len(worklog.Issue.Fields.Epic.Fields.Summary) == 0 {
			// Check if is an issue without epic
			worklogKey = worklog.Issue.Fields.Summary
			worklog.Issue.Fields.Epic.Fields.Summary = worklog.Issue.Fields.Summary
		}

		if val, ok := worklogGroup[worklogKey]; ok {
			// If a issue ID/KEY does not exist as map key, create with worklog data
			val.TimeSpentSeconds = val.TimeSpentSeconds + worklog.TimeSpentSeconds
			worklogGroup[worklogKey] = val
		} else {
			// else sum the TimeSpentSeconds field with existing map and new work
			worklogGroup[worklogKey] = models.Worklog(worklog)
		}

	}

	// Flatten the created map into an array
	result := []models.Worklog{}

	for _, worklog := range worklogGroup {
		result = append(result, worklog)
	}

	return result, nil
}
