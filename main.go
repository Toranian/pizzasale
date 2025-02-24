package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly"
)

const (
	PizzaURL = "https://www.thriftyfoods.com/product/pizzaristorante-thin-crust-spinach/00000_000000005833617000"
)

type Pizza struct {
	Name     string
	Price    float64
	ImageURL string
}

type Status struct {
	OnSale bool
	Price  float64
}

func (p Pizza) String() string {
	return fmt.Sprintf("%s | %f", p.Name, p.Price)
}

func getPizzas(url string) []Pizza {
	var pizzas []Pizza

	c := colly.NewCollector()
	c.OnHTML("div.grid", func(e *colly.HTMLElement) {

		imgURL := e.ChildAttr("img", "src")

		title := e.ChildText("a.js-ga-productname")

		price := e.ChildText("span.price")
		if imgURL == "" || title == "" || price == "" {
			return
		}

		cleanedPrice := strings.Replace(price, "$", "", -1)
		p, err := strconv.ParseFloat(cleanedPrice, 64)

		if err != nil {
			return
		}

		pizza := Pizza{
			Name:     title,
			Price:    p,
			ImageURL: imgURL,
		}
		pizzas = append(pizzas, pizza)
	})

	err := c.Visit(url)

	if err != nil {
		fmt.Println("Error visiting page:", err)
	}

	return pizzas
}

func main() {

	r := gin.Default()

	r.LoadHTMLGlob("templates/*")

	r.GET("/", func(c *gin.Context) {
		pizzas := getPizzas(PizzaURL)
		onSale := true
		if len(pizzas) > 0 && pizzas[0].Price < 5 {
			onSale = true
		}

		c.HTML(200, "index.html", gin.H{"pizzas": pizzas, "status": Status{OnSale: onSale, Price: pizzas[0].Price}})
	})

	r.Run()
}
