package server

//接口调用结果
type Result int

const DefaultVersion = 1

const (
	SUCCESS_INT    Result = 0 //成功, 加一个int以便与交易状态的区别
	FAILED_INT     Result = 1 //失败，加一个int以便与交易状态的区别
	NOTFUNDS       Result = 2 //余额不足
	GETDATA_FAIL   Result = 3 //从链上取到信息失败
	DOTS_EXCEPTION Result = 4 //组件异常，如果没有正常初始等
	GENERATE_FAIL  Result = 5 //生成数据出错
	PARAMETER_IN   Result = 6 //输入参数不正确

	TX_NOT_FINISH Result = 10 //交易没有完成，不能创建新的交易

	HTTP_ERROR     Result = 100 //http err when client call server, "{}"是http错误的内容
	DATABASE_ERROR Result = 200 //数据库出错
)

//交易状态
type Status string

const (
	None     Status = "" //没有初始化的状态
	Generate Status = "Generate"
	Cancel   Status = "Cancel"
	Sent     Status = "Sent"
	Success  Status = "Success"
	Failed   Status = "Failed"
)

type ResHeader struct {
	Version int32  `json:"version"` // 现在版本默认为 1
	Result  Result `json:"result"`  //
	Message string `json:"message"` //出错或提示的信息，
	Id      string `json:"id"`      //请求的id， 如果有订单id时，它就是订单id； 这个值就是请求的Id
	ExData  string `json:"exData"`  //扩展数据，服务端不做处理，只是原样返回
}

func NewResHeader() *ResHeader {
	return &ResHeader{
		Version: DefaultVersion,
		Result:  SUCCESS_INT,
		Id:      "",
	}
}

type ReqHeader struct {
	Version int32  `json:"version"`    // 现在版本默认为 1
	Id      string `json:"purchaseId"` //请求的id， 如果有订单id时，它就是订单id
	ExData  string `json:"exData"`     //扩展数据，服务端不做处理，只是原样返回
}

func NewReqHeader() *ReqHeader {
	return &ReqHeader{
		Version: DefaultVersion,
		Id:      "",
	}
}
