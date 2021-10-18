package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"regexp"

	"github.com/gocarina/gocsv"
	"github.com/hokaccha/go-prettyjson"
)

type BankStatement struct {
	Date        string  `csv:"Date"`
	Transaction string  `csv:"TransactionType"`
	Description string  `csv:"Description"`
	PaidOut     float64 `csv:"Paidout"`
	PaidIn      float64 `csv:"Paidin"`
	Balance     float64 `csv:"Balance"`
	Catagory    string
}

func main() {
	in, err := os.Open("statement-july.csv")
	if err != nil {
		panic(err)
	}
	defer in.Close()

	bankStatement := []*BankStatement{}

	if err := gocsv.UnmarshalFile(in, &bankStatement); err != nil {
		panic(err)
	}

	var total float64 = 0.0
	myRegex, err := regexp.Compile(`(ONE\sSTOP)|(MCDONALDS)|(TESCO)|(ASDA)|(CO\-OP)|(SAINSBURYS)|(UBER\s\*EATS)|(GREGGS)|(Olive\sTree)|(THE\sBAY\sVIEW\sINN)|(MORRISONS)|(ALLSTARSSP\*)|(TAUNTON\sDEANE\sCRICKET)`)

	for _, line := range bankStatement {
		found := myRegex.MatchString(line.Description)
		if err != nil {
			panic(err)
		}
		if found {
			total += line.PaidOut
			line.Catagory = "food"
		} else {
			line.Catagory = "unknown"
		}
	}

	file, err := json.MarshalIndent(bankStatement, "", "  ")
	if err != nil {
		panic(err)
	}

	_ = ioutil.WriteFile("statment_as_json.json", file, 0600)

	log(bankStatement)
	m := make(map[string]float64)
	m["food-total"] = math.Round(total*100) / 100
	log(m)

}

func log(data interface{}) {
	s, _ := prettyjson.Marshal(data)
	fmt.Println(string(s))
}
