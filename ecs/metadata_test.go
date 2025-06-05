package ecs

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestEmptyEnvVar(t *testing.T) {
	want := "environment variable {ECS_CONTAINER_METADATA_URI_V4} is not set"
	_, err := GetTaskMetadata()
	assert.Equal(t, err.Error(), want)
}

func TestFetchTaskMetadata(t *testing.T) {
	want := "mock_taskArn"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprintln(w, `{"TaskARN" : "mock_taskArn"}`)
	}))
	defer ts.Close()
	_ = os.Setenv("ECS_CONTAINER_METADATA_URI_V4", ts.URL)
	taskMeta, err := GetTaskMetadata()
	if err != nil {
		t.Log(err.Error())
	}
	assert.Equal(t, taskMeta.TaskARN, want)
}

func TestJsonParseFailure(t *testing.T) {
	want := "invalid character 'b' looking for beginning of object key string"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprintln(w, `{ bad_Json }`)
	}))
	defer ts.Close()
	err := os.Setenv("ECS_CONTAINER_METADATA_URI_V4", ts.URL)
	_, err = GetTaskMetadata()
	assert.Equal(t, err.Error(), want)
}

func TestUnmarshallEmpty(t *testing.T) {
	want := ""
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprintln(w, `{"field2" : "mock_task_arn"}`)
	}))
	defer ts.Close()
	err := os.Setenv("ECS_CONTAINER_METADATA_URI_V4", ts.URL)
	taskMeta, err := GetTaskMetadata()
	if err != nil {
		t.Log(err.Error())
	}
	assert.Equal(t, taskMeta.TaskARN, want)
}

func TestInstanceId(t *testing.T) {
	want := "c"
	instId := TaskMetadata{TaskARN: "a/b/c"}.TaskId()
	assert.Equal(t, instId, want)
}

func TestInstanceIdEmpty(t *testing.T) {
	want := ""
	instId := TaskMetadata{}.TaskId()
	assert.Equal(t, instId, want)
}
