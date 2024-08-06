.PHONY load:
	go run cmd/migrate/migrate.go
	go run cmd/import_from_jwc/import_from_jwc.go
	go run cmd/load/extend_teacher_profile.go
	cat ./log/logfile.log
	go run cmd/load/extend_training_plan.go

CLEAN_FILES=./log/* gorm.db

.PHONY: clean

clean:
	$(RM) gorm.db
	$(foreach file,$(wildcard $(CLEAN_FILES)),$(if $(wildcard $(file)),$(RM) $(file);))


.PHONY reset: clean load