package main

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/grafana/sobek"
)

var client = resty.New().SetHeader("User-Agent", userAgent)

func NewVM() (*sobek.Runtime, error) {
	vm := sobek.New()

	vm.Set("exit", func(msg string) {
		panic(vm.ToValue(msg))
	})

	console := vm.NewObject()
	if err := console.Set("log", func(v any) {
		fmt.Println(v)
	}); err != nil {
		return nil, fmt.Errorf("new console log: %w", err)
	}
	if err := vm.Set("console", console); err != nil {
		return nil, fmt.Errorf("set console: %w", err)
	}

	if err := vm.Set("get", func(url string, call sobek.FunctionCall) (sobek.Value, error) {
		var headers map[string]string
		if len(call.Arguments) > 0 {
			obj := call.Arguments[0].ToObject(vm)
			headers = convertMapString(obj.Get("headers").ToObject(vm))
		}

		r := client.R()
		r.SetHeaders(headers)

		resp, err := r.Get(url)
		if err != nil {
			return nil, fmt.Errorf("do request: %w", err)
		}

		return vm.ToValue(map[string]any{
			"status": resp.StatusCode(),
			"header": resp.Header(),
			"body":   string(resp.Body()),
		}), nil
	}); err != nil {
		return nil, fmt.Errorf("set http get: %w", err)
	}

	return vm, nil
}

func convertMapString(obj *sobek.Object) map[string]string {
	result := make(map[string]string)
	if obj == nil {
		return result
	}
	for _, key := range obj.Keys() {
		result[key] = obj.Get(key).String()
	}

	return result
}
