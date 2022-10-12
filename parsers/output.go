package parsers

import (
	"os"
	"path"

	"github.com/easilok/clockwork_parser/models"
)

func GenerateWorklogCSV(worklogs []models.Worklog, filename string, epicSummary bool) error {

	_ = os.Mkdir("export", os.ModePerm)

	f, err := os.Create(path.Join("export", filename))
	if err != nil {
		return err
	}

	defer f.Close()

	fileHeader := "Epic,Issue,Key,Type,HoursSpent,Assignee,"
	if epicSummary {
		fileHeader = "Epic,HoursSpent,"
	}

	f.WriteString(fileHeader + "\n")

	for _, worklog := range worklogs {
		var err error
		var line string
		if epicSummary {
			line, err = worklog.GetEpicCSVLine()
		} else {
			line, err = worklog.GetIssueCSVLine()
		}
		if err == nil {
			f.WriteString(line + "\n")
		}
	}

	f.Sync()

	return nil
}
