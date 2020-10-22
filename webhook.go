package main

import (
	"bytes"
	"fmt"
	"github.com/Logiase/gomirai/message"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"net/http"
	"path"
	"text/template"
)

type Tag map[string]string

type GrafanaWebhookRequest struct {
	DashboardID int `json:"dashboardId"`
	//EvalMatches []struct {
	//	Value  int    `json:"value"`
	//	Metric string `json:"metric"`
	//	Tags   Tag `json:"tags"`
	//} `json:"evalMatches"`
	ImageURL string `json:"imageUrl"`
	Message  string `json:"message" validate:"required"`
	OrgID    int    `json:"orgId"`
	PanelID  int    `json:"panelId" validate:"required"`
	RuleID   int    `json:"ruleId" validate:"required"`
	RuleName string `json:"ruleName" validate:"required"`
	RuleURL  string `json:"ruleUrl" validate:"required"`
	State    string `json:"state" validate:"required"`
	Tags     Tag `json:"tags"`
	Title    string `json:"title" validate:"required"`
}

func readTemplateFile(name string) string {
	t, _ := ioutil.ReadFile(path.Join("templates", fmt.Sprintf("%s.tmpl", name)))
	return string(t)
}

func webhookHandler(c echo.Context) error {
	r := new(GrafanaWebhookRequest)
	if err := c.Bind(r); err != nil {
		return responseError(http.StatusBadRequest, "bind request failed", err)
	}
	if err := c.Validate(r); err != nil {
		return responseError(http.StatusBadRequest, "request validation failed", err)
	}

	var messageTemplate string

	switch r.State {
	case GrafanaWebhookStateOK:
		messageTemplate = readTemplateFile("state-ok")
	case GrafanaWebhookStateAlerting:
		messageTemplate = readTemplateFile("state-alerting")
	default:
		return c.NoContent(http.StatusOK)
	}

	msg := bytes.NewBufferString("")
	err := template.Must(template.New("message").Parse(messageTemplate)).Execute(msg, r)
	if err != nil {
		return responseError(http.StatusInternalServerError, "failed to format message", err)
	}


	err = botSendGroupMessage(Conf.QQ.Group, 0, message.PlainMessage(msg.String()), message.ImageMessage("url", r.ImageURL))
	if err != nil {
		return responseError(http.StatusInternalServerError, "failed to send message after consecutive retries", err)
	}

	return c.NoContent(http.StatusAccepted)
}
