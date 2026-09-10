package entity

import "time"

type User struct {
	ID        int
	Name      string
	Email     string
	Password  string
	Role      string
	Phone     string
	CreatedAt time.Time
}

type UserMeasurement struct {
	ID                 int
	UserID             int
	Title              string
	HeightCM           float64
	ChestCircumference float64
	WaistCircumference float64
	CreatedAt          time.Time
}

type Fabric struct {
	ID            int
	Name          string
	PricePerMeter float64
}

type FabricPatternOption struct {
	FabricPatternID int
	FabricID        int
	PatternID       int
	PatternName     string
	StockCM         int
	PricePerMeter   float64
}

type SizeRequirement struct {
	ID         int
	Size       string
	RequiredCM int
}
