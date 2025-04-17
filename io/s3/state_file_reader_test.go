package s3

import (
	"context"
	"fmt"
	"github.com/reeceappling/goUtils/v2/io/awsclient"
	"github.com/stretchr/testify/assert"
	"testing"
)

var bucket = "tempBucket"
var path = "dir/state.json"

func BenchmarkStateFileReader(b *testing.B) {
	for idx := 1; idx < 64; idx <<= 1 {
		b.Run(fmt.Sprint(idx), func(b *testing.B) {
			ctx := context.Background()
			awsclient.SetupWithDefault(ctx)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := (&S3FileReader{
					Bucket: bucket,
				}).RaceReadN(ctx, path, idx)

				assert.NoError(b, err)
			}
		})
	}

}

func BenchmarkS3FileReader(b *testing.B) {
	ctx := context.Background()
	awsclient.SetupWithDefault(ctx)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := (&S3FileReader{
			Bucket: bucket,
		}).Read(ctx, path)

		assert.NoError(b, err)
	}
}
