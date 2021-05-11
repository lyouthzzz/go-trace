package attribute

type KeyValue struct {
	Key   string
	Value interface{}
}

func New(key string, val interface{}) KeyValue {
	return KeyValue{Key: key, Value: val}
}

func InputBody(i []byte) KeyValue {
	return KeyValue{Key: "input_body", Value: i}
}

func OutputBody(o []byte) KeyValue {
	return KeyValue{Key: "output_body", Value: o}
}

func GlobalTicketId(id string) KeyValue {
	return KeyValue{Key: "globalTicket", Value: id}
}

func MonitorId(id string) KeyValue {
	return KeyValue{Key: "monitorId", Value: id}
}

func ParentRpcId(id string) KeyValue {
	return KeyValue{Key: "parentRpcId", Value: id}
}

func RpcId(id string) KeyValue {
	return KeyValue{Key: "rpcId", Value: id}
}
