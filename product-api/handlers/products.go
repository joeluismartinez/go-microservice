package handlers

import (
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/joeluismartinez/go-microservice/product-api/data"
)

type Product struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Product {
	return &Product{l}
}

func (p *Product) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		p.Getproducts(rw, r)
		return
	}

	if r.Method == http.MethodPost {
		p.addProduct(rw, r)
		return
	}

	if r.Method == http.MethodPut {
		p.l.Println("PUT")
		reg := regexp.MustCompile(`/([0-9]+)`)
		g := reg.FindAllStringSubmatch(r.URL.Path, -1)
		if len(g[0]) != 2 {
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}
		idString := g[0][1]
		id, err := strconv.Atoi(idString)
		if err != nil {
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}

		p.updateProducts(id, rw, r)
	}

	rw.WriteHeader(http.StatusMethodNotAllowed)
}

func (p *Product) Getproducts(rw http.ResponseWriter, r *http.Request) {
	lp := data.GetProducts()
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal JSON", http.StatusInternalServerError)
	}
}

func (p *Product) addProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle post product")
	prod := &data.Product{}
	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshall JSON", http.StatusBadRequest)
	}
	data.AddProduct(prod)
}

func (p *Product) updateProducts(id int, rw http.ResponseWriter, r *http.Request) {
	prod := &data.Product{}
	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshall JSON", http.StatusBadRequest)
	}
	err = data.UpdateProduct(id, prod)
	if err == data.ErrorProductNotFound {
		http.Error(rw, "Product Not Found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(rw, "Prodcut Not Found", http.StatusInternalServerError)
		return
	}
}
