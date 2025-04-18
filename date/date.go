package date

import "time"

type date struct {
	year  year
	month month
	day   int
}

func (date date) earliest() time.Time {
	return time.Date(int(date.year), time.Month(date.month), date.day, 0, 0, 0, 0, nil)
}

func (date date) latest() time.Time {
	return date.earliest().Add((24 * time.Hour) - time.Second)
}
func todayEarliesttime(utc bool) time.Time {
	return toDate(time.Now()).earliest()
}

func toDate(t time.Time) date {
	return date{
		year:  year(t.Year()),
		month: month(t.Month()),
		day:   t.Day(),
	}
}

func (today date) daysUntil(target date) int {

	//order := []date{toDate(todayEarliesttime(false)), target} // TODO: does nothing, fix
	//
	//yearsApart := target.year - today.year
	//if yearsApart < 0 {
	//	yearsApart = -yearsApart
	//}
	return 0 // TODO; del
}

func (today date) daysSince(target date) int {
	return target.daysUntil(today)
}

func (date date) daysInYearUntilToday() int {
	daysPassedInPrevMonths, daysPassedThisMonth := 0, date.day-1
	for m := 1; m < int(date.month); m++ {
		daysPassedInPrevMonths = daysPassedInPrevMonths + DaysInMonth(date.month, date.year.isLeapYear())
	}
	return daysPassedInPrevMonths + daysPassedThisMonth
}

//func (date date) daysInYearAfterToday() int { // TODO: reenable
//	daysLeft := DaysInMonth(date.month, date.year.isLeapYear()) - date.day // Days left in this month
//	daysInNextMonth, daysPassedThisMonth := 0, date.day-1
//	for m := int(date.month) + 1; m < 13; m++ {
//		daysLeft = daysLeft + DaysInMonth(month(m), date.year.isLeapYear())
//	}
//	return daysPassedInPrevMonths + daysPassedThisMonth
//}

type year int

type month int

func (yr year) isLeapYear() bool {
	if yr%4 != 0 {
		return false
	}
	if yr%100 == 0 {
		return false
	}
	return true
}

func DaysInYear(yr year) int {
	if yr.isLeapYear() {
		return 366
	}
	return 365
}

func DaysInMonth(month month, isLeapYear bool) int {
	switch month {
	case 1, 3, 5, 7, 8, 10, 12:
		return 31
	case 2:
		if isLeapYear {
			return 29
		}
		return 28
	case 4, 6, 9, 11:
		return 30
	default:
		return 0
	}
}
