package entities

type Message struct {
	Value      float64 `json:"value"`
	SensorType string  `json:"sensor_type"`
	ID1        string  `json:"id1"`
	ID2        int     `json:"id2"`
	Timestamp  string  `json:"timestamp"`
	Key        string  `json:"key"`
}
