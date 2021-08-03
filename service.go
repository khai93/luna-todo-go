package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type ServiceData struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Version     string `json:"version"`
}

type LunaAuth struct {
	user string
	pass string
}

type BalancerOptions struct {
	weight int
}

type InstanceData struct {
	InstanceId      string `json:"instanceId"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	Version         string `json:"version"`
	Status          string `json:"status"`
	BalancerOptions `json:"balancerOptions"`
	Url             string `json:"url"`
}

var lunaAuth = LunaAuth{"admin", "secretphrase2"}
var lunaUrl = fmt.Sprintf("http://%s:%s@localhost:5000", lunaAuth.user, lunaAuth.pass)

func getServiceData() (*ServiceData, error) {
	serviceFile, err := os.Open("service.json")
	if err != nil {
		return nil, err
	}

	defer serviceFile.Close()
	fileBytes, err := ioutil.ReadAll(serviceFile)
	if err != nil {
		return nil, err
	}

	var data ServiceData
	if err := json.Unmarshal(fileBytes, &data); err != nil {
		return nil, err
	}

	return &data, nil
}

func StartService() {
	serviceData, err := getServiceData()
	if err != nil {
		log.Fatal(err)
	}

	instanceId := fmt.Sprintf("%s:localhost:4000", serviceData.Name)
	instanceUrl := fmt.Sprintf("%s/registry/v1/instances/%s", lunaUrl, instanceId)
	instanceData := InstanceData{
		InstanceId:      instanceId,
		Name:            serviceData.Name,
		Description:     serviceData.Description,
		Version:         serviceData.Version,
		Status:          "OK",
		BalancerOptions: BalancerOptions{},
		Url:             "http://localhost:4000",
	}

	client := &http.Client{}

	body, err := json.Marshal(instanceData)
	if err != nil {
		log.Fatal(err)
	}

	instanceDataBuffer := bytes.NewBuffer(body)

	// Send heartbeat to Luna
	req, err := http.NewRequest(http.MethodPut, instanceUrl, instanceDataBuffer)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		// Register Instance
		resp, err := http.Post(instanceUrl, "application/json", bytes.NewBuffer(body))
		if err != nil {
			log.Fatal(err)
		}

		if resp.StatusCode != 201 {
			log.Fatal(errors.New("Could not register instance."))
		}

		defer resp.Body.Close()

		log.Println("Registered Instance.")
	} else {
		log.Println("Sent Heartbeat.")
	}

	time.AfterFunc(15*time.Second, StartService)
}
