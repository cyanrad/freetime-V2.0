package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	lo "github.com/samber/lo"
)

type Period struct {
	day   int
	start int
	end   int
}

func main() {
	inputFiles, _ := initArgs()

	scheduleData := [][]Period{}

	for _, f := range inputFiles {
		scheduleData = append(scheduleData, parsePeriodCsvFile(f))
	}

	flattenedSchedule := lo.Flatten(scheduleData)
	byDayPeriodGroups := lo.GroupBy(flattenedSchedule, func(p Period) int { return p.day })

	for _, group := range byDayPeriodGroups {
		fmt.Println(compressPeriodGroup(group))
		// inverse the groups here as well
	}
}

func initArgs() ([]string, string) {
	inputFilesString := ""
	flag.StringVar(&inputFilesString, "f", "",
		"input schedule csv files (expected format: \"file1, file2\")")
	outputFile := ""
	flag.StringVar(&outputFile, "o", "", "output file")
	flag.Parse()

	inputFiles := strings.Split(inputFilesString, ", ")
	return inputFiles, outputFile
}

func parsePeriodCsvFile(fileName string) []Period {
	fReader, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}

	r := csv.NewReader(fReader)
	parsedValues, err := r.ReadAll()
	if err != nil {
		log.Println(err)
	}

	parsedData := []Period{}

	for _, val := range parsedValues {
		day, err := strconv.Atoi(val[0])
		if err != nil {
			log.Fatal(err)
		}
		start, err := strconv.Atoi(val[1])
		if err != nil {
			log.Fatal(err)
		}
		end, err := strconv.Atoi(val[2])
		if err != nil {
			log.Fatal(err)
		}

		parsedData = append(parsedData, Period{
			day:   day,
			start: start,
			end:   end,
		})
	}

	return parsedData
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
