package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type file struct {
	name string
	time time.Time
}

func parse(arg string) (f file, err error) {

	// remove path elements
	str := filepath.Base(arg)
	// remove extension
	str = str[:len(str)-len(filepath.Ext(str))]

	// parse for $(date +%Y%m%d%H%M%S)
	var t time.Time
	if t, err = time.Parse("20060102150405", str); err != nil {
		return
	}

	f = file{arg, t}
	return
}

type filesort []file

func (fs filesort) Len() int {
	return len(fs)
}

func (fs filesort) Less(i, j int) bool {
	return fs[i].time.After(fs[j].time)
}

func (fs *filesort) Swap(i, j int) {
	(*fs)[i], (*fs)[j] = (*fs)[j], (*fs)[i]
}

type moment struct {
	year  int
	month time.Month
	day   int
	hour  int
}

func (m moment) time() time.Time {
	return time.Date(m.year, m.month, m.day, m.hour, 0, 0, 0, time.UTC)
}

func (m moment) byHour() moment {
	if m.hour > 0 {
		m.hour--
		return m
	}
	m.hour = 23
	return m.byDay()
}

func (m moment) byDay() moment {
	if m.day > 1 {
		m.day--
		return m
	}
	m = m.byMonth()
	m.day = numDaysInMonth(m.year, m.month)
	return m
}

func (m moment) byMonth() moment {
	if m.month > 1 {
		m.month--
		return m
	}
	m.month = 12
	return m.byYear()
}

func (m moment) byYear() moment {
	m.year--
	return m
}

func numDaysInMonth(year int, month time.Month) int {
	switch month {
	case 2:
		if year%4 == 0 && (year%100 != 0 || year%400 == 0) {
			return 29
		}
		return 28
	case 4:
		return 30
	case 6:
		return 30
	case 9:
		return 30
	case 11:
		return 30
	}
	return 31
}

type policy interface {
	first(moment) moment
	next(moment) moment
}

type hourly struct{}

func (hourly) first(m moment) moment {
	return m
}

func (hourly) next(m moment) moment {
	return m.byHour()
}

type daily struct{}

func (daily) first(m moment) moment {
	m.hour = 0
	return m
}

func (daily) next(m moment) moment {
	return m.byDay()
}

type weekly struct{}

func (weekly) first(m moment) moment {
	m.hour = 0
	for d := m.time().Weekday() - time.Sunday; d > 0; d-- {
		m = m.byDay()
	}
	return m
}

func (weekly) next(m moment) moment {
	for i := 0; i < 7; i++ {
		m = m.byDay()
	}
	return m
}

type monthly struct{}

func (monthly) first(m moment) moment {
	m.hour = 0
	m.day = 1
	return m
}

func (monthly) next(m moment) moment {
	return m.byMonth()
}

type yearly struct{}

func (yearly) first(m moment) moment {
	m.hour = 0
	m.day = 1
	m.month = 1
	return m
}

func (yearly) next(m moment) moment {
	return m.byYear()
}

func selectFiles(s map[string]file, fs []file, p policy, q int) {

	switch {
	case len(fs) == 0:
		return
	case q < 1:
		return
	}

	f := fs[0]
	var m moment
	m.year, m.month, m.day = f.time.Date()
	m.hour = f.time.Hour()
	m = p.first(m)

	for i := 0; i < q; i++ {

		next := m.time()
		if !next.After(f.time) {
			s[f.name] = f
			if len(fs) == 1 {
				return
			}
			f, fs = fs[1], fs[1:]
		}
		for !next.After(f.time) {
			if len(fs) == 1 {
				return
			}
			f, fs = fs[1], fs[1:]
		}

		m = p.next(m)
	}
}

func main() {

	log.SetFlags(0)

	const factor = 3
	hours := flag.Int("hours", factor*24, "it keeps the latest backup for each hour, for the given number of hours")
	days := flag.Int("days", factor*7, "it keeps the latest backup for each day, for the given number of days")
	weeks := flag.Int("weeks", factor*4, "it keeps the latest backup for each week, for the given number of weeks")
	months := flag.Int("months", factor*12, "it keeps the latest backup for each month, for the given number of months")
	years := flag.Int("years", factor*5, "it keeps the latest backup for each year, for the given number of years")
	keep := flag.Bool("keep", false, "it shows the backups to keep instead of the ones to forget")
	flag.Parse()

	var fs filesort
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
		text := strings.TrimSpace(scanner.Text())
		if f, err := parse(text); err == nil {
			fs = append(fs, f)
		}
	}
	sort.Sort(&fs)

	s := make(map[string]file)
	selectFiles(s, fs, hourly{}, *hours)
	selectFiles(s, fs, daily{}, *days)
	selectFiles(s, fs, weekly{}, *weeks)
	selectFiles(s, fs, monthly{}, *months)
	selectFiles(s, fs, yearly{}, *years)

	for _, f := range fs {
		if _, ok := s[f.name]; ok == *keep {
			fmt.Println(f.name)
		}
	}
}
