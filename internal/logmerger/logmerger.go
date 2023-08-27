package logmerger

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
	"xiangzeli/logmerger/internal/errcode"
)

var (
	metaLineRegexp = regexp.MustCompile(`^(.*)\((\d+)\) (\d\d\d\d[/-]\d\d[/-]\d\d \d\d:\d\d:\d\d)$`)
	cst, _         = time.LoadLocation("Asia/Shanghai")
)

func mergeLogs(mainLogPath string, restLogPaths []string) ([]logItem, statistics, error) {
	stat := statistics{
		Total:   0,
		PerFile: make([]perFileStatistics, 0, len(restLogPaths)+1),
	}

	items, err := openAndRead(mainLogPath)
	if err != nil {
		return nil, statistics{}, err
	}
	stat.PerFile = append(stat.PerFile, perFileStatistics{
		FileName: filepath.Base(mainLogPath),
		Count:    len(items),
		IsMain:   true,
	})

	for _, path := range restLogPaths {
		i, err := openAndRead(path)
		if err != nil {
			return nil, statistics{}, err
		}
		items = append(items, i...)
		stat.PerFile = append(stat.PerFile, perFileStatistics{
			FileName: filepath.Base(path),
			Count:    len(i),
		})
	}

	sort.Sort(logItemByTime(items))
	stat.Total = len(items)
	sort.Sort(mainThenFn(stat.PerFile))

	return items, stat, nil
}

func openAndRead(path string) ([]logItem, error) {
	fd, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fd.Close()
	return readFromFile(fd)
}

func readFromFile(fd *os.File) ([]logItem, error) {
	fileName := fd.Name()
	r := bufio.NewReader(fd)

	var items []logItem
	var last *logItem

	for {
		line, err := r.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, errcode.ErrInputFileError(err)
		}
		line = strings.TrimSpace(line)

		if len(line) == 0 {
			// blank line
			continue
		}

		matches := metaLineRegexp.FindStringSubmatch(line)
		if len(matches) == 4 {
			// meta line
			matches[3] = strings.ReplaceAll(matches[3], "/", "-")
			t, _ := time.ParseInLocation("2006-01-02 15:04:05", matches[3], cst)
			items = append(items, logItem{
				UserID:   matches[2],
				Nickname: matches[1],
				Time:     t,
				Source:   fileName,
			})
			last = &items[len(items)-1]
		} else {
			// content line
			if last == nil {
				return nil, errcode.ErrReadError(fmt.Errorf("in file %s: content line before meta line", fileName))
			}
			last.Content += line + "\n"
		}
	}

	return items, nil
}
