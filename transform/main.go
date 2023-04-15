package main

import (
	"fmt"
)

type M = map[string]interface{}

func main() {
	data := M{
		"uid": "64256825b16d000cd6bcfb5d",
		//"values": M{
		//	"rids": []interface{}{
		//		"640e7c2c7d8a24d6f831e9bf",
		//		"640e7c2c7d8a24d6f831e9c0",
		//	},
		//},
		//"orders": []interface{}{
		//	M{"vid": "64256836b16d000cd6bcfb5f"},
		//	M{"vid": "642569d0b16d000cd6bcfb60"},
		//	M{"vid": "6425682cb16d000cd6bcfb5e"},
		//},
	}

	//format := M{
	//	"uid": "oid",
	//	//"values.rids.$": "oid",
	//	//"orders.$.vid": "oid",
	//}

	//for path, kind := range format {
	//	keys := strings.Split(path, ".")
	//	if err := transform(data, keys, kind); err != nil {
	//		panic(err)
	//	}
	//}

	fmt.Println(data)

}

//func transform(in interface{}, keys []string, kind interface{}) (err error) {
//	n := len(keys)
//	if n == 1 {
//		return
//	}
//	key := keys[0]
//	keys = keys[1:]
//	switch data := in.(type) {
//	case M:
//		if err = transform(data[key], keys, kind); err != nil {
//			return
//		}
//		break
//	case string:
//		switch kind {
//		case "oid":
//			var v primitive.ObjectID
//			if v, err = primitive.ObjectIDFromHex(data); err != nil {
//				return
//			}
//			fmt.Println(v)
//			break
//
//		}
//	}
//
//	return
//}

//func set(in interface{}, kind interface{}) (err error) {
//	switch kind {
//	case "oid":
//		var v primitive.ObjectID
//		if v, err = primitive.ObjectIDFromHex(in.(string)); err != nil {
//			return
//		}
//		in = &v
//		break
//
//	}
//	return
//}
