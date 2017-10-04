package entry

import (
	"strings"
	"time"
)

var logTimeFormat = "2006-01-02T15:04:05.000000Z"

type Entry struct {
	RequestTime time.Time
	Line        string
}

func ParseALBLog(line string) (*Entry, error) {
	s := strings.Split(line, " ")
	t, err := time.Parse(logTimeFormat, s[1])
	if err != nil {
		return nil, err
	}
	return &Entry{
		RequestTime: t,
		Line:        line,
	}, nil
}

func ParseELBLog(line string) (*Entry, error) {
	s := strings.Split(line, " ")
	t, err := time.Parse(logTimeFormat, s[0])
	if err != nil {
		return nil, err
	}
	return &Entry{
		RequestTime: t,
		Line:        line,
	}, nil
}
