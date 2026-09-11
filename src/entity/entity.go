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
	Unpaid PaymentStatusType = "unpaid"
)

type OrderStatusType string

const (
	Pending        OrderStatusType = "pending"
	InProgress     OrderStatusType = "in progress"
	WaitingPayment OrderStatusType = "waiting payment"
	Finished       OrderStatusType = "finished"
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
	Progress       string
	CreatedAt      time.Time
}

type AdminOrder struct {
	ID               int
	OrderCode        string
	CustomerName     string
	DeterminedSize   string
	CMUsed           int
	AssignedWorkerID *int
	Status           string
	Progress         string
	CreatedAt        time.Time
}

type Payment struct {
	ID           int
	OrderID      int
	OrderCode    string
	CustomerName string
	Amount       float64
	Status       string
	CreatedAt    time.Time
}

type AvailableWorker struct {
	ID   int
	Name string
}
type SalesReport struct {
	TotalOrders    int
	TotalRevenue   float64
	PaidOrders     int
	UnpaidOrders   int
	FinishedOrders int
}

type FabricPattern struct {
	ID          int
	FabricID    int
	FabricName  string
	PatternID   int
	PatternName string
	StockCM     int
	PricePerCM  float64
}

type Pattern struct {
	ID   int
	Name string
}
