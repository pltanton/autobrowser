.PHONY: clean
clean:
	rm -rf out

.PHONY: build-linux
build-linux:
	CGO_ENABLE=0 go build -C linux -o ../out/autobrowser cmd/autobrowser/main.go

.PHONY: build-macos
build-macos: clean
	mkdir -p "out/Autobrowser.app"
	cp macos/assets/* out/Autobrowser.app
	go build -C macos -o ../out/Autobrowser.app/autobrowser cmd/autobrowser/main.go

# only for reference, make for windows is _not really_ a thing
.PHONY: build-windows
build-windows:
	go build -C windows -o ../out/autobrowser cmd/autobrowser/main.go

.PHONY: install-macos
install-macos: build-macos
	cp -r out/Autobrowser.app ~/Applications