package packet

import (
	"errors"
	"fmt"
	"github.com/laconiz/metis/nets/codec"
	"hash/fnv"
	"reflect"
	"sort"
)

type MetaID uint32

type Meta struct {
	id    MetaID
	name  string
	typo  reflect.Type
	codec codec.Codec
}

func (meta *Meta) ID() MetaID {
	return meta.id
}

func (meta *Meta) Name() string {
	return meta.name
}

func (meta *Meta) Type() reflect.Type {
	return meta.typo
}

func (meta *Meta) Codec() codec.Codec {
	return meta.codec
}

func (meta *Meta) Encode(msg interface{}) ([]byte, error) {
	return meta.codec.Encode(msg)
}

func (meta *Meta) Decode(raw []byte) (interface{}, error) {
	msg := reflect.New(meta.typo).Interface()
	err := meta.codec.Decode(raw, msg)
	return msg, err
}

func Register(msg interface{}, codec codec.Codec) (*Meta, error) {

	typo := reflect.TypeOf(msg)
	if typo == nil {
		return nil, errors.New("register a nil message")
	}

	if typo.Kind() == reflect.Ptr {
		typo = typo.Elem()
	}

	name := typo.String()

	hash := fnv.New32()
	hash.Write([]byte(name))
	id := MetaID(hash.Sum32())

	return RegisterEx(id, name, msg, codec)
}

func RegisterEx(id MetaID, name string, msg interface{}, codec codec.Codec) (*Meta, error) {

	if meta := MetaByID(id); meta != nil {
		return nil, fmt.Errorf("conflict meta id: %s - %s", name, meta.name)
	}

	if meta := MetaByName(name); meta != nil {
		return nil, fmt.Errorf("conflict meta name: %s - %s", name, meta.name)
	}

	typo := reflect.TypeOf(msg)
	if typo == nil {
		return nil, fmt.Errorf("register a nil message: %s", name)
	}

	if typo.Kind() == reflect.Ptr {
		typo = typo.Elem()
	}

	if meta := MetaByType(typo); meta != nil {
		return nil, fmt.Errorf("conflict meta type: %s - %s", name, meta.name)
	}

	if codec == nil {
		return nil, fmt.Errorf("invalid codec: %s", name)
	}

	meta := &Meta{id: id, name: name, typo: typo, codec: codec}
	idMap[id] = meta
	nameMap[name] = meta
	typeMap[typo] = meta

	return meta, nil
}

func MetaByID(id MetaID) *Meta {
	return idMap[id]
}

func MetaByName(name string) *Meta {
	return nameMap[name]
}

func MetaByType(typo reflect.Type) *Meta {
	return typeMap[typo]
}

func MetaByMsg(msg interface{}) *Meta {

	typo := reflect.TypeOf(msg)
	if typo == nil {
		return nil
	}

	if typo.Kind() == reflect.Ptr {
		typo = typo.Elem()
	}

	return MetaByType(typo)
}

func MetaList() (ret []*Meta) {

	for _, meta := range nameMap {
		ret = append(ret, meta)
	}

	sort.Slice(ret, func(i, j int) bool {
		return ret[i].name < ret[j].name
	})

	return
}

var (
	idMap   = map[MetaID]*Meta{}
	nameMap = map[string]*Meta{}
	typeMap = map[reflect.Type]*Meta{}
)
