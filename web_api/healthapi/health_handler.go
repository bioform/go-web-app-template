package healthapi

import (
	"context"
	"net/http"

	"github.com/bioform/go-web-app-template/pkg/database"
)

type HealthOutput struct {
	Body struct {
		Database database.DBHealthStats
	}
	Status int
}

func HealthHandler(ctx context.Context, _ *struct{}) (*HealthOutput, error) {
	dbHealth := database.Health(ctx)

	result := &HealthOutput{}
	result.Body.Database = dbHealth

	if dbHealth.Status == "up" {
		result.Status = http.StatusOK
	} else {
		result.Status = http.StatusInternalServerError
	}

	return result, nil
}
