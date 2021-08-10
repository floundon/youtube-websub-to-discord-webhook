package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/floundon/youtube-websub-to-discord-webhook/config"
	"github.com/floundon/youtube-websub-to-discord-webhook/handler"
	"github.com/gin-gonic/gin"
	"github.com/guregu/dynamo"
	"log"
)

var ginLambda *ginadapter.GinLambda

func init() {
	r := gin.Default()

	db := dynamo.New(session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region: aws.String("ap-northeast-1"),
		},
	})))
	webSubHandler := &handler.WebSubHandler{
		YouTubeVideoDataTable: db.Table(config.Get().YouTubeVideoDataTableName),
	}

	r.GET("/websub/subscribe", webSubHandler.VerifySubscription)
	r.POST("/websub/subscribe", webSubHandler.ReceiveNotification)

	ginLambda = ginadapter.New(r)
}

func ginHandler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(ginHandler)
}
