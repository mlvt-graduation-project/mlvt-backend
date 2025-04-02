package entity

type MonitorMetric struct {
	Count     uint64 `json:"count" bson:"count"`
	Succeeded uint64 `json:"succeeded" bson:"succeeded"`
	Failed    uint64 `json:"failed" bson:"failed"`
}

// #region Pipeline
type MonitorPipeline struct {
	All MonitorMetric `json:"all" bson:"all"`
	TTS MonitorMetric `json:"tts" bson:"tts"`
	TTT MonitorMetric `json:"ttt" bson:"ttt"`
	STT MonitorMetric `json:"stt" bson:"stt"`
	LS  MonitorMetric `json:"ls" bson:"ls"`
	FP  MonitorMetric `json:"fp" bson:"fp"`
}

// #endregion

// #region Media

type MonitorDataType struct {
	Videos MonitorMetric `json:"videos" bson:"videos"`
	Audios MonitorMetric `json:"audios" bson:"audios"`
	Texts  MonitorMetric `json:"texts" bson:"texts"`
}

// #endregion

// region Time Period

type TimePeriodType string

const (
	TimePeriodDay  TimePeriodType = "day"
	TimePeriodWeek TimePeriodType = "week"
	TimePeriodYear TimePeriodType = "year"
)

// #endregion

// region Traffic

type MonitorTraffic struct {
	Cell  string `json:"cell" bson:"cell"`
	Value uint64 `json:"value" bson:"value"`
}

type MonitorTraffics struct {
	Count   uint64           `json:"count" bson:"count"`
	Traffic []MonitorTraffic `json:"traffic" bson:"traffic"`
}

// #endregion
