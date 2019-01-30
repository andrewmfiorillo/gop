build:
	go build -o ./bin/pgo ./cmd

crossbuild:
	gox -os="linux darwin windows" -arch="amd64 arm" -output="./bin/p_{{.OS}}_{{.Arch}}" ./cmd
