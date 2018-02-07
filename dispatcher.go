package jsonsock

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type JsonRequest struct {
	Id      int           `json:"id"`
	Jsonrpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

type JsonResponse struct {
	Id      int         `json:"id"`
	Jsonrpc string      `json:"jsonrpc"`
	Result  interface{} `json:"result"`
}

type Dispatcher struct {
	registry map[string]interface{}
}

func MarshalRequest(targetFunc string, funcArgs []interface{}) (string, error) {
	req := &JsonRequest{123, "2.0", targetFunc, funcArgs}
	resp, err := json.Marshal(req)
	if err != nil {
		return "", errors.New("Failed to convert!")
	}
	return string(resp), nil
}

func UnmarshalRequest(incomingRequest string) (JsonRequest, error) {
	var request JsonRequest
	err := json.Unmarshal([]byte(incomingRequest), &request)
	if err != nil {
		return request, errors.New("Failed to unmarshal request")
	}
	return request, nil
}

func MarshalResponse(response interface{}) (string, error) {
	jresp := &JsonResponse{123, "2.0", response}
	resp, err := json.Marshal(jresp)
	if err != nil {
		return "", errors.New("Failed to create response")
	}
	return string(resp), nil
}

func (dispatcher Dispatcher) RegisterImpl(namespace string, targetImpl interface{}) {
	dispatcher.registry[namespace] = targetImpl
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
		invalidTypeError := errors.New("Provided value type didn't match obj field type")
		return invalidTypeError
	}

	structFieldValue.Set(val)
	return nil
}

func fillStruct(s interface{}, m map[string]interface{}) error {
	for k, v := range m {
		err := setField(s, k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func CallFuncByName(targetImpl interface{}, funcName string, params ...interface{}) (out []reflect.Value, err error) {
	m := reflect.ValueOf(targetImpl).MethodByName(funcName)
	if !m.IsValid() {
		return make([]reflect.Value, 0), fmt.Errorf("Method not found \"%s\"", funcName)
	}
	in := make([]reflect.Value, len(params))
	for i, param := range params {
		switch reflect.ValueOf(param).Interface().(type) {
		case float64:
			in[i] = reflect.ValueOf(param)
		default:
			argType := m.Type().In(i)
			tStruct := reflect.New(argType).Interface()
			paramMap := reflect.ValueOf(param).Interface()
			err := fillStruct(tStruct, paramMap.(map[string]interface{}))
			if err == nil {
				in[i] = reflect.Indirect(reflect.ValueOf(tStruct))
			}
		}
	}
	out = m.Call(in)
	return
}

func (dispatcher Dispatcher) CallFunc(targetFunc string, funcArgs []interface{}) ([]reflect.Value, error) {
	s := strings.Split(targetFunc, ".")
	namespace, funcName := s[0], s[1]
	targetImpl := dispatcher.registry[namespace]
	return CallFuncByName(targetImpl, funcName, funcArgs...)
}
