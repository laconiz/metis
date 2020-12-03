package date

import (
	"fmt"
	"time"
)

func Now() Date {
	return Date{time: time.Now()}
}

func Parse(time time.Time) Date {
	return Date{time: time}
}

type Date struct {
	time time.Time
}

func (date *Date) Time() time.Time {
	return date.time
}

func (date Date) String() string {
	return date.time.Format(Layout)
}

func (date Date) MarshalJSON() ([]byte, error) {
	return []byte(`"` + date.time.Format(Layout) + `"`), nil
}

func (date Date) UnmarshalJSON(raw []byte) error {

	if len(raw) != len(Layout)+2 {
		return fmt.Errorf("invalid raw: %s", string(raw))
	}

	time, err := time.Parse(Layout, string(raw))
	if err != nil {
		return err
	}

	date.time = time
	return nil
}

func (date Date) Day() int {
	return date.time.Year()*10000 + (int(date.time.Month())+1)*100 + date.time.Day()
}

const Layout = "2006-01-02 15:04:05.000 Z0700"
