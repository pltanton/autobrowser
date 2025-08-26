module github.com/pltanton/autobrowser/macos

go 1.23.0

toolchain go1.24.4

require github.com/pltanton/autobrowser/common v0.0.0

require github.com/BurntSushi/toml v1.5.0 // indirect

replace github.com/pltanton/autobrowser/common => ../common
