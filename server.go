package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type product struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Skus        []skus `json:"skus"`
}

type skus struct {
	Sku            string  `json:"sku"`
	Price          float64 `json:"price"`
	Stock          int     `json:"stock"`
	VariationType  string  `json:"variation_type"`
	VariationValue string  `json:"variation_value"`
}

var products = []product{
	{
		ID: 1, Title: "Ar condicionado", Description: "Ar condicionado",
		Skus: []skus{
			{Sku: "CODIGO001", Price: 179.99, Stock: 55, VariationType: "voltage", VariationValue: "110V"},
			{Sku: "CODIGO002", Price: 239.99, Stock: 45, VariationType: "voltage", VariationValue: "220V"},
		},
	},
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/products", getAllHandler).Methods("GET")
	r.HandleFunc("/products/{id:[0-9]+}", getByIdHandler).Methods("GET")
	r.HandleFunc("/products", newProductHandler).Methods("POST")
	r.HandleFunc("/products/{id:[0-9]+}", updateProductHandler).Methods("PUT")
	r.HandleFunc("/products/{id:[0-9]+}/skus/{sku:[A-Za-z0-9]+}", updateSkuHandler).Methods("PUT")
	http.ListenAndServe("localhost:8080", r)
}

// GET /products
func getAllHandler(w http.ResponseWriter, r *http.Request) {
	data, err := json.Marshal(products)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// GET /products/{ID}
func getByIdHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err.Error())
		return
	}

	for _, value := range products {
		if id == value.ID {
			data, err := json.Marshal(value)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Println(err.Error())
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write(data)
			return
		}
	}
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(`{"error": true, "message": "Product Id not found!"}`))
}

// POST /products
func newProductHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var data product
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err.Error())
		return
	}

	dataToSend, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err.Error())
		return
	}

	products = append(products, data)
	w.WriteHeader(http.StatusOK)
	w.Write(dataToSend)

}

// PUT /products/{id}
func updateProductHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err.Error())
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err.Error())
		return
	}

	var updatedProduct product
	err = json.Unmarshal(body, &updatedProduct)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err.Error())
		return
	}

	for key, value := range products {
		if id == value.ID {
			products[key] = updatedProduct
			w.WriteHeader(http.StatusOK)
			w.Write(body)
			return
		}
	}

	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(`{"error": true, "message": "Product Id not found!"}`))
}

// PUT /products/{id}/skus/{sku}
func updateSkuHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	sku := vars["sku"]

	fmt.Println(sku, id)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err.Error())
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err.Error())
		return
	}

	var updatedSku skus
	err = json.Unmarshal(body, &updatedSku)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err.Error())
		return
	}

	for productKey, value := range products {
		if id == value.ID {
			for skuKey, skuInfo := range value.Skus {
				if sku == skuInfo.Sku {
					products[productKey].Skus[skuKey] = updatedSku
					w.WriteHeader(http.StatusOK)
					w.Write(body)
					return
				}
			}
		}
	}

	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(`{"error": true, "message": "Product Id not found!"}`))
}
