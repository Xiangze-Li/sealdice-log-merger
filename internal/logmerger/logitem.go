package logmerger

import (
	"sort"
	"time"
)

type logItem struct {
	UserID   string
	Nickname string
	Time     time.Time
	Content  string

	Source string
}

type logItemByTime []logItem

func (a logItemByTime) Len() int           { return len(a) }
func (a logItemByTime) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a logItemByTime) Less(i, j int) bool { return a[i].Time.Before(a[j].Time) }

// for compile time check
var _ sort.Interface = logItemByTime{}
