package v201809

// https://developers.google.com/adwords/api/docs/reference/v201809/AdGroupExtensionSettingService.Page
// Contains the results from a get call.
type Page struct {
	TotalNumEntries int    `xml:"https://adwords.google.com/api/adwords/cm/v201809 totalNumEntries,omitempty"`
	PageType        string `xml:"https://adwords.google.com/api/adwords/cm/v201809 Page.Type,omitempty"`
}
