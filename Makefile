compile: sms ses

sms:
	GOOS=linux go build -o bin/sms functions/sms/*.go

ses:
	GOOS=linux go build -o bin/ses functions/ses/*.go
