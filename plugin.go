package main

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/AlertFlow/runner/pkg/executions"
	"github.com/AlertFlow/runner/pkg/models"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type WaitPlugin struct{}

func (p *WaitPlugin) Init() models.Plugin {
	return models.Plugin{
		Name:    "Wait",
		Type:    "action",
		Version: "1.0.5",
		Creator: "JustNZ",
	}
}

func (p *WaitPlugin) Details() models.PluginDetails {
	params := []models.Param{
		{
			Key:         "WaitTime",
			Type:        "number",
			Default:     10,
			Required:    true,
			Description: "The time to wait in seconds",
		},
	}

	paramsJSON, err := json.Marshal(params)
	if err != nil {
		log.Error(err)
	}

	return models.PluginDetails{
		Action: models.ActionDetails{
			ID:          "wait",
			Name:        "Wait",
			Description: "Waits for a specified amount of time",
			Icon:        "solar:clock-circle-broken",
			Type:        "wait",
			Category:    "General",
			Function:    p.Execute,
			Params:      json.RawMessage(paramsJSON),
		},
	}
}

func (p *WaitPlugin) Execute(execution models.Execution, flow models.Flows, payload models.Payload, steps []models.ExecutionSteps, step models.ExecutionSteps, action models.Actions) (data map[string]interface{}, finished bool, canceled bool, no_pattern_match bool, failed bool) {
	// get the waittime from the action params
	waitTime := 10
	for _, param := range action.Params {
		if param.Key == "WaitTime" {
			waitTime, _ = strconv.Atoi(param.Value)
		}
	}

	err := executions.UpdateStep(execution.ID.String(), models.ExecutionSteps{
		ID:             step.ID,
		ActionID:       action.ID.String(),
		ActionMessages: []string{`Waiting for ` + strconv.Itoa(waitTime) + ` seconds`},
		Pending:        false,
		Paused:         true,
		StartedAt:      time.Now(),
	})
	if err != nil {
		return nil, false, false, false, true
	}

	executions.SetToPaused(execution)

	time.Sleep(time.Duration(waitTime) * time.Second)

	executions.SetToRunning(execution)

	err = executions.UpdateStep(execution.ID.String(), models.ExecutionSteps{
		ID:             step.ID,
		ActionMessages: []string{"Wait Action finished"},
		Paused:         false,
		Finished:       true,
		FinishedAt:     time.Now(),
	})
	if err != nil {
		return nil, false, false, false, true
	}

	return nil, true, false, false, false
}

func (p *WaitPlugin) Handle(context *gin.Context) {}

var Plugin WaitPlugin
