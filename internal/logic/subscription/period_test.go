package subscription

import (
	"fmt"
	"github.com/gogf/gf/v2/os/gtime"
	"testing"
	"time"
	"unibee/utility"
)

func TestPeriod(t *testing.T) {
	t.Run("Test for period", func(t *testing.T) {
		{
			fmt.Println(gtime.Timestamp())
			fmt.Println(gtime.NewFromTimeStamp(gtime.Timestamp()))
			today := gtime.NewFromTimeStamp(time.Date(2024, 1, 31, 12, 0, 0, 0, time.Local).Unix())
			d := today.Day()
			fmt.Println(d)
			fmt.Println(today)

			// nextMonthLastDay
			day2 := today.AddDate(0, 1, -today.Day()+1)
			day2 = day2.AddDate(0, 0, utility.MinInt(d, day2.EndOfMonth().Day())-1)
			fmt.Println(day2)

			// nextTwoMonthLastDay
			day3 := day2.AddDate(0, 1, -day2.Day()+1)
			day3 = day3.AddDate(0, 0, utility.MinInt(d, day3.EndOfMonth().Day())-1)
			fmt.Println(day3)
		}
	})
}
