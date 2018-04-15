package msg

import (
	_ "encoding/json"
)

type TFServer struct {
	Position []float64 `json:"position"`
	Rotation []float64 `json:"rotation"`
}
