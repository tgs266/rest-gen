package rumtime

type StandardError struct {
	ErrorId string                 `json:"errorId"`
	Code    string                 `json:"code"`
	Params  map[string]interface{} `json:"params"`
}
