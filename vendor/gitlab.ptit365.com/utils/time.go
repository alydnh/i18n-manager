package utils

import "time"

const DateLayout string = "2006-01-02"
const DateTimeLayout string = "2006-01-02 03:04:05"

func NormalizeDate(t time.Time) time.Time {
	s := t.Format(DateLayout)
	t, _ = time.Parse(DateLayout, s)
	return t
}

func FirstDayOfMonth(t time.Time) time.Time {
	return NormalizeDate(t.AddDate(0, 0, -t.Day()+1))
}

func ParseDate(dateString string, defaultDate time.Time) time.Time {
	if date, err := Date(dateString); nil != err {
		return defaultDate
	} else {
		return date
	}
}

func Date(date string) (time.Time, error) {
	return time.Parse(DateLayout, date)
}

func Datetime(t string) (time.Time, error) {
	return time.Parse(DateTimeLayout, t)
}

func ToDateString(t time.Time) string {
	return t.Format(DateLayout)
}
