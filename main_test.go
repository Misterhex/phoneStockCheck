package main

import (
	"io/ioutil"
	"strings"
	"testing"
	"time"
)

func TestShouldCheckIphone6Plus(t *testing.T) {
	CheckIphone6Plus()
}

func TestParseCheckIPhone6JsonResp(t *testing.T) {
	respBytes, _ := ioutil.ReadFile("iphone6plus_resp.json")
	checkIPhoneResponse := ParseCheckIPhone6JsonResp(respBytes)

	if strings.Contains(checkIPhoneResponse.Body.Content.Selected.PurchaseOptions.ShippingLead, "Currently Unavailable") {
		t.Fail()
	}
}

func TestShouldReallySendEmail(t *testing.T) {
	SendMail("send from test @ " + time.Now().String())
}
