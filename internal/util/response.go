package util

// Response 统一API返回结构体，网关全场景通用
type Response struct {
	Code    int         `json:"code"`              // 业务码:0成功,>0业务错误,5xx系统异常
	Message string      `json:"message"`           // 对外提示文案
	Data    interface{} `json:"data,omitempty"`    // 业务数据，为空时不输出该字段
	Details []string    `json:"details,omitempty"` // 详细错误信息，参数校验、多错误场景使用
	TraceId string      `json:"traceId,omitempty"` // 链路追踪ID，排查问题
}

// Page 分页公共data结构，查询分页接口直接套在Data里面
type Page struct {
	List     interface{} `json:"list"`      // 列表数据
	Total    int64       `json:"total"`     // 总条数
	PageNum  int         `json:"page_num"`  // 当前页码
	PageSize int         `json:"page_size"` // 每页大小
}

// ---------------- 成功返回 ----------------

// OK 成功，不带traceId
func OK(data interface{}) Response {
	return Response{
		Code:    0,
		Message: "ok",
		Data:    data,
	}
}

// OKWithTrace 成功，携带traceId（网关推荐使用）
func OKWithTrace(data interface{}, traceId string) Response {
	return Response{
		Code:    0,
		Message: "ok",
		Data:    data,
		TraceId: traceId,
	}
}

// OKPage 分页成功快捷方法
func OKPage(list interface{}, total int64, pageNum, pageSize int, traceId string) Response {
	return Response{
		Code:    0,
		Message: "ok",
		Data: Page{
			List:     list,
			Total:    total,
			PageNum:  pageNum,
			PageSize: pageSize,
		},
		TraceId: traceId,
	}
}

// ---------------- 业务错误返回 ----------------

// Error 自定义业务错误码+消息
func Error(code int, message string) Response {
	return Response{
		Code:    code,
		Message: message,
	}
}

// ErrorWithTrace 业务错误，带traceId
func ErrorWithTrace(code int, message string, traceId string) Response {
	return Response{
		Code:    code,
		Message: message,
		TraceId: traceId,
	}
}

// ErrorWithDetails 参数校验类错误，带多条详情
func ErrorWithDetails(code int, message string, details []string, traceId string) Response {
	return Response{
		Code:    code,
		Message: message,
		Details: details,
		TraceId: traceId,
	}
}

// ---------------- 常用预置快捷错误（可直接调用，不用每次传code） ----------------

// ErrBadRequest 参数错误 400
func ErrBadRequest(msg string, traceId string) Response {
	return ErrorWithTrace(400, msg, traceId)
}

// ErrUnauthorized 未登录/token失效 401
func ErrUnauthorized(msg string, traceId string) Response {
	return ErrorWithTrace(401, msg, traceId)
}

// ErrForbidden 权限不足 403
func ErrForbidden(msg string, traceId string) Response {
	return ErrorWithTrace(403, msg, traceId)
}

// ErrNotFound 资源不存在 404
func ErrNotFound(msg string, traceId string) Response {
	return ErrorWithTrace(404, msg, traceId)
}

// ErrServerInternal 服务内部异常500，panic捕获统一返回
func ErrServerInternal(traceId string) Response {
	return ErrorWithTrace(500, "服务器内部错误，请稍后重试", traceId)
}
