package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"jcourse_go/internal/dal"
	"jcourse_go/internal/infrastructure/entity"
	"jcourse_go/pkg/util/selenium-get"
)

func main() {
	_ = godotenv.Load()
	dal.InitDBClient()
	db := dal.GetDBClient()
	extend_teacher_data_path := "./data/teachers.json"
	log_file, err := os.OpenFile("./data/logfile.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer log_file.Close()
	defer log.SetOutput(os.Stdout)

	log.SetOutput(log_file)

	extend_teachers := seleniumget.LoadTeacherProfiles(extend_teacher_data_path)

	// to extend: email, profile_url, profile_desc, picture
	for _, t := range extend_teachers {
		// t.department 是全名，jwc是简称
		var teachers []entity.TeacherPO
		db.Model(entity.TeacherPO{}).Where("name = ?", t.Name).Find(&teachers)
		if len(teachers) == 1 {
			teachers[0].Email = t.Mail
			teachers[0].ProfileURL = t.ProfileUrl
			teachers[0].Biography = t.Biography
			teachers[0].Picture = t.HeadImage
			db.Save(&teachers[0])
			continue
		}
		// len == 0, no need to extend
		if len(teachers) > 1 {
			confirm := false
			for _, tt := range teachers {
				if tt.Department == t.Department {
					tt.Email = t.Mail
					tt.ProfileURL = t.ProfileUrl
					tt.Biography = t.Biography
					tt.Picture = t.HeadImage
					db.Save(&tt)
					confirm = true
					break
				}
			}
			if !confirm {
				log.Printf("name %s has multiple teachers, please extend by hand", t.Name)
			}
		}

	}

}
