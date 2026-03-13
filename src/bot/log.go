package bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

type line struct {
	Message   string
	Timestamp time.Time
	Severity  string
}

func (e *Bot) Write(p []byte) (int, error) {
	l := &line{
		Message:   string(p),
		Timestamp: time.Now().UTC(),
	}
	e.IRC().Cmd.Message(e.cfg.LogChan, fmt.Sprintf("[%s] %s: %s", l.Timestamp, l.Severity, l.Message))

	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(l)
	data := buf.Bytes()
	n, err := os.Stdout.Write(data)
	if err != nil {
		return n, err
	}
	if n != len(data) {
		return n, io.ErrShortWrite
	}
	return len(p), nil
}
