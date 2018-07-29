package models

import (
	"database/sql"
	_ "github.com/lib/pq" // github.com/lib/pq
	"golang.org/x/crypto/bcrypt"
	"log"
)

// NewMember is the struct for the member signup process
type NewSchedule struct {
	Name, Owner string
}

// CreateSchedule builds a new schedule with the creator as the Owner
func CreateSchedule(s *NewSchedule) error {

}

func GetSchedules(ownerId string) {
}

func GetSchedule(schduleId string)

func UpdateScheduleName(scheduleId string) {

}

func DeleteSchedule(scheduleId string) {

}

// Generate a link like google does with sheets or docs https://docs.google.com/spreadsheets/d/1Qm_7-QB9eZJBjK_mESb6Oy1kVzAiJgCp_rPp3c1zHrI/edit#gid=0
func generateShareLink() error {

}
