all: build

build:
	docker build \
		-t slack-swanson:latest \
		-t 485304306779.dkr.ecr.eu-west-1.amazonaws.com/slack-swanson:latest .

run:
	docker run \
		-v ~/.aws-lambda-rie:/aws-lambda \
		-p 9000:8080 \
		--entrypoint /aws-lambda/aws-lambda-rie \
		-e SLACK_TOKEN \
		slack-swanson:latest /main