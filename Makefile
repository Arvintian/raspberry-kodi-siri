RASPBERRY_IP = 192.168.199.110

build: cmd
	GOOS=linux GOARCH=arm GOARM=7 go build -o ./bin/rspberry-kodi-siri ./cmd
	scp -r $(shell pwd)/bin/rspberry-kodi-siri pi@$(RASPBERRY_IP):/home/pi/bin