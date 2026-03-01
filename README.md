# Payrail - Go Payment Framework

Payrail is a lightweight, developer-friendly **Go payment framework** that allows you to integrate multiple payment providers (like Lenco) with minimal effort. It handles **request validation**, **provider configuration**, and **response parsing**, so you can focus on building your application, not payment plumbing.

---

## Features

- Clean API for charging customers via mobile money
- Built-in payment status verification
- Automatic request validation
- Comprehensive error handling with error codes
- Configurable payment providers
- Structured responses with transaction details
- Extensible to multiple providers

---

## Installation

```bash
go get github.com/yourusername/payrail
```

---

## Quick Start

```go
package main

import (
    "fmt"
    "log"
    "payrail"
    "payrail/core"
)

func main() {
    // Initialize client
    client, err := payrail.NewClient("your-lenco-api-key", "lenco")
    if err != nil {
        log.Fatal(err)
    }

    // Request payment
    chargeReq := core.ChargeRequest{
        Amount:     0.50,
        Reference:  "order-12345",
        Phone:      "+260769312808",
        Operator:   core.OperatorMTN,
        Country:    core.CountryZM,
        Bearer:     core.BearerCustomer,
        CustomerID: "cus_456",
    }

    chargeResp, err := client.Charge(chargeReq)
    if err != nil {
        log.Fatal(err)
    }

    // Check payment status later
    status, err := client.Veryfi(chargeResp.Reference)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Payment Status: %s\n", status.Status)
}
```

---

## 1. Charge Endpoint

### Request Charge API

Create a charge to request payment from a customer via mobile money.

**Endpoint:** `POST https://api.lenco.co/access/v2/collections/mobile-money`

#### Request Parameters

| Field      | Type   | Required | Description |
|-----------|--------|----------|-------------|
| `operator` | string | Yes      | `"airtel"` or `"mtn"` |
| `bearer`   | string | Yes      | `"merchant"` or `"customer"` (who pays fees) |
| `amount`   | number | Yes      | Payment amount in local currency |
| `reference` | string | Yes     | Unique transaction reference (min 6 chars) |
| `country`  | string | Yes      | `"zm"` (Zambia) or `"mw"` (Malawi) |
| `phone`    | string | Yes      | Recipient phone (e.g., `+260769312808`) |

#### Go Code Example

```go
req := core.ChargeRequest{
    Amount:     50.00,                   // ZMW/MWK
    Reference:  "order-20260301-001",    // Your unique ID
    Phone:      "+260769312808",         // Customer phone
    Operator:   core.OperatorMTN,        // or OperatorAirtel
    Country:    core.CountryZM,          // or CountryMW
    Bearer:     core.BearerCustomer,     // or BearerMerchant
    CustomerID: "cus_78910",             // Your customer ID
}

resp, err := client.Charge(req)
if err != nil {
    log.Fatal(err)
}

fmt.Println("Status:", resp.Status)       // pending, successful, failed, pay-offline
fmt.Println("Transaction ID:", resp.TransactionID)
fmt.Println("Lenco Reference:", resp.Reference)
```

#### JSON Request Body

```json
{
  "operator": "mtn",
  "bearer": "customer",
  "amount": 50.00,
  "reference": "order-20260301-001",
  "country": "zm",
  "phone": "+260769312808"
}
```

#### Response Fields

```go
type ChargeResponse struct {
    Status        string // "pending" | "successful" | "failed" | "pay-offline"
    TransactionID string // Lenco's collection ID
    Reference     string // Lenco's internal reference for tracking
    RawResponse   []byte // Raw JSON response for debugging
}
```

#### Example Response

```json
{
  "status": true,
  "message": "Mobile money request initiated",
  "data": {
    "id": "coll_abc123xyz",
    "initiatedAt": "2026-03-01T14:30:00Z",
    "completedAt": null,
    "amount": "50.00",
    "fee": null,
    "bearer": "customer",
    "currency": "ZMW",
    "reference": "order-20260301-001",
    "lencoReference": "lref_789def456",
    "type": "mobile-money",
    "status": "pending",
    "source": "api",
    "reasonForFailure": null,
    "settlementStatus": null,
    "settlement": null,
    "mobileMoneyDetails": {
      "country": "zm",
      "phone": "+260769312808",
      "operator": "mtn",
      "accountName": "John Doe",
      "operatorTransactionId": "mtn_txn_ref_123"
    },
    "bankAccountDetails": null,
    "cardDetails": null
  }
}
```

---

## 2. Veryfi (Status Check) Endpoint

### Verify Payment Status

Query the current status of a payment request using its reference.

**Endpoint:** `GET https://api.lenco.co/access/v2/collections/status/{reference}`

#### Path Parameters

| Parameter | Type   | Required | Description |
|-----------|--------|----------|-------------|
| `reference` | string | Yes      | Lenco reference from charge response |

#### Go Code Example

```go
// After charging, verify the status
statusResp, err := client.Veryfi(chargeResp.Reference)
if err != nil {
    log.Fatal(err)
}

switch statusResp.Status {
case "pending":
    fmt.Println("Waiting for customer confirmation")
case "successful":
    fmt.Println("Payment received! Settlement status:", statusResp.SettlementStatus)
case "failed":
    fmt.Println("Payment failed:", statusResp.ReasonForFailure)
case "pay-offline":
    fmt.Println("Customer needs to complete payment offline")
case "3ds-auth-required":
    fmt.Println("3DS authentication required")
}
```

#### Response Fields

```go
type StatusResponse struct {
    Status           string // "pending" | "successful" | "failed" | "pay-offline" | "3ds-auth-required"
    TransactionID    string // Lenco's collection ID
    Reference        string // Lenco's reference for tracking
    Amount           string // Payment amount
    Currency         string // Currency code (e.g., "ZMW")
    Fee              string // Transaction fee
    SettlementStatus string // "pending" | "settled" | null
    ReasonForFailure string // Error message if failed
    RawResponse      []byte // Raw JSON for debugging
}
```

#### Example Response

```json
{
  "status": true,
  "message": "Collection found",
  "data": {
    "id": "coll_abc123xyz",
    "initiatedAt": "2026-03-01T14:30:00Z",
    "completedAt": "2026-03-01T14:35:00Z",
    "amount": "50.00",
    "fee": "0.50",
    "bearer": "customer",
    "currency": "ZMW",
    "reference": "order-20260301-001",
    "lencoReference": "lref_789def456",
    "type": "mobile-money",
    "status": "successful",
    "source": "api",
    "reasonForFailure": null,
    "settlementStatus": "pending",
    "settlement": {
      "id": "sett_xyz789",
      "amountSettled": "49.50",
      "currency": "ZMW",
      "createdAt": "2026-03-01T14:35:00Z",
      "settledAt": null,
      "status": "pending",
      "type": "next-day",
      "accountId": "acc_12345"
    },
    "mobileMoneyDetails": {
      "country": "zm",
      "phone": "+260769312808",
      "operator": "mtn",
      "accountName": "John Doe",
      "operatorTransactionId": "mtn_txn_ref_123"
    },
    "bankAccountDetails": null,
    "cardDetails": null
  }
}
```

---

## Error Handling

Payrail provides comprehensive error handling with structured error codes:

```go
import "payrail/core"

resp, err := client.Charge(req)
if err != nil {
    // Handle different error types
    if payErr, ok := err.(*core.PaymentError); ok {
        switch payErr.Code {
        case core.ErrCodeValidation:
            fmt.Println("Invalid request:", payErr.Message)
        case core.ErrCodeProvider:
            fmt.Println("Provider error:", payErr.Message)
        case core.ErrCodeNetwork:
            fmt.Println("Network error:", payErr.Message)
        }
    }
}
```

### Error Codes

| Code | Level | Meaning |
|------|-------|---------|
| `VALIDATION_ERROR` | Error | Request validation failed |
| `PROVIDER_ERROR` | Error | Provider API returned an error |
| `NETWORK_ERROR` | Error | Network connectivity issue |
| `PARSING_ERROR` | Error | Failed to parse response |
| `NOT_FOUND` | Error | Resource not found |
| `UNAUTHORIZED` | Error | API key invalid |
| `INTERNAL_ERROR` | Error | Unexpected framework error |

### Error Logger

Use the built-in error logger for structured logging:

```go
logger := core.NewErrorLogger()

// Log different severity levels
logger.LogErrorf(core.ErrorLevelInfo, "Payment initiated for reference %s", ref)
logger.LogErrorf(core.ErrorLevelWarning, "Slow response from provider")

// Log PaymentError with code
if payErr, ok := err.(*core.PaymentError); ok {
    logger.LogError(payErr)
}
```

---

## Complete Example

```go
package main

import (
    "fmt"
    "log"
    "payrail"
    "payrail/core"
    "time"
)

func main() {
    logger := core.NewErrorLogger()

    // 1. Initialize client
    client, err := payrail.NewClient("your-api-key", "lenco")
    if err != nil {
        logger.Fatalf("Failed to create client: %v", err)
    }

    // 2. Request charge
    charge := core.ChargeRequest{
        Amount:     100.00,
        Reference:  fmt.Sprintf("order-%d", time.Now().Unix()),
        Phone:      "+260769312808",
        Operator:   core.OperatorMTN,
        Country:    core.CountryZM,
        Bearer:     core.BearerCustomer,
        CustomerID: "cus_user_789",
    }

    resp, err := client.Charge(charge)
    if err != nil {
        logger.LogErrorf(core.ErrorLevelError, "Charge failed: %v", err)
        return
    }

    logger.LogErrorf(core.ErrorLevelInfo, "Charge created: %s (Status: %s)", 
        resp.TransactionID, resp.Status)

    // 3. Poll status (in production, use webhooks instead)
    for attempt := 0; attempt < 5; attempt++ {
        time.Sleep(2 * time.Second)

        status, err := client.Veryfi(resp.Reference)
        if err != nil {
            logger.LogErrorf(core.ErrorLevelWarning, "Status check failed: %v", err)
            continue
        }

        logger.LogErrorf(core.ErrorLevelInfo, "Payment status: %s", status.Status)

        if status.Status == "successful" {
            fmt.Printf("✓ Payment received! Amount: %s %s\n", 
                status.Amount, status.Currency)
            break
        } else if status.Status == "failed" {
            fmt.Printf("✗ Payment failed: %s\n", status.ReasonForFailure)
            break
        }
    }
}
```

---

## Request Validation

All requests are automatically validated before sending to the provider:

```go
type ChargeRequest struct {
    Amount    float64 // Must be > 0
    Reference string  // Required, minimum 6 characters
    Phone     string  // Required
    Operator  Operator // Required (airtel or mtn)
    Country   Country  // Required (zm or mw)
    Bearer    Bearer   // Required (merchant or customer)
    CustomerID string  // Required
}
```

Validation happens automatically in `client.Charge()`. Invalid requests return a validation error immediately without hitting the network.

---

## Adding More Providers

To add a new provider (e.g., Duo Payment):

1. Create a provider package:
```go
// provider/duo/charge.go
package duo

type DuoProvider struct {
    config Config
}

func (d *DuoProvider) Charge(req core.ChargeRequest) (*core.ChargeResponse, error) {
    // Implementation
}

func (d *DuoProvider) Veryfi(reference string) (*core.StatusResponse, error) {
    // Implementation
}
```

2. Wire it in `client.go`:
```go
case "duo":
    p := duo.NewProvider(duo.Config{APIKey: apiKey})
    return &Client{provider: p}, nil
```

3. Document the provider in README

---

## Why Payrail

✅ **Safe** - Automatic validation prevents invalid payments  
✅ **Simple** - One call to `Charge()` or `Veryfi()` handles all logic  
✅ **Structured** - Error codes and types for programmatic handling  
✅ **Extensible** - Add providers without breaking existing code  
✅ **Developer-friendly** - Clean API, clear documentation
