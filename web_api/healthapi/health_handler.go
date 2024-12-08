package healthapi

import (
	"context"
	"net/http"

	"github.com/bioform/go-web-app-template/pkg/database"
	"github.com/bioform/go-web-app-template/pkg/mail"
)

type HealthOutput struct {
	Body struct {
		Dependencies map[string]interface{}
	}
	Status int
}

// Generalize health checks for multiple dependencies
type Dependency struct {
	Name  string
	Check func(ctx context.Context) (interface{}, error)
}

func HealthHandler(ctx context.Context, _ *struct{}) (*HealthOutput, error) {
	dependencies := []Dependency{
		{
			Name: "Database",
			Check: func(ctx context.Context) (interface{}, error) {
				// Replace with actual database health check logic
				dbHealth := database.Health(ctx)
				return dbHealth, nil // Return the entire dbHealth object
			},
		},
		{
			Name: "SMTP",
			Check: func(ctx context.Context) (interface{}, error) {
				// Replace with actual SMTP health check logic
				smtpHealth := map[string]string{
					"status": "down",
				}

				if err := mail.Client().DialAndSendWithContext(ctx); err == nil {
					smtpHealth["status"] = "up"
				}

				return smtpHealth, nil
			},
		},
		// Add more dependencies as needed
	}

	result := &HealthOutput{}
	result.Body.Dependencies = make(map[string]interface{})
	allUp := true

	for _, dep := range dependencies {
		data, err := dep.Check(ctx)
		if err != nil {
			allUp = false
		}

		// Determine status from data
		status := extractStatus(data)
		if status != "up" {
			allUp = false
		}

		result.Body.Dependencies[dep.Name] = data
	}

	if allUp {
		result.Status = http.StatusOK
	} else {
		result.Status = http.StatusInternalServerError
	}

	return result, nil
}

func extractStatus(data interface{}) string {
	switch v := data.(type) {
	case map[string]string:
		return v["status"]
	case map[string]interface{}:
		if s, ok := v["status"].(string); ok {
			return s
		}
	case database.DBHealthStats:
		return v.Status
	}
	return "unknown"
}
