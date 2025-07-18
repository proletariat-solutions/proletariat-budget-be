package integration_tests

import (
	"ghorkov32/proletariat-budget-be/integration_tests/utils"
	"ghorkov32/proletariat-budget-be/internal/core/domain"
	"ghorkov32/proletariat-budget-be/openapi"
	"net/http"
)

func (s *Suite) TestExpenditures() {
	s.T().Log("Starting TestExpenditures")

	// Setup test data
	testMember := s.createTestHouseholdMember()
	testAccount := s.createTestAccount(
		&testMember,
		"150",
	)
	testCategory := s.createTestCategory(openapi.CategoryTypeExpenditure)

	s.Run(
		"Create expenditure successfully",
		func() {
			expenditureReq := s.createTestExpenditureRequest(
				&testAccount.Id,
				&testCategory,
			)

			apiResponse, err := s.createExpenditureRequest(expenditureReq)
			s.handleErr(
				err,
				"error while making request",
			)

			s.Equal(
				http.StatusCreated,
				apiResponse.StatusCode,
			)

			var expenditure openapi.Expenditure
			s.decodeResponse(
				apiResponse,
				&expenditure,
			)

			// Verify the created expenditure
			s.NotEmpty(expenditure.Id)
			s.Equal(
				expenditureReq.Amount,
				expenditure.Amount,
			)
			s.Equal(
				expenditureReq.Currency,
				expenditure.Currency,
			)
			s.Equal(
				expenditureReq.Description,
				expenditure.Description,
			)
			s.Equal(
				expenditureReq.AccountId,
				expenditure.AccountId,
			)
			s.Equal(
				expenditureReq.Category.Id,
				expenditure.Category.Id,
			)
			s.Equal(
				*expenditureReq.Declared,
				*expenditure.Declared,
			)
			s.Equal(
				*expenditureReq.Planned,
				*expenditure.Planned,
			)
			s.NotNil(expenditure.CreatedAt)
			s.NotNil(expenditure.UpdatedAt)
		},
	)

	s.Run(
		"Create expenditure with non-existent account",
		func() {
			expenditureReq := s.createTestExpenditureRequest(
				utils.StringPtr("999999"),
				&testCategory,
			)

			apiResponse, err := s.createExpenditureRequest(expenditureReq)
			s.handleErr(
				err,
				"error while making request",
			)

			s.assertHttpError(
				apiResponse,
				http.StatusBadRequest,
				domain.ErrAccountNotFound.Error(),
			)
		},
	)

	s.Run(
		"Create expenditure with non-existent category",
		func() {
			expenditureReq := s.createTestExpenditureRequest(
				&testAccount.Id,
				&testCategory,
			)
			expenditureReq.Category.Id = "999999"

			apiResponse, err := s.createExpenditureRequest(expenditureReq)
			s.handleErr(
				err,
				"error while making request",
			)

			s.assertHttpError(
				apiResponse,
				http.StatusBadRequest,
				domain.ErrCategoryNotFound.Error(),
			)
		},
	)

	s.Run(
		"Create expenditure with inactive account",
		func() {
			// Create an inactive account
			inactiveAccount := s.createTestAccount(
				&testMember,
				"150",
			)
			_, err := s.deactivateAccountRequest(inactiveAccount.Id)
			s.handleErr(
				err,
				"error while deactivating account",
			)

			expenditureReq := s.createTestExpenditureRequest(
				&inactiveAccount.Id,
				&testCategory,
			)

			apiResponse, err := s.createExpenditureRequest(expenditureReq)
			s.handleErr(
				err,
				"error while making request",
			)

			s.assertHttpError(
				apiResponse,
				http.StatusConflict,
				domain.ErrAccountInactive.Error(),
			)
		},
	)

	s.Run(
		"Create expenditure with inactive category",
		func() {
			// Create an inactive category
			inactiveCategory := s.createTestCategory(openapi.CategoryTypeExpenditure)
			_, err := s.deactivateCategory(inactiveCategory.Id)
			s.handleErr(
				err,
				"error while deactivating category",
			)

			expenditureReq := s.createTestExpenditureRequest(
				&testAccount.Id,
				&inactiveCategory,
			)

			apiResponse, err := s.createExpenditureRequest(expenditureReq)
			s.handleErr(
				err,
				"error while making request",
			)

			s.assertHttpError(
				apiResponse,
				http.StatusConflict,
				domain.ErrCategoryInactive.Error(),
			)
		},
	)

	s.Run(
		"Create expenditure with insufficient balance",
		func() {
			// Create an account with low balance
			lowBalanceAccount := s.createTestAccountWithBalance(
				&testMember,
				"150",
				10.0,
			)

			expenditureReq := s.createTestExpenditureRequest(
				&lowBalanceAccount.Id,
				&testCategory,
			)
			expenditureReq.Amount = 1000.0 // Amount higher than account balance

			apiResponse, err := s.createExpenditureRequest(expenditureReq)
			s.handleErr(
				err,
				"error while making request",
			)

			s.assertHttpError(
				apiResponse,
				http.StatusConflict,
				domain.ErrInsufficientBalance.Error(),
			)
		},
	)

	s.Run(
		"Create expenditure with non-existent tag",
		func() {
			expenditureReq := s.createTestExpenditureRequest(
				&testAccount.Id,
				&testCategory,
			)
			expenditureReq.Tags = &[]openapi.Tag{
				{
					Id:      "999999",
					Name:    "Non-existent Tag",
					TagType: openapi.TagTypeExpenditure,
				},
			}

			apiResponse, err := s.createExpenditureRequest(expenditureReq)
			s.handleErr(
				err,
				"error while making request",
			)

			s.assertHttpError(
				apiResponse,
				http.StatusBadRequest,
				domain.ErrTagNotFound.Error(),
			)
		},
	)

	s.Run(
		"Create expenditure with valid tags",
		func() {
			// Create a test tag
			testTag := s.createTestTag(openapi.TagTypeExpenditure)

			expenditureReq := s.createTestExpenditureRequest(
				&testAccount.Id,
				&testCategory,
			)
			expenditureReq.Tags = &[]openapi.Tag{testTag}

			apiResponse, err := s.createExpenditureRequest(expenditureReq)
			s.handleErr(
				err,
				"error while making request",
			)

			s.Equal(
				http.StatusCreated,
				apiResponse.StatusCode,
			)

			var expenditure openapi.Expenditure
			s.decodeResponse(
				apiResponse,
				&expenditure,
			)

			// Verify the tags are included
			s.NotNil(expenditure.Tags)
			s.Equal(
				1,
				len(*expenditure.Tags),
			)
			s.Equal(
				testTag.Id,
				(*expenditure.Tags)[0].Id,
			)
		},
	)

	s.Run(
		"Create expenditure with minimal required fields",
		func() {
			expenditureReq := &openapi.ExpenditureRequest{
				AccountId: testAccount.Id,
				Amount:    50.0,
				Category:  testCategory,
				Currency:  "150",
			}

			apiResponse, err := s.createExpenditureRequest(expenditureReq)
			s.handleErr(
				err,
				"error while making request",
			)

			s.Equal(
				http.StatusCreated,
				apiResponse.StatusCode,
			)

			var expenditure openapi.Expenditure
			s.decodeResponse(
				apiResponse,
				&expenditure,
			)

			s.NotEmpty(expenditure.Id)
			s.Equal(
				expenditureReq.Amount,
				expenditure.Amount,
			)
			s.Equal(
				expenditureReq.Currency,
				expenditure.Currency,
			)
		},
	)

	s.Run(
		"Create expenditure with all optional fields",
		func() {
			apiResponse, err := s.createTag(
				&openapi.TagRequest{
					BackgroundColor: utils.StringPtr("#FFFFFF"),
					Color:           utils.StringPtr("#000000"),
					Description:     utils.StringPtr("Test tag with all fields"),
					Name:            "Test",
					TagType:         openapi.TagTypeExpenditure,
				},
			)
			s.handleErr(
				err,
				"error while making request",
			)

			s.Equal(
				http.StatusCreated,
				apiResponse.StatusCode,
			)

			var testTag openapi.Tag
			s.decodeResponse(
				apiResponse,
				&testTag,
			)
			expenditureReq := &openapi.ExpenditureRequest{
				AccountId:   testAccount.Id,
				Amount:      75.50,
				Category:    testCategory,
				Currency:    "150",
				Declared:    utils.BoolPtr(true),
				Description: "Test expenditure with all fields",
				Planned:     utils.BoolPtr(false),
				Tags:        &[]openapi.Tag{testTag},
			}

			apiResponse, err = s.createExpenditureRequest(expenditureReq)
			s.handleErr(
				err,
				"error while making request",
			)

			s.Equal(
				http.StatusCreated,
				apiResponse.StatusCode,
			)

			var expenditure openapi.Expenditure
			s.decodeResponse(
				apiResponse,
				&expenditure,
			)

			s.NotEmpty(expenditure.Id)
			s.Equal(
				expenditureReq.Amount,
				expenditure.Amount,
			)
			s.Equal(
				expenditureReq.Currency,
				expenditure.Currency,
			)
			s.Equal(
				expenditureReq.Description,
				expenditure.Description,
			)
			s.Equal(
				*expenditureReq.Declared,
				*expenditure.Declared,
			)
			s.Equal(
				*expenditureReq.Planned,
				*expenditure.Planned,
			)
			s.NotNil(expenditure.Tags)
			s.Equal(
				1,
				len(*expenditure.Tags),
			)
		},
	)
}

// Helper function to create expenditure request
func (s *Suite) createExpenditureRequest(expenditureReq *openapi.ExpenditureRequest) (
	*http.Response,
	error,
) {
	body, errBodyPrepare := utils.PrepareRequestBody(expenditureReq)
	s.handleErr(
		errBodyPrepare,
		"error while preparing request body",
	)

	req, errReq := http.NewRequestWithContext(
		s.ctx,
		http.MethodPost,
		"http://localhost:9091/expenditures",
		body,
	)
	s.handleErr(
		errReq,
		"error while creating request",
	)

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	client := &http.Client{}
	apiResponse, err := client.Do(req)
	s.handleErr(
		err,
		"error while making request",
	)

	return apiResponse, nil
}

// Helper function to create test expenditure request
func (s *Suite) createTestExpenditureRequest(
	account *string,
	category *openapi.Category,
) *openapi.ExpenditureRequest {
	return &openapi.ExpenditureRequest{
		AccountId:   *account,
		Amount:      100.50,
		Category:    *category,
		Currency:    "150",
		Declared:    utils.BoolPtr(true),
		Description: "Test expenditure for integration testing",
		Planned:     utils.BoolPtr(false),
	}
}

// Helper function to create test account with specific balance
func (s *Suite) createTestAccountWithBalance(
	owner *openapi.HouseholdMember,
	currency string,
	balance float32,
) openapi.Account {
	accountReq := &openapi.AccountRequest{
		AccountInformation: utils.StringPtr("Test Account Information"),
		AccountNumber:      utils.StringPtr("1234567890"),
		Active:             utils.BoolPtr(true),
		Currency:           currency,
		Description:        utils.StringPtr("Test account with specific balance"),
		InitialBalance:     balance,
		Institution:        utils.StringPtr("Test Bank"),
		Name:               "Test Account",
		Owner:              owner,
		Type:               "bank",
	}

	body, err := utils.PrepareRequestBody(accountReq)
	s.handleErr(
		err,
		"error while preparing account request body",
	)

	req, err := http.NewRequestWithContext(
		s.ctx,
		http.MethodPost,
		"http://localhost:9091/accounts",
		body,
	)
	s.handleErr(
		err,
		"error while creating account request",
	)

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	client := &http.Client{}
	apiResponse, err := client.Do(req)
	s.handleErr(
		err,
		"error while making account request",
	)

	defer apiResponse.Body.Close()

	var account openapi.Account
	s.decodeResponse(
		apiResponse,
		&account,
	)

	return account
}

// Helper function to create test account (using existing pattern)
func (s *Suite) createTestAccount(
	owner *openapi.HouseholdMember,
	currency string,
) openapi.Account {
	return s.createTestAccountWithBalance(
		owner,
		currency,
		1000.0,
	)
}

// Helper function to create test tag
func (s *Suite) createTestTag(tagType openapi.TagType) openapi.Tag {
	tagReq := &openapi.TagRequest{
		Name:    "Test Tag",
		TagType: tagType,
	}

	body, err := utils.PrepareRequestBody(tagReq)
	s.handleErr(
		err,
		"error while preparing tag request body",
	)

	req, err := http.NewRequestWithContext(
		s.ctx,
		http.MethodPost,
		"http://localhost:9091/tags",
		body,
	)
	s.handleErr(
		err,
		"error while creating tag request",
	)

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	client := &http.Client{}
	apiResponse, err := client.Do(req)
	s.handleErr(
		err,
		"error while making tag request",
	)

	defer apiResponse.Body.Close()

	var tag openapi.Tag
	s.decodeResponse(
		apiResponse,
		&tag,
	)

	return tag
}

// Helper
