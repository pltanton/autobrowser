BUILD_ENVPARAMS:=CGO_ENABLE=0

clear:
	rm -rf build

.PHONY: build-linux
build-linux:
	$(BUILD_ENVPARAMS) go build -o build/autobrowser cmd/autobrowser-linux/main.go

build-macos: build-linux
	mkdir "build/Autobrowser.app"
	cp build/autobrowser build/Autobrowser.app/autobrowser-bin
	cp macos/* build/Autobrowser.app