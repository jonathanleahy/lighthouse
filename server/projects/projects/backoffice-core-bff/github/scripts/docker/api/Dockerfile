# build stage
FROM 231406070346.dkr.ecr.us-east-1.amazonaws.com/golden-images-builder:golang-1-21 as builder
WORKDIR /go/src/github.com/pismo/backoffice-core-bff/
COPY ./ /go/src/github.com/pismo/backoffice-core-bff/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags "-X main.gitVersion=$(git rev-parse HEAD ) -w -extldflags "-static""  -o application cmd/api/main.go && mv application /application

# run stage
FROM 231406070346.dkr.ecr.us-east-1.amazonaws.com/golden-images-runner:binary-0-0-1 as runner
COPY --from=builder --chown=999 /application /application
EXPOSE 8080
