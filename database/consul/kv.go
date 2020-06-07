// consul键值对操作

package consul

import (
	"errors"
	"github.com/hashicorp/consul/api"
	"github.com/laconiz/metis/utils/json"
)

var ErrNotFound = errors.New("key not found")

type KV struct {
	kv *api.KV
}

func (kv *KV) KV() *api.KV {
	return kv.kv
}

func (kv *KV) Get(key string, value interface{}) (bool, error) {

	pair, _, err := kv.kv.Get(key, nil)
	if err != nil {
		return false, err
	}

	if pair == nil {
		return false, nil
	}

	err = json.Unmarshal(pair.Value, value)
	return true, err
}

func (kv *KV) Set(key string, value interface{}) error {

	raw, err := json.Marshal(value)
	if err != nil {
		return err
	}

	pair := &api.KVPair{Key: key, Value: raw}
	_, err = kv.kv.Put(pair, nil)
	return err
}

func (kv *KV) Delete(key string) error {
	_, err := kv.kv.Delete(key, nil)
	return err
}

func (kv *KV) List(prefix string) (api.KVPairs, error) {
	pairs, _, err := kv.kv.List(prefix, nil)
	return pairs, err
}

func (kv *KV) ListLoose(prefix string, receiver interface{}) error {

	pairs, err := kv.List(prefix)
	if err != nil {
		return err
	}

	return Pairs(prefix, pairs, receiver, Loose)
}

func (kv *KV) ListStrict(prefix string, receiver interface{}) error {

	pairs, err := kv.List(prefix)
	if err != nil {
		return err
	}

	return Pairs(prefix, pairs, receiver, Strict)
}
