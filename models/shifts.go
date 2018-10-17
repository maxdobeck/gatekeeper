package models

import (
	"log"
)

// Shift contains the core values
type Shift struct {
	ID, Title, Start, End, Stop, minEnrollees, Schedule string
	Days                                                []string
}

// CreateShift builds a new schedule with the creator as the Owner
func CreateShift(s *Shift) error {
	_, err := Db.Query("INSERT INTO shifts(Title, Start, End, Stop, minEnrollees, Days) VALUES ($1,$2,$3,$4,$5,$6,$7)", s.Title, s.Start, s.End, s.Stop, s.minEnrollees, s.Days, s.Schedule)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("Shift Created: ", s)
	return err
}
