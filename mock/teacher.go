package mock

import (
	"jcourse_go/model/po"

	"gorm.io/gorm/clause"

	"syreclabs.com/go/faker"
)

func MockTeachers(gen MockDBGenerator, n int) ([]po.TeacherPO, error) {
	teachers := make([]po.TeacherPO, n)
	titles := []string{"教授", "副教授", "讲师"}
	codes, err := GenerateUniqueSet(n, func() string {
		return faker.Code().Ean8()
	})
	if err != nil {
		return nil, err
	}
	for i := 0; i < n; i++ {
		teachers[i] = po.TeacherPO{
			Name:       faker.Name().FirstName(), // TODO 中文支持?
			Code:       codes[i],
			Department: GenDepartment(gen),
			Email:      faker.Internet().Email(),
			Title:      titles[gen.Rand.Intn(len(titles))],
			Pinyin:     faker.Lorem().Characters(10),
			PinyinAbbr: faker.Lorem().Characters(3),
		}
	}
	err = gen.db.Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(teachers, gen.batchSize).Error
	return teachers, err
}
