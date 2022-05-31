package errs

import "net/http"

// *******************
// type
// *******************

type Err struct {
	StatusCode int         `json:"-"`
	Code       int         `json:"code"`
	Msg        string      `json:"msg"`
	Detail     interface{} `json:"detail"`
}

// *******************
// function
// *******************

func (e Err) Ret(detail ...interface{}) (int, Err) {
	if len(detail) > 0 {
		e.Detail = detail[0]
	}
	return e.StatusCode, e
}

// *******************
// error
// *******************

var ERR1000 = Err{
	StatusCode: http.StatusInternalServerError,
	Code:       1000,
	Msg:        "internal server error",
}

var ERR2000 = Err{
	StatusCode: http.StatusBadRequest,
	Code:       2000,
	Msg:        "illegal request parameter",
}

var ERR2001 = Err{
	StatusCode: http.StatusBadRequest,
	Code:       2001,
	Msg:        "illegal request header",
}

var ERR2002 = Err{
	StatusCode: http.StatusNotFound,
	Code:       2002,
	Msg:        "resource not found",
}

var ERR2003 = Err{
	StatusCode: http.StatusUnauthorized,
	Code:       2003,
	Msg:        "unauthorized",
}

var ERR2004 = Err{
	StatusCode: http.StatusUnauthorized,
	Code:       2004,
	Msg:        "generate token error",
}

var ERR2005 = Err{
	StatusCode: http.StatusUnauthorized,
	Code:       2005,
	Msg:        "parse token error",
}

var ERR2006 = Err{
	StatusCode: http.StatusForbidden,
	Code:       2006,
	Msg:        "forbidden",
}
