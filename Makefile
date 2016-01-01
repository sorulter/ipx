GIT=`git rev-list HEAD -n 1|cut -c 1-7`

all: dependencies current linux
	@printf "\e[31m Build Complete.\e[0m\r\n"
	@printf "\e[31m Compress Use Upx.\e[0m\r\n"
	@upx --best --backup ./bin/ipx
	@ls -al ./bin/*

dependencies:
	@printf "\e[33;47m Update Dependencies \e[0m\r\n"
	@go get github.com/gocubes/config
	@go get github.com/lessos/lessgo/logger
	@go get github.com/shadowsocks/shadowsocks-go/shadowsocks

current:
	@printf "\e[34m Build Current OS Binary \e[0m\r\n"
	@go build -ldflags "-s -w -X main.Git=$(GIT)" -o ./bin/ipx ./src/*.go

linux:
	@printf "\e[34m Build Linux Binary \e[0m\r\n"
	@GOARCH=amd64 GOOS=linux go build -ldflags "-s -w -X main.Git=$(GIT)" -o ./bin/ipx-64x-linux ./src/*.go

run: dependencies current
	@ls -l ./bin/*
	@printf "\e[32;47m RUN Server \e[0m"
	./bin/ipx -log_dir="./logs"
