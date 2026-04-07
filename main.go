package main

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/rs/cors"
	"html/template"
	"net/http"
	"os"
	"strings"
)

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
	Id   int    `json:"id"`
}

type PageData struct {
	Rates map[string]string
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/people", func(w http.ResponseWriter, r *http.Request) {
		people := []Person{
			{Name: "Person 1", Age: 19, Id: 1},
			{Name: "Person 2", Age: 20, Id: 2},
			{Name: "Person 3", Age: 21, Id: 3},
			{Name: "Person 4", Age: 22, Id: 4},
			{Name: "Person 5", Age: 23, Id: 5},
			{Name: "Person 6", Age: 24, Id: 6},
			{Name: "Person 7", Age: 25, Id: 7},
			{Name: "Person 8", Age: 26, Id: 8},
			{Name: "Person 9", Age: 27, Id: 9},
			{Name: "Person 10", Age: 28, Id: 10},
			{Name: "Person 11", Age: 29, Id: 11},
			{Name: "Person 12", Age: 30, Id: 12},
			{Name: "Person 13", Age: 31, Id: 13},
			{Name: "Person 14", Age: 32, Id: 14},
			{Name: "Person 15", Age: 33, Id: 15},
			{Name: "Person 16", Age: 34, Id: 16},
			{Name: "Person 17", Age: 35, Id: 17},
			{Name: "Person 18", Age: 36, Id: 18},
			{Name: "Person 19", Age: 37, Id: 19},
			{Name: "Person 20", Age: 38, Id: 20},
			{Name: "Person 21", Age: 39, Id: 21},
			{Name: "Person 22", Age: 40, Id: 22},
			{Name: "Person 23", Age: 41, Id: 23},
			{Name: "Person 24", Age: 42, Id: 24},
			{Name: "Person 25", Age: 43, Id: 25},
			{Name: "Person 26", Age: 44, Id: 26},
			{Name: "Person 27", Age: 45, Id: 27},
			{Name: "Person 28", Age: 46, Id: 28},
			{Name: "Person 29", Age: 47, Id: 29},
			{Name: "Person 30", Age: 48, Id: 30},
			{Name: "Person 31", Age: 49, Id: 31},
			{Name: "Person 32", Age: 50, Id: 32},
			{Name: "Person 33", Age: 51, Id: 33},
			{Name: "Person 34", Age: 52, Id: 34},
			{Name: "Person 35", Age: 53, Id: 35},
			{Name: "Person 36", Age: 54, Id: 36},
			{Name: "Person 37", Age: 55, Id: 37},
			{Name: "Person 38", Age: 56, Id: 38},
			{Name: "Person 39", Age: 57, Id: 39},
			{Name: "Person 40", Age: 58, Id: 40},
			{Name: "Person 41", Age: 59, Id: 41},
			{Name: "Person 42", Age: 60, Id: 42},
			{Name: "Person 43", Age: 61, Id: 43},
			{Name: "Person 44", Age: 62, Id: 44},
			{Name: "Person 45", Age: 63, Id: 45},
			{Name: "Person 46", Age: 64, Id: 46},
			{Name: "Person 47", Age: 65, Id: 47},
			{Name: "Person 48", Age: 66, Id: 48},
			{Name: "Person 49", Age: 67, Id: 49},
			{Name: "Person 50", Age: 68, Id: 50},
			{Name: "Person 51", Age: 69, Id: 51},
			{Name: "Person 52", Age: 70, Id: 52},
			{Name: "Person 53", Age: 18, Id: 53},
			{Name: "Person 54", Age: 19, Id: 54},
			{Name: "Person 55", Age: 20, Id: 55},
			{Name: "Person 56", Age: 21, Id: 56},
			{Name: "Person 57", Age: 22, Id: 57},
			{Name: "Person 58", Age: 23, Id: 58},
			{Name: "Person 59", Age: 24, Id: 59},
			{Name: "Person 60", Age: 25, Id: 60},
			{Name: "Person 61", Age: 26, Id: 61},
			{Name: "Person 62", Age: 27, Id: 62},
			{Name: "Person 63", Age: 28, Id: 63},
			{Name: "Person 64", Age: 29, Id: 64},
			{Name: "Person 65", Age: 30, Id: 65},
			{Name: "Person 66", Age: 31, Id: 66},
			{Name: "Person 67", Age: 32, Id: 67},
			{Name: "Person 68", Age: 33, Id: 68},
			{Name: "Person 69", Age: 34, Id: 69},
			{Name: "Person 70", Age: 35, Id: 70},
			{Name: "Person 71", Age: 36, Id: 71},
			{Name: "Person 72", Age: 37, Id: 72},
			{Name: "Person 73", Age: 38, Id: 73},
			{Name: "Person 74", Age: 39, Id: 74},
			{Name: "Person 75", Age: 40, Id: 75},
			{Name: "Person 76", Age: 41, Id: 76},
			{Name: "Person 77", Age: 42, Id: 77},
			{Name: "Person 78", Age: 43, Id: 78},
			{Name: "Person 79", Age: 44, Id: 79},
			{Name: "Person 80", Age: 45, Id: 80},
			{Name: "Person 81", Age: 46, Id: 81},
			{Name: "Person 82", Age: 47, Id: 82},
			{Name: "Person 83", Age: 48, Id: 83},
			{Name: "Person 84", Age: 49, Id: 84},
			{Name: "Person 85", Age: 50, Id: 85},
			{Name: "Person 86", Age: 51, Id: 86},
			{Name: "Person 87", Age: 52, Id: 87},
			{Name: "Person 88", Age: 53, Id: 88},
			{Name: "Person 89", Age: 54, Id: 89},
			{Name: "Person 90", Age: 55, Id: 90},
			{Name: "Person 91", Age: 56, Id: 91},
			{Name: "Person 92", Age: 57, Id: 92},
			{Name: "Person 93", Age: 58, Id: 93},
			{Name: "Person 94", Age: 59, Id: 94},
			{Name: "Person 95", Age: 60, Id: 95},
			{Name: "Person 96", Age: 61, Id: 96},
			{Name: "Person 97", Age: 62, Id: 97},
			{Name: "Person 98", Age: 63, Id: 98},
			{Name: "Person 99", Age: 64, Id: 99},
			{Name: "Person 100", Age: 65, Id: 100},
			{Name: "Person 101", Age: 66, Id: 101},
			{Name: "Person 102", Age: 67, Id: 102},
			{Name: "Person 103", Age: 68, Id: 103},
			{Name: "Person 104", Age: 69, Id: 104},
			{Name: "Person 105", Age: 70, Id: 105},
			{Name: "Person 106", Age: 18, Id: 106},
			{Name: "Person 107", Age: 19, Id: 107},
			{Name: "Person 108", Age: 20, Id: 108},
			{Name: "Person 109", Age: 21, Id: 109},
			{Name: "Person 110", Age: 22, Id: 110},
			{Name: "Person 111", Age: 23, Id: 111},
			{Name: "Person 112", Age: 24, Id: 112},
			{Name: "Person 113", Age: 25, Id: 113},
			{Name: "Person 114", Age: 26, Id: 114},
			{Name: "Person 115", Age: 27, Id: 115},
			{Name: "Person 116", Age: 28, Id: 116},
			{Name: "Person 117", Age: 29, Id: 117},
			{Name: "Person 118", Age: 30, Id: 118},
			{Name: "Person 119", Age: 31, Id: 119},
			{Name: "Person 120", Age: 32, Id: 120},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(people)
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		rates := scrapeBCV()
		tmpl := `
		<!DOCTYPE html>
		<html>
		<head>
			<title>Tasas BCV</title>
			<style>
				body { font-family: sans-serif; display: flex; justify-content: center; background: #f4f4f9; }
				table { border-collapse: collapse; width: 300px; background: white; margin-top: 50px; box-shadow: 0 4px 8px rgba(0,0,0,0.1); }
				th, td { border: 1px solid #ddd; padding: 12px; text-align: left; }
				th { background-color: #004a99; color: white; }
				tr:nth-child(even) { background-color: #f2f2f2; }
			</style>
		</head>
		<body>
			<table>
				<tr><th>Moneda</th><th>Precio (Bs.)</th></tr>
				{{range $coin, $price := .Rates}}
				<tr>
					<td><strong>{{$coin | printf "%s" | u}}</strong></td>
					<td>{{$price}}</td>
				</tr>
				{{end}}
			</table>
		</body>
		</html>`

		funcMap := template.FuncMap{"u": strings.ToUpper}
		t, _ := template.New("webpage").Funcs(funcMap).Parse(tmpl)

		data := PageData{Rates: rates}
		t.Execute(w, data)
	})
	handler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	}).Handler(mux)
	fmt.Println("Servidor corriendo en http://localhost:8080")
	http.ListenAndServe(":"+port, handler)
}
func scrapeBCV() map[string]string {
	rates := make(map[string]string)
	rates["DEBUG"] = "Debug"
	c := colly.NewCollector()
	c.UserAgent = "Mozilla/5.0 (iPhone; CPU iPhone OS 16_5 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.5 Mobile/15E148 Safari/604.1"
	c.OnHTML("#euro, #yuan, #lira, #rublo, #dolar", func(e *colly.HTMLElement) {
		idCoin := e.Attr("id")
		price := strings.TrimSpace(e.ChildText(".centrado strong"))
		rates[idCoin] = price
	})

	c.Visit("https://www.bcv.org.ve/")
	return rates
}
