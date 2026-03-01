package main

import (
	"fmt"
	"log"
	"payrail"
	"payrail/core"
)

func main() {
	// Initialize client with Lenco provider
	client, err := payrail.NewClient("f423b2e4be7e45bf63e378e12baf5d56d0781ddc094aa2884764759a9bc68984", "lenco")
	if err != nil {
		log.Fatal(err)
	}

	// Test 1: Charge a customer
	fmt.Println("=== Testing Charge Endpoint ===")
	chargeReq := core.ChargeRequest{
		Amount:     0.1,
		Country:    core.CountryZM,
		Reference:  "test-ref-12345",
		Phone:      "+260769312808",
		Operator:   core.OperatorMTN,
		Bearer:     core.BearerCustomer,
		CustomerID: "cus_12345",
	}

	chargeResp, err := client.Charge(chargeReq)
	if err != nil {
		log.Printf("Charge error: %v\n", err)
	} else {
		fmt.Println("Charge Status:", chargeResp.Status)
		fmt.Println("Transaction ID:", chargeResp.TransactionID)
		fmt.Println("Lenco Reference:", chargeResp.Reference)
	}

	// Test 2: Verify payment status
	fmt.Println("\n=== Testing Veryfi (Status Check) Endpoint ===")
	if chargeResp != nil && chargeResp.Reference != "" {
		statusResp, err := client.Veryfi(chargeResp.Reference)
		if err != nil {
			log.Printf("Veryfi error: %v\n", err)
		} else {
			fmt.Println("Status:", statusResp.Status)
			fmt.Println("Transaction ID:", statusResp.TransactionID)
			fmt.Println("Reference:", statusResp.Reference)
		}
	}
}
