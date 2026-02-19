package lenco

import (
	"fmt"
	"io"
	"net/http"
	"payrail/core"
	"strings"
)

type LencoProvider struct {
}

func (l *LencoProvider) Charge(req core.ChargeRequest) (*core.ChargeResponse, error) {

	url := "https://api.lenco.co/access/v2/collections/mobile-money"

	payload := strings.NewReader("{\"operator\":\"airtel\",\"bearer\":\"merchant\"}")

	lencoReq, _ := http.NewRequest("POST", url, payload)

	lencoReq.Header.Add("accept", "application/json")
	lencoReq.Header.Add("content-type", "application/json")
	lencoReq.Header.Add("Authorization", "Bearer xo+CAiijrIy9XvZCYyhjrv0fpSAL6CfU8CgA+up1NXqK")

	res, _ := http.DefaultClient.Do(lencoReq)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	fmt.Println(string(body))

	return &core.ChargeResponse{
		Status:    "pending",
		Reference: req.Reference,
	}, nil
}
