package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"strings"
	"time"
)

func main() {

	ticker := time.NewTicker(time.Minute * 5)

	go func() {
		for _ = range ticker.C {
			fmt.Println("tick")

			go func() {
				CheckMi3()
			}()

			iphones := []string{
				"http://store.apple.com/sg/buyFlowSelectionSummary/IPHONE6P?node=home/shop_iphone/family/iphone6&step=select&option.dimensionScreensize=5_5inch&option.dimensionColor=gold&option.dimensionCapacity=64gb&option.carrierModel=UNLOCKED%2FWW&carrierPolicyType=UNLOCKED",
				"http://store.apple.com/sg/buyFlowSelectionSummary/IPHONE6P?node=home/shop_iphone/family/iphone6&step=select&option.dimensionScreensize=5_5inch&option.dimensionColor=silver&option.dimensionCapacity=64gb&option.carrierModel=UNLOCKED%2FWW&carrierPolicyType=UNLOCKED",
				"http://store.apple.com/sg/buyFlowSelectionSummary/IPHONE6P?node=home/shop_iphone/family/iphone6&step=select&option.dimensionScreensize=5_5inch&option.dimensionColor=space_gray&option.dimensionCapacity=64gb&option.carrierModel=UNLOCKED%2FWW&carrierPolicyType=UNLOCKED",
				"http://store.apple.com/sg/buyFlowSelectionSummary/IPHONE6?node=home/shop_iphone/family/iphone6&step=select&option.dimensionScreensize=4_7inch&option.dimensionColor=gold&option.dimensionCapacity=64gb&option.carrierModel=UNLOCKED%2FWW&carrierPolicyType=UNLOCKED",
				"http://store.apple.com/sg/buyFlowSelectionSummary/IPHONE6?node=home/shop_iphone/family/iphone6&step=select&option.dimensionScreensize=4_7inch&option.dimensionColor=gold&option.dimensionCapacity=64gb&option.carrierModel=UNLOCKED%2FWW&carrierPolicyType=UNLOCKED",
			}

			for _, url := range iphones {
				go func() {
					CheckIPhone6(url)
				}()
			}
		}
	}()

	<-make(chan int)
}

func CheckMi3() {
	resp, _ := http.Get("http://www.mi.com/sg/mi3/")

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	html := string(body)

	if strings.Contains(html, "Out of Stock") {
		SendMail("mi3 available.")
	}
}

type CheckIPhoneResponse struct {
	Head Head
	Body Body
}

type Head struct {
	Status string
	Data   Data
}

type Data struct {
}

type PurchaseOptions struct {
	ShippingPrice string
	Price         string
	ShippingLead  string
	Financing     string
	Promotions    string
	IsBuyable     bool
}

type Selected struct {
	PurchaseOptions PurchaseOptions
	PartNumber      string
	ProductImage    string
	ProductTitle    string
}

type Content struct {
	PageTitle string
	PageUrl   string
	Selected  Selected
}

type Body struct {
	Content Content
}

type Configuration struct {
	SmtpUsername  string
	SmtpPasswd    string
	EmailRecipent string
}

func CheckIPhone6(productUrl string) {

	resp, _ := http.Get(productUrl)
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	checkIPhoneResponse := ParseCheckIPhone6JsonResp(body)

	if strings.Contains(checkIPhoneResponse.Body.Content.Selected.PurchaseOptions.ShippingLead, "Currently unavailable") {
		message := fmt.Sprintf("\n %v \n %v", checkIPhoneResponse.Body.Content.Selected.ProductTitle, checkIPhoneResponse.Body.Content.Selected.PurchaseOptions.ShippingLead)
		SendMail(message)
	}
}

func ParseCheckIPhone6JsonResp(jsonBytes []byte) CheckIPhoneResponse {
	var m CheckIPhoneResponse

	err := json.Unmarshal(jsonBytes, &m)
	if err != nil {
		log.Fatal(err)
	}

	return m
}

func SendMail(body string) {

	file, _ := os.Open("conf.json")
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Printf("#%v\n", configuration)

	auth := smtp.PlainAuth("", configuration.SmtpUsername, configuration.SmtpPasswd, "smtp.live.com")

	to := []string{configuration.EmailRecipent}
	msg := []byte(body)
	err = smtp.SendMail("smtp.live.com:587", auth, configuration.SmtpUsername, to, msg)
	if err != nil {
		log.Fatal(err)
	}
}
