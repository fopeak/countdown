package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"time"

	tm "github.com/buger/goterm"
	"github.com/fatih/color"
)

// Dates @ dates in dates.json
type Dates struct {
	Dates []Date `json:"dates"`
}

// Date @ date data in dates.json
type Date struct {
	Name      string `json:"name"`
	Timestamp int64  `json:"timestamp"`
}

const (
	day  = time.Minute * 60 * 24
	year = 365 * day
)

func main() {
	// arguments
	args := os.Args[1:]
	// check arguments
	if len(args) == 0 {
		fmt.Println("'countdown c' to countdown")
		fmt.Println("'countdown add name unixtimestamp' to add new countdown")
		return
	}
	var file *os.File
	// check if dates.json exists
	if _, err := os.Stat("dates.json"); err == nil {
		file, _ = os.OpenFile("dates.json", os.O_RDWR, 0777)
	} else {
		file, _ = os.Create("dates.json")
	}
	switch args[0] {
	case "c":
		countdown(file)
		return
	case "add":
		if len(args) == 3 {
			arg2, err := strconv.ParseInt(args[2], 10, 64)
			if err != nil {
				panic(err)
			}
			var date = Date{args[1], arg2}
			add(file, &date)
		} else {
			fmt.Println("set countdown name and timestamp like this, countdown add endoftheworld 1536393600")
		}
		return
	}
}

func add(file *os.File, date *Date) {
	var dates Dates
	old, _ := ioutil.ReadAll(file)
	json.Unmarshal(old, &dates)
	dates.Dates = append(dates.Dates, *date)
	json, err := json.Marshal(dates)
	if err != nil {
		panic(err.Error)
	}
	//a, err := file.Write(json)
	if err != nil {
		fmt.Println(err)
	}
	ioutil.WriteFile("dates.json", json, 0644)
}

func countdown(file *os.File) {
	// clear screen
	tm.Clear()

	var dates, _ = ioutil.ReadAll(file)
	var timestamps Dates

	// unmarshal dates
	json.Unmarshal(dates, &timestamps)
	// update screen per second
	ticker := time.NewTicker(time.Second)
	// color set
	color.Set(color.FgHiBlue)
	for range ticker.C {
		tm.MoveCursor(1, 1)
		for i := range timestamps.Dates {
			tm.Println(timestamps.Dates[i].Name+":", humanizeDuration(time.Until(time.Unix(timestamps.Dates[i].Timestamp, 0))))
		}
		tm.Flush() // Call it every time at the end of rendering
	}
	ticker.Stop()
}

// humanizeDuration humanizes time.Duration output to a meaningful value,
// golang's default ``time.Duration`` output is badly formatted and unreadable.
func humanizeDuration(duration time.Duration) string {
	if duration.Seconds() < 60.0 {
		return fmt.Sprintf("%d seconds", int64(duration.Seconds()))
	}
	if duration.Minutes() < 60.0 {
		remainingSeconds := math.Mod(duration.Seconds(), 60)
		return fmt.Sprintf("%d minutes %d seconds", int64(duration.Minutes()), int64(remainingSeconds))
	}
	if duration.Hours() < 24.0 {
		remainingMinutes := math.Mod(duration.Minutes(), 60)
		remainingSeconds := math.Mod(duration.Seconds(), 60)
		return fmt.Sprintf("%d hours %d minutes %d seconds",
			int64(duration.Hours()), int64(remainingMinutes), int64(remainingSeconds))
	}
	remainingHours := math.Mod(duration.Hours(), 24)
	remainingMinutes := math.Mod(duration.Minutes(), 60)
	remainingSeconds := math.Mod(duration.Seconds(), 60)
	return fmt.Sprintf("%d days %d hours %d minutes %d seconds",
		int64(duration.Hours()/24), int64(remainingHours),
		int64(remainingMinutes), int64(remainingSeconds))
}
