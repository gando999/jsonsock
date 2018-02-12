package jsonsock

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type Dispatcher interface {
	RegisterImpl(namespace string, targetImpl interface{})
	CallFunc(targetFunc string, funcArgs []interface{}) DispatcherResult
}

type Registry interface {
	RegisterImpl(identifier string, impl interface{})
	FindTargetImpl(identifier string) interface{}
}

type DispatcherResult interface {
	GetResult() (interface{}, error)
}

type CallResult struct {
	actual []reflect.Value
	err    error
}

func (callResult CallResult) GetResult() (interface{}, error) {
	return callResult.actual[0].Interface(), callResult.err
}

type DynamicDispatcher struct {
	registry Registry
}

type NamespaceBasedRegistry struct {
	namespaceMap map[string]interface{}
}

func CreateRegistry() Registry {
	registry := new(NamespaceBasedRegistry)
	registry.namespaceMap = make(map[string]interface{})
	return registry
}

func (registry NamespaceBasedRegistry) FindTargetImpl(identifier string) interface{} {
	return registry.namespaceMap[identifier]
}

func (registry NamespaceBasedRegistry) RegisterImpl(identifier string, impl interface{}) {
	registry.namespaceMap[identifier] = impl
}

func CreateDispatcher() Dispatcher {
	dispatcher := new(DynamicDispatcher)
	dispatcher.registry = CreateRegistry()
	return dispatcher
}

func (dispatcher DynamicDispatcher) RegisterImpl(namespace string, targetImpl interface{}) {
	dispatcher.registry.RegisterImpl(namespace, targetImpl)
}

func (dispatcher DynamicDispatcher) CallFunc(targetFunc string, funcArgs []interface{}) DispatcherResult {
	s := strings.Split(targetFunc, ".")
	namespace, funcName := s[0], s[1]
	targetImpl := dispatcher.registry.FindTargetImpl(namespace)
	rValues, err := CallFuncByName(targetImpl, funcName, funcArgs...)
	return CallResult{rValues, err}
}

func setField(obj interface{}, name string, value interface{}) error {
	structValue := reflect.ValueOf(obj).Elem()
	structFieldValue := structValue.FieldByName(name)
	if !structFieldValue.IsValid() {
		return fmt.Errorf("No such field: %s in obj", name)
	}
	if !structFieldValue.CanSet() {
		return fmt.Errorf("Cannot set %s field value", name)
	}
	structFieldType := structFieldValue.Type()
	val := reflect.ValueOf(value)
	if structFieldType != val.Type() {
		switch val.Interface().(type) {
		case map[string]interface{}:
			tStruct := reflect.New(structFieldType).Interface()
			err := FillStruct(tStruct, val.Interface().(map[string]interface{}))
			if err != nil {
				return err
			}
			structFieldValue.Set(reflect.Indirect(reflect.ValueOf(tStruct)))
			return nil
		default:
			invalidTypeError := errors.New("Provided value type didn't match obj field type")
			return invalidTypeError
		}
	} else {
		structFieldValue.Set(val)
		return nil
	}
}

func FillStruct(s interface{}, m map[string]interface{}) error {
	for k, v := range m {
		err := setField(s, k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func CreateParameter(param interface{}, argType reflect.Type) reflect.Value {
	switch v := reflect.ValueOf(param).Interface().(type) {
	case float64, string, int:
		if argType.Kind() == reflect.Ptr {
			p := reflect.New(reflect.TypeOf(v))
			p.Elem().Set(reflect.ValueOf(param))
			return p
		}
		return reflect.ValueOf(param)
	case map[string]interface{}:
		tStruct := reflect.New(argType).Interface()
		paramMap := reflect.ValueOf(param).Interface()
		err := FillStruct(tStruct, paramMap.(map[string]interface{})) //check error
		if err != nil {
			return reflect.Zero(argType)
		}
		return reflect.Indirect(reflect.ValueOf(tStruct))
	default:
		fmt.Printf("Unknown type %T!\n", v)
		return reflect.Zero(argType)
	}
}

func CallFuncByName(targetImpl interface{}, funcName string, params ...interface{}) (out []reflect.Value, err error) {
	m := reflect.ValueOf(targetImpl).MethodByName(funcName)
	if !m.IsValid() {
		return make([]reflect.Value, 0), fmt.Errorf("Method not found \"%s\"", funcName)
	}
	in := make([]reflect.Value, len(params))
	for i, param := range params {
		switch reflect.ValueOf(param).Interface().(type) {
		case []interface{}:
			argType := m.Type().In(i)
			tList := reflect.Indirect(reflect.New(argType))
			paramList := reflect.ValueOf(param).Interface().([]interface{})
			for _, sliceElement := range paramList {
				tList.Set(reflect.Append(tList, CreateParameter(sliceElement, argType.Elem())))
			}
			in[i] = tList
		default:
			argType := m.Type().In(i)
			in[i] = CreateParameter(param, argType)
		}
	}
	out = m.Call(in)
	return
}
