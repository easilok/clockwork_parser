package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/easilok/clockwork_parser/models"
)

const CLOCKWORK_API_WORKLOGS = "https://api.clockwork.report/v1/worklogs"

func GetWorklogs(token string, startDate string, endDate string, project string, author string) ([]models.Worklog, error) {

	requestArgs := "?starting_at=%s&ending_at=%s&expand=issues,authors"

	if len(project) > 0 {
		requestArgs = requestArgs + "&project_keys[]=%s"
	} else {
		// This allows for appending an empty string with Sprintf
		requestArgs = requestArgs + "%s"
	}

	if len(author) > 0 {
		requestArgs = requestArgs + "&user_query=%s"
	} else {
		// This allows for appending an empty string with Sprintf
		requestArgs = requestArgs + "%s"
	}

	requestUrl := fmt.Sprintf(
		CLOCKWORK_API_WORKLOGS+requestArgs,
		startDate,
		endDate,
		project,
		author,
	)

	req, err := http.NewRequest(http.MethodGet, requestUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Token "+token)

	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 299 {
		customError := fmt.Errorf("RESPONSE FAILED STATUS %d", res.StatusCode)
		return nil, customError
	}
	if err != nil {
		return nil, err
	}

	var worklogs []models.Worklog

	err = json.Unmarshal(body, &worklogs)
	if err != nil {
		return nil, err
	}

	return worklogs, nil
}
