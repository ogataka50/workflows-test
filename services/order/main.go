package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

type Order struct {
	Unit     uint64 `json:"unit"`
	Price    uint64 `json:"price"`
	Unstable bool   `json:"unstable"`
}

type Result struct {
	Status string `json:"status"`
}

func create(w http.ResponseWriter, req *http.Request) {
	var o Order
	err := json.NewDecoder(req.Body).Decode(&o)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println("order accepted, unit : " + strconv.Itoa(int(o.Unit)) + " price : $" + strconv.Itoa(int(o.Price)))

	err = executeWorkflow(o)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("executed workflows")

	res, err := json.Marshal(Result{
		Status: "accepted",
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(res)
}

func update(w http.ResponseWriter, req *http.Request) {
	fmt.Println("order updated")

	res, err := json.Marshal(Result{
		Status: "updated",
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(res)
}

func void(w http.ResponseWriter, req *http.Request) {
	fmt.Println("order voided")

	res, err := json.Marshal(Result{
		Status: "voided",
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(res)
}

func main() {
	http.HandleFunc("/create", create)
	http.HandleFunc("/update", update)
	http.HandleFunc("/void", void)

	http.ListenAndServe(":8080", nil)
}

func getToken() (string, error) {
	url := "http://metadata.google.internal/computeMetadata/v1/instance/service-accounts/default/token?scopes=https://www.googleapis.com/auth/cloud-platform"
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Metadata-Flavor", "Google")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	type st struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int64  `json:"expires_in"`
		TokenType   string `json:"token_type"`
	}

	var s st
	err = decoder.Decode(&s)
	if err != nil {
		return "", err
	}

	fmt.Println("response body:", s)

	return s.AccessToken, nil
}

func executeWorkflow(o Order) error {
	token, err := getToken()
	if err != nil {
		return err
	}

	projectId := os.Getenv("PROJECT_ID")

	url := "https://workflowexecutions.googleapis.com/v1/projects/" + projectId + "/locations/us-central1/workflows/workflow-test/executions"
	var jsonStr = []byte(`{"argument":"{\"unit\":` + strconv.Itoa(int(o.Unit)) + `, \"price\":` + strconv.Itoa(int(o.Price)) + `, \"unstable\":` + strconv.FormatBool(o.Unstable) + `}"}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
