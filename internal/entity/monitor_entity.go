package entity

type MonitorMetric struct {
	Count     uint64 `json:"count" bson:"count"`
	Succeeded uint64 `json:"succeeded" bson:"succeeded"`
	Failed    uint64 `json:"failed" bson:"failed"`
}

type MonitorPipeline struct {
	All MonitorMetric `json:"all" bson:"all"`
	TTS MonitorMetric `json:"tts" bson:"tts"`
	TTT MonitorMetric `json:"ttt" bson:"ttt"`
	STT MonitorMetric `json:"stt" bson:"stt"`
	LS  MonitorMetric `json:"ls" bson:"ls"`
	FP  MonitorMetric `json:"fp" bson:"fp"`
}

type MonitorDataType struct {
	Videos MonitorMetric `json:"videos" bson:"videos"`
	Audios MonitorMetric `json:"audios" bson:"audios"`
	Texts  MonitorMetric `json:"texts" bson:"texts"`
}
