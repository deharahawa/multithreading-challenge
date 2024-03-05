package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type BrasilAPICEPResponse struct {
	CEP          string `json:"cep"`
	City         string `json:"city"`
	Message      string `json:"message"`
	Neighborhood string `json:"neighborhood"`
	Service      string `json:"service"`
	State        string `json:"state"`
	Street       string `json:"street"`
}

type ViaCEPResponse struct {
	Bairro      string `json:"bairro"`
	CEP         string `json:"cep"`
	Complemento string `json:"complemento"`
	DDD         string `json:"ddd"`
	GIA         string `json:"gia"`
	IBGE        string `json:"ibge"`
	Localidade  string `json:"localidade"`
	Logradouro  string `json:"logradouro"`
	SIAFI       string `json:"siafi"`
	UF          string `json:"uf"`
}

func callAPI(url string, response interface{}) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	json.NewDecoder(res.Body).Decode(&response)
	return nil
}

func callAPIAndSendResponse(url string, response interface{}, messageFormat string, apiRes chan<- string) error {
	err := callAPI(url, response)
	if err != nil {
		return err
	}

	switch r := response.(type) {
	case *BrasilAPICEPResponse:
		if r.Message != "" {
			errMsg := fmt.Sprintf("brasil API error: %s", r.Message)
			return fmt.Errorf(errMsg)
		}
		apiRes <- fmt.Sprintf(messageFormat, r.Street)
	case *ViaCEPResponse:
		if r.Logradouro == "" {
			errMsg := "via CEP API error: CEP not found"
			return fmt.Errorf(errMsg)
		}
		apiRes <- fmt.Sprintf(messageFormat, r.Logradouro)
	}
	return nil
}

func CallBrasilAPI(cep string, apiRes chan<- string) error {
	brasilAPIURL := fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cep)
	var response BrasilAPICEPResponse
	return callAPIAndSendResponse(brasilAPIURL, &response, "Received first from BrasilAPI: %s", apiRes)
}

func CallViaCEPAPI(cep string, apiRes chan<- string) error {
	viaCEPURL := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep)
	var response ViaCEPResponse
	return callAPIAndSendResponse(viaCEPURL, &response, "Received first from ViaCEP: %s", apiRes)
}
