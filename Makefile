SRCS=cmd/lambda/main.go

build: template.yml $(SRCS)
	sam build

.aws-sam: build

package: .aws-sam
	sam package --s3-bucket floundon-sam-package > packaged.yml

packaged.yml: package

deploy: packaged.yml
	sam deploy --template-file packaged.yml --stack-name youtube-websub-to-discord-webhook-stack --capabilities CAPABILITY_IAM

tidy:
	go mod tidy
	go mod verify