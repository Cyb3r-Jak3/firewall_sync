lint:
	go vet
	golint -set_exit_status


scan:
	gosec -no-fail -fmt sarif -out security.sarif ./...