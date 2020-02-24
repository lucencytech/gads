package v201809

import (
	"encoding/xml"
)

type OfflineCallConversionService struct {
	Auth
}

const (
	uploadCallConversionAction = "ADD"
)

type OfflineCallConversionFeed struct {
	CallerID               string  `xml:"callerId"`
	CallStartTime          string  `xml:"callStartTime"`
	ConversionName         string  `xml:"conversionName"`
	ConversionTime         string  `xml:"conversionTime,omitempty"`
	ConversionValue        float64 `xml:"conversionValue,omitempty"`
	ConversionCurrencyCode string  `xml:"conversionCurrencyCode,omitempty"`
}

type OfflineCallConversionOperations map[string][]OfflineCallConversionFeed

func NewOfflineCallConversionService(auth *Auth) *OfflineCallConversionService {
	return &OfflineCallConversionService{Auth: *auth}
}

func (o *OfflineCallConversionService) Mutate(conversionOperations OfflineCallConversionOperations) (conversion []OfflineCallConversionFeed, error error) {
	type callConversionOperation struct {
		Action     string                    `xml:"operator"`
		Conversion OfflineCallConversionFeed `xml:"operand"`
	}
	var ops []callConversionOperation

	for _, conversion := range conversionOperations {
		for _, con := range conversion {
			ops = append(ops,
				callConversionOperation{
					Action:     uploadCallConversionAction,
					Conversion: con,
				},
			)
		}
	}
	mutation := struct {
		XMLName xml.Name
		Ops     []callConversionOperation `xml:"operations"`
	}{
		XMLName: xml.Name{
			Space: baseUrl,
			Local: "mutate",
		},
		Ops: ops}
	respBody, err := o.Auth.request(offlineCallConversionFeedServiceUrl, "mutate", mutation)
	if err != nil {
		return conversion, err
	}
	mutateResp := struct {
		Conversions []OfflineCallConversionFeed `xml:"rval>value"`
	}{}
	err = xml.Unmarshal([]byte(respBody), &mutateResp)
	if err != nil {
		return conversion, err
	}
	return mutateResp.Conversions, nil
}
