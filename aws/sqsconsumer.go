package aws


import (
	"fmt"
	"time"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type MessageProcessor struct {
	client   sqs.SQS
	sqsQueue string
}

var (
	// MaxNumberOfMessage at one poll
	MaxNumberOfMessage int64 = 10
	// WaitTimeSecond for each poll
	WaitTimeSecond int64 = 20
)


/*func main() {

	queueUrl := "https://sqs.us-west-2.amazonaws.com/891551238374/ami-build-notifications"
	region := "us-west-2"
	readSqsMessage(queueUrl, region)
}*/

func readSqsMessage( queueUrl string, region string) {
	client, err := sqsClient(region)
	fatalfIfErr("Error: Failed to get the client: %v", err)
	mp := &MessageProcessor{*client, queueUrl}
	mp.pollQueue()
}

func (mp *MessageProcessor) pollQueue() {
	for {
		fmt.Println("Long polling for a message... (Will wait for 10s)")

		// Fetch some messages, hide it for 10 seconds.
		params := &sqs.ReceiveMessageInput{
			MaxNumberOfMessages: aws.Int64(MaxNumberOfMessage),
			QueueUrl:            aws.String(mp.sqsQueue),
			WaitTimeSeconds:     aws.Int64(WaitTimeSecond),
		}
		resp, err := mp.client.ReceiveMessage(params)

		if err != nil {
			fatalfIfErr("Error: Failed to read the message: %v", err)
		}

		// Sleep a little if we didn't find any messages in this poll.
		if len(resp.Messages) < 1 {
			fmt.Println("No messages on queue. Will sleep for 30s, then long poll again.")
			time.Sleep(30 * time.Second)
		}

		// Iterate over the messages we received, and dispatch a processor for each.
		for _, message := range resp.Messages {
			go mp.processMessage(message)
		}
	}
}

func (mp *MessageProcessor) processMessage(message *sqs.Message) {
	fmt.Println(*message.Body)
	params := &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(mp.sqsQueue),
		ReceiptHandle: message.ReceiptHandle,
	}
	_, err := mp.client.DeleteMessage(params)
	fatalfIfErr("Error: while deleting message from queue: %v", err)
}

func sqsClient(region string) (*sqs.SQS, error) {
	awsSession, err := session.NewSession(aws.NewConfig().WithRegion(region))
	fatalfIfErr("Error: Failed to fetch aws session: %v", err)
	//noinspection ALL
	sqsClientxx := sqs.New(awsSession)
	return sqsClientxx, err
}
