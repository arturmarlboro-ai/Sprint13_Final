package api

import (
	"errors"
	"fmt"
	"log"
	"slices"
	"strconv"
	"strings"
	"time"
)

func NextDate(now time.Time, dstart, repeat string) (string, error) {
	fmt.Println(now.Format("20060102"), dstart, repeat)

	taskTime, err := time.Parse("20060102", dstart)
	if err != nil {
		log.Print("dstart parse err: ", err)
		return "", fmt.Errorf("dstart parse err %w", err)
	}

	var repModeTimes []string

	repModeTimes = strings.SplitN(repeat, " ", 3) //m 15,16 10,11

	fmt.Println(repModeTimes)

	if len(repModeTimes) > 3 {
		log.Print("task repeat wrong format")
	}

	var repTimesString []string

	repMode := strings.ToLower(strings.TrimSpace(repModeTimes[0]))

	if len(repModeTimes) != 1 {
		repTimesString = strings.Split(repModeTimes[1], ",")
	}

	var repTimes []int

	for _, n := range repTimesString {
		v, err := strconv.Atoi(n)
		if err != nil {
			log.Print(err)
			return "", fmt.Errorf("%w", err)
		}
		repTimes = append(repTimes, v)
	}

	var nextdate time.Time

	switch repMode {
	case "d": //d 2
		for _, v := range repTimes {
			if v > 400 {
				log.Print("max interval d 400 exceeded")
				return "", errors.New("max interval d 400 exceeded")
			}
		}
		if len(repTimes) != 1 || repTimes[0] <= 0 {
			log.Print("mode 'd' needs one positive number")
			return "", errors.New("mode 'd' needs one positive number")
		}

		fmt.Println("Got D!", repTimes)

		if taskTime.After(now) {
			nextdate = taskTime.AddDate(0, 0, repTimes[0])
		} else {
			for taskTime.Before(now) {
				taskTime = taskTime.AddDate(0, 0, repTimes[0])
			}
			nextdate = taskTime
		}

	case "w": //w 1,2,3,4,5
		fmt.Println("Got W!", repTimes)

		if len(repTimes) == 0 {
			log.Print("mode 'w' needs up to 7 reps")
			return "", errors.New("mode 'w' needs reps")
		}

		wdCheckIn := make(map[time.Weekday]bool)

		for _, d := range repTimes {
			if d < 1 || d > 7 {
				log.Print("wrong repTimes format")
				return "", fmt.Errorf("repTimes format 1-7, got %d", d)
			}
			var wday time.Weekday

			if d == 7 {
				wday = time.Sunday
			} else {
				wday = time.Weekday(d)
			}
			wdCheckIn[wday] = true
		}
		taskWeekBegin := now.AddDate(0, 0, 1)
		taskWeekEnd := taskWeekBegin.AddDate(0, 0, 8)

		for taskWeekBegin.Before(taskWeekEnd) {
			if wdCheckIn[taskWeekBegin.Weekday()] {
				nextdate = taskWeekBegin
				break
			}
			taskWeekBegin = taskWeekBegin.AddDate(0, 0, 1)
		}
		if nextdate.IsZero() {
			return "", errors.New("nextdate w zero")
		}

	case "m": //m 15,16 10,11
		const (
			lastDayMonth     = -1
			sec2lastDayMonth = -2
		)

		var repDatesStr []string
		var repMonthsStr []string

		if len(repModeTimes) == 3 {
			repDatesStr = strings.Split(repModeTimes[1], ",")
			repMonthsStr = strings.Split(repModeTimes[2], ",")
		}
		if len(repModeTimes) == 2 {
			repDatesStr = strings.Split(repModeTimes[1], ",")
		}
		if len(repModeTimes) < 2 {
			log.Print("wrong m mode sizes")
			return "", errors.New("wrong m mode size")
		}

		repDays := make([]int, len(repDatesStr))

		for i, v := range repDatesStr {
			val, err := strconv.Atoi(v)
			if err != nil {
				log.Print("fail conv m repDates:", err)
				return "", fmt.Errorf("fail conv m repDates: %w", err)
			}
			if val == lastDayMonth || val == sec2lastDayMonth {
				repDays[i] = val
			} else if val < 1 || val > 31 {
				return "", fmt.Errorf("days must be 1–31, got: %d", val)
			} else {
				repDays[i] = val
			}
		}
		slices.Sort(repDays)

		var repMonths []int

		if len(repMonthsStr) > 0 {
			repMonths = make([]int, len(repMonthsStr))
			for i, v := range repMonthsStr {
				val, err := strconv.Atoi(v)
				if err != nil {
					log.Print(err)
					return "", fmt.Errorf("fail conv m repMonths: %w", err)
				}
				if val < 1 || val > 12 {
					return "", fmt.Errorf("months must be 1–12, got: %d", val)
				}
				repMonths[i] = val
			}
			slices.Sort(repMonths)

		} else {
			repMonths = make([]int, 12)
			for m := 1; m <= 12; m++ {
				repMonths[m-1] = m
			}
		}
		//var dueDay [32]bool
		var dueMonth [13]bool

		// for _, d := range repDays {
		// 	dueDay[d] = true
		// }
		for _, m := range repMonths {
			dueMonth[m] = true
		}

		nextPossibleDueDate := now
		maxRange := 400
		count := 0

		for {
			nextPossibleDueDate = nextPossibleDueDate.AddDate(0, 0, 1)
			count++
			if count > maxRange {
				log.Print("mode m: no nextdate found, out of range 400 days")
				return "", errors.New("mode m: no nextdate found, out of range 400 days")
			}

			y, m, d := nextPossibleDueDate.Date()

			if !dueMonth[m] {
				continue
			}
			ok := false
			for _, ruleDay := range repDays {
				targetDay := 0

				if ruleDay == lastDayMonth {
					targetDay = time.Date(y, m+1, 0, 0, 0, 0, 0, nextPossibleDueDate.Location()).Day()

				} else if ruleDay == sec2lastDayMonth {
					lastday := time.Date(y, m+1, 0, 0, 0, 0, 0, nextPossibleDueDate.Location()).Day()
					targetDay = lastday - 1
				} else {
					targetDay = ruleDay
				}

				if d == targetDay {
					ok = true
					break
				}
			}
			if ok {
				nextdate = nextPossibleDueDate
				break
			}
		}

	case "y":
		fmt.Println("Got Y!", repTimes)
		newtaskTime := taskTime
		if taskTime.After(now) {
			nextdate = taskTime.AddDate(1, 0, 0)
		} else {
			for newtaskTime.Before(now) {
				newtaskTime = newtaskTime.AddDate(1, 0, 0)
			}
			nextdate = newtaskTime
		}

	default:
		log.Printf("wrong mode: must be d, w, m, y, got %v", repModeTimes)
		return "", errors.New("wrong mode: must be d, w, m, y")

	}

	fmt.Println(now.Format("20060102"), dstart, repeat)

	return nextdate.Format("20060102"), nil
}
