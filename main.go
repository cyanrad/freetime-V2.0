package main

import (
	"io"
	"log"
	"os"
	"sort"

	lo "github.com/samber/lo"
)

type Period struct {
	day   int
	start int
	end   int
}

func main() {
	inputFiles, outputFile := initArgs()

	scheduleData := []Period{}

	for _, f := range inputFiles {
		csvFile, err := os.Open(f)
		if err != nil {
			log.Fatal(err)
		}

		scheduleData = append(scheduleData, csvToPeriods(csvFile)...)

		csvFile.Close()
	}

	byDayPeriodGroups := lo.GroupBy(scheduleData, func(p Period) int { return p.day })

	// initializing the writer of the output free time data (file|stdout)
	var w io.Writer
	if outputFile != "" {
		f, err := os.OpenFile(outputFile, os.O_WRONLY|os.O_CREATE, 0222)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		w = f
	} else {
		w = os.Stdout
	}

	freeTime := []Period{}
	for _, group := range byDayPeriodGroups {
		freeTime = append(
			freeTime,
			reverse(compressPeriodGroup(group))...,
		)
	}

	periodsToCsv(freeTime, w)
}

func compressPeriodGroup(group []Period) []Period {
	start, end := unzipPeriods(group)
	sort.Ints(start)
	sort.Ints(end)

	compressed := []Period{} // the final compressed list
	periodStart := -1        // the time a compressed preiod starts
	day := group[0].day      // just the group day
	counter := 0             // counter to check when a compressed preiod is over

	for i, j := 0, 0; i < len(start); {
		if start[i] < end[j] {
			if periodStart == -1 {
				periodStart = start[i]
			}
			counter++
			i++
		} else if start[i] == end[j] {
			i++
			j++
		} else {
			counter--

			if counter == 0 {
				compressed = append(compressed, Period{
					day:   day,
					start: periodStart,
					end:   end[j],
				})
				periodStart = -1 // resetting period start
			}

			j++
		}
	}

	// compressing any remaining periods
	if counter > 0 {
		compressed = append(compressed, Period{
			day:   day,
			start: periodStart,
			end:   end[len(end)-1],
		})
	}

	return compressed
}

func unzipPeriods(periods []Period) ([]int, []int) {
	start := []int{}
	end := []int{}

	for _, p := range periods {
		start = append(start, p.start)
		end = append(end, p.end)
	}

	return start, end
}

func reverse(group []Period) []Period {
	freeTime := []Period{}
	start := 0
	i := 0

	if group[0].start == 0 {
		start = group[0].end
		i++
	}

	for ; i < len(group); i++ {
		freeTime = append(freeTime, Period{
			day:   group[i].day,
			start: start,
			end:   group[i].start,
		})

		start = group[i].end
	}

	i--
	if start != 2460 {
		freeTime = append(freeTime, Period{
			day:   group[i].day,
			start: start,
			end:   2460,
		})
	}

	return freeTime
}
