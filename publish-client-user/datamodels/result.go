package datamodels
type Result struct {
	Status bool
	Msg    string
	Code   string
	Data   interface{}
	Token  string
}
type ResultPage struct {
	Status    bool
	Msg       string
	Code      string
	PageCount int
	Data      interface{}
}

/*
d为通用返回数据，e为可选参数，e不为空时返回失败result
 */
func NewResult(d interface{}, e ... string) *Result {
	if len(e) > 0 {
		return &Result{false, e[0], "200", nil, ""}
	}
	return &Result{true, "success", "200", d, ""}
}
func GetResult(result *Result,d interface{}, e ... string)*Result{
	if len(e)>0 {
		result.Code="200"
		result.Status=false
		result.Msg=e[0]
		result.Data=nil
	}else{
		result.Code="200"
		result.Status=true
		result.Msg="success"
		result.Data=d
	}
	return nil
}
/*
d为通用返回数据，e为可选参数，e不为空时返回失败result
 */
func NewResultPage(d interface{}, pageCount int, e ... string) *ResultPage {
	if len(e) > 0 {
		return &ResultPage{false, e[0], "200", 0, ""}
	}
	return &ResultPage{true, "success", "200", pageCount, d}
}
func GetResultPage(result *ResultPage,d interface{},pageCount int, e ... string)*ResultPage{
	if len(e)>0 {
		result.Code="200"
		result.Status=false
		result.Msg=e[0]
		result.Data=nil
	}else{
		result.Code="200"
		result.Status=true
		result.Msg="success"
		result.Data=d
		result.PageCount = pageCount
	}
	return nil
}