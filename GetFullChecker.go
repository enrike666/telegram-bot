package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

type GetFullChecker interface {
	GetFullCheck(CheckInfo) (*Check, error)
}

type TelegramGetFullChecker struct {
	apiURL     string
	httpClient HttpDoer
}

func NewTelegramGetFullChecker(apiURL string, httpClient HttpDoer) *TelegramGetFullChecker {
	return &TelegramGetFullChecker{apiURL: apiURL, httpClient: httpClient}
}

func (t *TelegramGetFullChecker) GetFullCheck(checkInfo CheckInfo) (*Check, error) {
	data := url.Values{}
	data.Set("fn", checkInfo.FN)
	data.Set("fd", checkInfo.FD)
	data.Set("fp", checkInfo.FP)
	data.Set("n", checkInfo.N)
	data.Set("s", checkInfo.Sum)

	dateTime, err := checkInfo.GetTime()
	if err != nil {
		return nil, err
	}

	data.Set("t", dateTime.Format("01.06.2006 15:04"))

	requestBody := bytes.NewBufferString(data.Encode())
	request, err := http.NewRequest(http.MethodPost, t.apiURL, requestBody)
	if err != nil {
		return nil, err
	}

	request.Header.Add(contectTypeHeader, contectTypeXWWWWFormUrlencoded)

	contentLength := strconv.Itoa(len(data.Encode()))
	request.Header.Add(contentLengthHeader, contentLength)

	resp, err := t.httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	defer func() {
		errClose := resp.Body.Close()
		if errClose != nil {
			log.Println(errClose)
		}
	}()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: code: %d, body: %s", ErrUnexpectedResponse, resp.StatusCode, string(body))
	}

	var checkResponse Response
	err = json.Unmarshal(body, &checkResponse)
	if err != nil {
		var errResponse ErrResponse

		err = json.Unmarshal(body, &errResponse)
		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("%w: %s", ErrBadRequest, errResponse)
	}

	return &checkResponse.Data.Json, nil
}

type ErrResponse struct {
	Code int    `json:"code"`
	Data string `json:"data"`
}

func (e ErrResponse) String() string {
	return fmt.Sprintf("code: %d, message: %s", e.Code, e.Data)
}

var ErrBadRequest = errors.New("bad request from proverkacheka.com")
var ErrUnexpectedResponse = errors.New("unexpected response code returned")
