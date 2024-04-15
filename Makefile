.PHONY: clean
clean:
	rm -rf build

.PHONY: build-linux
build-linux:
	CGO_ENABLE=0 go build -C linux -o build/autobrowser-linux cmd/autobrowser/main.go

.PHONY: build-macos
build-macos:
	go build -C macos -o build/autobrowser-mac cmd/autobrowser/main.go
	mkdir -p "build/Autobrowser.app"
	mv build/autobrowser-mac build/Autobrowser.app/autobrowser-bin
	cp macos/assets/* build/Autobrowser.app

.PHONY: install-macos
install-macos: build-macos
	cp -r build/Autobrowser.app ~/Applications