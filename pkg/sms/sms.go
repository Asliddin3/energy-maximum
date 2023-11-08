package sms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

type Sms struct {
	Login    string
	Password string
	Url      string
	Port     string
	Sender   string
}
type SmsRequest struct {
	Messages []SmsReqMessage `json:"messages"`
}

type SmsReqMessage struct {
	Recipient string `json:"recipient"`
	MessageID string `json:"message-id"`
	Sms       struct {
		Originator string `json:"originator"`
		Content    []struct {
			Text string `json:"text"`
		} `json:"content"`
	} `json:"sms"`
}

func NewSmsSender() *Sms {
	return &Sms{
		Login:    "",
		Password: "",
		Url:      "",
		Port:     "",
		Sender:   "",
	}
}

func (s *Sms) SendCode(phone string, text string) error {
	id := uuid.NewString()
	id = "mxb" + string([]rune(id)[:10])
	reqData := SmsRequest{
		Messages: []SmsReqMessage{
			{
				Recipient: phone,
				MessageID: id,
				Sms: struct {
					Originator string `json:"originator"`
					Content    []struct {
						Text string `json:"text"`
					} `json:"content"`
				}{
					Originator: s.Sender,
					Content: []struct {
						Text string `json:"text"`
					}{
						{
							Text: text,
						},
					},
				},
			},
		},
	}
	data, err := json.Marshal(&reqData)
	if err != nil {
		return err
	}
	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(data)
	client := &http.Client{}
	req, err := http.NewRequest("POST", s.Url, payloadBuf)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Basic %s:%s", s.Login, s.Password))
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("Content-Type", "application/json")
	res, e := client.Do(req)
	if e != nil {
		return e
	}

	defer res.Body.Close()
	fmt.Println("MESSAGE sended ")
	return nil
}
