build:
	go build -o terraform-provider-phpipam_${TRAVIS_TAG}

test:
	go test -v $(shell go list ./... | grep -v /vendor/) 

testacc:
	TF_ACC=1 go test -v ./plugin/providers/phpipam -run="TestAcc"

release: release_bump release_build

release_bump:
	scripts/release_bump.sh

release_build:
	scripts/release_build.sh

clean:
	rm -rf pkg/
