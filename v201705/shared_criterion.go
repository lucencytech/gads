package v201705

import (
	"encoding/xml"
	"fmt"
)

type SharedCriterionService struct {
	Auth
}

func NewSharedCriterionService(auth *Auth) *SharedCriterionService {
	return &SharedCriterionService{Auth: *auth}
}

type SharedCriterion struct {
	Id        int64     `xml:"sharedSetId,omitempty"`
	Negative  bool      `xml:"negative,omitempty"`
	Criterion Criterion `xml:"criterion,omitempty"`
}

func (s SharedCriterionService) Get(selector Selector) (sharedCriteria []SharedCriterion, totalCount int64, err error) {
	selector.XMLName = xml.Name{baseUrl, "selector"}
	respBody, err := s.Auth.request(
		sharedCriterionServiceUrl,
		"get",
		struct {
			XMLName xml.Name
			Sel     Selector
		}{
			XMLName: xml.Name{
				Space: baseUrl,
				Local: "get",
			},
			Sel: selector,
		},
	)
	if err != nil {
		return sharedCriteria, totalCount, err
	}
	getResp := struct {
		Size           int64             `xml:"rval>totalNumEntries"`
		SharedCriteria []SharedCriterion `xml:"rval>entries"`
	}{}
	err = xml.Unmarshal([]byte(respBody), &getResp)
	if err != nil {
		return sharedCriteria, totalCount, err
	}
	return getResp.SharedCriteria, getResp.Size, err
}

func (s *SharedCriterion) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	for token, err := dec.Token(); err == nil; token, err = dec.Token() {
		if err != nil {
			return err
		}
		switch start := token.(type) {
		case xml.StartElement:
			tag := start.Name.Local
			switch tag {
			case "sharedSetId":
				if err := dec.DecodeElement(&s.Id, &start); err != nil {
					return err
				}
			case "negative":
				if err := dec.DecodeElement(&s.Negative, &start); err != nil {
					return err
				}
			case "criterion":
				criterion, err := criterionUnmarshalXML(dec, start)
				if err != nil {
					return err
				}
				s.Criterion = criterion
			default:
				return fmt.Errorf("unknown BiddableAdGroupCriterion field %s", tag)
			}
		}
	}
	return nil
}
