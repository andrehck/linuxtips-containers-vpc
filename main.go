package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func main() {
	var toogleReceiveMessage bool
	// Defina as URLs das filas de origem e destino
	sourceQueueURL := "https://sqs.us-east-1.amazonaws.com/014419883644/teste"
	destinationQueueURL := "https://sqs.us-east-1.amazonaws.com/014419883644/teste2"

	// Carregar a configuração padrão da AWS
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		log.Fatalf("Erro ao carregar a configuração da AWS: %v", err)
	}

	// Crie o cliente SQS
	svc := sqs.NewFromConfig(cfg)

	for {
		toogleReceiveMessage = true
		fmt.Println("O toogle está", toogleReceiveMessage)
		if toogleReceiveMessage == true {
			// Recebe a mensagem da fila de origem
			receiveMsgOutput, err := svc.ReceiveMessage(context.TODO(), &sqs.ReceiveMessageInput{
				QueueUrl:            aws.String(sourceQueueURL),
				MaxNumberOfMessages: 10, // Número máximo de mensagens a receber
				WaitTimeSeconds:     5,  // Tempo de espera para o recebimento (Long Polling)
			})

			if err != nil {
				log.Fatalf("Erro ao receber mensagem: %v", err)
			}

			// Verifica se há mensagens na fila
			if len(receiveMsgOutput.Messages) == 0 {
				fmt.Println("Nenhuma mensagem recebida")
				continue
			}

			for _, message := range receiveMsgOutput.Messages {
				// Copia a mensagem para a fila de destino
				_, err := svc.SendMessage(context.TODO(), &sqs.SendMessageInput{
					QueueUrl:    aws.String(destinationQueueURL),
					MessageBody: message.Body,
				})

				if err != nil {
					log.Fatalf("Erro ao enviar mensagem para a fila de destino: %v", err)
				}

				// Exclui a mensagem da fila de origem após o envio para a fila de destino
				_, err = svc.DeleteMessage(context.TODO(), &sqs.DeleteMessageInput{
					QueueUrl:      aws.String(sourceQueueURL),
					ReceiptHandle: message.ReceiptHandle,
				})

				if err != nil {
					log.Fatalf("Erro ao excluir mensagem da fila de origem: %v", err)
				}

				fmt.Printf("Mensagem copiada e excluída: %s\n", *message.Body)
			}
		}
	}
}
