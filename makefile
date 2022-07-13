dev:
	go run main.go
	
build:
	env CC=/opt/gcc-linaro-7.5.0-2019.12-x86_64_arm-linux-gnueabihf/bin/arm-linux-gnueabihf-gcc CXX=/opt/gcc-linaro-7.5.0-2019.12-x86_64_arm-linux-gnueabihf/bin/arm-linux-gnueabihf-g++ \
    CGO_ENABLED=1 GOOS=linux GOARCH=arm \
    go build -v
