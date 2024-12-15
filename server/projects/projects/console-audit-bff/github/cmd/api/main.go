package main

import (
	"context"

	_ "github.com/pismo/console-audit-bff/docs/openapi"
	"github.com/pismo/console-audit-bff/internal/app"
	"github.com/pismo/psm-sdk/psm"
	"github.com/pismo/psm-sdk/psm/attribute/value"
	"github.com/pismo/psm-sdk/psm/factory"
	"github.com/pismo/psm-sdk/psm/schema/data"
)

func main() {
	psm.Configure("console-audit-bff", factory.Default())
	defer psm.Shutdown()
	psm.Start()

	application, err := app.NewApplication(psm.HttpClient())
	if err != nil {
		psm.Logger().Panic(context.Background(), "Failed to build application", data.LogAttributes{
			"err": value.Plain(err),
		})
	}

	err = application.Start()
	if err != nil {
		psm.Logger().Warn(context.Background(), "Failed to start application", data.LogAttributes{
			"err": value.Plain(err),
		})
	}

}

