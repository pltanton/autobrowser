BUILD_ENVPARAMS:=CGO_ENABLE=0

clean:
	rm -rf build

.PHONY: build-linux
build-linux:
	$(BUILD_ENVPARAMS) go build -o build/autobrowser cmd/autobrowser-linux/main.go

build-macos:
	go build -C cmd/autobrowser-macos -o ../../build/autobrowser .
	mkdir -p "build/Autobrowser.app"
	mv build/autobrowser build/Autobrowser.app/autobrowser-bin
	cp macos/* build/Autobrowser.app

install-macos: build-macos
	cp -r build/Autobrowser.app ~/Applications