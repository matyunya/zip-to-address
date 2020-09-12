package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/asaskevich/govalidator"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

type Address struct {
	Prefecture string `json:"prefecture"`
	City       string `json:"city"`
	Town       string `json:"town"`
}

var (
	m       map[string]Address
	csvPath *string
)

func init() {
	csvPath = flag.String("f", "./csv/out.csv", "path to zip codes file")
	flag.Parse()
}

func main() {
	err := readCsv(*csvPath)
	if err != nil {
		panic(err)
	}

	router := fasthttprouter.New()
	router.GET("/", handler)

	fmt.Println("Server is listening at 8050. Try http://localhost:8050/?z=4010502")
	err = fasthttp.ListenAndServe(":8050", router.Handler)
	if err != nil {
		log.Fatalf("Fatal error: %s", err)
	}
}

func readCsv(filename string) error {

	csvFile, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer csvFile.Close()

	csvData, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		return err
	}

	m = make(map[string]Address)

	for _, line := range csvData {
		m[line[2]] = Address{
			Prefecture: line[6],
			City:       line[7],
			Town:       line[8],
		}
	}

	return nil
}

func handler(ctx *fasthttp.RequestCtx) {
	index := string(ctx.QueryArgs().Peek("z"))

	ctx.Response.Header.SetContentType("application/json")

	valid := govalidator.IsNumeric(index) && len(index) > 0 && len(index) == 7
	if !valid {
		jsonData, _ := json.Marshal(map[string]string{"error": "Index must contain 7 digits"})
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.Write(jsonData)
		return
	}

	jsonData, err := json.Marshal(m[index])
	if err != nil {
		log.Println(err)
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.Write(jsonData)
}
