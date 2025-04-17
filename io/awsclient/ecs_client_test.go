package awsclient

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	ecsUtils "github.com/reeceappling/goUtils/v2/ecs"
	"github.com/reeceappling/goUtils/v2/io/awsclient/mocks"
	"github.com/reeceappling/goUtils/v2/utils"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var clusterName = "some-cluster"
var taskArn = "arn-stuff-goes-here/some-task-id"

func TestStopThisTask(t *testing.T) {
	ctx := context.Background()
	ts := setupMetadataServer(ctx)
	defer ts.Close()

	t.Run("should stop the task", func(t *testing.T) {
		mockEcsClient := mocks.NewEcsClient(t)
		mockEcsClient.On("StopTask", mock.Anything, &ecs.StopTaskInput{
			Cluster: &clusterName,
			Task:    utils.Pointer("a-task-id"),
		}).Return(nil, nil)
		ctx := context.WithValue(ctx, EcsClientKey, mockEcsClient)
		StopThisTask(ctx)
	})
}

func setupMetadataServer(ctx context.Context) *httptest.Server {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		bs, err := json.Marshal(ecsUtils.TaskMetadata{
			Cluster: clusterName,
			TaskARN: taskArn,
		})
		must(err)
		_, err = fmt.Fprintln(w, string(bs))
		must(err)
	}))
	must(os.Setenv("ECS_CONTAINER_METADATA_URI_V4", ts.URL))
	return ts
}
func must(err error) {
	if err != nil {
		panic(err)
	}
}
