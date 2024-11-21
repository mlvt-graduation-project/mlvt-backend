package entity

type StatusEntity string

const (
	StatusRaw        StatusEntity = "raw"
	StatusProcessing StatusEntity = "processing"
	StatusSucceeded  StatusEntity = "succeeded"
	StatusFailed     StatusEntity = "failed"
)
