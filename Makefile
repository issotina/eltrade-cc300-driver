APP = MECEF.driver
VER = $(shell git describe --tags)
osx:
	go build -o eltrade
	mkdir "bin/macOS-x64/application"
	mv ./eltrade ./bin/macOS-x64/application/eltrade
	cd bin/macOS-x64 && yes N | ./build-macos-x64.sh $(APP) $(VER)

win32:
win64:
linux:
all: osx win32 win64 linux

