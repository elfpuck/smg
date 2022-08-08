package registry

import (
	"fmt"
	"path"
	"strings"
	"sync"

	"github.com/hashicorp/consul/api"
)

type registryConsul struct {
	client *api.KV
	prefix string
	once   sync.Once
}

func (r *registryConsul) Init(config *Config) error {
	var err error
	var cli *api.Client
	r.once.Do(func() {
		cli, err = api.NewClient(&api.Config{
			Address: config.Address,
			Token:   config.Token,
		})
		if err != nil {
			return
		}
		r.client = cli.KV()
		r.prefix = config.Prefix
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *registryConsul) addPrefix(key string) string {
	return strings.TrimPrefix(path.Join(r.prefix, key), "/")
}

func (r *registryConsul) delPrefix(key string) string {
	return strings.TrimPrefix(strings.TrimPrefix(key, r.prefix), "/")
}

func (r *registryConsul) Get(key string) ([]byte, error) {
	key = r.addPrefix(key)
	kv, _, err := r.client.Get(key, nil)
	if err != nil {
		return nil, err
	}
	if kv == nil {
		return nil, fmt.Errorf("key ( %s ) was not found", key)
	}
	return kv.Value, nil
}

func (r *registryConsul) Set(key string, value []byte) error {
	key = r.addPrefix(key)
	kv := &api.KVPair{
		Key:   key,
		Value: value,
	}
	_, err := r.client.Put(kv, nil)
	return err
}

func (r *registryConsul) Delete(key string) error {
	key = r.addPrefix(key)
	_, err := r.client.Delete(key, nil)
	return err
}

func (r *registryConsul) List(prefix string) (KVPairs, error) {
	prefix = r.addPrefix(prefix)
	pairs, _, err := r.client.List(prefix, nil)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	ret := make(KVPairs, len(pairs))
	for i, kv := range pairs {
		ret[i] = &KVPair{Key: r.delPrefix(kv.Key), Value: kv.Value}
	}
	return ret, nil
}
