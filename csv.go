package main

import (
	"encoding/csv"
	"io"
	"log"
	"strconv"
)

func csvToPeriods(r io.Reader) []Period {
	cr := csv.NewReader(r)
	parsedValues, err := cr.ReadAll()
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

func periodsToCsv(periods []Period, w io.Writer) {
	cw := csv.NewWriter(w)

	headers := []string{"day", "start", "end"}
	if err := cw.Write(headers); err != nil {
		log.Fatal(err)
	}

	for _, p := range periods {
		row := []string{
			strconv.Itoa(p.day),
			strconv.Itoa(p.start),
			strconv.Itoa(p.end),
		}

		if err := cw.Write(row); err != nil {
			log.Fatal(err)
		}
	}

	cw.Flush()
}
