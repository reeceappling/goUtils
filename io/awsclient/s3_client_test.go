package awsclient

//func TestDefaultS3Client(t *testing.T) { // TODO: reenable later
//	ctx := context.Background()
//	t.Run("StandardizeError", func(t *testing.T) {
//
//		t.Run("not an error", func(t *testing.T) {
//			assert.Nil(t, StandardizeError(ctx, nil))
//		})
//
//		t.Run("404", func(t *testing.T) {
//			assert.Equal(t,
//				errorreference.ErrorNotFound,
//				StandardizeError(ctx, utils.Pointer(types.NoSuchKey{})),
//			)
//		})
//
//		t.Run("429", func(t *testing.T) {
//			assert.Equal(t,
//				errorreference.ErrorSlowDown,
//				StandardizeError(ctx, errorreference.ErrorSlowDown), // these never seem to work, so we needed the more generic version below
//			)
//
//			assert.Equal(t,
//				errorreference.ErrorSlowDown,
//				StandardizeError(ctx, errors.New("operation error S3: GetObject, failed to get rate limit token, retry quota exceeded, 3 available, 5 requested")),
//			)
//		})
//
//		t.Run("empty s3 bucket", func(t *testing.T) {
//			assert.Equal(t,
//				ErrorUndefinedS3Bucket,
//				StandardizeError(ctx, errors.New("operation error S3: GetObject, input member Bucket must not be empty")),
//			)
//		})
//
//		t.Run("empty s3 key", func(t *testing.T) {
//			assert.Equal(t,
//				ErrorUndefinedS3Key,
//				StandardizeError(ctx, errors.New("operation error S3: GetObject, input member Key must not be empty")),
//			)
//		})
//	})
//}
