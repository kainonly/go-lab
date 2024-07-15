package main

import (
	"fmt"
	"github.com/sony/sonyflake"
)

func main() {
	st := sonyflake.NewSonyflake(sonyflake.Settings{})
	id, err := st.NextID()
	if err != nil {
		panic(err)
	}
	fmt.Println(id)
}
