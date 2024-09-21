package subscription

import (
	"fmt"
	"github.com/gogf/gf/v2/os/gtime"
	"testing"
	"time"
	"unibee/utility"
)

func TestPeriod(t *testing.T) {
	//ctx := context.Background()
	t.Run("Test for period", func(t *testing.T) {
		//time := gtime.NewFromTimeStamp(1717137911)
		//fmt.Println(time.Month())
		//fmt.Println(time)
		//fmt.Println(time.StartOfMonth())
		//fmt.Println(time.Timestamp() - time.StartOfMonth().Timestamp())
		//time = time.AddDate(0, 1, 0)
		//fmt.Println(time.Month())
		//fmt.Println(time)
		//fmt.Println(time.StartOfMonth())
		//fmt.Println(time.Timestamp() - time.StartOfMonth().Timestamp())
		//today := time.Date
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
		//{
		//	today := time.Date(2022, 3, 31, 0, 0, 0, 0, time.Local)
		//	d := today.AddDate(0, -1, 0)
		//	fmt.Println(d.Format("20060102"))
		//	// 20220303
		//
		//	today = time.Date(2022, 3, 31, 0, 0, 0, 0, time.Local)
		//	d = today.AddDate(0, 1, 0)
		//	fmt.Println(d.Format("20060102"))
		//	// 20220501
		//
		//	today = time.Date(2022, 10, 31, 0, 0, 0, 0, time.Local)
		//	d = today.AddDate(0, -1, 0)
		//	fmt.Println(d.Format("20060102"))
		//	// 20221001
		//
		//	today = time.Date(2022, 10, 31, 0, 0, 0, 0, time.Local)
		//	d = today.AddDate(0, 1, 0)
		//	fmt.Println(d.Format("20060102"))
		//	// 20221201
		//}
	})
}
