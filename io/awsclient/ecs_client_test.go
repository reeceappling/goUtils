package awsclient

//var clusterName = "some-cluster" // TODO: REENABLE LATER
//var taskArn = "arn-str-goes-here/a-task-id"
//
//func TestStopThisTask(t *testing.T) {
//	ctx := context.Background()
//	ts := setupMetadataServer(ctx)
//	defer ts.Close()
//
//	t.Run("should stop the task", func(t *testing.T) {
//		mockEcsClient := mocks.NewEcsClient(t)
//		mockEcsClient.On("StopTask", mock.Anything, &ecs.StopTaskInput{
//			Cluster: &clusterName,
//			Task:    utils.Pointer("a-task-id"),
//		}).Return(nil, nil)
//		ctx := context.WithValue(ctx, EcsClientKey, mockEcsClient)
//		StopThisTask(ctx)
//	})
//}
//
//func setupMetadataServer(ctx context.Context) *httptest.Server {
//	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		w.Header().Set("Content-Type", "application/json")
//		bs, err := json.Marshal(ecsUtils.TaskMetadata{
//			Cluster: clusterName,
//			TaskARN: taskArn,
//		})
//		if err != nil {
//			panic("marshal fail: " + err.Error())
//		}
//		_, err = fmt.Fprintln(w, string(bs))
//		if err != nil {
//			panic("print fail: " + err.Error())
//		}
//	}))
//	must(os.Setenv("ECS_CONTAINER_METADATA_URI_V4", ts.URL))
//	return ts
//}
//func must(err error) {
//	if err != nil {
//		panic(err)
//	}
//}
