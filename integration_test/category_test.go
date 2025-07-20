package integration_test

import (
	"net/http"

	"ghorkov32/proletariat-budget-be/integration_test/utils"
	"ghorkov32/proletariat-budget-be/internal/core/domain"
	"ghorkov32/proletariat-budget-be/openapi"
)

func (s *Suite) TestCategory() {
	s.T().Log("Starting TestCategory")

	var deletableCategoryId string

	s.Run(
		"Create a category successfully",
		func() {
			categoryType := openapi.CategoryTypeExpenditure
			categoryReq := &openapi.CategoryRequest{
				Name:            "Test Category",
				Description:     "Test category description",
				CategoryType:    &categoryType,
				Color:           utils.StringPtr("#FF0000"),
				BackgroundColor: utils.StringPtr("#00FF00"),
			}

			apiResponse := s.createCategory(categoryReq)

			s.Equal(
				http.StatusCreated,
				apiResponse.StatusCode,
			)

			var createdCategory openapi.Category
			s.decodeResponse(
				apiResponse,
				&createdCategory,
			)

			s.Equal(
				categoryReq.Name,
				createdCategory.Name,
			)
			s.Equal(
				categoryReq.Description,
				createdCategory.Description,
			)
			s.Equal(
				openapi.CategoryTypeExpenditure,
				*createdCategory.CategoryType,
			)
			s.Equal(
				true,
				*createdCategory.Active,
			)
			s.NotEmpty(createdCategory.Id)

			deletableCategoryId = createdCategory.Id
		},
	)

	s.Run(
		"Create a category with minimal data",
		func() {
			categoryType := openapi.CategoryTypeIngress
			categoryReq := &openapi.CategoryRequest{
				Name:            "Minimal Category",
				CategoryType:    &categoryType,
				Color:           utils.StringPtr("#FF0000"),
				BackgroundColor: utils.StringPtr("#00FF00"),
			}

			apiResponse := s.createCategory(categoryReq)

			s.Equal(
				http.StatusCreated,
				apiResponse.StatusCode,
			)

			var createdCategory openapi.Category
			s.decodeResponse(
				apiResponse,
				&createdCategory,
			)

			s.Equal(
				categoryReq.Name,
				createdCategory.Name,
			)
			s.Equal(
				openapi.CategoryTypeIngress,
				*createdCategory.CategoryType,
			)
			s.NotEmpty(createdCategory.Id)
		},
	)

	s.Run(
		"List categories without filter",
		func() {
			// Create a few test categories first
			_ = s.createTestCategory(openapi.CategoryTypeIngress)
			_ = s.createTestCategory(openapi.CategoryTypeExpenditure)

			apiResponse, err := s.listCategories(nil)
			s.handleErr(
				err,
				"error while making request",
			)

			s.Equal(
				http.StatusOK,
				apiResponse.StatusCode,
			)

			var categoryList openapi.ListCategories200JSONResponse
			s.decodeResponse(
				apiResponse,
				&categoryList,
			)

			s.NotNil(categoryList.Categories)
			s.Greater(
				len(*categoryList.Categories),
				0,
			)

			// Verify structure of returned categories
			for _, category := range *categoryList.Categories {
				s.NotEmpty(category.Id)
				s.NotEmpty(category.Name)
				s.NotNil(category.CategoryType)
				s.NotNil(category.Active)
			}
		},
	)

	s.Run(
		"List categories filtered by CategoryType",
		func() {
			// Create categories of different types
			_ = s.createTestCategory(openapi.CategoryTypeIngress)
			_ = s.createTestCategory(openapi.CategoryTypeExpenditure)

			categoryType := openapi.CategoryTypeExpenditure
			apiResponse, err := s.listCategories(&categoryType)
			s.handleErr(
				err,
				"error while making request",
			)

			s.Equal(
				http.StatusOK,
				apiResponse.StatusCode,
			)

			var categoryList openapi.ListCategories200JSONResponse
			s.decodeResponse(
				apiResponse,
				&categoryList,
			)

			s.NotNil(categoryList.Categories)
			s.Greater(
				len(*categoryList.Categories),
				0,
			)

			// Verify all returned categories have the correct CategoryType
			for _, category := range *categoryList.Categories {
				s.Equal(
					openapi.CategoryTypeExpenditure,
					*category.CategoryType,
				)
			}
		},
	)

	s.Run(
		"Update a category successfully",
		func() {
			categoryType := openapi.CategoryTypeExpenditure
			category := s.createTestCategory(categoryType)

			updateReq := &openapi.CategoryRequest{
				Name:            "Updated Category Name",
				Description:     "Updated description",
				CategoryType:    &categoryType,
				Color:           utils.StringPtr("#FF0000"),
				BackgroundColor: utils.StringPtr("#00FF00"),
			}

			apiResponse, err := s.updateCategory(
				category.Id,
				updateReq,
			)
			s.handleErr(
				err,
				"error while making request",
			)

			s.Equal(
				http.StatusOK,
				apiResponse.StatusCode,
			)

			var updatedCategory openapi.Category
			s.decodeResponse(
				apiResponse,
				&updatedCategory,
			)

			s.Equal(
				updateReq.Name,
				updatedCategory.Name,
			)
			s.Equal(
				updateReq.Description,
				updatedCategory.Description,
			)
			s.Equal(
				category.Id,
				updatedCategory.Id,
			)
		},
	)

	s.Run(
		"Update a non-existent category",
		func() {
			categoryType := openapi.CategoryTypeExpenditure
			updateReq := &openapi.CategoryRequest{
				Name:            "Non-existent Category",
				CategoryType:    &categoryType,
				Color:           utils.StringPtr("#FF0000"),
				BackgroundColor: utils.StringPtr("#00FF00"),
			}

			apiResponse, err := s.updateCategory(
				"non-existent-id",
				updateReq,
			)
			s.handleErr(
				err,
				"error while making request",
			)

			s.assertHttpError(
				apiResponse,
				http.StatusNotFound,
				domain.ErrCategoryNotFound.Error(),
			)
		},
	)

	s.Run(
		"Activate a category successfully",
		func() {
			// Create and deactivate a category first
			categoryType := openapi.CategoryTypeExpenditure
			category := s.createTestCategory(categoryType)

			deactivateResponse, err := s.deactivateCategory(category.Id)
			s.handleErr(
				err,
				"error while deactivating category",
			)
			s.Equal(
				http.StatusNoContent,
				deactivateResponse.StatusCode,
			)

			// Now activate it
			apiResponse, err := s.activateCategory(category.Id)
			s.handleErr(
				err,
				"error while making request",
			)

			s.Equal(
				http.StatusNoContent,
				apiResponse.StatusCode,
			)
		},
	)

	s.Run(
		"Activate a non-existent category",
		func() {
			apiResponse, err := s.activateCategory("non-existent-id")
			s.handleErr(
				err,
				"error while making request",
			)

			s.assertHttpError(
				apiResponse,
				http.StatusNotFound,
				domain.ErrCategoryNotFound.Error(),
			)
		},
	)

	s.Run(
		"Activate an already active category",
		func() {
			categoryType := openapi.CategoryTypeExpenditure
			category := s.createTestCategory(categoryType)

			apiResponse, err := s.activateCategory(category.Id)
			s.handleErr(
				err,
				"error while making request",
			)

			s.assertHttpError(
				apiResponse,
				http.StatusBadRequest,
				domain.ErrCategoryAlreadyActive.Error(),
			)
		},
	)

	s.Run(
		"Deactivate a category successfully",
		func() {
			categoryType := openapi.CategoryTypeExpenditure
			category := s.createTestCategory(categoryType)

			apiResponse, err := s.deactivateCategory(category.Id)
			s.handleErr(
				err,
				"error while making request",
			)

			s.Equal(
				http.StatusNoContent,
				apiResponse.StatusCode,
			)
		},
	)

	s.Run(
		"Deactivate a non-existent category",
		func() {
			apiResponse, err := s.deactivateCategory("non-existent-id")
			s.handleErr(
				err,
				"error while making request",
			)

			s.assertHttpError(
				apiResponse,
				http.StatusNotFound,
				domain.ErrCategoryNotFound.Error(),
			)
		},
	)

	s.Run(
		"Deactivate an already inactive category",
		func() {
			categoryType := openapi.CategoryTypeExpenditure
			category := s.createTestCategory(categoryType)

			// Deactivate first time
			deactivateResponse, err := s.deactivateCategory(category.Id)
			s.handleErr(
				err,
				"error while deactivating category",
			)
			s.Equal(
				http.StatusNoContent,
				deactivateResponse.StatusCode,
			)

			// Try to deactivate again
			apiResponse, err := s.deactivateCategory(category.Id)
			s.handleErr(
				err,
				"error while making request",
			)

			s.assertHttpError(
				apiResponse,
				http.StatusBadRequest,
				domain.ErrCategoryAlreadyInactive.Error(),
			)
		},
	)

	s.Run(
		"Delete a category successfully",
		func() {
			apiResponse, err := s.deleteCategory(deletableCategoryId)
			s.handleErr(
				err,
				"error while making request",
			)

			s.Equal(
				http.StatusNoContent,
				apiResponse.StatusCode,
			)

			// Now check if the category is deleted. Using list because i didn't made a single get
			getResponse, err := s.listCategories(nil)
			s.handleErr(
				err,
				"error while making request",
			)
			s.Equal(
				http.StatusOK,
				getResponse.StatusCode,
			)
			var categories openapi.ListCategories200JSONResponse
			s.decodeResponse(
				getResponse,
				&categories,
			)
			for _, category := range *categories.Categories {
				if category.Id == deletableCategoryId {
					s.Fail("Deleted category still exists in the list")
				}
			}
		},
	)

	s.Run(
		"Delete a non-existent category",
		func() {
			apiResponse, err := s.deleteCategory("999")
			s.handleErr(
				err,
				"error while making request",
			)

			s.assertHttpError(
				apiResponse,
				http.StatusNotFound,
				domain.ErrCategoryNotFound.Error(),
			)
		},
	)

	s.Run(
		"Delete a category used in expenditure",
		func() {
			// This test assumes there's a category with ID "1" that's used in expenditures
			// You might need to adjust this based on your mock data
			apiResponse, err := s.deleteCategory("1")
			s.handleErr(
				err,
				"error while making request",
			)

			s.assertHttpError(
				apiResponse,
				http.StatusBadRequest,
				domain.ErrCategoryUsedInExpenditure.Error(),
			)
		},
	)
}

// Helper function to create a test category
func (s *Suite) createTestCategory(categoryType openapi.CategoryType) openapi.Category {
	categoryReq := &openapi.CategoryRequest{
		Name:            "Test Category " + string(categoryType),
		Description:     "Test category description",
		CategoryType:    &categoryType,
		Color:           utils.StringPtr("#FF0000"),
		BackgroundColor: utils.StringPtr("#00FF00"),
	}

	category, err := s.createCategoryAndReturn(categoryReq)
	s.handleErr(
		err,
		"error while creating test category",
	)

	return category
}

// Function to create a category and return the response
func (s *Suite) createCategory(categoryReq *openapi.CategoryRequest) *http.Response {
	body, errBodyPrepare := utils.PrepareRequestBody(categoryReq)
	s.handleErr(
		errBodyPrepare,
		"error while preparing request body",
	)

	req, errReq := http.NewRequestWithContext(
		s.ctx,
		http.MethodPost,
		"http://localhost:9091/categories",
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

	return apiResponse
}

// Function to create a category and return the created category object
func (s *Suite) createCategoryAndReturn(categoryReq *openapi.CategoryRequest) (
	openapi.Category,
	error,
) {
	apiResponse := s.createCategory(categoryReq)
	defer apiResponse.Body.Close()

	var category openapi.Category
	s.decodeResponse(
		apiResponse,
		&category,
	)

	return category, nil
}

// Function to list categories
func (s *Suite) listCategories(categoryType *openapi.CategoryType) (
	*http.Response,
	error,
) {
	req, errReq := http.NewRequestWithContext(
		s.ctx,
		http.MethodGet,
		"http://localhost:9091/categories",
		nil,
	)
	s.handleErr(
		errReq,
		"error while creating request",
	)

	if categoryType != nil {
		q := req.URL.Query()
		q.Set(
			"type",
			string(*categoryType),
		)
		req.URL.RawQuery = q.Encode()
	}

	client := &http.Client{}
	apiResponse, err := client.Do(req)
	s.handleErr(
		err,
		"error while making request",
	)

	return apiResponse, nil
}

// Function to update a category
func (s *Suite) updateCategory(
	id string,
	categoryReq *openapi.CategoryRequest,
) (
	*http.Response,
	error,
) {
	body, errBodyPrepare := utils.PrepareRequestBody(categoryReq)
	s.handleErr(
		errBodyPrepare,
		"error while preparing request body",
	)

	req, errReq := http.NewRequestWithContext(
		s.ctx,
		http.MethodPut,
		"http://localhost:9091/categories/"+id,
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

// Function to delete a category
func (s *Suite) deleteCategory(id string) (
	*http.Response,
	error,
) {
	req, errReq := http.NewRequestWithContext(
		s.ctx,
		http.MethodDelete,
		"http://localhost:9091/categories/"+id,
		nil,
	)
	s.handleErr(
		errReq,
		"error while creating request",
	)

	client := &http.Client{}
	apiResponse, err := client.Do(req)
	s.handleErr(
		err,
		"error while making request",
	)

	return apiResponse, nil
}

// Function to activate a category
func (s *Suite) activateCategory(id string) (
	*http.Response,
	error,
) {
	req, errReq := http.NewRequestWithContext(
		s.ctx,
		http.MethodPatch,
		"http://localhost:9091/categories/"+id+"/activate",
		nil,
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

// Function to deactivate a category
func (s *Suite) deactivateCategory(id string) (
	*http.Response,
	error,
) {
	req, errReq := http.NewRequestWithContext(
		s.ctx,
		http.MethodPatch,
		"http://localhost:9091/categories/"+id+"/deactivate",
		nil,
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
