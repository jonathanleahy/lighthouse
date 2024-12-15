package main

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/pismo/backoffice-core-bff/internal/app/adapter/web/handler"
	"github.com/pismo/backoffice-core-bff/internal/app/domain/service"
	"github.com/pismo/backoffice-core-bff/internal/app/infrastructure/aws/sns"
	"github.com/pismo/backoffice-core-bff/internal/app/infrastructure/env"
	"github.com/pismo/backoffice-core-bff/internal/app/infrastructure/logger"
	"github.com/pismo/backoffice-core-bff/internal/app/infrastructure/network"
	"github.com/pismo/backoffice-core-bff/internal/app/infrastructure/opentelemetry"
	"github.com/pismo/backoffice-core-bff/internal/app/infrastructure/server"
	"go.opentelemetry.io/otel/metric/global"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func main() {
	logger.Init()
	defer logger.Defer()
	ctx := context.Background()
	setTracer(ctx)
	setMetric(ctx)
	start(setSns())
}

func start(auditService service.AuditService) {
	httpService := setHttpService(auditService)
	handler.SetResolver(service.NewResolver(httpService))
	err := <-server.Start(httpService)
	if err == nil {
		logger.Warn("Graceful Shutdown", "", "", nil)
	} else {
		logger.Panic("Application down", "", "", logger.Fields{"error": err})
	}
	auditService.WaitFinish()
}

func setSns() service.AuditService {
	awsSession, err := session.NewSession(&aws.Config{Region: aws.String(env.GetEnvWithDefaultAsString(env.AwsRegion, env.DefaultAwsRegion))})
	if err != nil {
		zap.L().Error("failed to start SNS service", zap.NamedError("Attributes.error", err))
	}
	eventSender := sns.NewAwsSNSClient(env.GetEnvWithDefaultAsString(env.SnsConsoleAudit, env.DefaultSnsConsoleAudit), awsSession)
	return service.NewAuditService(eventSender)
}

func setHttpService(auditService service.AuditService) network.HttpService {
	duration := env.GetEnvWithDefaultAsInt(env.HttpTimeout, env.DefaultHttpTimeout)
	timeout := time.Second * time.Duration(duration)
	client := &http.Client{Timeout: timeout}
	return network.NewHttpService(client, auditService)
}

const AttributesError = "Attributes.error"

func setTracer(ctx context.Context) {
	tracerProvider := opentelemetry.InitTracer(ctx)
	defer func() {
		if err := tracerProvider.Shutdown(ctx); err != nil {
			zap.L().Error("error on opentelemetry", zap.NamedError(AttributesError, err))
		}
	}()
}

func setMetric(ctx context.Context) {
	metricController := opentelemetry.InitMetric(ctx)
	defer func() {
		_ = metricController.Stop(ctx)
	}()
	global.SetMeterProvider(metricController.MeterProvider())

	//if err := runtime.Start(runtime.WithMinimumReadMemStatsInterval(time.Second)); err != nil {
	//	zap.L().Error("failed to start runtime instrumentation:", zap.NamedError(AttributesError, err))
	//}
}

