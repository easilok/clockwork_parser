package main

import (
	"fmt"
	"log"
	"os"

	"github.com/easilok/clockwork_parser/api"
	"github.com/easilok/clockwork_parser/parsers"
	"github.com/joho/godotenv"
)

const CLOCKWORK_API_WORKLOGS = "https://api.clockwork.report/v1/worklogs"

func main() {
	fmt.Println("Starting Program")

	err := godotenv.Load()

	if err != nil {
		log.Fatal(err)
	}

	apiToken, ok := os.LookupEnv("CLOCKWORK_API_TOKEN")
	if !ok {
		log.Fatal("Missing Clockwork API token")
	}

	project, ok := os.LookupEnv("CLOCKWORK_API_PROJECT")
	if !ok {
		project = "LT"
	}

	author := os.Getenv("CLOCKWORK_API_AUTHOR")

	startDate, ok := os.LookupEnv("CLOCKWORK_API_START_DATE")
	if !ok {
		// TODO - make it today month start
		startDate = "2022-10-01"
	}

	endDate, ok := os.LookupEnv("CLOCKWORK_API_START_END")
	if !ok {
		// TODO - make it today month end
		endDate = "2022-10-09"
	}

	worklogs, err := api.GetWorklogs(apiToken, startDate, endDate, project, author)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Fetched logs from clockwork API:")
	for _, worklog := range worklogs {
		// fmt.Printf("Log %s by %s with %ds spent: %s - %s, assigned to %s, from epic %s\n", worklog.Id, worklog.Author.Name, worklog.TimeSpentSeconds, worklog.Issue.Id, worklog.Issue.Fields.Summary, worklog.Issue.Fields.Assignee.Name, worklog.Issue.Fields.Epic.Fields.Summary)
		fmt.Printf(
			"Log by %s with %ds spent: Issue: %s (%s) - %s (%s), assigned to %s, from epic %s\n",
			worklog.Author.Name,
			worklog.TimeSpentSeconds,
			worklog.Issue.Id,
			worklog.Issue.Key,
			worklog.Issue.Fields.Summary,
			worklog.Issue.Fields.IssueType.Name,
			worklog.Issue.Fields.Assignee.Name,
			worklog.Issue.Fields.Epic.Fields.Summary,
		)
	}

	groupedWorklogs, err := parsers.GroupWorklogsByIssue(worklogs)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("")
	fmt.Println("Grouped worklogs:")
	for _, worklog := range groupedWorklogs {
		// fmt.Printf("Log %s by %s with %ds spent: %s - %s, assigned to %s, from epic %s\n", worklog.Id, worklog.Author.Name, worklog.TimeSpentSeconds, worklog.Issue.Id, worklog.Issue.Fields.Summary, worklog.Issue.Fields.Assignee.Name, worklog.Issue.Fields.Epic.Fields.Summary)
		fmt.Printf(
			"Log by %s with %ds spent: Issue: %s (%s) - %s (%s), assigned to %s, from epic %s\n",
			worklog.Author.Name,
			worklog.TimeSpentSeconds,
			worklog.Issue.Id,
			worklog.Issue.Key,
			worklog.Issue.Fields.Summary,
			worklog.Issue.Fields.IssueType.Name,
			worklog.Issue.Fields.Assignee.Name,
			worklog.Issue.Fields.Epic.Fields.Summary,
		)
	}

	epicWorklogs, err := parsers.GroupWorklogsByEpics(worklogs)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("")
	fmt.Println("Grouped epics:")
	for _, worklog := range epicWorklogs {
		// fmt.Printf("Log %s by %s with %ds spent: %s - %s, assigned to %s, from epic %s\n", worklog.Id, worklog.Author.Name, worklog.TimeSpentSeconds, worklog.Issue.Id, worklog.Issue.Fields.Summary, worklog.Issue.Fields.Assignee.Name, worklog.Issue.Fields.Epic.Fields.Summary)
		fmt.Printf(
			"Log by %s with %ds spent from Epic: %s\n",
			worklog.Author.Name,
			worklog.TimeSpentSeconds,
			worklog.Issue.Fields.Epic.Fields.Summary,
		)
	}

	err = parsers.GenerateWorklogCSV(groupedWorklogs, "test_issue.csv", false)
	if err != nil {
		log.Fatal(err)
	}

	err = parsers.GenerateWorklogCSV(epicWorklogs, "test_epic.csv", true)
	if err != nil {
		log.Fatal(err)
	}
}
