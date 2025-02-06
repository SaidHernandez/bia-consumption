package model

type AggregatedConsumption struct {
	Period             []string  `json:"period"`
	ActiveEnergy       []float64 `json:"active"`
	ReactiveInductive  []float64 `json:"reactive_inductive"`
	ReactiveCapacitive []float64 `json:"reactive_capacitive"`
	ExportedEnergy     []float64 `json:"exported"`
}
