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
	file, err := os.Open("statement-sept.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	bank_statement := []*BankStatement{}

	if err := gocsv.UnmarshalFile(file, &bank_statement); err != nil {
		panic(err)
	}

	var total float64
	my_regex, err := regexp.Compile(`(ONE\sSTOP)|(MCDONALDS)|(TESCO)|(ASDA)|(CO\-OP)|(SAINSBURYS)|(UBER\s\*EATS)|(GREGGS)|(Olive\sTree)|(THE\sBAY\sVIEW\sINN)|(MORRISONS)|(ALLSTARSSP\*)|(TAUNTON\sDEANE\sCRICKET)`)
	if err != nil {
		panic(err)
	}

	for _, line := range bank_statement {
		matched := my_regex.MatchString(line.Description)
		if err != nil {
			panic(err)
		}
		if matched {
			total += line.PaidOut
			line.Catagory = "food"
		} else {
			line.Catagory = "unknown"
		}
	}

	json_file, err := json.MarshalIndent(bank_statement, "", "  ")
	if err != nil {
		panic(err)
	}

	_ = ioutil.WriteFile("statment_as_json.json", json_file, 0600)

	log(bank_statement)
	m := make(map[string]float64)
	m["food-total"] = math.Round(total*100) / 100
	log(m)

}

func log(data interface{}) {
	s, _ := prettyjson.Marshal(data)
	fmt.Println(string(s))
}
