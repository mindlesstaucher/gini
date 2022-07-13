dev:
	go run main.go

arm:
	env CC=/opt/gcc-linaro-7.5.0-2019.12-x86_64_arm-linux-gnueabihf/bin/arm-linux-gnueabihf-gcc GXX=/opt/gcc-linaro-7.5.0-2019.12-x86_64_arm-linux-gnueabihf/bin/arm-linux-gnueabihf-g++ CGO_ENABLED=1 GOOS=linux GOARCH=arm go build -v -o gini-arm-linux .

	
	
build:
	env CC=/opt/gcc-linaro-7.5.0-2019.12-x86_64_arm-linux-gnueabihf/bin/arm-linux-gnueabihf-gcc CXX=/opt/gcc-linaro-7.5.0-2019.12-x86_64_arm-linux-gnueabihf/bin/arm-linux-gnueabihf-g++ \
    CGO_ENABLED=1 GOOS=linux GOARCH=arm \
    go build -v
