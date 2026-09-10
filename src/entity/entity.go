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

type PaymentStatusType string

const (
	Paid   PaymentStatusType = "paid"
	unpaid PaymentStatusType = "unpaid"
)

type OrderStatusType string

const (
	Pending  OrderStatusType = "pending"
	Verified OrderStatusType = "verified"
	Rejected OrderStatusType = "rejected"
)

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
type CustomerOrder struct {
	ID             int
	OrderCode      string
	DeterminedSize string
	TotalPrice     float64
	PaymentStatus  PaymentStatusType
	Status         OrderStatusType
	CreatedAt      time.Time
}
