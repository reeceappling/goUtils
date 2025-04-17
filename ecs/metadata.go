package ecs

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"strings"
)

func GetTaskMetadata() (TaskMetadata, error) {
	v4Uri, isSet := os.LookupEnv("ECS_CONTAINER_METADATA_URI_V4")

	if !isSet {
		return TaskMetadata{}, errors.New("environment variable {ECS_CONTAINER_METADATA_URI_V4} is not set")
	}

	fetch, err := http.Get(v4Uri + "/task")
	if err != nil {
		return TaskMetadata{}, err
	}
	defer fetch.Body.Close()
	jsonBody, err := io.ReadAll(fetch.Body)
	if err != nil {
		return TaskMetadata{}, err
	}
	taskMetadata := TaskMetadata{}
	err = json.Unmarshal(jsonBody, &taskMetadata)
	return taskMetadata, err
}

type TaskMetadata struct {
	Cluster string
	TaskARN string
}

func (tm TaskMetadata) TaskId() string {
	pieces := strings.Split(tm.TaskARN, "/")
	return pieces[len(pieces)-1]
}
