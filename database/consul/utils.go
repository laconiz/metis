package consul

import (
	"errors"
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/laconiz/metis/utils/json"
	"reflect"
	"strings"
)

const (
	Loose  = false
	Strict = true
)

func Pairs(prefix string, pairs api.KVPairs, receiver interface{}, strict bool) error {

	typo := reflect.TypeOf(receiver)
	if typo == nil || typo.Kind() != reflect.Ptr {
		return errors.New("receiver must be *map[string]Any")
	}
	typo = typo.Elem()

	if typo.Kind() != reflect.Map {
		return errors.New("receiver must be *map[string]Any")
	}

	value := reflect.ValueOf(receiver).Elem()
	value.Set(reflect.MakeMap(typo))

	for _, pair := range pairs {

		key := strings.Replace(pair.Key, prefix, "", 1)

		elem := reflect.New(typo.Elem())
		err := json.Unmarshal(pair.Value, elem.Interface())
		if err == nil {
			value.SetMapIndex(reflect.ValueOf(key), elem)
			continue
		}

		if strict {
			return fmt.Errorf("unmarshal %s:%s error: %w", pair.Key, string(pair.Value), err)
		}
	}

	return nil
}
