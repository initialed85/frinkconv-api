package frinkconv_server

import (
	"encoding/json"
	"fmt"
	"log"
)

func getConvertRequestFromJSON(requestJSON []byte) (ConvertRequest, error) {
	convertRequest := ConvertRequest{}

	log.Printf("<<< %#+v", string(requestJSON))

	err := json.Unmarshal(requestJSON, &convertRequest)
	if err != nil {
		return ConvertRequest{}, err
	}

	log.Printf("<<< %#+v", convertRequest)

	return convertRequest, nil
}

func getConvertErrorResponseJSON(givenErr error) []byte {
	response := ConvertErrorResponse{
		Error: givenErr.Error(),
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		return []byte(fmt.Sprintf("{\"error\": %#+v}", err.Error()))
	}

	log.Printf(">>> %#+v", string(responseJSON))

	return responseJSON
}

func getConvertSuccessResponseJSON(destinationValue float64) []byte {
	response := ConvertSuccessResponse{
		DestinationValue: destinationValue,
	}

	log.Printf(">>> %#+v", response)

	responseJSON, err := json.Marshal(response)
	if err != nil {
		return []byte(fmt.Sprintf("{\"error\": %#+v}", err.Error()))
	}

	log.Printf(">>> %#+v", string(responseJSON))

	return responseJSON
}

func getConvertRequestsFromJSON(requestJSON []byte) ([]ConvertRequest, error) {
	convertRequests := make([]ConvertRequest, 0)

	log.Printf("<<< %#+v", string(requestJSON))

	err := json.Unmarshal(requestJSON, &convertRequests)
	if err != nil {
		return []ConvertRequest{}, err
	}

	log.Printf("<<< %#+v", convertRequests)

	return convertRequests, nil
}

func getConvertResponsesJSON(response []interface{}) []byte {
	log.Printf(">>> %#+v", response)

	responseJSON, err := json.Marshal(response)
	if err != nil {
		return []byte(fmt.Sprintf("{\"error\": %#+v}", err.Error()))
	}

	log.Printf(">>> %#+v", string(responseJSON))

	return responseJSON
}
