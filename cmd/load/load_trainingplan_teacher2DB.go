package load

import (
	"jcourse_go/dal"
	seleniumget "jcourse_go/util/selenium-get"
)

func main() {
	dal.InitDBClient()
	// seleniumget.LoadTrainingPlan2DB("./util/selenium-get/data/trainingPlan.txt", dal.GetDBClient())
	seleniumget.LoadTeacherProfile2DB("./util/selenium-get/data/teachers.json", dal.GetDBClient())
}
