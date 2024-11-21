package mock

import (
	"jcourse_go/dal"
	"jcourse_go/repository"
	"math/rand"

	"gorm.io/gorm"
)

type MockDBGenerator struct {
	db        *gorm.DB
	batchSize int
	Rand      *rand.Rand
}

type MockDBParams struct {
	Seed      int64
	BatchSize int
	Users     int
	Courses   int
	Reviews   int
	Teachers  int
}

// ClearAllTables 清空数据库中的所有表(PG限定)
func ClearAllTables(db *gorm.DB) error {
	tables := []string{}
	err := db.Raw("SELECT tablename FROM pg_tables WHERE schemaname = 'public'").Scan(&tables).Error
	if err != nil {
		return err
	}

	for _, table := range tables {
		err = db.Migrator().DropTable(table)
		if err != nil {
			return err
		}
	}
	return nil
}

func MockDB(params MockDBParams) (*gorm.DB, error) {
	randSource := SetSeed(params.Seed)
	err := dal.InitTestDBClient()
	if err != nil {
		return nil, err
	}
	err = ClearAllTables(dal.GetDBClient())
	if err != nil {
		return nil, err
	}
	gen := MockDBGenerator{
		db:        dal.GetDBClient(),
		batchSize: params.BatchSize,
		Rand:      rand.New(randSource),
	}
	err = repository.Migrate(gen.db)
	if err != nil {
		return nil, err
	}
	users, err := MockUsers(gen, params.Users)
	if err != nil {
		return nil, err
	}
	teachers, err := MockTeachers(gen, params.Teachers)
	if err != nil {
		return nil, err
	}
	baseCourses, err := MockBaseCourses(gen, params.Courses)
	if err != nil {
		return nil, err
	}
	courses, err := MockCourses(gen, params.Courses, baseCourses, teachers)
	if err != nil {
		return nil, err
	}
	_, err = MockOfferedCourses(gen, params.Courses, courses, teachers)
	if err != nil {
		return nil, err
	}
	_, err = MockReviews(gen, params.Reviews, courses, users)
	if err != nil {
		return nil, err
	}
	return gen.db, nil
}
