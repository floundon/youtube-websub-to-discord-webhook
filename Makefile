SRCS=cmd/lambda/main.go

build: template.yml $(SRCS)
	sam build

.aws-sam: build

package: .aws-sam
	sam package --s3-bucket $(SAM_PACKAGE_BUCKET) > packaged.yml

packaged.yml: package

deploy: packaged.yml
	sam deploy --template-file packaged.yml --stack-name $(DEPLOY_STACK_NAME) --capabilities CAPABILITY_IAM

tidy:
	go mod tidy
	go mod verify