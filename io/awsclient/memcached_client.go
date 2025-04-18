package awsclient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	ssmv2 "github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/bradfitz/gomemcache/memcache"
	"gocloud.dev/runtimevar"
	"gocloud.dev/runtimevar/awsparamstore"
	"os"
	"sync"
	"time"
)

//go:generate mockery --name MemcachedClient
type MemcachedClient interface {
	Get(key string) (item *memcache.Item, err error)
	Set(item *memcache.Item) error
}

type UninitializedMemcachedClient struct {
}

func (cache UninitializedMemcachedClient) Get(key string) (item *memcache.Item, err error) {
	return nil, errors.New("tile cache is not initialized")
}
func (cache UninitializedMemcachedClient) Set(item *memcache.Item) error {
	return nil
}

var memcachedClient MemcachedClient = UninitializedMemcachedClient{}
var endpointsMut = &sync.RWMutex{}

func getClient() MemcachedClient {
	endpointsMut.RLock()
	defer endpointsMut.RUnlock()
	return memcachedClient
}
func setClient(endpoints []string) {
	newClient := memcache.New(endpoints...)
	newClient.Timeout = memcachedTimeout

	endpointsMut.Lock()
	defer endpointsMut.Unlock()
	memcachedClient = newClient
}

const memcachedTimeout = 5 * time.Second

func init() {
	varName := "/temp/memcached-addresses"
	options := awsparamstore.Options{
		WaitDuration: time.Minute * 2,
	}
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "could not init cache, not using cache")
	}
	client := ssmv2.NewFromConfig(cfg)
	go func() {
		if v, err := awsparamstore.OpenVariableV2(client, varName, runtimevar.StringDecoder, &options); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "failed to watch memcached endpoint param")
		} else {
			defer v.Close() //nolint:errcheck
			for {
				snapShot, err := v.Watch(context.Background())
				if err != nil {
					continue
				}
				endpointsString, ok := snapShot.Value.(string)
				if !ok {
					continue
				}
				var parsedEndpoints []string
				err = json.Unmarshal([]byte(endpointsString), &parsedEndpoints)
				if err != nil {
					continue
				}
				setClient(parsedEndpoints)
			}
		}
	}()
}

var MemcachedClientKey = "memcached-client-key"

func GetMemcachedClient(ctx context.Context) (context.Context, MemcachedClient) {
	if existingClient, ok := ctx.Value(MemcachedClientKey).(MemcachedClient); ok {
		return ctx, existingClient
	}
	mc := getClient()
	return context.WithValue(ctx, MemcachedClientKey, mc), mc
}
