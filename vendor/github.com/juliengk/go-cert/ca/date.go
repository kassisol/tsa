package ca

import (
	"fmt"
	"strconv"
	"time"
)

type CertDate struct {
	Now    time.Time
	Expire time.Time
}

func CreateDate(month int) CertDate {
	now := time.Now().UTC()

	return CertDate{
		Now:    now,
		Expire: now.AddDate(0, month, 0),
	}
}

func (cd *CertDate) ExpireDateString() string {
	year := strconv.Itoa(cd.Expire.Year())
	month := convertNumber(int(cd.Expire.Month()))
	day := convertNumber(cd.Expire.Day())

	return fmt.Sprintf("%s-%s-%s", year, month, day)
}

func ExpireDiffDays(notafter time.Time) int {
	days := 1

	now := time.Now().UTC()
	diff := notafter.Sub(now)

	hours := int(diff.Hours())

	if hours > 24 {
		days = hours / 24
	}

	return days
}

func DatabaseDateTimeFormat(datetime time.Time) string {
	year := strconv.Itoa(datetime.Year())
	month := convertNumber(int(datetime.Month()))
	day := convertNumber(datetime.Day())
	hour := convertNumber(datetime.Hour())
	minute := convertNumber(datetime.Minute())
	second := convertNumber(datetime.Second())

	return fmt.Sprintf("%s%s%s%s%s%sZ", year, month, day, hour, minute, second)
}

func convertNumber(number int) string {
	result := strconv.Itoa(number)

	return fmt.Sprintf("%02s", result)
}
