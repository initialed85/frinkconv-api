package frinkconv_server

type ConvertRequest struct {
	SourceValue      float64 `json:"source_value"`
	SourceUnits      string  `json:"source_units"`
	DestinationUnits string  `json:"destination_units"`
}

type ConvertErrorResponse struct {
	Error string `json:"error,omitempty"`
}

type ConvertSuccessResponse struct {
	DestinationValue float64 `json:"destination_value"`
}
