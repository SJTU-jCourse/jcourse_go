export LOCAL_MODULE := jcourse_go

.PHONY :lint
lint:
	@go fmt ./...
	@goimports -local $(LOCAL_MODULE) -w $$(find . -type f -name '*.go')
	@go mod tidy
