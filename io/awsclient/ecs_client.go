package awsclient

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	ecsUtils "github.com/reeceappling/goUtils/v2/ecs"
	"github.com/reeceappling/goUtils/v2/logging"
	"github.com/reeceappling/goUtils/v2/utils"
	"os"
)

const EcsClientKey = "ecs-client-key"

//go:generate mockery --name EcsClient
type EcsClient interface {
	StopTask(ctx context.Context, params *ecs.StopTaskInput, optFns ...func(*ecs.Options)) (*ecs.StopTaskOutput, error)
}

func GetEcsClient(ctx context.Context) (context.Context, EcsClient, error) {
	if existingClient, ok := ctx.Value(EcsClientKey).(EcsClient); ok {
		return ctx, existingClient, nil
	}
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, nil, err
	}
	var client EcsClient = ecs.NewFromConfig(cfg) // ensures we don't add something that doesn't work with the interface
	return context.WithValue(ctx, EcsClientKey, client), client, nil
}

func StopThisTask(ctx context.Context) {
	_, client, err := GetEcsClient(ctx)
	hardStopTaskOnError(ctx, err)
	md, err := ecsUtils.GetTaskMetadata()
	hardStopTaskOnError(ctx, err)
	_, err = client.StopTask(ctx, &ecs.StopTaskInput{
		Cluster: utils.Pointer(md.Cluster),
		Task:    utils.Pointer(md.TaskId()),
	})
	hardStopTaskOnError(ctx, err)
}

func hardStopTaskOnError(ctx context.Context, err error) {
	if err != nil {
		log := logging.GetSugaredLogger(ctx)
		log.Errorw("Error during task shutdown, hard stopping task", "error", err) // could just use log.Fatal here, but that doesn't seem to sync?
		_ = log.Sync()
		os.Exit(42)
	}
}
