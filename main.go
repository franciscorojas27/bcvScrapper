package main

import (
	"crypto/tls"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"

	"github.com/gocolly/colly"
	"github.com/rs/cors"
)

type PageData struct {
	Rates map[string]string
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	mux := http.NewServeMux()

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
	c.WithTransport(&http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    })
	c.UserAgent = "Mozilla/5.0 (iPhone; CPU iPhone OS 16_5 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.5 Mobile/15E148 Safari/604.1"
	c.OnHTML("#euro, #yuan, #lira, #rublo, #dolar", func(e *colly.HTMLElement) {
		idCoin := e.Attr("id")
		price := strings.TrimSpace(e.ChildText(".centrado strong"))
		rates[idCoin] = price
	})

	c.Visit("https://www.bcv.org.ve/")
	return rates
}
