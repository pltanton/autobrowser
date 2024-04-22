.PHONY: clean
clean:
	rm -rf build

.PHONY: build-linux
build-linux:
	CGO_ENABLE=0 go build -C linux -o ../out/autobrowser cmd/autobrowser/main.go

.PHONY: build-macos
build-macos: clean
	mkdir -p "build/Autobrowser.app"
	cp macos/assets/* build/Autobrowser.app
	go build -C macos -o ../out/Autobrowser.App/autobrowser cmd/autobrowser/main.go 

.PHONY: install-macos
install-macos: build-macos
	cp -r out/Autobrowser.app ~/Applications