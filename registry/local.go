package registry

import (
	"errors"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
)

type registryLocal struct {
	prefix string
	once   sync.Once
}

func absDir(fp string) string {
	if path.IsAbs(fp) {
		return fp
	}
	homeFlag := "~"
	if strings.HasPrefix(fp, homeFlag) {
		homeUrl, _ := os.UserHomeDir()
		return path.Join(homeUrl, strings.TrimLeft(fp, homeFlag))
	} else {
		absDp, _ := filepath.Abs(fp)
		return absDp
	}
}

func (r *registryLocal) addPrefix(key string) string {
	return path.Join(r.prefix, key)
}

func (r *registryLocal) delPrefix(key string) string {
	if r.prefix == key {
		return path.Base(key)
	}
	return strings.TrimPrefix(strings.TrimPrefix(key, r.prefix), "/")
}

func (r *registryLocal) Init(config *Config) error {
	var err error
	r.once.Do(func() {
		r.prefix = absDir(config.Prefix)
		_, err = os.Stat(r.prefix)
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *registryLocal) Get(key string) ([]byte, error) {
	absFp := path.Join(r.prefix, key)
	return ioutil.ReadFile(absFp)
}

func (r *registryLocal) Set(key string, value []byte) error {
	key = r.addPrefix(key)
	absDp := absDir(path.Dir(key))
	if _, err := os.Stat(absDp); errors.Is(err, fs.ErrNotExist) {
		err := os.MkdirAll(absDp, os.ModePerm)
		if err != nil {
			return err
		}
	}
	file, err := os.OpenFile(key, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(value)
	return err
}

func (r *registryLocal) Delete(key string) error {
	key = r.addPrefix(key)
	return os.Remove(key)
}

func (r *registryLocal) List(prefix string) (KVPairs, error) {
	ret := KVPairs{}
	prefix = r.addPrefix(prefix)
	err := filepath.Walk(prefix, func(path string, info fs.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		if b, err := ioutil.ReadFile(path); err == nil {
			ret = append(ret, &KVPair{
				Key:   r.delPrefix(path),
				Value: b,
			})
		}
		return nil
	})
	return ret, err
}
