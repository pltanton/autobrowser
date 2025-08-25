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

.PHONY: install-macos
install-macos: build-macos
	cp -r out/Autobrowser.app ~/Applications

.PHONY: build-macos-and-install-dev
build-macos-and-install-dev: clean
	mkdir -p "out/AutobrowserDev.app"
	cp macos/assets/AppIcon.icns out/AutobrowserDev.app
	sed -e 's/Autobrowser/AutobrowserDev/g' \
		-e 's/autobrowser/autobrowser-dev/g' \
		-e 's/com\.pltanton\.autobrowser/dev.pltanton.autobrowser.dev/g' \
		macos/assets/Info.plist > out/AutobrowserDev.app/Info.plist
	go build -C macos -o ../out/AutobrowserDev.app/autobrowser-dev cmd/autobrowser/main.go
	cp -r out/AutobrowserDev.app ~/Applications
