all: wechat update control

wechat: main.go */*.go */*/*.go
	go build -o wechat main.go

update: script/updateAccessToken.go
	go build -o update script/updateAccessToken.go

control: script/control.go util/*.go
	go build -o control script/control.go

.PHONY: clean run 

clean:
	rm -rf wechat update control

run:
	./control restart
