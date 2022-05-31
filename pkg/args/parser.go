package args

import "github.com/gin-gonic/gin"

type ReqData struct {
	Params map[string]string
	Query  map[string][]string
	Body   map[string]interface{}
}

func ParseGinCtx(ctx *gin.Context) ReqData {
	params := make(map[string]string)
	for _, p := range ctx.Params {
		params[p.Key] = p.Value
	}
	var body map[string]interface{}
	ctx.ShouldBindJSON(&body)
	return ReqData{
		Params: params,
		Query:  ctx.Request.URL.Query(),
		Body:   body,
	}
}
