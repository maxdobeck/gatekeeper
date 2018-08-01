package models

import (
	"database/sql"
	"log"
)

// NewMember is the struct for the member signup process
type Schedule struct {
	Title, Owner string
}

// CreateSchedule builds a new schedule with the creator as the Owner
func CreateSchedule(s *Schedule) error {
	_, err := Db.Query("INSERT INTO schedules(title, owner_id) VALUES ($1,$2)", s.Title, s.Owner)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("Schedule Created: ", s)
	return err
}

func GetSchedules(memberId string) ([]*Schedule, error) {
	rows, err := Db.Query("SELECT id, title FROM schedules WHERE owner_id = $1;", memberId)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()
	s := make([]*Schedule, 0)
	for rows.Next() {
		var id, title string
		err := rows.Scan(id, title)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		s = append(s, &Schedule{id, title})
	}
	log.Println("Array of all schedules owned by: ", memberId, s)
	return s, err
}

/*
func GetSchedule(schduleId string) {

}
*/
// UpdateScheduleTitle will change the title of the specificed schedule
func UpdateScheduleTitle(scheduleId string, newTitle string) error {
	_, err := Db.Query("UPDATE schedules SET title = $2 WHERE id = $1", scheduleId, newTitle)
	if err == sql.ErrNoRows {
		return err
	}
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

/*func DeleteSchedule(scheduleId string) {

}

// Generate a link like google does with sheets or docs https://docs.google.com/spreadsheets/d/1Qm_7-QB9eZJBjK_mESb6Oy1kVzAiJgCp_rPp3c1zHrI/edit#gid=0
func generateShareLink() error {

}
*/
