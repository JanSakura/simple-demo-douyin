package models

// ResponseStatus 状态码，每个接口对应的都不太一样，这个是基础的；用gin.H{}写请求状态太多，抽象出来
type ResponseStatus struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}
