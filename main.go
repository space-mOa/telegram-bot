package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"io/ioutil"
	"net/http"
)

type Telegram struct {
	Url        string
	Clovek     string
	Metoda     string
	CilovaUrl  string
	KonecnaUrl string
}

type Balicek struct {
	Odkazy []string
	Nazev  []string
	Text   []string
}

func (t *Telegram) odesliZpravu(zpravy []string) {
	adresa := t.Url + t.Metoda + t.Clovek + "&text="
	for i, link := range zpravy {
		if i <= 1 {
			continue
		}
		t.KonecnaUrl = t.KonecnaUrl + link + "%0A" + "%0A"
	}
	adresa = adresa + t.KonecnaUrl
	odpoved, err := http.Get(adresa)
	if err != nil {
		fmt.Println(" Nepodařilo se poslat pozadavek. \n error: ", err)
	}
	defer odpoved.Body.Close()
	sbirka, err := ioutil.ReadAll(odpoved.Body)
	if err != nil {
		fmt.Println(" Nepodařilo se rozbalit odpověd. \n error: ", err)
	}
	fmt.Println("sbirka: ", string(sbirka))
}

func main() {
	fmt.Println("AHOJ")
	Bot := Telegram{Url: "vaše URL",
		Clovek: "chat_id= uživatele", Metoda: "sendMessage?", CilovaUrl: "http://www.kvic.cz/aktuality/8/1/Pracovni_prilezitosti"}
	var Dopisy Balicek
	css_polozka := ".h2Link"
	c := colly.NewCollector()
	c.OnHTML(css_polozka, func(e *colly.HTMLElement) {
		// atribut html
		link := e.Attr("href")
		fmt.Printf("Link found: %q -> %v\n\n", e.Text, link)
		Dopisy.Nazev = append(Dopisy.Nazev, e.Text)
		c.Visit(e.Request.AbsoluteURL(link))
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Prohledávám: ", r.URL.String())
		Dopisy.Odkazy = append(Dopisy.Odkazy, r.URL.String())
	})
	c.Visit(Bot.CilovaUrl)
	fmt.Println("\nLinky: ", Dopisy.Odkazy[1], "\nLinky: ", Dopisy.Odkazy[2], "\n")
	Bot.odesliZpravu(Dopisy.Odkazy)
}
