BUILD_ENVPARAMS:=CGO_ENABLE=0

.PHONY: build-linux
build-linux:
	$(BUILD_ENVPARAMS) go build -o build/autobrowser cmd/linux/main.go