package services

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type JiraService struct {
	auth     string
	endpoint string
	headers  *map[string]string
	jql      *map[string]interface{}
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

func NewJiraService(auth, endpoint, jql string) *JiraService {
	return &JiraService{
		auth:     auth,
		endpoint: endpoint,
		jql: &map[string]interface{}{
			"jql":        jql,
			"maxResults": 20,
		},
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
		"Authorization": "Basic " + js.auth,
	}
	return js
}

func (js *JiraService) SetJql(jql *map[string]interface{}) *JiraService {
	js.jql = jql
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
		http.MethodPost, js.endpoint, bytes.NewBuffer(*body),
	)
	if err != nil {
		return nil, err
	}
	for key, value := range *js.headers {
		req.Header.Add(key, value)
	}
	return http.DefaultClient.Do(req)
}
