# push - simple apis to send push notifications

Supports SMS and SES channels


Sample calls:

<pre>
curl -X POST -H "x-api-key:xxxx" -d '{"sender_email":"sender@domain", "to":["dest@somewhere"], "subject":"stuff","message":"yo"}' https://myapiid.execute-api.us-east-1.amazonaws.com/dev/push/api/v1/emailmessage

curl -X POST -H "x-api-key:xxxx" -d '{"email_address":"someone@somewhere"}' https://myapiid.execute-api.us-east-1.amazonaws.com/dev/push/api/v1/regaddress
</pre>