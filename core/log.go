package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/fatih/color"
)

const (
	FATAL     = 5
	ERROR     = 4
	IMPORTANT = 3
	WARN      = 2
	INFO      = 1
	DEBUG     = 0
)

var LogColors = map[int]*color.Color{
	FATAL:     color.New(color.FgRed).Add(color.Bold),
	ERROR:     color.New(color.FgRed),
	WARN:      color.New(color.FgYellow),
	IMPORTANT: color.New(),
	DEBUG:     color.New(color.Faint),
}

type Logger struct {
	sync.Mutex

	debug  bool
	silent bool
}

type DiscMessage struct {
	Embeds []struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		URL         string `json:"url"`
		Color       int    `json:"color"`
	} `json:"embed"`
}

func (l *Logger) SetDebug(d bool) {
	l.debug = d
}

func (l *Logger) SetSilent(d bool) {
	l.silent = d
}

func (l *Logger) Log(level int, format string, args ...interface{}) {
	l.Lock()
	defer l.Unlock()

	if level == DEBUG && !l.debug {
		return
	}

	if l.silent && level < IMPORTANT {
		return
	}

	if c, ok := LogColors[level]; ok {
		c.Printf(format+"\n", args...)
	} else {
		fmt.Printf(format+"\n", args...)
	}

	// if level > INFO && session.Config.SlackWebhook != "" {
	// 	values := map[string]string{"text": fmt.Sprintf(format+"\n", args...)}
	// 	jsonValue, _ := json.Marshal(values)
	// 	http.Post(session.Config.SlackWebhook, "application/json", bytes.NewBuffer(jsonValue))
	// }

	if level == FATAL {
		os.Exit(1)
	}
}

func (l *Logger) Fatal(format string, args ...interface{}) {
	l.Log(FATAL, format, args...)
}

func (l *Logger) Error(format string, args ...interface{}) {
	l.Log(ERROR, format, args...)
}

func (l *Logger) Warn(format string, args ...interface{}) {
	l.Log(WARN, format, args...)
}

func (l *Logger) Important(format string, args ...interface{}) {
	l.Log(IMPORTANT, format, args...)
}

func (l *Logger) Info(format string, args ...interface{}) {
	l.Log(INFO, format, args...)
}

func (l *Logger) Debug(format string, args ...interface{}) {
	l.Log(DEBUG, format, args...)
}

func (l *Logger) Discord(title string, description string, URL string,) {
	var values DiscMessage
	values.Embeds.Title = title
	values.Embeds.Description = "```" + description + "```"
	values.Embeds.URL = URL
	values.Embeds.Color = 6545520
	jsonValue, err := json.Marshal(values)
	if err != nil {
		fmt.Println(err.Error())
	}
	http.Post(session.Config.SlackWebhook, "application/json", bytes.NewBuffer(jsonValue))
}
