package model

type Users struct {
	User_name    string
	Email        string
	Credit_limit int
	Spent        int
}

type Merchants struct {
	Merchant_name  string
	Email          string
	Discount       float64
	Total_discount float64
}

type User_report struct {
	User_name string
	Spent     int
}

type User_spent struct {
	Balance int
	Spent   int
}

type Merchant_discount struct {
	Discount       float64
	Total_discount float64
}
