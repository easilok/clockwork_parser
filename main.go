package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/easilok/clockwork_parser/api"
	"github.com/easilok/clockwork_parser/parsers"
	"github.com/joho/godotenv"
)

const CLOCKWORK_API_WORKLOGS = "https://api.clockwork.report/v1/worklogs"

func getCurrentMonthLimits() (time.Time, time.Time) {
	now := time.Now()
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()

	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

	return firstOfMonth, lastOfMonth
}

func main() {
	fmt.Println("Starting Program")

	defaultStartDate, defaultEndDate := getCurrentMonthLimits()

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
		// Using a default project for lima team
		project = "LT"
	}

	author := os.Getenv("CLOCKWORK_API_AUTHOR")

	startDate, ok := os.LookupEnv("CLOCKWORK_API_START_DATE")
	if !ok {
		startDate = defaultStartDate.Format("2006-01-02")
	}

	endDate, ok := os.LookupEnv("CLOCKWORK_API_END_DATE")
	if !ok {
		endDate = defaultEndDate.Format("2006-01-02")
	}

	fmt.Printf("Fetching clockwork logs for %s, from %s to %s\n\n", author, startDate, endDate)
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

	fmt.Println("")
	fmt.Println("Creating CSV files reports.")
	fileTimestamp := defaultStartDate.Format("200601")
	startDateTime, err := time.Parse("2006-01-02", startDate)
	if err == nil {
		fileTimestamp = startDateTime.Format("200601")
	}
	// Create a gorouting waiting group with two waits
	wg := new(sync.WaitGroup)
	wg.Add(2)
	go parsers.GenerateWorklogCSV(
		wg,
		groupedWorklogs,
		fmt.Sprintf(
			"%s_test_issue.csv",
			fileTimestamp,
		),
		false,
	)

	go parsers.GenerateWorklogCSV(
		wg,
		epicWorklogs,
		fmt.Sprintf(
			"%s_test_epic.csv",
			fileTimestamp,
		),
		true,
	)

	fmt.Println("CSV files created \"export\" folder.")
}
