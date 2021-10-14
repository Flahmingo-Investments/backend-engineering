package commons

import (
	"fmt"
	"time"
)

func TestingTime(textTime string) {
	TestingNow, _ = time.Parse(time.RFC3339, fmt.Sprintf("2020-01-01T%s:00.000Z", textTime))
}
