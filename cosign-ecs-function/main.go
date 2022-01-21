package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
)

func handler(event events.CloudWatchEvent) {

	var eventDetail Detail
	err := json.Unmarshal(event.Detail, &eventDetail)
	if err != nil {
		log.Fatalf("[ERROR] %v error during event unmarshalling: %v", event.ID, err)
	}

	lambdaEvent := LambdaEvent{
		Version:    event.Version,
		ID:         event.ID,
		DetailType: event.DetailType,
		Source:     event.Source,
		Account:    event.AccountID,
		Time:       event.Time,
		Region:     event.Region,
		Resources:  event.Resources,
		Detail:     eventDetail,
	}

	log.Printf("Cluster: %v\n", lambdaEvent.Detail.ClusterArn)
	log.Printf("taskArn: %v\n", lambdaEvent.Detail.TaskArn)
	log.Printf("taskDefinitionArn: %v\n", lambdaEvent.Detail.TaskDefinitionArn)
	log.Printf("accountId: %v\n", lambdaEvent.Account)

	var verified bool
	for i := 0; i < len(lambdaEvent.Detail.Containers); i++ {
		log.Printf("container Image %v : %v", i, lambdaEvent.Detail.Containers[i].Image)
		verified, err = Verify(lambdaEvent.Detail.Containers[i].Image, lambdaEvent.Region, lambdaEvent.Account)
		if err != nil {
			log.Printf("[ERROR] %v error Verifing image: %v %v", event.ID, verified, err)
		}
		if !verified {
			log.Println("NOT VERIFIED")
		} else {
			log.Println("VERIFIED")
		}
	}
}

func main() {
	lambda.Start(handler)
}
