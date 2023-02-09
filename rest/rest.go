package rest

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	// ErrvalueNotStructOrPointer value is not struct or pointer / value 不是结构体或指针
	ErrValueNotStructOrPointer = errors.New("value is not struct or pointer")
	ErrInvalidType             = errors.New("invalid type")
)

func parserMethods(gt any) ([]Methodor, error) {
	mrs := []Methodor{}

	if e := rangeStructMethod(gt, func(v reflect.Value, m reflect.Method) error {
		if f, ok := v.Interface().(func() Groupor); ok {
			tmrs, e := f().export("/", nil)
			if e != nil {
				return e
			}

			mrs = append(mrs, tmrs...)
		}
		return nil
	}); e != nil {
		return nil, e
	}

	return mrs, nil
}

func Test(v any) {
	mrs, e := parserMethods(v)
	if e != nil {
		if errors.Is(e, ErrValueNotStructOrPointer) {
			panic("参数必须是结构体或结构体指针")
		}

		panic(fmt.Sprintf("parserGroups: %s", e))
	}

	for i := 0; i < len(mrs); i++ {
		m := mrs[i]

		fmt.Printf("method => %s\n", m.exportString())
	}
}
