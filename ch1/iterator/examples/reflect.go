package main

import (
	"fmt"
	"reflect"
)

func main() {
	// NOTE: 传统反射
	{
		u := User{Name: "江湖十年", Age: 18}
		v := reflect.ValueOf(u)
		t := v.Type()

		// 传统反射：遍历结构体字段
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			val := v.Field(i)
			fmt.Printf("Field: %s = %v\n", field.Name, val.Interface())
		}

		// 传统反射：遍历函数入参类型
		funcType := reflect.TypeOf(ProcessUser)
		for i := 0; i < funcType.NumIn(); i++ {
			inType := funcType.In(i)
			fmt.Printf("Arg %d Type: %v\n", i, inType)
		}
	}

	// NOTE: Go 1.26 迭代器
	{
		u := User{Name: "江湖十年", Age: 18}
		v := reflect.ValueOf(u)

		// Go 1.26 迭代器：遍历结构体字段与值
		// v.Fields() 返回 iter.Seq2[reflect.StructField, reflect.Value]
		for field, val := range v.Fields() {
			fmt.Printf("Field: %s = %v\n", field.Name, val.Interface())
		}

		// Go 1.26 迭代器：遍历函数入参类型
		// funcType.Ins() 返回 iter.Seq[reflect.Type]
		funcType := reflect.TypeOf(ProcessUser)
		for inType := range funcType.Ins() {
			fmt.Printf("Type: %v\n", inType)
		}
	}
}

type User struct {
	Name string
	Age  int
}

func ProcessUser(name string, age int) User {
	return User{name, age}
}
