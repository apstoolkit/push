compile: sms

sms:
	GOOS=linux go build -o bin/sms functions/sms/*.go

