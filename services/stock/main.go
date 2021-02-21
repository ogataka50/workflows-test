package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

type Reserve struct {
	Unit     uint64 `json:"unit"`
	Unstable bool   `json:"unstable"`
}

type Result struct {
	Status string `json:"status"`
	Unit   uint64 `json:"unit"`
}

func reserve(w http.ResponseWriter, req *http.Request) {
	var r Reserve
	err := json.NewDecoder(req.Body).Decode(&r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println("try to reserve stock")

	rand.Seed(time.Now().UnixNano())
	if r.Unstable && rand.Intn(5) < 4 {
		fmt.Println("stock service is unstable")
		http.Error(w, "stock service is unstable", http.StatusInternalServerError)
		return
	}

	var result Result
	if r.Unit > 30 {
		fmt.Println("no stock")
		result = Result{
			Status: "no stock",
		}
	} else {
		fmt.Println("reserved stock")
		result = Result{
			Status: "reserved",
			Unit:   r.Unit,
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

func cancelReserve(w http.ResponseWriter, req *http.Request) {
	fmt.Println("cancel reserve stock")

	res, err := json.Marshal(Result{
		Status: "cancel reserved",
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(res)
}

func main() {
	http.HandleFunc("/reserve", reserve)
	http.HandleFunc("/cancelReserve", cancelReserve)

	http.ListenAndServe(":8080", nil)
}
