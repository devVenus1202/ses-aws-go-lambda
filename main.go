package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

var (
	EmailNotProvided   = errors.New("no email provided")
	MessageNotProvided = errors.New("no message provided")
)

const (
	Success = "success"
	Error   = "error"
)

type LogData struct {
	UserIdentity       UserIdentity `json:"userIdentity"`
	EventSource        string       `json:"eventSource"`
	EventName          string       `json:"eventName"`
	AwsRegion          string       `json:"awsRegion"`
	SourceIPAddress    string       `json:"sourceIPAddress"`
	UserAgent          string       `json:"userAgent"`
	ErrorCode          string       `json:"errorCode"`
	ErrorMessage       string       `json:"errorMessage"`
	EventType          string       `json:"eventType"`
	RequestID          string       `json:"requestID"`
	RecipientAccountID string       `json:"recipientAccountId"`
}

type UserIdentity struct {
	Type        string `json:"type"`
	PrincipalID string `json:"principalId"`
	Arn         string `json:"arn"`
	AccountID   string `json:"accountId"`
	AccessKeyID string `json:"accessKeyId"`
}

type ResponseMessage struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

var toEmail string
var fromEmail string
var subject string
var emailClient *ses.SES
var regionName string

func init() {
	toEmail = os.Getenv("TO_EMAIL")
	subject = os.Getenv("SUBJECT")
	fromEmail = os.Getenv("FROM_EMAIL")
	regionName = os.Getenv("REGION")

	if len(subject) < 0 {
		subject = "Message from website"
	}

	emailClient = ses.New(session.New(), aws.NewConfig().WithRegion(regionName))
}

func HandleRequest(ctx context.Context, logEvent events.CloudwatchLogsEvent) (string, error) {

	d, err := logEvent.AWSLogs.Parse()
	if err != nil {
		return fmt.Sprintf("1-Error %s!", err), nil
	}

	var message_html string
	for _, event := range d.LogEvents {

		var logData LogData
		err := json.Unmarshal([]byte(event.Message), &logData)
		if err != nil {
			return fmt.Sprintf("2-Error %s!", err), nil
		}
		message_html = message_html + "<b>" + logData.UserIdentity.AccessKeyID + "</b><br>"
		message_html = message_html + "<h3> EventSource:" + logData.EventSource + "</h3>"
		message_html = message_html + "<p><strong><label> User Identity Arn: </label></strong> " + logData.UserIdentity.Arn + "</p>"
		message_html = message_html + "<p><strong><label> User Identity AccountId: </label> </strong>" + logData.UserIdentity.AccountID + "</p>"
		message_html = message_html + "<p><strong><label> User Identity AccessKeyId: </label> </strong> " + logData.UserIdentity.AccessKeyID + "</p>"
		message_html = message_html + "<p><strong><label> Event Name: </label> </strong>" + logData.EventName + "</p>"
		message_html = message_html + "<p><strong><label> Event Type: </label> </strong>" + logData.EventType + "</p>"
		message_html = message_html + "<p><strong><label> IP Address: </label> </strong>" + logData.SourceIPAddress + "</p>"
		message_html = message_html + "<p><strong><label> AWS Region: </label> </strong>" + logData.AwsRegion + "</p>"
		message_html = message_html + "<p><strong><label> User Agent: </label> </strong> " + logData.UserAgent + "</p>"
		message_html = message_html + "<p><strong><label> Event Type: </label> </strong>" + logData.EventType + "</p>"
		message_html = message_html + "<p><strong><label> Recipient AccountId: </label> </strong>" + logData.RecipientAccountID + "</p>"
		message_html = message_html + "<p><strong><label> Error Message: </label></strong>" + logData.ErrorMessage + "</p>"
		message_html = message_html + "<p><strong><label> Error Code: </label></strong>" + logData.ErrorMessage + "</p>"
		message_html = message_html + "<br/>"
	}

	emailParams := &ses.SendEmailInput{
		Message: &ses.Message{
			Body: &ses.Body{
				Text: &ses.Content{
					Data: aws.String("From " + fromEmail + "<br/>" + message_html),
				},
				Html: &ses.Content{
					Data: aws.String("From " + fromEmail + "<br/><br/>" + message_html),
				},
			},
			Subject: &ses.Content{
				Data: aws.String(subject),
			},
		},
		Destination: &ses.Destination{
			ToAddresses: []*string{aws.String(toEmail)},
		},
		Source: aws.String(toEmail),
	}

	_, err = emailClient.SendEmail(emailParams)

	if err != nil {
		return fmt.Sprintf("Error %s!", err), nil
	}

	successResponse, err := json.Marshal(ResponseMessage{Success, "Message is sent"})
	return fmt.Sprintf("Success %s!", successResponse), nil

}

func main() {
	lambda.Start(HandleRequest)
}
