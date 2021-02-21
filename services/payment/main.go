package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type Pay struct {
	Price    uint64 `json:"price"`
	Unstable bool   `json:"unstable"`
}

type Result struct {
	Status string `json:"status"`
	Price  uint64 `json:"price"`
}

func authorize(w http.ResponseWriter, req *http.Request) {
	var p Pay
	err := json.NewDecoder(req.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println("start authorize payment")

	rand.Seed(time.Now().UnixNano())
	if p.Unstable && rand.Intn(5) < 4 {
		fmt.Println("payment service is unstable")
		http.Error(w, "payment service is unstable", http.StatusInternalServerError)
		return
	}

	var result Result
	if p.Price > 100000 {
		fmt.Println("no money... : $" + strconv.Itoa(int(p.Price)))
		result = Result{
			Status: "no money",
		}
	} else {
		fmt.Println("authorize payment : $" + strconv.Itoa(int(p.Price)))
		result = Result{
			Status: "authorized",
			Price:  p.Price,
		}
	}

	res, err := json.Marshal(result)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(res)
}

func main() {

	http.HandleFunc("/authorize", authorize)

	http.ListenAndServe(":8080", nil)
}
