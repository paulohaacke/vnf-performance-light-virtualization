package main

import (
	"os"
	"fmt"
	"time"

	vegeta "./vegeta/lib"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/plotutil"
)

func boxplot(title string, yLabel string, data ...interface{}) {
	p, err := plot.New()
	if err != nil {
		panic(err)
	}
	p.Title.Text = title
	p.Y.Label.Text = yLabel


	err = plotutil.AddBoxPlots(p, vg.Points(40), data...)

	if err != nil {
		panic(err)
	}

	if err := p.Save(4*vg.Inch, 4*vg.Inch, "boxplot.png"); err != nil {
		panic(err)
	}
}

func carga(url string, rate vegeta.Pacer, duration time.Duration) plotter.Values  {
	targeter := vegeta.NewStaticTargeter(vegeta.Target{
		Method: "GET",
		URL:    url,
	})
	attacker := vegeta.NewAttacker()

	var results vegeta.Results
	for res := range attacker.Attack(targeter, rate, duration, "Carga") {
		results.Add(res)
	}

	results.Close()

	data := make(plotter.Values, len(results))
	for i,res:=range results {
		data[i] = float64(res.Latency)/1000000
	}
	return data

}

func main() {
	rate := vegeta.Rate{Freq: 2, Per: time.Second}
	duration := 3 * time.Second
	url := "http://10.0.3.187/500kb.html"

	fmt.Print("Inicio...\n")
	os.Setenv("HTTP_PROXY", "http://10.0.3.188:3128")
	fmt.Println("Carga 0... ", os.Getenv("HTTP_PROXY"))
	data0 := carga(url, rate, duration)
	os.Setenv("HTTP_PROXY", "http://10.0.3.106:3128")
	fmt.Println("Carga 1... ", os.Getenv("HTTP_PROXY"))
	data1 := carga(url, rate, duration)
	os.Setenv("HTTP_PROXY", "http://10.0.3.188:3128")
	fmt.Println("Carga 2... ", os.Getenv("HTTP_PROXY"))
	data2 := carga(url, rate, duration)
	os.Setenv("HTTP_PROXY", "http://10.0.3.106:3128")
	fmt.Println("Carga 3... ", os.Getenv("HTTP_PROXY"))
	data3 := carga(url, rate, duration)
	fmt.Print("Plot...\n")
	data := []interface{}{"Nativo", data0, "Docker", data1, "LXC", data2, "LXD", data3}
	boxplot("TITULO", "YLABEL", data...)
	fmt.Print("Fim.\n")
}


