package test

import (
	"io/ioutil"
	"strings"

	"github.com/google/logger"
)

// QueriesTest - Struct with json query's GraphQL to request
type QueriesTest struct {
	Query []string
}

// MountRequests is responsible to prepare son's to request
func MountRequests(limits string) (querys QueriesTest) {
	limitsArray := strings.Split(limits, ",")
	jsonBytes, err := ioutil.ReadFile("./resources/domains/*.json")
	if err != nil {
		logger.Fatalf("error to load domains: ", err)
	}

	querys = QueriesTest{}
	for i, _ := range limitsArray {
		jsonReplaced := strings.Replace(string(jsonBytes), "$lines", limitsArray[i], 1)
		querys.Query = append(querys.Query, jsonReplaced)
	}

	return querys
}
