module github.com/pltanton/autobrowser/macos

go 1.23.0

toolchain go1.24.4

require (
	github.com/pltanton/autobrowser/common v0.0.0
	golang.org/x/term v0.34.0
)

require golang.org/x/sys v0.35.0 // indirect

replace github.com/pltanton/autobrowser/common => ../common
