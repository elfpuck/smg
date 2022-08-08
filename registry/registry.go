package registry

type KVPairs []*KVPair

type KVPair struct {
	Key   string
	Value []byte
}

type RemoteProvider interface {
	Init(config *Config) error
	Get(key string) ([]byte, error)
	Set(key string, value []byte) error
	List(prefix string) (KVPairs, error)
	Delete(key string) error
}

type Config struct {
	Provider string
	Address  string
	Token    string
	Prefix   string
}

func AddRemoteProvider(config *Config) (RemoteProvider, error) {
	var rp RemoteProvider
	switch config.Provider {
	case "consul":
		rp = new(registryConsul)
	case "local":
		rp = new(registryLocal)
	}
	if rp != nil {
		if err := rp.Init(config); err != nil {
			return rp, err
		}
	}
	return rp, nil
}
