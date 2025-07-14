package integration_tests

import (
	"encoding/json"
	"ghorkov32/proletariat-budget-be/integration_tests/utils"
	"ghorkov32/proletariat-budget-be/internal/core/domain"
	"ghorkov32/proletariat-budget-be/internal/core/port"
	"ghorkov32/proletariat-budget-be/openapi"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"time"
)

const accountResourceURL = "http://localhost:9091/accounts"

func (s *Suite) TestAccount() {
	s.T().Log("Starting TestAccount")
	s.Run("Create account successfully", func() {
		householdMember := s.createTestHouseholdMember()
		account := s.createTestAccountRequest(&householdMember, "150")

		apiResponse, err := s.makeAccountRequest(account)
		s.handleErr(err, "error while making request")
		defer apiResponse.Body.Close()

		var createdAccount openapi.Account
		s.decodeResponse(apiResponse, &createdAccount)

		// Verify that the account was created
		s.Equal(http.StatusCreated, apiResponse.StatusCode)
		s.Equal(account.AccountInformation, createdAccount.AccountInformation)
		s.Equal(account.AccountNumber, createdAccount.AccountNumber)
		s.Equal(account.Active, createdAccount.Active)
		s.Equal(account.Currency, createdAccount.Currency)
	})

	s.Run("Create account with inactive member", func() {
		householdMember := s.createTestHouseholdMember()
		_, errDeactivate := s.deactivateHouseholdMember(householdMember.Id)
		s.handleErr(errDeactivate, "error while deactivating member")
		account := s.createTestAccountRequest(&householdMember, "150")

		apiResponse, err := s.makeAccountRequest(account)
		s.handleErr(err, "error while making request")

		s.assertHttpError(apiResponse, http.StatusBadRequest, domain.ErrMemberInactive.Error())
	})

	s.Run("Create account and get 400 error to non existant member", func() {
		nonExistentMember := &openapi.HouseholdMember{
			Id:        "0",
			Active:    utils.BoolPtr(true),
			FirstName: "John",
			LastName:  "Doe",
			Nickname:  utils.StringPtr("Johnny"),
			Role:      "primary",
		}
		account := s.createTestAccountRequest(nonExistentMember, "150")

		s.testAccountCreationError(account, http.StatusBadRequest, domain.ErrMemberNotFound.Error())
	})

	s.Run("Create account and get 400 error to non existant currency", func() {
		householdMember := s.createTestHouseholdMember()
		account := s.createTestAccountRequest(&householdMember, "999")

		s.testAccountCreationError(account, http.StatusBadRequest, domain.ErrInvalidCurrency.Error())
	})

	s.Run("Delete an account with no transactions", func() {
		householdMember := s.createTestHouseholdMember()
		account := s.createTestAccountRequest(&householdMember, "150")

		apiResponse, err := s.makeAccountRequest(account)
		s.handleErr(err, "error while making request")
		defer apiResponse.Body.Close()

		var createdAccount openapi.Account
		s.decodeResponse(apiResponse, &createdAccount)

		apiResponse, err = s.deleteAccountRequest(createdAccount.Id)
		s.handleErr(err, "error while making request")

		s.Equal(http.StatusNoContent, apiResponse.StatusCode)
	})

	s.Run("Delete an account with transactions", func() {
		apiResponse, err := s.deleteAccountRequest("1")
		s.handleErr(err, "error while making request")

		s.assertHttpError(apiResponse, http.StatusConflict, domain.ErrAccountHasTransactions.Error())
	})

	s.Run("Check if account can be deleted OK", func() {
		apiResponse, err := s.checkCanDeleteAccount("6")
		s.handleErr(err, "error while making request")

		var response openapi.CanDelete
		s.decodeResponse(apiResponse, &response)
		s.Equal(http.StatusOK, apiResponse.StatusCode)
		s.Equal(true, response.CanDelete)
		s.Equal("", *response.Reason)
	})

	s.Run("Check if account can be deleted - false due having transactions", func() {
		apiResponse, err := s.checkCanDeleteAccount("1")
		s.handleErr(err, "error while making request")

		var response openapi.CanDelete
		s.decodeResponse(apiResponse, &response)
		s.Equal(http.StatusOK, apiResponse.StatusCode)
		s.Equal(false, response.CanDelete)
		s.Equal("Account has transactions", *response.Reason)
	})

	s.Run("Delete an account that does not exist", func() {
		apiResponse, err := s.deleteAccountRequest("0")
		s.handleErr(err, "error while making request")

		s.assertHttpError(apiResponse, http.StatusNotFound, domain.ErrAccountNotFound.Error())
	})

	s.Run("Deactivate an account", func() {

		householdMember := s.createTestHouseholdMember()
		account := s.createTestAccountRequest(&householdMember, "150")

		apiResponse, err := s.makeAccountRequest(account)
		s.handleErr(err, "error while making request")
		defer apiResponse.Body.Close()

		var createdAccount openapi.Account
		s.decodeResponse(apiResponse, &createdAccount)

		apiResponse, err = s.deactivateAccountRequest(createdAccount.Id)
		s.handleErr(err, "error while making request")

		s.Equal(http.StatusNoContent, apiResponse.StatusCode)
	})

	s.Run("Deactivate an already deactivated account", func() {

		householdMember := s.createTestHouseholdMember()
		account := s.createTestAccountRequest(&householdMember, "150")

		apiResponse, err := s.makeAccountRequest(account)
		s.handleErr(err, "error while making request")
		defer apiResponse.Body.Close()

		var createdAccount openapi.Account
		s.decodeResponse(apiResponse, &createdAccount)

		apiResponse, err = s.deactivateAccountRequest(createdAccount.Id)
		s.handleErr(err, "error while making request")

		apiResponse, err = s.deactivateAccountRequest(createdAccount.Id)
		s.handleErr(err, "error while making request")

		s.assertHttpError(apiResponse, http.StatusBadRequest, domain.ErrAccountAlreadyInactive.Error())
	})

	s.Run("Deactivate an account that does not exist", func() {
		apiResponse, err := s.deactivateAccountRequest("0")
		s.handleErr(err, "error while making request")

		s.assertHttpError(apiResponse, http.StatusNotFound, domain.ErrAccountNotFound.Error())
	})

	s.Run("Reactivate an account", func() {

		householdMember := s.createTestHouseholdMember()
		account := s.createTestAccountRequest(&householdMember, "150")

		apiResponse, err := s.makeAccountRequest(account)
		s.handleErr(err, "error while making request")
		defer apiResponse.Body.Close()

		var createdAccount openapi.Account
		s.decodeResponse(apiResponse, &createdAccount)

		apiResponse, err = s.deactivateAccountRequest(createdAccount.Id)
		s.handleErr(err, "error while making request")

		apiResponse, err = s.activateAccountRequest(createdAccount.Id)
		s.handleErr(err, "error while making request")

		s.Equal(http.StatusOK, apiResponse.StatusCode)
	})

	s.Run("Activate an already active account", func() {

		householdMember := s.createTestHouseholdMember()
		account := s.createTestAccountRequest(&householdMember, "150")

		apiResponse, err := s.makeAccountRequest(account)
		s.handleErr(err, "error while making request")
		defer apiResponse.Body.Close()

		var createdAccount openapi.Account
		s.decodeResponse(apiResponse, &createdAccount)

		apiResponse, err = s.activateAccountRequest(createdAccount.Id)
		s.handleErr(err, "error while making request")

		s.assertHttpError(apiResponse, http.StatusBadRequest, domain.ErrAccountAlreadyActive.Error())
	})

	s.Run("Reactivate an account that does not exist", func() {
		apiResponse, err := s.activateAccountRequest("0")
		s.handleErr(err, "error while making request")

		s.assertHttpError(apiResponse, http.StatusNotFound, domain.ErrAccountNotFound.Error())
	})

	s.Run("Get account by ID", func() {
		householdMember := s.createTestHouseholdMember()
		account := s.createTestAccountRequest(&householdMember, "150")

		apiResponse, err := s.makeAccountRequest(account)
		s.handleErr(err, "error while making request")
		defer apiResponse.Body.Close()

		var createdAccount openapi.Account
		s.decodeResponse(apiResponse, &createdAccount)

		apiResponse, err = s.getAccountRequest(createdAccount.Id)
		s.handleErr(err, "error while making request")

		var retrievedAccount openapi.Account
		s.decodeResponse(apiResponse, &retrievedAccount)
		s.handleErr(err, "error while decoding response body")
		s.Equal(http.StatusOK, apiResponse.StatusCode)
		s.Equal(createdAccount, retrievedAccount)
	})

	s.Run("Get account by ID that does not exist", func() {
		apiResponse, err := s.getAccountRequest("0")
		s.handleErr(err, "error while making request")

		s.assertHttpError(apiResponse, http.StatusNotFound, domain.ErrAccountNotFound.Error())
	})

	s.Run("Update account", func() {
		householdMember := s.createTestHouseholdMember()
		accountRequest := s.createTestAccountRequest(&householdMember, "150")

		apiResponse, err := s.makeAccountRequest(accountRequest)
		s.handleErr(err, "error while making request")
		defer apiResponse.Body.Close()

		var createdAccount openapi.Account
		s.decodeResponse(apiResponse, &createdAccount)

		updatedAccount := createdAccount

		updatedAccount.Description = utils.StringPtr("Updated description")
		updatedAccount.Owner = &householdMember

		apiResponse, err = s.updateAccountRequest(&updatedAccount)
		s.handleErr(err, "error while making request")

		var retrievedAccount openapi.Account
		s.decodeResponse(apiResponse, &retrievedAccount)

		s.Equal(http.StatusOK, apiResponse.StatusCode)
		s.assertAccountsEquals(&updatedAccount, &retrievedAccount)
	})

	s.Run("Update a non existant account", func() {

		owner := &openapi.HouseholdMember{
			Id:        "0",
			FirstName: "John",
			LastName:  "Doe",
			Nickname:  utils.StringPtr("Johnny"),
			Role:      "owner",
		}

		updatedAccount := openapi.Account{
			AccountInformation: utils.StringPtr("Test Account Information"),
			AccountNumber:      utils.StringPtr("1234567890"),
			Active:             utils.BoolPtr(true),
			Currency:           "150",
			Description:        utils.StringPtr("Test savings account for integration testing"),
			InitialBalance:     1000.50,
			Institution:        utils.StringPtr("Test Bank of America"),
			Name:               "Test Savings Account",
			Owner:              owner,
			Type:               "bank",
			Id:                 "0",
		}
		updatedAccount.Description = utils.StringPtr("Updated description")
		updatedAccount.Owner = &openapi.HouseholdMember{
			Id: "0",
		}
		apiResponse, err := s.updateAccountRequest(&updatedAccount)
		s.handleErr(err, "error while making request")
		s.assertHttpError(apiResponse, http.StatusNotFound, domain.ErrAccountNotFound.Error())

	})

	s.Run("Update account and get 400 error to non existant member", func() {
		householdMember := s.createTestHouseholdMember()
		account := s.createTestAccountRequest(&householdMember, "150")

		apiResponse, err := s.makeAccountRequest(account)
		s.handleErr(err, "error while making request")
		defer apiResponse.Body.Close()

		var createdAccount openapi.Account
		s.decodeResponse(apiResponse, &createdAccount)

		updatedAccount := createdAccount

		updatedAccount.Owner = &openapi.HouseholdMember{
			Id: "0",
		}
		apiResponse, err = s.updateAccountRequest(&updatedAccount)
		s.handleErr(err, "error while making request")
		s.assertHttpError(apiResponse, http.StatusBadRequest, domain.ErrMemberNotFound.Error())

	})

	s.Run("Update account and get 400 error to non existant currency", func() {
		householdMember := s.createTestHouseholdMember()
		account := s.createTestAccountRequest(&householdMember, "150")

		apiResponse, err := s.makeAccountRequest(account)
		s.handleErr(err, "error while making request")
		defer apiResponse.Body.Close()

		var createdAccount openapi.Account
		s.decodeResponse(apiResponse, &createdAccount)

		updatedAccount := createdAccount

		updatedAccount.Currency = "999"
		apiResponse, err = s.updateAccountRequest(&updatedAccount)
		s.handleErr(err, "error while making request")
		s.assertHttpError(apiResponse, http.StatusBadRequest, domain.ErrInvalidCurrency.Error())
	})

	s.Run("Update account with invalid data", func() {
		householdMember := s.createTestHouseholdMember()
		account := s.createTestAccountRequest(&householdMember, "150")

		apiResponse, err := s.makeAccountRequest(account)
		s.handleErr(err, "error while making request")
		defer apiResponse.Body.Close()

		var createdAccount openapi.Account
		s.decodeResponse(apiResponse, &createdAccount)

		updatedAccount := createdAccount

		updatedAccount.Currency = "ASD"
		apiResponse, err = s.updateAccountRequest(&updatedAccount)
		s.handleErr(err, "error while making request")
		s.assertHttpError(apiResponse, http.StatusBadRequest, port.ErrInvalidDataFormat.Error())
	})

	s.Run("List accounts with no filter", func() {
		apiResponse := s.listAccountRequest(openapi.ListAccountsParams{
			Type:     nil,
			Currency: nil,
			Active:   nil,
			Limit:    utils.IntPtr(5),
			Offset:   utils.IntPtr(0),
		})

		var accounts openapi.AccountList
		s.decodeResponse(apiResponse, &accounts)

		s.Equal(http.StatusOK, apiResponse.StatusCode)
		s.Equal(len(*accounts.Accounts), 5)
	})

	s.Run("List accounts with pagination", func() {
		apiResponse := s.listAccountRequest(openapi.ListAccountsParams{
			Type:     nil,
			Currency: nil,
			Active:   nil,
			Limit:    utils.IntPtr(2),
			Offset:   utils.IntPtr(5),
		})

		var accounts openapi.AccountList
		s.decodeResponse(apiResponse, &accounts)

		s.Equal(http.StatusOK, apiResponse.StatusCode)
		s.Equal(len(*accounts.Accounts), 2)
	})

	s.Run("List accounts with currency filter", func() {
		filters := openapi.ListAccountsParams{
			Type:     nil,
			Currency: utils.StringPtr("150"),
			Active:   nil,
			Limit:    utils.IntPtr(5),
			Offset:   utils.IntPtr(0),
		}
		apiResponse := s.listAccountRequest(filters)

		var accounts openapi.AccountList
		s.decodeResponse(apiResponse, &accounts)

		s.Equal(http.StatusOK, apiResponse.StatusCode)
		s.Equal(len(*accounts.Accounts), 5)
		for _, account := range *accounts.Accounts {
			s.Equal(*filters.Currency, account.Currency)
		}
	})

	s.Run("List accounts with account type filter", func() {

		filters := openapi.ListAccountsParams{
			Type:     utils.StringPtr("bank"),
			Currency: nil,
			Active:   nil,
			Limit:    utils.IntPtr(5),
			Offset:   utils.IntPtr(0),
		}
		apiResponse := s.listAccountRequest(filters)

		var accounts openapi.AccountList
		s.decodeResponse(apiResponse, &accounts)

		s.Equal(http.StatusOK, apiResponse.StatusCode)
		s.Equal(len(*accounts.Accounts), 5)
		for _, account := range *accounts.Accounts {
			s.Equal(openapi.AccountTypeBank, account.Type)
		}
	})

	s.Run("Empty account list", func() {
		filters := openapi.ListAccountsParams{
			Type:     nil,
			Currency: utils.StringPtr("151"),
			Active:   nil,
			Limit:    utils.IntPtr(5),
			Offset:   utils.IntPtr(0),
		}
		apiResponse := s.listAccountRequest(filters)

		var accounts openapi.AccountList
		s.decodeResponse(apiResponse, &accounts)

		s.Equal(http.StatusOK, apiResponse.StatusCode)
		s.Equal(len(*accounts.Accounts), 0)
	})

	s.Run("List active accounts", func() {

		filters := openapi.ListAccountsParams{
			Type:     nil,
			Currency: nil,
			Active:   utils.BoolPtr(true),
			Limit:    utils.IntPtr(5),
			Offset:   utils.IntPtr(0),
		}
		apiResponse := s.listAccountRequest(filters)

		var accounts openapi.AccountList
		s.decodeResponse(apiResponse, &accounts)

		s.Equal(http.StatusOK, apiResponse.StatusCode)
		s.Equal(len(*accounts.Accounts), 5)
		for _, account := range *accounts.Accounts {
			s.Equal(*filters.Active, *account.Active)
		}
	})

	s.Run("List inactive accounts", func() {

		filters := openapi.ListAccountsParams{
			Type:     nil,
			Currency: nil,
			Active:   utils.BoolPtr(false),
			Limit:    utils.IntPtr(5),
			Offset:   utils.IntPtr(0),
		}
		apiResponse := s.listAccountRequest(filters)

		var accounts openapi.AccountList
		s.decodeResponse(apiResponse, &accounts)

		s.Equal(http.StatusOK, apiResponse.StatusCode)
		s.Equal(2, len(*accounts.Accounts))
		for _, account := range *accounts.Accounts {
			s.Equal(*filters.Active, *account.Active)
		}
	})

}

func (s *Suite) assertAccountsEquals(expected, actual *openapi.Account) {
	s.Equal(expected.AccountNumber, actual.AccountNumber)
	s.Equal(expected.Active, actual.Active)
	s.Equal(expected.Currency, actual.Currency)
	s.Equal(expected.Description, actual.Description)
	s.Equal(expected.InitialBalance, actual.InitialBalance)
	s.Equal(expected.Institution, actual.Institution)
	s.Equal(expected.Name, actual.Name)
	s.Equal(expected.Owner.Id, actual.Owner.Id)
	s.Equal(expected.Type, actual.Type)
	s.Equal(expected.Id, actual.Id)
	s.Equal(expected.CreatedAt, actual.CreatedAt)
	s.WithinDuration(expected.UpdatedAt, actual.UpdatedAt, time.Second*1)
}

// Helper for error test cases
func (s *Suite) testAccountCreationError(account *openapi.AccountRequest, expectedStatus int, expectedMessage string) {
	apiResponse, err := s.makeAccountRequest(account)
	s.handleErr(err, "error while making request")

	s.assertHttpError(apiResponse, expectedStatus, expectedMessage)
}

func (s *Suite) assertHttpError(response *http.Response, expectedStatus int, expectedMessage string) {
	var errorResponse openapi.Error
	s.decodeResponse(response, &errorResponse)
	s.Equal(expectedStatus, response.StatusCode)
	s.Equal(expectedMessage, errorResponse.Message)
}

// Helper method to make account creation requests
func (s *Suite) makeAccountRequest(account *openapi.AccountRequest) (*http.Response, error) {
	body, err := utils.PrepareRequestBody(account)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(
		s.ctx,
		http.MethodPost,
		accountResourceURL,
		body,
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	return client.Do(req)
}
func (s *Suite) updateAccountRequest(account *openapi.Account) (*http.Response, error) {
	body, err := utils.PrepareRequestBody(*account)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(
		s.ctx,
		http.MethodPut,
		accountResourceURL+"/"+url.PathEscape(account.Id),
		body,
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	return client.Do(req)
}
func (s *Suite) deleteAccountRequest(id string) (*http.Response, error) {

	req, err := http.NewRequestWithContext(
		s.ctx,
		http.MethodDelete,
		accountResourceURL+"/"+url.PathEscape(id),
		http.NoBody,
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	return client.Do(req)
}

func (s *Suite) checkCanDeleteAccount(id string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(
		s.ctx,
		http.MethodGet,
		accountResourceURL+"/"+path.Join(url.PathEscape(id), "can-delete"),
		http.NoBody,
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	return client.Do(req)
}

func (s *Suite) deactivateAccountRequest(id string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(
		s.ctx,
		http.MethodPatch,
		accountResourceURL+"/"+path.Join(url.PathEscape(id), "deactivate"),
		http.NoBody,
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	return client.Do(req)
}

func (s *Suite) activateAccountRequest(id string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(
		s.ctx,
		http.MethodPatch,
		accountResourceURL+"/"+path.Join(url.PathEscape(id), "activate"),
		http.NoBody,
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	return client.Do(req)
}

func (s *Suite) getAccountRequest(id string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(
		s.ctx,
		http.MethodGet,
		accountResourceURL+"/"+url.PathEscape(id),
		nil,
	)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	apiResponse, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	return apiResponse, nil
}

func (s *Suite) listAccountRequest(params openapi.ListAccountsParams) *http.Response {

	req, err := http.NewRequestWithContext(
		s.ctx,
		http.MethodGet,
		accountResourceURL,
		nil,
	)
	s.handleErr(err, "error while creating request")

	q := req.URL.Query()
	if params.Limit != nil {
		q.Set("limit", strconv.Itoa(*params.Limit))
	}
	if params.Offset != nil {
		q.Set("offset", strconv.Itoa(*params.Offset))
	}
	if params.Currency != nil {
		q.Set("currency", *params.Currency)
	}
	if params.Type != nil {
		q.Set("sort_by", *params.Type)
	}
	if params.Active != nil {
		q.Set("active", strconv.FormatBool(*params.Active))
	}
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	apiResponse, err := client.Do(req)
	s.handleErr(err, "error while making request")

	return apiResponse
}

// Helper method to create a test household member
func (s *Suite) createTestHouseholdMember() openapi.HouseholdMember {
	householdMemberReq := &openapi.HouseholdMemberRequest{
		Active:    utils.BoolPtr(true),
		FirstName: "John",
		LastName:  "Doe",
		Nickname:  utils.StringPtr("Johnny"),
		Role:      "primary",
	}

	householdMember, err := s.createHouseholdMember(householdMemberReq)
	s.handleErr(err, "error while creating household member")
	return householdMember
}

// Helper method to create a test account request
func (s *Suite) createTestAccountRequest(owner *openapi.HouseholdMember, currency string) *openapi.AccountRequest {
	return &openapi.AccountRequest{
		AccountInformation: utils.StringPtr("Test Account Information"),
		AccountNumber:      utils.StringPtr("1234567890"),
		Active:             utils.BoolPtr(true),
		Currency:           currency,
		Description:        utils.StringPtr("Test savings account for integration testing"),
		InitialBalance:     1000.50,
		Institution:        utils.StringPtr("Test Bank of America"),
		Name:               "Test Savings Account",
		Owner:              owner,
		Type:               "bank",
	}
}

// Generic helper to decode response body
func (s *Suite) decodeResponse(response *http.Response, target interface{}) {
	defer response.Body.Close()
	err := json.NewDecoder(response.Body).Decode(target)
	s.handleErr(err, "error while decoding response body")
}
