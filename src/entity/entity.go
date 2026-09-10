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
	ID         int
	Name       string
	PricePerCM float64
}

type FabricPatternOption struct {
	FabricPatternID int
	FabricID        int
	PatternID       int
	PatternName     string
	StockCM         int
	PricePerCM      float64
}

type OrderSummary struct {
	FabricName  string
	PatternName string
	Size        string
	RequiredCM  int
	PricePerCM  float64
	TotalPrice  float64
}
