package model

import (
	"encoding/json"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/pulpfree/auth"
)

// Logger struct
type Logger struct {
	Body *bodyStruct
}

type bodyStruct struct {
	AppID    string      `json:"appID"`
	AppType  string      `json:"appType"`
	Context  string      `json:"context"`
	Messages []msgStruct `json:"messages"`
	SentAt   string      `json:"sentAt"`
}

type msgStruct struct {
	Level    int    `json:"level"`
	Message  string `json:"message"`
	Location string `json:"location"`
}

type record struct {
	AppID      bson.ObjectId `bson:"appID"`
	AppType    string        `bson:"appType"`
	Context    string        `bson:"context"`
	Level      int           `bson:"level"`
	Location   interface{}   `bson:"location"`
	Message    interface{}   `bson:"message"`
	ReceivedAt time.Time     `bson:"receivedAt"`
	SentAt     time.Time     `bson:"sentAt"`
}

type records struct {
	Records []record
}

// Record function
func (l *Logger) Record(db *auth.DB) error {

	const longForm = "2006-01-02T15:04:05Z"
	sentTime, _ := time.Parse(longForm, l.Body.SentAt)

	col := db.DB.C("logs")
	ts := time.Now()
	appID := bson.ObjectIdHex(l.Body.AppID)

	for _, m := range l.Body.Messages {

		var msg interface{}
		if err := json.Unmarshal([]byte(m.Message), &msg); err != nil {
			return err
		}

		var location interface{}
		if m.Location != "" {
			if err := json.Unmarshal([]byte(m.Location), &location); err != nil {
				return err
			}
		}

		r := &record{
			AppID:      appID,
			AppType:    l.Body.AppType,
			Context:    l.Body.Context,
			Message:    msg,
			Level:      m.Level,
			Location:   location,
			ReceivedAt: ts,
			SentAt:     sentTime,
		}
		if err := col.Insert(r); err != nil {
			return err
		}
	}
	return nil
}
