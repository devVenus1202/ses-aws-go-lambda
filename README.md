# ses-aws-go-lambda

Run following commands:
GOOS=linux go build main.go
zip ./main.zip main

In AWS console, create Go 1.x lambda and upload zip file.

To have a test you need to set test event.
On "Configue test event", 
1. Select "Create new test event"
2. Select "Amazon Cloudwatch Logs" for "Event template"
3. Type event name in "Event name"
4. Copy and Past following content in Editor.
{
  "awslogs": {
    "data": "H4sIAAAAAAAAAHWPwQqCQBCGX0Xm7EFtK+smZBEUgXoLCdMhFtKV3akI8d0bLYmibvPPN3wz00CJxmQnTO41whwWQRIctmEcB6sQbFC3CjW3XW8kxpOpP+OC22d1Wml1qZkQGtoMsScxaczKN3plG8zlaHIta5KqWsozoTYw3/djzwhpLwivWFGHGpAFe7DL68JlBUk+l7KSN7tCOEJ4M3/qOI49vMHj+zCKdlFqLaU2ZHV2a4Ct/an0/ivdX8oYc1UVX860fQDQiMdxRQEAAA=="
  }
}

5. Set environment value
TO_EMAIL "receiver@go.com"
FROM_EMAILSUBJECT "receiver@go.com"
REGION "us-east-1"
SUBJECT "Test"
