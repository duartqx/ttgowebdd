package services

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type JiraService struct {
	base64EncodedAuth string
	searchEndpoint    string
	assignee          string
	headers           *map[string]string
	jql               *map[string]interface{}
}

type jiraName struct {
	Name string `json:"name"`
}

type jiraResponse struct {
	Issues []struct {
		Key    string `json:"name"`
		Fields struct {
			IssueType jiraName   `json:"issuetype"`
			Sprint    []jiraName `json:"customfield_10020"`
			Status    jiraName   `json:"status"`
			Summary   string     `json:"summary"`
		} `json:"fields"`
	} `json:"issues"`
}

func NewJiraService(auth, endpoint, assignee string) *JiraService {
	return &JiraService{
		base64EncodedAuth: auth,
		searchEndpoint:    endpoint,
		assignee:          assignee,
	}
}

func (js *JiraService) SetHeaders(headers *map[string]string) *JiraService {
	js.headers = headers
	return js
}

func (js *JiraService) SetDefaultHeaders() *JiraService {
	js.headers = &map[string]string{
		"Accept":        "application/json",
		"Content-Type":  "application/json",
		"Authorization": "Basic" + js.base64EncodedAuth,
	}
	return js
}

func (js *JiraService) SetJql(jql *map[string]interface{}) *JiraService {
	js.jql = jql
	return js
}

func (js *JiraService) SetDefaultJql() *JiraService {

	js.jql = &map[string]interface{}{
		"jql": `
			project = "AJ" and assignee = ` + js.assignee + ` 
			and status = Open ORDER BY created DESC
	`,
		"maxResults": 20,
	}

	return js
}

func (js JiraService) Query() (*jiraResponse, error) {
	body, err := json.Marshal(js.jql)
	if err != nil {
		return nil, err
	}

	resp, err := js.post(&body)
	if err != nil {
		return nil, err
	}

	var jiraResponse jiraResponse
	if err := json.NewDecoder(resp.Body).Decode(&jiraResponse); err != nil {
		return nil, err
	}

	return &jiraResponse, nil
}

func (js JiraService) post(body *[]byte) (*http.Response, error) {
	req, err := http.NewRequest(
		http.MethodPost, js.searchEndpoint, bytes.NewBuffer(*body),
	)
	if err != nil {
		return nil, err
	}
	for key, value := range *js.headers {
		req.Header.Add(key, value)
	}
	return http.DefaultClient.Do(req)
}
