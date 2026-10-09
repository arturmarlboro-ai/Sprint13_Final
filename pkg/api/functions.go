package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	"Sprint13_Final/pkg/db"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func checkDate(task *db.Task) error {
	now := time.Now()

	if task.Date == "" {
		task.Date = now.Format(DateFormat)
	}

	t, err := time.Parse("20060102", task.Date)
	if err != nil {
		return fmt.Errorf("checkDate parseing error: %w", err)
	}

	if task.Repeat != "" {
		repMode := strings.SplitN(task.Repeat, " ", 3)
		a := strings.ToLower(repMode[0])
		if a != "d" && a != "w" && a != "m" && a != "y" {
			return errors.New("wrong mode task.Repeat. must be d , w, m, y. got: " + a)
		}

		next, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return fmt.Errorf("checkDate error: %w", err)
		}

		today := now.Format(DateFormat)

		if t.Format(DateFormat) < today {
			task.Date = next
		}
	} else {

		today := now.Format(DateFormat)
		if t.Format(DateFormat) < today {
			task.Date = today
		}
	}
	return nil
}

func writeJson(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(data)
}

func NextDate(now time.Time, dstart, repeat string) (string, error) {
	fmt.Println(now.Format("20060102"), dstart, repeat)

	taskTime, err := time.Parse(DateFormat, dstart)
	if err != nil {
		log.Print("dstart parse err: ", err)
		return "", fmt.Errorf("dstart parse err %w", err)
	}

	// if taskTime.After(now) {
	// 	return taskTime.Format(DateFormat), nil
	// }
	if repeat == "" {
		if !taskTime.After(now) {
			return "", fmt.Errorf("date in the past, repeat rule is empty")
		}
		return dstart, nil
	}
	repModeTimes := strings.SplitN(repeat, " ", 3) //m 15,16 10,11
	repMode := strings.ToLower(strings.TrimSpace(repModeTimes[0]))
	var repTimes []int

	if len(repModeTimes) > 1 {
		for _, n := range strings.Split(repModeTimes[1], ",") {
			v, err := strconv.Atoi(n)
			if err != nil {
				log.Print("atoi error: ", err)
				return "", fmt.Errorf("atoi error: %w", err)
			}
			repTimes = append(repTimes, v)
		}
	}

	var nextdate time.Time

	switch repMode {
	case "d": //d 2
		if len(repTimes) != 1 || repTimes[0] <= 0 {
			log.Print("mode 'd' needs one positive number")
			return "", errors.New("mode 'd' needs one positive number")
		}

		if repTimes[0] > 400 {
			return "", errors.New("max interval d 400 exceeded")
		}

		next := taskTime.AddDate(0, 0, repTimes[0])

		for !next.After(now) {
			next = next.AddDate(0, 0, repTimes[0])
		}

		nextdate = next

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
			if !nextPossibleDueDate.After(now) {
				continue
			}
			if !nextPossibleDueDate.After(taskTime) {
				continue
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

		next := taskTime.AddDate(1, 0, 0)
		for !next.After(now) {
			next = next.AddDate(1, 0, 0)
		}
		nextdate = next

	default:
		log.Printf("wrong mode: must be d, w, m, y, got %v", repModeTimes)
		return "", errors.New("wrong mode: must be d, w, m, y")

	}

	return nextdate.Format(DateFormat), nil
}
