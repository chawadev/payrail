# Payrail - Go Payment Framework

Lightweight payment framework for integrating mobile money providers (Lenco) into your Go applications.

---

## Installation

```bash
go get github.com/chawadev/payrail
```

---

## Setup

```go
import (
    "github.com/chawadev/payrail"
    "github.com/chawadev/payrail/payrailCore"
)

// Create a client
client, err := payrail.NewClient("your-lenco-api-key", "lenco")
if err != nil {
    log.Fatal(err)
}
```

---

## 1. Request Payment (Charge)

### What You Send

```go
req := core.ChargeRequest{
    Amount:     100.00,                    // Payment amount
    Reference:  "order-INV-2024-001",      // Your unique ID (6+ chars)
    Phone:      "+260769312808",          // Customer phone
    Operator:   core.OperatorMTN,         // OperatorMTN or OperatorAirtel
    Country:    core.CountryZM,           // CountryZM or CountryMW
    Bearer:     core.BearerCustomer,      // Who pays fee: BearerCustomer or BearerMerchant
    CustomerID: "customer-123",           // Your internal customer ID
}

resp, err := client.Charge(req)
if err != nil {
    log.Fatal(err)
}
```

### What You Get Back

**Success Response:**
```go
ChargeResponse{
    Status:        "pending",                      // pending, successful, failed, pay-offline
    TransactionID: "49b9628d-b5cd-41c8-bd1e-98b6591360f5",  // Lenco's ID
    Reference:     "2606007027",                  // Use this for status checks
    RawResponse:   // Raw JSON bytes for debugging
}
```

**Print the response:**
```go
fmt.Println("Status:", resp.Status)
fmt.Println("Transaction ID:", resp.TransactionID)
fmt.Println("Lenco Reference:", resp.Reference)
```

### Field Guide

| Field | Options | Example |
|-------|---------|---------|
| `Amount` | Any positive number | `100.00`, `0.50` |
| `Reference` | 6+ unique alphanumeric | `order-12345`, `INV-2024-001` |
| `Phone` | International format | `+260769312808`, `+265989123456` |
| `Operator` | `OperatorMTN` `OperatorAirtel` | `core.OperatorMTN` |
| `Country` | `CountryZM` `CountryMW` | `core.CountryZM` |
| `Bearer` | `BearerCustomer` `BearerMerchant` | `core.BearerCustomer` |
| `CustomerID` | Your system's customer ID | `cus_abc123` |

---

## 2. Check Payment Status (Veryfi)

### What You Send

```go
// Use the Lenco Reference from the charge response
statusResp, err := client.Veryfi("2606007027")
if err != nil {
    log.Fatal(err)
}
```

### What You Get Back

```go
StatusResponse{
    Status:           "successful",     // pending, successful, failed, pay-offline, 3ds-auth-required
    TransactionID:    "49b9628d-b5cd-41c8-bd1e-98b6591360f5",
    Reference:        "2606007027",
    Amount:           "100.00",
    Currency:         "ZMW",
    Fee:              "2.50",
    SettlementStatus: "pending",        // pending, settled, null
    ReasonForFailure: "",               // Only if status is failed
    RawResponse:      // Raw JSON for debugging
}
```

**Check status:**
```go
switch statusResp.Status {
case "pending":
    fmt.Println("Waiting for customer...")
case "successful":
    fmt.Println("Payment received! Amount:", statusResp.Amount, statusResp.Currency)
case "failed":
    fmt.Println("Payment failed:", statusResp.ReasonForFailure)
}
```

---

## Complete Example

```go
package main

import (
	"fmt"
	"log"
	"time"
	"github.com/chawadev/payrail"
	"github.com/chawadev/payrail/payrailCore"
)

func main() {
	// Initialize
	client, err := payrail.NewClient("your-api-key", "lenco")
	if err != nil {
		log.Fatal(err)
	}

	// Request payment
	charge := core.ChargeRequest{
		Amount:     250.00,
		Reference:  "order-2024-001",
		Phone:      "+260769312808",
		Operator:   core.OperatorMTN,
		Country:    core.CountryZM,
		Bearer:     core.BearerCustomer,
		CustomerID: "user_456",
	}

	chargeResp, err := client.Charge(charge)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("✓ Charge requested\n")
	fmt.Printf("  Status: %s\n", chargeResp.Status)
	fmt.Printf("  Transaction: %s\n", chargeResp.TransactionID)

	// Wait a few seconds
	time.Sleep(3 * time.Second)

	// Check status
	statusResp, err := client.Veryfi(chargeResp.Reference)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\n✓ Status checked\n")
	fmt.Printf("  Current Status: %s\n", statusResp.Status)
	fmt.Printf("  Amount: %s %s\n", statusResp.Amount, statusResp.Currency)

	if statusResp.Status == "successful" {
		fmt.Printf("  Settlement: %s\n", statusResp.SettlementStatus)
	}
}
```

**Output:**
```
✓ Charge requested
  Status: pending
  Transaction: 49b9628d-b5cd-41c8-bd1e-98b6591360f5

✓ Status checked
  Current Status: successful
  Amount: 250.00 ZMW
  Settlement: pending
```

---

## Error Handling

```go
resp, err := client.Charge(req)
if err != nil {
    // Handle error
    if payErr, ok := err.(*core.PaymentError); ok {
        fmt.Printf("Error Code: %s\n", payErr.Code)
        fmt.Printf("Message: %s\n", payErr.Message)

        switch payErr.Code {
        case core.ErrCodeValidation:
            fmt.Println("Request validation failed")
        case core.ErrCodeNetwork:
            fmt.Println("Network error - retry later")
        case core.ErrCodeProvider:
            fmt.Println("Provider error")
        }
    }
}
```

### Common Errors

| Error | Cause | Fix |
|-------|-------|-----|
| `customer ID is required` | Missing CustomerID | Add `CustomerID: "your_id"` |
| `reference too short` | Reference < 6 characters | Use longer reference |
| `amount must be greater than zero` | Amount ≤ 0 | Set `Amount > 0` |
| `Duplicate reference` | Reference already used | Use unique reference |

---

## Validation

All requests are validated before sending:
- ✓ Amount > 0
- ✓ Reference length ≥ 6 chars
- ✓ CustomerID not empty
- ✓ Phone format valid
- ✓ Operator is valid
- ✓ Country is valid
- ✓ Bearer is valid

Invalid requests fail instantly without network calls.

---

## Status Values

**Payment Status:**
- `pending` - Awaiting customer confirmation
- `successful` - Payment received
- `failed` - Payment declined
- `pay-offline` - Customer needs to pay manually
- `3ds-auth-required` - Additional authentication needed

**Settlement Status:**
- `pending` - Funds not yet settled
- `settled` - Funds received
- `null` - Not applicable (still processing)

---

## Tips

1. **Always use unique references** - Duplicate references will be rejected
2. **Store the Transaction ID** - Keep `TransactionID` for your records
3. **Use Veryfi for polling** - Check status periodically or use webhooks
4. **Handle all statuses** - Some payments go to "pay-offline" state
5. **Retry on network errors** - Implement backoff for transient failures

---

## What Lenco Actually Receives

When you call `client.Charge()`, these fields are sent to Lenco:

```json
{
  "operator": "mtn",
  "bearer": "customer",
  "amount": 100.00,
  "reference": "order-2024-001",
  "country": "zm",
  "phone": "+260769312808"
}
```

When you call `client.Veryfi()`, a GET request is sent:
```
GET /access/v2/collections/status/2606007027
Authorization: Bearer your-api-key
```

---

## Support

For issues or questions:
1. Check error messages and codes
2. Verify all required fields are present
3. Ensure API key is valid
4. Check phone number format (must include country code)
