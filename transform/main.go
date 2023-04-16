package main

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
)

type M = map[string]interface{}

func main() {
	data := M{
		"meta": M{
			"uid": "64256825b16d000cd6bcfb5d",
		},
		"values": M{
			"rids": []interface{}{
				"640e7c2c7d8a24d6f831e9bf",
				"640e7c2c7d8a24d6f831e9c0",
			},
		},
		"orders": []interface{}{
			M{"vid": "64256836b16d000cd6bcfb5f"},
			M{"vid": "642569d0b16d000cd6bcfb60"},
			M{"vid": "6425682cb16d000cd6bcfb5e"},
		},
		"snapshot": M{
			"products": []interface{}{
				M{"ids": []interface{}{
					"640e7c2c7d8a24d6f831e9bf",
					"640e7c2c7d8a24d6f831e9c0",
				}},
				M{"ids": []interface{}{
					"640e7c2c7d8a24d6f831e9bf",
					"640e7c2c7d8a24d6f831e9c0",
				}},
				M{"ids": []interface{}{
					"640e7c2c7d8a24d6f831e9bf",
					"640e7c2c7d8a24d6f831e9c0",
				}},
			},
		},
	}

	format := M{
		"meta.uid":                "oid",
		"values.rids":             "oids",
		"orders.$.vid":            "oid",
		"snapshot.products.$.ids": "oids",
	}

	for path, kind := range format {
		keys := strings.Split(path, ".")
		if err := transform(data, keys, kind); err != nil {
			panic(err)
		}
	}

	fmt.Println(data)
}

func transform(data M, keys []string, kind interface{}) (err error) {
	var cursor interface{}
	cursor = data
	n := len(keys) - 1
	for i, key := range keys[:n] {
		if key == "$" {
			for _, value := range cursor.([]interface{}) {
				if err = transform(value.(M), keys[i+1:], kind); err != nil {
					return
				}
			}
			return
		}
		cursor = cursor.(M)[key]
	}
	key := keys[n]
	switch kind {
	case "oid":
		if cursor.(M)[key], err = primitive.ObjectIDFromHex(cursor.(M)[key].(string)); err != nil {
			return
		}
		break
	case "oids":
		oids := cursor.(M)[key].([]interface{})
		for i, id := range oids {
			if oids[i], err = primitive.ObjectIDFromHex(id.(string)); err != nil {
				return
			}
		}
		break
	}
	return
}
