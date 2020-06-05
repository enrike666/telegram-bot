package main

import (
	"encoding/json"
	"github.com/gorilla/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

var validCheckInfo = CheckInfoFromBot{
	DateTimeString: "20200601T2122",
	Sum:            "213.07",
	FN:             "9282440300649733",
	FD:             "12416",
	FP:             "2858316733",
	N:              "1",
}

type MockHadler struct {
}

func (m MockHadler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	validResponse := Response{
		Data: Data{
			Json: Check{
				OrganizationAddress: "390037,62,РЯЗАНСКАЯ ОБЛАСТЬ,РЯЗАНЬ Г,ЗУБКОВОЙ УЛ,СТРОЕНИЕ 25А",
				OrganizationName:    "Агроторг ООО",
				Products: []Product{
					{
						Name:     "*3220501 GREENF.Чай BARB.GARDEN  25х1,5г",
						Price:    5999,
						Quantity: 1,
						Sum:      5999,
					},
					{
						Name:     "3683200 ЛОЖКАРЕВ Котлеты кур.в пан.500г",
						Price:    13909,
						Quantity: 1,
						Sum:      13909,
					},
					{
						Name:     "3221397 NESC.Нап.МЯГКИЙ коф.3в1 16г",
						Price:    1399,
						Quantity: 1,
						Sum:      1399,
					},
				},
				Sum: 21307,
			},
		},
	}

	requestBytes, err := ioutil.ReadAll(request.Body)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	bodyRequestValues, err := url.ParseQuery(string(requestBytes))
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	var decoder = schema.NewDecoder()
	var checkInfoRequest CheckInfoRequest

	err = decoder.Decode(&checkInfoRequest, bodyRequestValues)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	validCheckInfoRequest, err := NewCheckInfoRequestBasedCheckInfoFromBot(validCheckInfo)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if checkInfoRequest == *validCheckInfoRequest {
		responseBody, err := json.Marshal(&validResponse)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		writer.WriteHeader(http.StatusOK)
		_, _ = writer.Write(responseBody)
		return
	}

	if checkInfoRequest.FD == "" {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	//invalidResponse := struct {
	//
	//}{}
	//invalidResponse.Data.Json.OrganizationAddress = "спортивная"
	//invalidResponse.Data.Json.OrganizationName = "ддд"
	//invalidResponse.Data.Json.Sum = "строка"

	//responseBody, err := json.Marshal(&invalidResponse)
	//if err != nil {
	//	writer.WriteHeader(http.StatusInternalServerError)
	//	return
	//}

	//writer.WriteHeader(http.StatusOK)
	//_, _ = writer.Write(responseBody)
}

func TestGetFullCheck(t *testing.T) {
	server := httptest.NewServer(MockHadler{})
	defer server.Close()

	testCases := []struct {
		name      string
		apiUrl    string
		checkInfo CheckInfoFromBot
		wantCheck *Check
		wantErr   bool
	}{
		{
			name:      "positiv",
			apiUrl:    server.URL,
			checkInfo: validCheckInfo,
			wantCheck: &Check{
				OrganizationAddress: "390037,62,РЯЗАНСКАЯ ОБЛАСТЬ,РЯЗАНЬ Г,ЗУБКОВОЙ УЛ,СТРОЕНИЕ 25А",
				OrganizationName:    "Агроторг ООО",
				Products: []Product{
					{
						Name:     "*3220501 GREENF.Чай BARB.GARDEN  25х1,5г",
						Price:    5999,
						Quantity: 1,
						Sum:      5999,
					},
					{
						Name:     "3683200 ЛОЖКАРЕВ Котлеты кур.в пан.500г",
						Price:    13909,
						Quantity: 1,
						Sum:      13909,
					},
					{
						Name:     "3221397 NESC.Нап.МЯГКИЙ коф.3в1 16г",
						Price:    1399,
						Quantity: 1,
						Sum:      1399,
					},
				},
				Sum: 21307,
			},
			wantErr: false,
		},
		{
			name:   "invalid check DateTime",
			apiUrl: "https://proverkacheka.com/check/get",
			checkInfo: CheckInfoFromBot{
				DateTimeString: "0601T2122",
				Sum:            "213.07",
				FN:             "9282440300649733",
				FD:             "12416",
				FP:             "2858316733",
				N:              "1",
			},
			wantErr: true,
		},
		{
			name:      "api unavailable",
			apiUrl:    "https://invalidUrl/check/get",
			checkInfo: validCheckInfo,
			wantErr:   true,
		},
		{
			name:   "unexpected response code returned",
			apiUrl: server.URL,
			checkInfo: CheckInfoFromBot{
				DateTimeString: "20200601T2122",
				Sum:            "213.07",
				FN:             "9282440300649733",
				FD:             "",
				FP:             "2858316733",
				N:              "1",
			},
			wantErr: true,
		},
		{
			name:   "bad request string field FD",
			apiUrl: "https://proverkacheka.com/check/get",
			checkInfo: CheckInfoFromBot{
				DateTimeString: "20200601T2122",
				Sum:            "213.07",
				FN:             "9282440300649733",
				FD:             "may",
				FP:             "2858316733",
				N:              "1",
			},
			wantErr: true,
		},
		//{  ЕГО ПРОВЕРЯТЬ АНОНИМНОЙ СТРУКТУРОЙ
		//	name:   "invalid response",
		//	apiUrl: server.URL,
		//	checkInfo: validCheckInfo,
		//	wantErr: true,
		//},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			getFullChecker := NewTelegramGetFullChecker(tt.apiUrl, &http.Client{})
			actualCheck, err := getFullChecker.GetFullCheck(tt.checkInfo)
			if tt.wantErr {
				t.Log(err)
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantCheck, actualCheck)
			}
		})
	}

}
