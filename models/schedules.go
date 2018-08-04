package models

import (
	"database/sql"
	"log"
)

// NewMember is the struct for the member signup process
type Schedule struct {
	Id, Title, Owner string
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

// GetSchedules gets all schedules based on the provided member id
func GetSchedules(memberId string) ([]Schedule, error) {
	rows, err := Db.Query("SELECT id, title, owner_id FROM schedules WHERE owner_id = $1;", memberId)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()
	s := make([]Schedule, 0)
	for rows.Next() {
		var id, title, owner_id string
		err := rows.Scan(&id, &title, &owner_id)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		s = append(s, Schedule{id, title, owner_id})
	}
	log.Println("Array of all schedules owned by: ", memberId, s)
	return s, err
}

// GetScheduleById will obtain the schedule based on the provided schedule id
func GetScheduleById(scheduleId string) (Schedule, error) {
	var s Schedule
	row, err := Db.Query("SELECT id, title, owner_id FROM schedules WHERE id = $1;", scheduleId)
	if err != nil {
		log.Println(err)
		return s, err
	}
	if err == sql.ErrNoRows {
		log.Println("No record found for: ", scheduleId)
		return s, err
	}
	defer row.Close()
	for row.Next() {
		var id, title, owner_id string
		err := row.Scan(&id, &title, &owner_id)
		if err != nil || id != scheduleId {
			log.Println(err)
			return s, err
		}
		s.Title = title
		s.Owner = owner_id
		s.Id = id
	}
	log.Println("Schedule found: ", s)
	return s, err
}

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
