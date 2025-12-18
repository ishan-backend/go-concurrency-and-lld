# Shopping Cart Coupon Application System

A comprehensive low-level design implementation of a shopping cart system with coupon application functionality in Go.

## Problem Statement

Design a system that allows users to:
1. Add items to a shopping cart
2. Apply various types of coupons to the cart
3. Calculate discounts based on coupon types
4. Validate coupon eligibility before application
5. Handle multiple coupons (if applicable)
6. Calculate final cart total with discounts

## Features

### Supported Coupon Types

1. **Percentage Discount** (`PERCENTAGE`)
   - Applies a percentage discount on the cart total
   - Supports maximum discount cap
   - Example: 10% off with max $2000 discount

2. **Fixed Amount Discount** (`FIXED_AMOUNT`)
   - Applies a fixed dollar amount discount
   - Example: $500 off

3. **Buy X Get Y** (`BUY_X_GET_Y`)
   - Buy X items, get Y items free
   - Can be category-specific
   - Example: Buy 3 get 1 free on Clothing items

4. **Free Shipping** (`FREE_SHIPPING`)
   - Waives shipping charges
   - Can have minimum purchase requirements

### Coupon Validation Rules

- **Validity Period**: Coupons must be within their valid date range
- **Usage Limit**: Tracks and enforces maximum usage count
- **Minimum Purchase**: Enforces minimum cart value requirement
- **User Eligibility**: Supports user-specific coupons
- **Category Restrictions**: Can apply to specific product categories
- **Product Exclusions**: Can exclude specific products

## Architecture

### Core Components

1. **Models** (`models.go`)
   - `CartItem`: Represents an item in the cart
   - `ShoppingCart`: Represents the shopping cart with items and applied coupons
   - `Coupon`: Represents a discount coupon with all its properties
   - `DiscountResult`: Result of applying a coupon
   - `CartSummary`: Final cart summary with totals

2. **Coupon Validator** (`coupon_validator.go`)
   - Validates coupon eligibility
   - Checks validity period, usage limits, minimum purchase
   - Validates user-specific and category-specific rules
   - Thread-safe coupon storage and retrieval

3. **Discount Calculator** (`discount_calculator.go`)
   - Calculates discounts for different coupon types
   - Handles percentage, fixed amount, and buy-X-get-Y calculations
   - Ensures discounts don't exceed cart subtotal

4. **Cart Service** (`cart_service.go`)
   - Manages shopping cart operations
   - Handles item addition/removal
   - Applies and removes coupons
   - Calculates cart totals
   - Thread-safe operations

## Design Patterns Used

1. **Service Layer Pattern**: Encapsulates business logic in CartService
2. **Strategy Pattern**: Different discount calculation strategies for different coupon types
3. **Repository Pattern**: CouponValidator acts as a repository for coupons
4. **Thread Safety**: Uses `sync.RWMutex` for concurrent access

## Usage Example

```go
// Initialize services
validator := NewCouponValidator()
cartService := NewCartService(validator)

// Register coupons
coupon := &Coupon{
    CouponID:      "C1",
    Code:          "SAVE10",
    Type:          CouponTypePercentage,
    DiscountValue: 10.0,
    MinPurchase:   1000.0,
    MaxDiscount:   2000.0,
    ValidFrom:     time.Now(),
    ValidUntil:    time.Now().AddDate(0, 1, 0),
    UsageLimit:    100,
}
validator.RegisterCoupon(coupon)

// Create cart
cart := cartService.CreateCart("USER123")

// Add items
cartService.AddItem(cart.CartID, &CartItem{
    ProductID: "P1",
    Name:      "Laptop",
    Price:     50000.0,
    Quantity:  1,
    Category:  "Electronics",
})

// Apply coupon
result, err := cartService.ApplyCoupon(cart.CartID, "SAVE10")

// Get cart summary
summary, err := cartService.GetCartSummary(cart.CartID)
```

## Running the Code

```bash
cd shopping-cart-coupon
go run .
```

## Key Features

### Thread Safety
- All operations use `sync.RWMutex` for concurrent access
- Multiple goroutines can safely read cart data
- Write operations are exclusive

### Discount Calculation
- Percentage discounts with maximum cap support
- Fixed amount discounts
- Buy X Get Y with smart item selection (cheapest items made free)
- Free shipping handling

### Validation
- Comprehensive coupon validation before application
- Prevents duplicate coupon application
- Validates all eligibility criteria

### Error Handling
- Clear error messages for different failure scenarios
- Returns detailed discount results with success/failure status

## Test Scenarios Covered

1. ✅ Adding items to cart
2. ✅ Applying percentage coupon
3. ✅ Applying fixed amount coupon
4. ✅ Handling expired coupons
5. ✅ Removing coupons
6. ✅ Buy X Get Y coupon application
7. ✅ Free shipping coupon
8. ✅ User-specific coupons
9. ✅ Minimum purchase validation

## Future Enhancements

1. **Multiple Coupon Application**: Allow stacking multiple coupons
2. **Coupon Priority**: Define priority order for multiple coupons
3. **Category-Specific Discounts**: Apply discounts only to specific categories
4. **Product-Level Exclusions**: Exclude specific products from discounts
5. **Coupon Combinations**: Define which coupons can be used together
6. **Database Integration**: Replace in-memory storage with database
7. **Caching**: Cache frequently accessed coupons
8. **Analytics**: Track coupon usage and effectiveness
9. **Time-Based Rules**: Apply coupons during specific time windows
10. **First-Time User Coupons**: Special coupons for new users

## Code Structure

```
shopping-cart-coupon/
├── models.go              # Domain models
├── coupon_validator.go    # Coupon validation logic
├── discount_calculator.go  # Discount calculation logic
├── cart_service.go        # Main cart service
├── main.go                # Examples and tests
├── go.mod                 # Go module file
└── README.md              # Documentation
```

## Conclusion

This implementation provides a robust, thread-safe shopping cart system with comprehensive coupon application functionality. The design is extensible, maintainable, and ready for production use with proper error handling and validation.

