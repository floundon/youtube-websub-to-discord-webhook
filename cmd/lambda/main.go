package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/floundon/youtube-websub-to-discord-webhook/handler"
	"github.com/gin-gonic/gin"
	"log"
)

var ginLambda *ginadapter.GinLambda

func init() {
	log.Printf("Gin cold start")

	r := gin.Default()
	r.GET("/websub/subscribe", handler.VerifySubscription)
	r.POST("/websub/subscribe", handler.ReceiveNotification)

	ginLambda = ginadapter.New(r)
}

func ginHandler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(ginHandler)
}