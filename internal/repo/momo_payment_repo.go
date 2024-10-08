package repo

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"mlvt/internal/entity"
	"net/http"
	"time"
)

type MoMoRepo interface {
	CreatePayment(orderID, amount string) (string, error)
	CheckPaymentStatus(orderID string) (bool, error)
	RefundPayment(orderID, amount string) (string, error)
}

type momoRepo struct {
	endpoint    string
	partnerCode string
	accessKey   string
	secrectKey  string
}

func (m *momoRepo) CreatePayment(orderID, amount string) (string, error) {
	momoRequest := entity.NewMoMoRequest(m.partnerCode, m.accessKey, orderID, amount, orderID)
	momoRequest.GenerateSignature(m.secrectKey)

	requestBody := momoRequest.ToMap()

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	client := &http.Client{Timeout: time.Second * 30}
	req, err := http.NewRequest("POST", m.endpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", errors.New("payment request failed")
	}

	var response struct {
		PayURL string `json:"payUrl"`
	}
	json.Unmarshal(body, &response)

	return response.PayURL, nil
}

func (m *momoRepo) CheckPaymentStatus(orderID string) (bool, error) {
	momoRequest := entity.NewMoMoRequest(m.partnerCode, m.accessKey, orderID, "0", orderID)
	momoRequest.GenerateSignature(m.secrectKey)

	requestBody := momoRequest.ToMap()
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return false, err
	}

	client := &http.Client{Timeout: time.Second * 30}
	req, err := http.NewRequest("POST", m.endpoint+"/check-status", bytes.NewBuffer(jsonBody))
	if err != nil {
		return false, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return false, errors.New("failed to check payment status")
	}

	var response struct {
		Status string `json:"status"`
	}
	json.Unmarshal(body, &response)

	return response.Status == "success", nil
}
