package logmerger

import (
	"sort"
	"time"
)

type LogItem struct {
	UserID   string
	Nickname string
	Time     time.Time
	Content  string

	Source string
}

type LogItemByTime []LogItem

func (a LogItemByTime) Len() int           { return len(a) }
func (a LogItemByTime) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a LogItemByTime) Less(i, j int) bool { return a[i].Time.Before(a[j].Time) }

// for compile time check
var _ sort.Interface = LogItemByTime{}
