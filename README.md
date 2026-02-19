# Payrail - Go Payment Framework

Payrail is a lightweight, developer-friendly **Go payment framework** that allows you to integrate multiple payment providers (like Lenco) with minimal effort.  
It handles **request validation**, **provider configuration**, and **response parsing**, so you can focus on building your application, not payment plumbing.

---

## Features

- Clean API for charging customers
- Automatic request validation
- Configurable payment providers
- Structured response with transaction status
- Extensible to multiple providers in the future

---

## Installation

```bash
go get github.com/yourusername/payrail


## Usage

### 1. Initialize a Client

The client stores provider configuration (API key, base URL, etc.).
Set the provider at initialization:

```go
package main

import (
    "fmt"
    "github.com/yourusername/payrail"
)

func main() {
    // Initialize client for Lenco
    client := payrail.New("my-backend-api-key", "lenco")
}
```

* **Provider is set at client level** for safety and simplicity.
* Later, you can create another client for a different provider if needed.

---

### 2. Create a Charge Request

```go
req := payrail.ChargeRequest{
    Amount:      100,           // Payment amount
    Currency:    "ZMW",         // ISO 3-letter currency code
    CustomerID:  "cus_123",     // Unique customer identifier
    Reference:   "txn_001",     // Unique transaction reference
    Description: "Ticket purchase", // Optional
}
```

**Automatic validation includes:**

* Amount > 0
* Currency must be 3 letters (ISO code)
* CustomerID is required
* Reference is required and minimum 6 characters

---

### 3. Charge a Customer

```go
resp, err := client.Charge(req)
if err != nil {
    fmt.Println("Error:", err)
    return
}

fmt.Println("Transaction ID:", resp.TransactionID)
fmt.Println("Status:", resp.Status)
```

* Validation happens automatically before the request is sent.
* The client sends the request to the configured provider (Lenco in this case).
* `ChargeResponse` contains:

```go
type ChargeResponse struct {
    Status        string // "success", "pending", "failed"
    TransactionID string // Provider transaction ID
    Reference     string // Echoed from request
    RawResponse   []byte // Raw provider response for debugging
}
```

---

### 4. Adding Support for Multiple Providers

If your app needs multiple payment providers:

```go
lencoClient := payrail.New("lenco-key", "lenco")
duoClient := payrail.New("duo-key", "duo")
```

* Each client instance is bound to one provider.
* Requests are routed automatically to the correct backend configuration.

---

### 5. Why Use Payrail

* **Safety:** Automatic request validation prevents invalid payments.
* **Simplicity:** One call to `client.Charge()` handles all provider logic.
* **Extensible:** Add new providers without breaking existing code.
* **Developer-friendly:** Clean structs, clear API, minimal boilerplate.

```

---

If you want, the **next step is to write `client.go` and `Charge()`** exactly to match this README so developers can use the framework as described.  

Do you want me to do that next?
```
