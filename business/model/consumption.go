package model

import "time"

type Consumption struct {
	ID                 string    `gorm:"primaryKey"`
	MeterID            int       `gorm:"index"`
	Date               time.Time `gorm:"index"`
	ActiveEnergy       float64
	ReactiveInductive  float64
	ReactiveCapacitive float64
	ExportedEnergy     float64
}
