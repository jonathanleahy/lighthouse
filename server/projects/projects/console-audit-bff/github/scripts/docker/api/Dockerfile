# build stage
FROM 231406070346.dkr.ecr.us-east-1.amazonaws.com/golden-images-builder:golang-1-21 as builder

ARG GIT_ACCESS_TOKEN_USER
ARG GIT_ACCESS_TOKEN
ARG GOPRIVATE=github.com/pismo/*

# Add access to private libs on pismo's github
RUN apk add git && \
    git config --global url."https://${GIT_ACCESS_TOKEN_USER}:${GIT_ACCESS_TOKEN}@github.com".insteadOf https://github.com

WORKDIR /go/src/github.com/pismo/console-audit-bff/
COPY ./ /go/src/github.com/pismo/console-audit-bff/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags "-X main.gitVersion=$(git rev-parse HEAD ) -w -extldflags "-static""  -o application cmd/api/main.go && mv application /application

# run stage
FROM 231406070346.dkr.ecr.us-east-1.amazonaws.com/golden-images-runner:binary-0-0-1 as runner
COPY --from=builder --chown=999 /application /application
EXPOSE 8080
