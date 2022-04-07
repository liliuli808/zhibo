package utils

import (
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/structpb"
)

func IsEmptyAnyMap(a *anypb.Any) bool {
	if a == nil {
		return true
	}
	// 判断
	return len(UnmarshalAny(a).(map[string]interface{})) == 0
}

func UnmarshalAny(a *anypb.Any) interface{} {
	value := &structpb.Value{}
	_ = anypb.UnmarshalTo(a, value, proto.UnmarshalOptions{})
	return value.AsInterface()
}

func MarshalAny(v interface{}) *anypb.Any {
	// 尝试
	v = ProtoValueCheck(v)
	// 底层转换
	m, err := structpb.NewValue(v)
	if err != nil {
		panic(err)
	}
	// 转换
	any, err := anypb.New(m)
	// 判断
	if err != nil {
		panic(err)
	}
	return any
}

func ProtoValueCheck(v interface{}) interface{} {
	// 兼容int8
	switch v.(type) {
	case int8:
		v = int(v.(int8))
	case int16:
		v = int32(v.(int16))
	case uint8:
		v = int(v.(uint8))
	case uint16:
		v = int(v.(uint16))
	case int:
		v = v.(int)
	case map[string]interface{}:
		// 遍历
		for key, value := range v.(map[string]interface{}) {
			v.(map[string]interface{})[key] = ProtoValueCheck(value)
		}
	case []interface{}:
		// 遍历
		for i, value := range v.([]interface{}) {
			v.([]interface{})[i] = ProtoValueCheck(value)
		}
	}
	// 返回
	return v
}
