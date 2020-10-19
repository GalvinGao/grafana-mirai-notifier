package mirai

import (
	"bytes"
	"encoding/json"
	"github.com/Logiase/gomirai/message"
	"io"
	"net/http"
	"path"
	"time"
)

const sessionTimeout = time.Minute * 30

type Session struct {
	sessionKey string
	issuedAt time.Time
}

func (s Session) IsValid() bool {
	return s.sessionKey != "" && s.issuedAt.Add(sessionTimeout).After(time.Now())
}

type Bot struct {
	authKey string
	relayAddress string

	session *Session

	client *http.Client
}

func NewBot(authKey string, relayAddress string) *Bot {
	return &Bot{
		authKey: authKey,
		relayAddress: relayAddress,

		client: &http.Client{
			Timeout: time.Minute,
		},
	}
}

func (b Bot) sendRawRequest(method string, requestPath string, body io.Reader) (*http.Response, error) {
	url := path.Join(b.relayAddress, requestPath)
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return &http.Response{}, err
	}
	resp, err := b.client.Do(req)
	if err != nil {
		return &http.Response{}, err
	}
	return resp, nil
}

func jsoned(v interface{}) io.Reader {
	buf := bytes.NewBufferString("")
	if err := json.NewEncoder(buf).Encode(v); err != nil {
		panic(err)
	}
	return buf
}

type AuthResponse struct {
	Session string `json:"session"`
}

func (b Bot) getSessionKey() string {
	if b.session.IsValid() {
		return b.session.sessionKey
	}

	resp, err := b.sendRawRequest("POST", "/auth", jsoned(map[string]string{
		"authKey": b.authKey,
	}))
	if err != nil {
		panic(err)
	}

	var authResponse AuthResponse
	err = json.NewDecoder(resp.Body).Decode(&authResponse)
	if err != nil {
		panic(err)
	}

	b.session = &Session{
		sessionKey: authResponse.Session,
		issuedAt:   time.Now(),
	}

	return b.session.sessionKey
}

func (b Bot) SendGroupMessage(groupId int, messageChain message.Message)  {

}