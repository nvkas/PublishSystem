package tool

import "Publish_System/datamodels"

func NewResult(data interface{}, c string, m ...string) *datamodels.Result {
	r := &datamodels.Result{Data: data, Code: c}

	if e, ok := data.(error); ok {
		if m == nil {
			r.Msg = e.Error()
		}
	} else {
		r.Status = true
		r.Msg = "SUCCESS"
	}
	if len(m) > 0 {
		r.Msg = m[0]
	}

	return r
}
