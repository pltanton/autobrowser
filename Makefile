.PHONY: clean
clean:
	rm -rf build

.PHONY: build-linux
build-linux:
	CGO_ENABLE=0 go build -C linux -o build/autobrowser-linux cmd/autobrowser/main.go

.PHONY: build-macos
build-macos: clean
	mkdir -p "build/Autobrowser.app"
	cp macos/assets/* build/Autobrowser.app
	go build -C macos -o ../build/Autobrowser.App/autobrowser cmd/autobrowser/main.go 

.PHONY: install-macos
install-macos: build-macos
	cp -r build/Autobrowser.app ~/Applications