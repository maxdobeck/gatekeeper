package models

import (
	"log"
	"strconv"
	"strings"
)

// Shift contains the start, end, title, days worked, stop date, minimum enrollees, and owning schedule
type Shift struct {
	ID, Title, Start, End, Stop, MinEnrollees, Schedule string
	Days                                                [7]string
}

// CreateShift builds a new shift and attaches it to a schedule
func CreateShift(s *Shift) error {
	d := week(s.Days)
	min, _ := strconv.Atoi(s.MinEnrollees)
	_, err := Db.Query("INSERT INTO shifts(title, start_time, end_time, stop_date, min_enrollees, schedule_id, sun, mon, tue, wed, thu, fri, sat) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)",
		s.Title, s.Start, s.End, s.Stop, min, s.Schedule, d[0], d[1], d[2], d[3], d[4], d[5], d[6])
	if err != nil {
		log.Println("Error creating shift with schedule: ", err, s)
		return err
	}
	log.Println("Shift Created: ", s)
	return err
}

func week(d [7]string) [7]bool {
	var w [7]bool
	sw := []string{"sunday", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday"}
	for i, v := range sw {
		for _, d := range d {
			if strings.EqualFold(v, d) {
				w[i] = true
				break
			} else {
				w[i] = false
			}
		}
	}
	return w
}
