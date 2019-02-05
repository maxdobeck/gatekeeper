package models

import (
	"database/sql"
	log "github.com/sirupsen/logrus"
)

// Schedule contains the core values
type Schedule struct {
	ID, Title, OwnerID string
}

// CreateSchedule builds a new schedule with the creator as the Owner
func CreateSchedule(s *Schedule) error {
	_, err := Db.Query("INSERT INTO schedules(title, owner_id) VALUES ($1,$2)", s.Title, s.OwnerID)
	if err != nil {
		log.Info(err)
		return err
	}
	log.Info("Schedule Created: ", s)
	return err
}

// GetSchedules gets all schedules based on the provided member id
func GetSchedules(memberID string) ([]Schedule, error) {
	rows, err := Db.Query("SELECT id, title, owner_id FROM schedules WHERE owner_id = $1;", memberID)
	if err != nil {
		log.Info(err)
		return nil, err
	}
	defer rows.Close()
	s := make([]Schedule, 0)
	for rows.Next() {
		var id, title, ownerID string
		err := rows.Scan(&id, &title, &ownerID)
		if err != nil {
			log.Info(err)
			return nil, err
		}
		s = append(s, Schedule{id, title, ownerID})
	}
	log.Info("Array of all schedules owned by: ", memberID, s)
	return s, err
}

// GetScheduleByID will obtain the schedule based on the provided schedule id
func GetScheduleByID(scheduleID string) (Schedule, error) {
	var s Schedule
	row, err := Db.Query("SELECT id, title, owner_id FROM schedules WHERE id = $1;", scheduleID)
	if err != nil {
		log.Info(err)
		return s, err
	}
	if err == sql.ErrNoRows {
		log.Info("No record found for: ", scheduleID)
		return s, err
	}
	defer row.Close()
	for row.Next() {
		var id, title, ownerID string
		err := row.Scan(&id, &title, &ownerID)
		if err != nil || id != scheduleID {
			log.Info(err)
			return s, err
		}
		s.Title = title
		s.OwnerID = ownerID
		s.ID = id
	}
	log.Info("Schedule found: ", s)
	return s, err
}

// UpdateScheduleTitle will change the title of the specificed schedule
func UpdateScheduleTitle(scheduleID string, newTitle string) error {
	_, err := Db.Query("UPDATE schedules SET title = $2 WHERE id = $1", scheduleID, newTitle)
	if err == sql.ErrNoRows {
		return err
	}
	if err != nil {
		log.Info(err)
		return err
	}
	return nil
}

// DeleteSchedule will delete a schedule based on the id of the schedule
func DeleteSchedule(scheduleID string) error {
	_, err := Db.Query("DELETE FROM schedules WHERE id = $1;", scheduleID)
	if err != nil {
		log.Info("Schedule could not be deleted by id: ", scheduleID)
		return err
	}
	log.Info("Schedule deleted: ", scheduleID)
	return nil
}

/*// Generate a link like google does with sheets or docs
like so: https://docs.google.com/spreadsheets/d/1Qm_7-QB9eZJBjK_mESb6Oy1kVzAiJgCp_rPp3c1zHrI/edit#gid=0
func generateShareLink() error {
}
*/
