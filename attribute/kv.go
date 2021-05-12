package attribute

type Body []byte

func (b Body) String() string {
	return string(b)
}

type KeyValue struct {
	Key   string
	Value interface{}
}

func KV(key string, val interface{}) KeyValue {
	return KeyValue{Key: key, Value: val}
}

func InputBody(i Body) KeyValue {
	return KeyValue{Key: "input", Value: i}
}

func OutputBody(o Body) KeyValue {
	return KeyValue{Key: "output", Value: o}
}

func RpcType(rpcType string) KeyValue {
	return KV("rpcType", rpcType)
}
