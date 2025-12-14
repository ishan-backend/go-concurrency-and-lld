package main

// UserType represents the type of user
type UserType string

const (
	UserTypeRegular UserType = "Regular"
	UserTypePremium UserType = "Premium"
)

// ProductType represents the type of product
type ProductType string

const (
	ProductTypeGroceries  ProductType = "Groceries"
	ProductTypeBeverages ProductType = "Beverages"
	ProductTypeElectronics ProductType = "Electronics"
)

// PaymentMethod represents the payment method
type PaymentMethod string

const (
	PaymentMethodCreditCard PaymentMethod = "CreditCard"
	PaymentMethodDebitCard  PaymentMethod = "DebitCard"
	PaymentMethodUPI        PaymentMethod = "UPI"
	PaymentMethodCash       PaymentMethod = "Cash"
)

// User represents a user in the system
type User struct {
	UserID        string
	UserType      UserType
	PastPurchases []string
}

// CartItem represents an item in the cart
type CartItem struct {
	ProductID   string
	ProductType ProductType
	Quantity    int
	Price       float64
}

// Cart represents a shopping cart
type Cart struct {
	CartID string
	Items  []CartItem
}

// GetTotalAmount calculates the total amount of the cart
func (c *Cart) GetTotalAmount() float64 {
	total := 0.0
	for _, item := range c.Items {
		total += item.Price * float64(item.Quantity)
	}
	return total
}

// Order represents an order
type Order struct {
	OrderID     string
	User        *User
	Cart        *Cart
	TotalAmount float64
	PaymentMethod PaymentMethod
}

// Coupon represents a discount coupon
type Coupon struct {
	CouponID       string
	DiscountAmount float64
	Rules          []Rule
	RulesFormula   string // Logical formula, e.g., "AND(0, 1)" or "OR(0, 1)"
}

