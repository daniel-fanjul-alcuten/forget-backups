package main

import (
	"fmt"
	"sort"
	"testing"
	"time"
)

func TestParse(t *testing.T) {

	format := "20060102150405"

	if _, err := parse("foo"); err == nil {
		t.Error(err)
	}

	if v, err := parse("12340102030405"); err != nil {
		t.Error(err)
	} else if s := v.time.Format(format); s != "12340102030405" {
		t.Error(s)
	}

	if v, err := parse("12340102030405.tar"); err != nil {
		t.Error(err)
	} else if s := v.time.Format(format); s != "12340102030405" {
		t.Error(s)
	}

	if v, err := parse("/12340102030405.tar"); err != nil {
		t.Error(err)
	} else if s := v.time.Format(format); s != "12340102030405" {
		t.Error(s)
	}

	if v, err := parse("foo/12340102030405.tar"); err != nil {
		t.Error(err)
	} else if s := v.time.Format(format); s != "12340102030405" {
		t.Error(s)
	}

	if v, err := parse("foo/12340102030405.bundle.git"); err != nil {
		t.Error(err)
	} else if s := v.time.Format(format); s != "12340102030405" {
		t.Error(s)
	}
}

func TestMomentByMinute(t *testing.T) {

	m := moment{2007, 5, 18, 19, 20}
	if n := m.byMinute(); n != (moment{2007, 5, 18, 19, 19}) {
		t.Error(n)
	}

	m = moment{2007, 5, 18, 19, 0}
	if n := m.byMinute(); n != (moment{2007, 5, 18, 18, 59}) {
		t.Error(n)
	}

	m = moment{2007, 5, 1, 0, 0}
	if n := m.byMinute(); n != (moment{2007, 4, 30, 23, 59}) {
		t.Error(n)
	}

	m = moment{2007, 1, 1, 0, 0}
	if n := m.byMinute(); n != (moment{2006, 12, 31, 23, 59}) {
		t.Error(n)
	}
}

func TestMomentByHour(t *testing.T) {

	m := moment{2007, 5, 18, 19, 20}
	if n := m.byHour(); n != (moment{2007, 5, 18, 18, 20}) {
		t.Error(n)
	}

	m = moment{2007, 5, 18, 0, 20}
	if n := m.byHour(); n != (moment{2007, 5, 17, 23, 20}) {
		t.Error(n)
	}

	m = moment{2007, 5, 1, 0, 20}
	if n := m.byHour(); n != (moment{2007, 4, 30, 23, 20}) {
		t.Error(n)
	}

	m = moment{2007, 1, 1, 0, 20}
	if n := m.byHour(); n != (moment{2006, 12, 31, 23, 20}) {
		t.Error(n)
	}
}

func TestMomentByDay(t *testing.T) {

	m := moment{2007, 5, 18, 19, 20}
	if n := m.byDay(); n != (moment{2007, 5, 17, 19, 20}) {
		t.Error(n)
	}

	m = moment{2007, 5, 1, 19, 20}
	if n := m.byDay(); n != (moment{2007, 4, 30, 19, 20}) {
		t.Error(n)
	}

	m = moment{2007, 1, 1, 19, 20}
	if n := m.byDay(); n != (moment{2006, 12, 31, 19, 20}) {
		t.Error(n)
	}
}

func TestMomentByMonth(t *testing.T) {

	m := moment{2007, 5, 18, 19, 20}
	if n := m.byMonth(); n != (moment{2007, 4, 18, 19, 20}) {
		t.Error(n)
	}

	m = moment{2007, 1, 18, 19, 20}
	if n := m.byMonth(); n != (moment{2006, 12, 18, 19, 20}) {
		t.Error(n)
	}
}

func TestMomentByYear(t *testing.T) {

	m := moment{2007, 1, 18, 19, 20}
	if n := m.byYear(); n != (moment{2006, 1, 18, 19, 20}) {
		t.Error(n)
	}
}

func TestNumDaysInMonth(t *testing.T) {

	type st struct {
		year     int
		month    time.Month
		expected int
	}

	for _, s := range []st{
		{2007, 1, 31},
		{2007, 2, 28},
		{2008, 2, 29},
		{2007, 3, 31},
		{2007, 4, 30},
		{2007, 5, 31},
		{2007, 6, 30},
		{2007, 7, 31},
		{2007, 8, 31},
		{2007, 9, 30},
		{2007, 10, 31},
		{2007, 11, 30},
		{2007, 12, 31},
	} {
		if n := numDaysInMonth(s.year, s.month); n != s.expected {
			t.Errorf("%v, %v, %v: %v", s.year, s.month, s.expected, n)
		}
	}
}

func TestSelectFiles_Zero(t *testing.T) {
	var fs filesort
	fs = append(fs, file{"foo", time.Date(2007, 5, 18, 19, 0, 0, 0, time.UTC)})
	s := make(map[string]file)
	selectFiles(s, fs, hourly{}, 0)
	if l := len(s); l != 0 {
		t.Fatal(l)
	}
}

func TestSelectFiles_One(t *testing.T) {
	var fs filesort
	s := make(map[string]file)
	selectFiles(s, fs, hourly{}, 1)
	if l := len(s); l != 0 {
		t.Fatal(l)
	}
}

func TestSelectFiles_Short(t *testing.T) {

	const nfiles = 5
	fs := make(filesort, 0, nfiles)
	year := 2007
	month := 5
	day := 18
	for hour := 19; hour < 24; hour++ {
		name := fmt.Sprintf("%04d/%02d/%02d/%02d", year, int(month), day, hour)
		fs = append(fs, file{name, time.Date(year, time.Month(month), day, hour, 1, 2, 3, time.UTC)})
	}
	if l := len(fs); l != nfiles {
		t.Error(l)
	}
	sort.Sort(&fs)

	const sfiles = 5
	s := make(map[string]file, sfiles)
	selectFiles(s, fs, hourly{}, 28)
	selectFiles(s, fs, daily{}, 8)
	selectFiles(s, fs, weekly{}, 5)
	selectFiles(s, fs, monthly{}, 13)
	selectFiles(s, fs, yearly{}, 3)
	if l := len(s); l != sfiles {
		t.Fatal(l)
	}

	actual := make([]string, 0, sfiles)
	for k, _ := range s {
		actual = append(actual, k)
	}
	sort.Strings(actual)

	expected := [sfiles]string{
		// hourly
		"2007/05/18/19",
		"2007/05/18/20",
		"2007/05/18/21",
		"2007/05/18/22",
		"2007/05/18/23",
	}
	for k, a := range actual {
		if e := expected[k]; a != e {
			t.Errorf("%v %v: %v", k, e, a)
		}
	}
}

func TestSelectFiles_Long(t *testing.T) {

	const nfiles = 315576
	fs := make(filesort, 0, nfiles)
	for year := 2007; year < 2007+4; year++ {
		for month := time.January; month < 13; month++ {
			nd := numDaysInMonth(year, month)
			for day := 1; day <= nd; day++ {
				for hour := 0; hour < 24; hour++ {
					for minute := 0; minute < 60; minute += 7 {
						name := fmt.Sprintf("%04d/%02d/%02d/%02d/%02d", year, int(month), day, hour, minute)
						fs = append(fs, file{name, time.Date(year, month, day, hour, minute, 2, 3, time.UTC)})
					}
				}
			}
		}
	}
	if l := len(fs); l != nfiles {
		t.Error(l)
	}
	sort.Sort(&fs)

	const sfiles = 59
	s := make(map[string]file, sfiles)
	selectFiles(s, fs, minutely{}, 70)
	selectFiles(s, fs, hourly{}, 28)
	selectFiles(s, fs, daily{}, 8)
	selectFiles(s, fs, weekly{}, 5)
	selectFiles(s, fs, monthly{}, 13)
	selectFiles(s, fs, yearly{}, 3)
	if l := len(s); l != sfiles {
		t.Fatal(l)
	}

	actual := make([]string, 0, sfiles)
	for k, _ := range s {
		actual = append(actual, k)
	}
	sort.Strings(actual)

	expected := [sfiles]string{
		// yearly
		"2008/12/31/23/56",
		// monthly
		"2009/12/31/23/56",
		"2010/01/31/23/56",
		"2010/02/28/23/56",
		"2010/03/31/23/56",
		"2010/04/30/23/56",
		"2010/05/31/23/56",
		"2010/06/30/23/56",
		"2010/07/31/23/56",
		"2010/08/31/23/56",
		"2010/09/30/23/56",
		"2010/10/31/23/56",
		"2010/11/30/23/56",
		// weekly
		"2010/12/04/23/56",
		"2010/12/11/23/56",
		"2010/12/18/23/56",
		// daily
		"2010/12/24/23/56",
		"2010/12/25/23/56",
		"2010/12/26/23/56",
		"2010/12/27/23/56",
		"2010/12/28/23/56",
		"2010/12/29/23/56",
		// hourly
		"2010/12/30/20/56",
		"2010/12/30/21/56",
		"2010/12/30/22/56",
		"2010/12/30/23/56",
		"2010/12/31/00/56",
		"2010/12/31/01/56",
		"2010/12/31/02/56",
		"2010/12/31/03/56",
		"2010/12/31/04/56",
		"2010/12/31/05/56",
		"2010/12/31/06/56",
		"2010/12/31/07/56",
		"2010/12/31/08/56",
		"2010/12/31/09/56",
		"2010/12/31/10/56",
		"2010/12/31/11/56",
		"2010/12/31/12/56",
		"2010/12/31/13/56",
		"2010/12/31/14/56",
		"2010/12/31/15/56",
		"2010/12/31/16/56",
		"2010/12/31/17/56",
		"2010/12/31/18/56",
		"2010/12/31/19/56",
		"2010/12/31/20/56",
		"2010/12/31/21/56",
		// minutely
		"2010/12/31/22/49",
		"2010/12/31/22/56",
		"2010/12/31/23/00",
		"2010/12/31/23/07",
		"2010/12/31/23/14",
		"2010/12/31/23/21",
		"2010/12/31/23/28",
		"2010/12/31/23/35",
		"2010/12/31/23/42",
		"2010/12/31/23/49",
		"2010/12/31/23/56",
	}
	for k, a := range actual {
		if e := expected[k]; a != e {
			t.Errorf("%v %v: %v", k, e, a)
		}
	}
}
