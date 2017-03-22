test:
	go test -v $(shell go list ./... | grep -v /vendor/) 

testacc:
	TF_ACC=1 go test -v ./plugin/providers/phpipam -run="TestAcc"

build: deps
	gox -osarch="linux/amd64 windows/amd64 darwin/amd64" \
	-output="pkg/{{.OS}}_{{.Arch}}/terraform-provider-phpipam" .

deps:
	go get -u github.com/mitchellh/gox

clean:
	rm -rf pkg/
