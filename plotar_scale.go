package main

import (
	"os"
	"fmt"
	"time"
	"io"
	"os/signal"
	"flag"

	vegeta "./vegeta/lib"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/plotutil"
)

func scaleplot(data ...interface{}) {
	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	p.Title.Text = "Plotutil example"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	err = plotutil.AddLinePoints(p, data...)
	if err != nil {
		panic(err)
	}

	// Save the plot to a PNG file.
	if err := p.Save(4*vg.Inch, 4*vg.Inch, "points.png"); err != nil {
		panic(err)
	}
}

func boxplot(outputFile string, title string, yLabel string, data ...interface{}) {
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

	if err := p.Save(vg.Length(1+(len(data)/2))*vg.Inch, 4*vg.Inch, outputFile); err != nil {
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

func file(name string, create bool) (*os.File, error) {
	switch name {
	case "stdin":
		return os.Stdin, nil
	case "stdout":
		return os.Stdout, nil
	default:
		if create {
			return os.Create(name)
		}
		return os.Open(name)
	}
}

func get_data_from_file(filename string) plotter.Values {
	rc, err := file(filename, false)
	dec := vegeta.DecoderFor(rc)

	var r vegeta.Result
	dec.Decode(&r)

	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, os.Interrupt)

	//rc, _ := report.(vegeta.Closer)

	var results vegeta.Results
decode:
	for {
		select {
		case <-sigch:
			break decode
		default:
			var r vegeta.Result
			if err = dec.Decode(&r); err != nil {
				if err == io.EOF {
					break decode
				}
			}

			results.Add(&r)
		}
	}


	results.Close()

	data := make(plotter.Values, len(results))
	for i,res:=range results {
		data[i] = float64(res.Latency)/1000000
	}
	return data
}

func main() {
	//rate := vegeta.Rate{Freq: 2, Per: time.Second}
	//duration := 3 * time.Second
	//url := "http://10.0.3.187/500kb.html"

	//files := []string{"results.gob"}
	//dec, mc, err := decoder(files)
	//defer mc.Close()
	/*if err != nil {
		return err
	}*/


	//titlePtr := flag.String("title", "SEM TITULO", "Titulo do grafico")
	//yLabelPtr := flag.String("ylabel", "SEM LABEL", "Label do Eixo Y")
	//outputPtr := flag.String("output", "boxplot.png", "Nome do arquivo de Saida")
	flag.Parse()

	var data []interface{}
	for i, arg := range flag.Args() {
		if i%2 == 0 {
			fmt.Println("NAME: ", arg)
			data = append(data, arg)
		} else {
			fmt.Println("FILENAME: ", arg)
			data = append(data, get_data_from_file(arg))
		}
	}


	//boxplot(*outputPtr, *titlePtr, *yLabelPtr, data...)
	scaleplot(data...)


/*
	fmt.Print("Inicio...\n")
	os.Setenv("HTTP_PROXY", "http://10.0.3.188:3128")
	fmt.Println("Carga 0... ", os.Getenv("HTTP_PROXY"))
	data0 := carga(url, rate, duraion)
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
*/
}


