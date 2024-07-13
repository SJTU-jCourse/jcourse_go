package seleniumget

import (
	"testing"
)

func TestLoadTrainingPlan2DB(t *testing.T){
	LoadTrainingPlan2DB("./data/trainingPlan.txt", nil)
}
func TestLoadTeacherProfile2DB(t *testing.T){
	LoadTeacherProfile2DB("./data/teachers.json", nil)
}