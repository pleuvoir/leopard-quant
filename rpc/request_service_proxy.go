package rpc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gookit/color"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"leopard-quant/common/model"
	"reflect"
)

func RequestHandlerProxy(requestService RequestService, request model.RequestMessage) model.ResultMessage {
	r, err := requestHandlerProxy(requestService, request)
	if err == nil {
		return r
	}
	return model.NewFail(int32(codes.Unknown), request.Method, err.Error())
}

// 反射请求
func requestHandlerProxy(requestService RequestService, request model.RequestMessage) (r model.ResultMessage, err error) {

	defer func() {
		if r := recover(); r != nil {
			color.Redln(fmt.Sprintf("socketIO请求出现错误，%s", r))
			err = errors.New(fmt.Sprintf("%s", r))
		}
	}()

	tp := reflect.TypeOf(requestService)
	methodName := request.Method
	method, ok := tp.MethodByName(methodName)
	if !ok {
		return model.NewFail(int32(codes.NotFound), methodName, "Method not find"), nil
	}

	parameter := method.Type.In(2)
	req := reflect.New(parameter.Elem()).Interface()

	in := make([]reflect.Value, 0)
	ctx := context.Background()
	in = append(in, reflect.ValueOf(ctx))
	in = append(in, reflect.ValueOf(req))

	// 接口返回有错误
	call := reflect.ValueOf(&requestService).MethodByName(methodName).Call(in)
	if call[1].Interface() != nil {
		e := call[1].Interface().(error)
		st, _ := status.FromError(e)
		return model.NewFail(int32(st.Code()), methodName, st.Message()), nil
	}

	// 返回结果转json异常
	payload, err := json.Marshal(call[0].Interface())
	if err != nil {
		return model.NewFail(int32(codes.Aborted), methodName, "response data convert json error"), nil
	}

	// 业务正常返回
	return model.NewSuccess(methodName, string(payload)), nil
}
