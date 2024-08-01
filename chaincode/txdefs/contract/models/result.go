package models

type Result struct {
	Success  bool
	Feedback string
	Data     map[string]interface{}
	Assets   []map[string]interface{}
}
