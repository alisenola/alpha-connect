package logging

import (
	"reflect"
	"time"
)

func NewMapEncoder() MapEncoder {
	return MapEncoder{M: make(map[string]interface{})}
}

type MapEncoder struct {
	M map[string]interface{}
}

func (me MapEncoder) EncodeBool(key string, val bool) {
	me.M[key] = val //= fmt.Sprintf("%t", val)
}

func (me MapEncoder) EncodeFloat64(key string, val float64) {
	me.M[key] = val
	//me.M[key] = fmt.Sprintf("%f", val)
}

func (me MapEncoder) EncodeInt(key string, val int) {
	me.M[key] = val
	//me.M[key] = fmt.Sprintf("%d", val)
}

func (me MapEncoder) EncodeInt64(key string, val int64) {
	me.M[key] = val
	//me.M[key] = fmt.Sprintf("%d", val)
}

func (me MapEncoder) EncodeDuration(key string, val time.Duration) {
	me.M[key] = val
	//me.M[key] = fmt.Sprintf("%s", val)
}

func (me MapEncoder) EncodeUint(key string, val uint) {
	me.M[key] = val
	//me.M[key] = fmt.Sprintf("%d", val)
}

func (me MapEncoder) EncodeUint64(key string, val uint64) {
	me.M[key] = val
	//me.M[key] = fmt.Sprintf("%d", val)
}

func (me MapEncoder) EncodeString(key string, val string) {
	me.M[key] = val
	//me.M[key] = fmt.Sprintf("%s", val)
}

func (me MapEncoder) EncodeObject(key string, val interface{}) {
	me.M[key] = val
	//me.M[key] = fmt.Sprintf("%q", val)
}

func (me MapEncoder) EncodeType(key string, val reflect.Type) {
	me.M[key] = val
	//me.M[key] = fmt.Sprintf("%v", val)
}
