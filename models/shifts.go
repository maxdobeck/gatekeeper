package models

import (
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

// Shift contains the start, end, title, days worked, stop date, minimum enrollees, and owning schedule
type Shift struct {
	ID, Title, Start, End, Stop, MinEnrollees, Schedule string
	Days                                                [7]string
}

//ShiftPayload is the shift struct sent to the frontend
type ShiftPayload struct {
	ID, Title, Start, End, Stop, MinEnrollees, Schedule, Created string
	Sun, Mon, Tue, Wed, Thur, Fri, Sat                           bool
}

// CreateShift builds a new shift and attaches it to a schedule
func CreateShift(s *Shift) error {
	d := week(s.Days)
	min, _ := strconv.Atoi(s.MinEnrollees)
	_, err := Db.Query("INSERT INTO shifts(title, start_time, end_time, stop_date, min_enrollees, schedule_id, sun, mon, tue, wed, thu, fri, sat) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)",
		s.Title, s.Start, s.End, s.Stop, min, s.Schedule, d[0], d[1], d[2], d[3], d[4], d[5], d[6])
	if err != nil {
		log.Info("Error creating shift with schedule: ", err, s)
		return err
	}
	log.Info("Shift Created: ", s)
	return err
}

//GetShifts obtains all shifts linked to the supplied schedule
func GetShifts(scheduleID string) ([]ShiftPayload, error) {
	rows, err := Db.Query("SELECT * FROM shifts WHERE schedule_id = $1;", scheduleID)
	if err != nil {
		log.Info(err)
		return nil, err
	}
	defer rows.Close()
	s := make([]ShiftPayload, 0)
	for rows.Next() {
		var ID, Title, StartTime, EndTime, StopDate, MinEnrollees, Schedule, created string
		var sun, mon, tue, wed, thu, fri, sat bool
		err := rows.Scan(&ID, &Schedule, &Title, &MinEnrollees, &StartTime, &EndTime, &StopDate, &sun, &mon, &tue, &wed, &thu, &fri, &sat, &created)
		if err != nil {
			log.Info(err)
			return nil, err
		}
		s = append(s, ShiftPayload{ID, Title, StartTime, EndTime, StopDate, MinEnrollees, Schedule, created, sun, mon, tue, wed, thu, fri, sat})
	}
	log.Info("Array of all shifts owned by this schedule: ", scheduleID, s)
	return s, err
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
