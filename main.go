package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
)

// .
var ThreadCount = 100

func main() {
	filePath := flag.String("s", "", "entrypoints filename")
	topic := flag.String("t", "", "topic arn")
	newTopic := flag.String("n", "", "new topic name")
	protocol := flag.String("p", "application", "sns protocol")
	region := flag.String("r", "ap-northeast-1", "aws region")
	flag.Parse()
	if *filePath == "" || (*topic == "" && *newTopic == "") || *protocol == "" || *region == "" {
		flag.PrintDefaults()
		return
	}
	if *newTopic != "" {
		svc := sns.New(&aws.Config{Region: aws.String(*region)})
		params := &sns.CreateTopicInput{Name: aws.String(*newTopic)}
		out, err := svc.CreateTopic(params)
		if err != nil {
			log.Println(err.Error())
			return
		}
		topic = out.TopicArn
		log.Println("new arn: " + *topic)
	}

	var finCh = make(chan int)
	var endpointCh = make(chan string)
	for i := 0; i < ThreadCount; i++ {
		go func() {
			svc := sns.New(&aws.Config{Region: aws.String(*region)})
			for {
				select {
				case endpoint := <-endpointCh:
					params := &sns.SubscribeInput{
						Protocol: aws.String(*protocol), // Required
						TopicArn: aws.String(*topic),    // Required
						Endpoint: aws.String(endpoint),
					}
					_, err := svc.Subscribe(params)

					if err != nil {
						if !strings.Contains(err.Error(), "Endpoint does not exist") {
							log.Println(err.Error())
						}
					}
				case <-finCh:
					break
				}
			}
		}()
	}

	fp, err := os.Open(*filePath)
	if err != nil {
		panic(err)
	}
	defer fp.Close()
	scanner := bufio.NewScanner(fp)
	num := 0
	for scanner.Scan() {
		endpointCh <- scanner.Text()
		num++
		if num%1000 == 0 {
			log.Printf("%d\n", num)
		}
	}
	for i := 0; i < ThreadCount; i++ {
		finCh <- 0
	}
	log.Println("fin all")
	log.Println("arn: " + *topic)
}
