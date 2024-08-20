package utils

import "time"

func Timestamp() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func TimeNowStringVietNam() string {

	location, _ := time.LoadLocation("Asia/Bangkok")

	// this should give you time in location
	t := time.Now().In(location).Format("2006-01-02T15:04:05.000Z")

	return t
}

func TimeNowVietNam() time.Time {

	location, _ := time.LoadLocation("Asia/Bangkok")

	// this should give you time in location
	t := time.Now().In(location)

	return t
}

func StringToTimestamp(timestamp string) (time.Time, error) {
	return time.Parse("2006-01-02T15:04:05Z", timestamp)
}

func TimestampToString(timestamp time.Time) string {
	return timestamp.Format("2006-01-02T15:04:05")
}
