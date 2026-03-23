package helpers

import "net/http"

type Result[T any] interface {
	isResult()
}

type Success[T any] struct {
	Data T
}

func (s Success[T]) isResult() {}

type Error struct {
	Err  error
	Code int
}

func (e Error) Error() string {
	return e.Err.Error()
}

func (e Error) isResult() {}

func Try[T any](val T, err error, codes ...int) Result[T] {
	return ToResult(val, err, codes...)
}

func TryVoid(err error, codes ...int) Result[struct{}] {
	return ToResultVoid(err, codes...)
}

func Bind[T any, U any](r Result[T], next func(T) Result[U]) Result[U] {
	switch v := r.(type) {
	case Success[T]:
		return next(v.Data)
	case Error:
		return v
	default:
		panic("invalid result type")
	}
}

func Then[U any](r Result[struct{}], next func() Result[U]) Result[U] {
	switch v := r.(type) {
	case Success[struct{}]:
		_ = v
		return next()
	case Error:
		return v
	default:
		panic("invalid result type")
	}
}

func BindVoid[T any](r Result[T], next func(T) error) Result[struct{}] {
	return Bind(r, func(v T) Result[struct{}] {
		return ToResultVoid(next(v))
	})
}

func ToResult[T any](val T, err error, codes ...int) Result[T] {
	if err != nil {
		code := 500
		if len(codes) > 0 {
			code = codes[0]
		}
		return Error{Err: err, Code: code}
	}
	return Success[T]{Data: val}
}

func ToResultVoid(err error, codes ...int) Result[struct{}] {
	if err != nil {
		code := 500
		if len(codes) > 0 {
			code = codes[0]
		}
		return Error{Err: err, Code: code}
	}
	return Success[struct{}]{Data: struct{}{}}
}

func Match[T any](r Result[T], onSuccess func(T), onError func(Error)) {
	switch v := r.(type) {
	case Success[T]:
		onSuccess(v.Data)
	case Error:
		onError(v)
	default:
		panic("invalid result type")
	}
}

func ErrorCodeToHttp(code int) int {
	switch code {
	case 400:
		return http.StatusBadRequest
	case 401:
		return http.StatusUnauthorized
	case 403:
		return http.StatusForbidden
	case 404:
		return http.StatusNotFound
	case 500:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}
