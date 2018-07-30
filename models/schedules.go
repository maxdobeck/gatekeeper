package models

import (
	"log"
)

// NewMember is the struct for the member signup process
type NewSchedule struct {
	Title, Owner string
}

// CreateSchedule builds a new schedule with the creator as the Owner
func CreateSchedule(s *NewSchedule) error {
	_, err := Db.Query("INSERT INTO schedules(title, owner_id) VALUES ($1,$2)", s.Title, s.Owner)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("Schedule Created: ", s)
	return err
}

/*func GetSchedules(ownerId string) {
}

func GetSchedule(schduleId string)

func UpdateScheduleName(scheduleId string) {

}

func DeleteSchedule(scheduleId string) {

}

// Generate a link like google does with sheets or docs https://docs.google.com/spreadsheets/d/1Qm_7-QB9eZJBjK_mESb6Oy1kVzAiJgCp_rPp3c1zHrI/edit#gid=0
func generateShareLink() error {

}
*/
