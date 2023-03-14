package reflect

import (
	"reflect"
	"testing"
)

type Person struct {
	Age int
}

func (x *Person) SetAge(v int) {
	x.Age = v
}

func (x *Person) GetAge() int {
	return x.Age
}

func TestCall(t *testing.T) {
	p := new(Person)
	r := reflect.ValueOf(p)
	r.MethodByName("SetAge").
		Call([]reflect.Value{reflect.ValueOf(30)})
	t.Log(p.Age)
	rr := r.MethodByName("GetAge").Call([]reflect.Value{})
	for _, v := range rr {
		t.Log(v.Interface())
	}
}
