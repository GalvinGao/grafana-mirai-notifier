package main

import (
	"bytes"
	"github.com/Logiase/gomirai/message"
	"github.com/labstack/echo"
	"net/http"
	"text/template"
)

const MessageTemplate = `> New Alert State from Grafana
> State: {{.State}}
> {{.Message}}

Details:
  - Source: {{.RuleName}} (id:{{.RuleID}}, panel:{{.PanelID}})
  - For more, visit {{.RuleURL}}`

type Tag map[string]string

type GrafanaWebhookRequest struct {
	DashboardID int `json:"dashboardId"`
	EvalMatches []struct {
		Value  int    `json:"value"`
		Metric string `json:"metric"`
		Tags   Tag `json:"tags"`
	} `json:"evalMatches"`
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

func webhookHandler(c echo.Context) error {
	r := new(GrafanaWebhookRequest)
	if err := c.Bind(r); err != nil {
		return responseError(http.StatusBadRequest, "bind request failed", err)
	}
	if err := c.Validate(r); err != nil {
		return responseError(http.StatusBadRequest, "request validation failed", err)
	}

	msg := bytes.NewBufferString("")
	err := template.Must(template.New("message").Parse(MessageTemplate)).Execute(msg, r)
	if err != nil {
		return responseError(http.StatusInternalServerError, "failed to format message", err)
	}

	if _, err := Bot.SendGroupMessage(Conf.QQ.Group, 0, message.PlainMessage(msg.String())); err != nil {
		return responseError(http.StatusInternalServerError, "failed to send message", err)
	}
	if _, err := Bot.SendGroupMessage(Conf.QQ.Group, 0, message.ImageMessage("url", r.ImageURL)); err != nil {
		return responseError(http.StatusInternalServerError, "failed to send message", err)
	}

	return c.NoContent(http.StatusOK)
}
