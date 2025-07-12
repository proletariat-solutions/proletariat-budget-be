package integration_tests

import (
	"encoding/json"
	"ghorkov32/proletariat-budget-be/integration_tests/utils"
	"ghorkov32/proletariat-budget-be/openapi"
	"net/http"
)

func (s *Suite) createHouseholdMember(householdMemberReq *openapi.HouseholdMemberRequest) (openapi.HouseholdMember, error) {
	body, errBodyPrepare := utils.PrepareRequestBody(householdMemberReq)

	s.handleErr(errBodyPrepare, "error while preparing request body")

	req, errReq := http.NewRequestWithContext(
		s.ctx,
		http.MethodPost,
		"http://localhost:9091/household-members",
		body,
	)
	s.handleErr(errReq, "error while creating request")

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	apiResponse, err := client.Do(req)
	s.handleErr(err, "error while making request")

	defer apiResponse.Body.Close()

	var bodyBytes []byte

	_, errRead := apiResponse.Body.Read(bodyBytes)
	s.handleErr(errRead, "error while reading response body")

	var member openapi.HouseholdMember
	errDecode := json.NewDecoder(apiResponse.Body).Decode(&member)
	s.handleErr(errDecode, "error while decoding response body")

	return member, nil
}
