package main

import (
	"fmt"
	"log"
	"payrail"
	"payrail/core"
)

func main() {
	client, err := payrail.NewClient("f423b2e4be7e45bf63e378e12baf5d56d0781ddc094aa2884764759a9bc68984", "lenco")
	if err != nil {
		log.Fatal(err)
	}

	req := core.ChargeRequest{
		Amount:    0.1,
		Country:   core.CountryZM,
		Reference: "test-ref-123",
		Phone:     "+260769312808",
		Operator:  core.OperatorMTN,
		Bearer:    core.BearerCustomer,
	}

	resp, err := client.Charge(req)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Status:", resp.Status)
	fmt.Println("Reference:", resp.Reference)
	fmt.Println("Raw:", string(resp.RawResponse))
}
