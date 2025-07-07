package entity

type StatusEntity string
const (
	StatusRaw        StatusEntity = "raw"
	StatusProcessing StatusEntity = "processing"
	StatusSucceeded  StatusEntity = "succeeded"
	StatusFailed     StatusEntity = "failed"
)

type MediaType string
const (
	MediaTypeVideo MediaType = "video"
	MediaTypeAudio MediaType = "audio"
	MediaTypeText  MediaType = "text"
)

type ProgressType string
const (
	ProgressTypeTTT ProgressType = "ttt"
	ProgressTypeSTT ProgressType = "stt"
	ProgressTypeTTS ProgressType = "tts"
	ProgressTypeLS  ProgressType = "ls"
	ProgressTypeFP  ProgressType = "fp" // Full pipeline
)