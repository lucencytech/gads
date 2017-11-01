package main

import (
	"fmt"
	//"reflect"

	"io/ioutil"
	"net/http"

	"os"
	"encoding/csv"
	gads "github.com/Getsidecar/gads/v201705"
	"github.com/Getsidecar/sidecar-go-utils/config"

	"strings"
)

func getReport(auth *gads.Auth, headers []string) (collection []map[string]string) {
	fmt.Println("getting report with auth:", auth)
	rds := gads.NewReportDownloadService(auth)

	rd := gads.ReportDefinition{
		ReportName: "ReportNameGoesHere",
		ReportType: "SHOPPING_PERFORMANCE_REPORT",
		DateRangeType: "YESTERDAY",
		DownloadFormat: "CSV",
		Selector: gads.Selector{
			Fields: headers,
			// Predicates: []gads.Predicate{
			// 	{"Status", "EQUALS", []string{"ENABLED"}},
			// 	{"AdvertisingChannelType", "EQUALS", []string{"SHOPPING"}},
			// },
			// Paging: &paging,
		},
	}

	collection, _ = rds.Get(rd)
	return collection
}

func getAWQLResult(auth *gads.Auth, query string) ([]map[string]string) {
	rds := gads.NewReportDownloadService(auth)
	report, err := rds.AWQL(query, "CSV")
	if err != nil {
		fmt.Println("Error in AWQL Query: ", err)
		return nil
	}

	return report
}

func writeReportToCsv(filename string, report []map[string]string) {
	file, _ := os.Create(filename)
	writer := csv.NewWriter(file)
	defer writer.Flush()

	var returnHeaders []string
	for _, value := range report[0:1] {
		for key, _ := range value {
			returnHeaders = append(returnHeaders, key)
		}

		writer.Write(returnHeaders)
	}

	for _, line := range report {
		var lineList []string
		for _, header := range returnHeaders {
			lineList = append(lineList, line[header])
		}

		writer.Write(lineList)
	}
}

func writeFieldExclusionsToCsv(filename string, reportName string, auth *gads.Auth) {
	rds := gads.NewReportDefinitionService(auth)
	reportFields, _ := rds.GetReportFields(reportName)
	// fmt.Println("reportFields:", reportFields)
	var fieldExclusions  []map[string]string
	for _, field := range reportFields {
		fieldExclusion := make(map[string]string)
		fieldExclusion["fieldName"] = field.FieldName
		fieldExclusion["fieldExclusions"] = strings.Join(field.ExclusiveFields, ";")
		fieldExclusions = append(fieldExclusions, fieldExclusion)
	}
	fmt.Println("fieldExclusions:", fieldExclusions)

	writeReportToCsv(filename, fieldExclusions)
}

func writeQueryReportToCsv(awqlFilename string, csvFilename string, auth *gads.Auth) {
	queryBytes, _ := ioutil.ReadFile(awqlFilename)
	report := getAWQLResult(auth, string(queryBytes))
	writeReportToCsv(csvFilename, report)
}

func main() {
	authConfig, err := gads.NewCredentialsFromFile("config.json")

	if err != nil {
		panic(err)
	}

	c := &http.Client{}
	configClient := config.ConfigStoreClient{
		HttpClient: c,
		ReadAll:    ioutil.ReadAll,
		BaseUrl:    "https://config.sidecartechnologies.com:4000",
		Username:   "root",
		Password:   "tkw2yWejYMqXm9y",
	}

	clientConfigs, err := configClient.GetClients()

	if err != nil {
		panic(err)
	}
	// f, _ := os.Create("test.csv")
	// w := csv.NewWriter(f)
	// defer w.Flush()
	// w.Comma = '\t'
	// defer f.Close()
	for _, client := range clientConfigs {
		if client.Status != "active" {
			//fmt.Printf("Skipping %s due to inactive flag...\n", client.Shortname)
			continue
		}
		if client.Shortname != "moosejaw" {
			//fmt.Printf("Skipping %s to focus on moosejaw...\n", client.Shortname)
			continue
		}
		fmt.Printf("Running %s...\n", client.Shortname)
		authConfig.Auth.CustomerId = client.Accounts.Adwords.AccountId

		// headers := []string{
		// 	"AccountDescriptiveName",
		// 	"CampaignId",
		// 	"AdGroupId",
		// 	"Cost",
		// 	"Clicks",
		// 	"Impressions",
		// 	"Conversions",
		// 	"ConversionValue",
		// 	"OfferId",
		// 	"ExternalCustomerId",
		// 	"Date",
		// 	"AdGroupName",
		// 	"Device",
		// }
		//
		// // For using Report Download Service
		// report := getReport(&authConfig.Auth, headers)
		// writeReportToCsv("result.csv", report)


		//For using AWQL
		writeQueryReportToCsv("report1.awql", "report1.csv", &authConfig.Auth)
		writeQueryReportToCsv("report2.awql", "report2.csv", &authConfig.Auth)
		writeQueryReportToCsv("report3.awql", "report3.csv", &authConfig.Auth)


		// writeFieldExclusionsToCsv("field-exclusions.CAMPAIGN_PERFORMANCE_REPORT.csv", "CAMPAIGN_PERFORMANCE_REPORT", &authConfig.Auth)


	}
	// w.Flush()
}
