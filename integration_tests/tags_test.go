package integration_tests

import (
	"fmt"
	"ghorkov32/proletariat-budget-be/integration_tests/utils"
	"ghorkov32/proletariat-budget-be/internal/core/domain"
	"ghorkov32/proletariat-budget-be/openapi"
	"net/http"
)

func (s *Suite) TestTags() {
	s.T().Log("Starting TestTags")

	s.Run(
		"Create a tag successfully",
		func() {
			tagRequest := &openapi.TagRequest{
				Name:            "Test Tag",
				Description:     utils.StringPtr("A test tag for integration testing"),
				TagType:         openapi.TagTypeExpenditure,
				Color:           utils.StringPtr("#FF0000"),
				BackgroundColor: utils.StringPtr("#00FF00"),
			}

			apiResponse, err := s.createTag(tagRequest)
			s.handleErr(
				err,
				"error while making request",
			)

			s.Equal(
				http.StatusCreated,
				apiResponse.StatusCode,
			)

			var createdTag openapi.Tag
			s.decodeResponse(
				apiResponse,
				&createdTag,
			)

			s.NotEmpty(createdTag.Id)
			s.Equal(
				tagRequest.Name,
				createdTag.Name,
			)
			s.Equal(
				*tagRequest.Description,
				*createdTag.Description,
			)
			s.Equal(
				tagRequest.TagType,
				createdTag.TagType,
			)
			s.Equal(
				*tagRequest.Color,
				*createdTag.Color,
			)
			s.Equal(
				*tagRequest.BackgroundColor,
				*createdTag.BackgroundColor,
			)
			s.NotNil(createdTag.CreatedAt)
			s.NotNil(createdTag.UpdatedAt)
		},
	)

	s.Run(
		"Create a tag with minimal required fields",
		func() {
			tagRequest := &openapi.TagRequest{
				Name:    "Minimal Tag",
				TagType: openapi.TagTypeIngress,
			}

			apiResponse, err := s.createTag(tagRequest)
			s.handleErr(
				err,
				"error while making request",
			)

			s.Equal(
				http.StatusCreated,
				apiResponse.StatusCode,
			)

			var createdTag openapi.Tag
			s.decodeResponse(
				apiResponse,
				&createdTag,
			)

			s.NotEmpty(createdTag.Id)
			s.Equal(
				tagRequest.Name,
				createdTag.Name,
			)
			s.Equal(
				tagRequest.TagType,
				createdTag.TagType,
			)
			s.NotNil(createdTag.CreatedAt)
			s.NotNil(createdTag.UpdatedAt)
		},
	)

	s.Run(
		"Create a tag with empty name should fail",
		func() {
			tagRequest := &openapi.TagRequest{
				Name:    "",
				TagType: openapi.TagTypeIngress,
			}

			apiResponse, err := s.createTag(tagRequest)
			s.handleErr(
				err,
				"error while making request",
			)

			// Expecting validation error for empty name
			s.Equal(
				http.StatusBadRequest,
				apiResponse.StatusCode,
			)
		},
	)

	s.Run(
		"Create a tag with invalid type should fail",
		func() {
			tagRequest := &openapi.TagRequest{
				Name:    "Invalid Type Tag",
				TagType: "invalid_type",
			}

			apiResponse, err := s.createTag(tagRequest)
			s.handleErr(
				err,
				"error while making request",
			)

			// Expecting validation error for invalid type
			s.Equal(
				http.StatusBadRequest,
				apiResponse.StatusCode,
			)
		},
	)

	s.Run(
		"Create a tag with duplicate name should handle gracefully",
		func() {
			tagRequest := &openapi.TagRequest{
				Name:        "Duplicate Tag",
				Description: utils.StringPtr("First tag with this name"),
				TagType:     openapi.TagTypeExpenditure,
			}

			// Create first tag
			apiResponse1, err := s.createTag(tagRequest)
			s.handleErr(
				err,
				"error while making first request",
			)

			s.Equal(
				http.StatusCreated,
				apiResponse1.StatusCode,
			)

			// Try to create second tag with same name
			tagRequest2 := &openapi.TagRequest{
				Name:        "Duplicate Tag",
				Description: utils.StringPtr("Second tag with same name"),
				TagType:     openapi.TagTypeExpenditure,
			}

			apiResponse2, err := s.createTag(tagRequest2)
			s.handleErr(
				err,
				"error while making second request",
			)

			// Depending on business rules, this might succeed or fail
			// If duplicates are allowed, expect 201; if not, expect 400 or 409
			s.assertHttpError(
				apiResponse2,
				http.StatusConflict,
				domain.ErrTagAlreadyExists.Error(),
			)
		},
	)

	s.Run(
		"Create a tag with very long name should handle appropriately",
		func() {
			longName := "This is a very long tag name that might exceed the maximum allowed length for tag names in the database schema and should be handled appropriately by the validation layer"

			tagRequest := &openapi.TagRequest{
				Name:    longName,
				TagType: openapi.TagTypeExpenditure,
			}

			apiResponse, err := s.createTag(tagRequest)
			s.handleErr(
				err,
				"error while making request",
			)

			// Expecting either success (if length is within limits) or validation error
			s.True(
				apiResponse.StatusCode == http.StatusCreated ||
					apiResponse.StatusCode == http.StatusBadRequest,
				"Response should be either success or validation error",
			)
		},
	)

	s.Run(
		"Create a tag with very long description should handle appropriately",
		func() {
			longDescription := "This is a very long description that might exceed the maximum allowed length for tag descriptions in the database schema. It contains multiple sentences and should test the validation limits properly. The system should either accept it if within limits or reject it with appropriate error message."

			tagRequest := &openapi.TagRequest{
				Name:        "Long Description Tag",
				Description: &longDescription,
				TagType:     openapi.TagTypeExpenditure,
			}

			apiResponse, err := s.createTag(tagRequest)
			s.handleErr(
				err,
				"error while making request",
			)

			// Expecting either success (if length is within limits) or validation error
			s.True(
				apiResponse.StatusCode == http.StatusCreated ||
					apiResponse.StatusCode == http.StatusBadRequest,
				"Response should be either success or validation error",
			)
		},
	)

	s.Run(
		"Create a tag with special characters in name",
		func() {
			tagRequest := &openapi.TagRequest{
				Name:        "Special-Tag_123!@#",
				Description: utils.StringPtr("Tag with special characters"),
				TagType:     openapi.TagTypeExpenditure,
			}

			apiResponse, err := s.createTag(tagRequest)
			s.handleErr(
				err,
				"error while making request",
			)

			// Depending on validation rules, this might succeed or fail
			s.True(
				apiResponse.StatusCode == http.StatusCreated ||
					apiResponse.StatusCode == http.StatusBadRequest,
				"Response should be either success or validation error",
			)

			if apiResponse.StatusCode == http.StatusCreated {
				var createdTag openapi.Tag
				s.decodeResponse(
					apiResponse,
					&createdTag,
				)
				s.Equal(
					tagRequest.Name,
					createdTag.Name,
				)
			}
		},
	)

	s.Run(
		"Create a tag with unicode characters",
		func() {
			tagRequest := &openapi.TagRequest{
				Name:        "Тег на русском 🏷️",
				Description: utils.StringPtr("Tag with unicode characters including emoji"),
				TagType:     openapi.TagTypeExpenditure,
			}

			apiResponse, err := s.createTag(tagRequest)
			s.handleErr(
				err,
				"error while making request",
			)

			// Should handle unicode properly
			s.True(
				apiResponse.StatusCode == http.StatusCreated ||
					apiResponse.StatusCode == http.StatusBadRequest,
				"Response should be either success or validation error",
			)

			if apiResponse.StatusCode == http.StatusCreated {
				var createdTag openapi.Tag
				s.decodeResponse(
					apiResponse,
					&createdTag,
				)
				s.Equal(
					tagRequest.Name,
					createdTag.Name,
				)
			}
		},
	)

	s.Run(
		"Delete an existing tag successfully",
		func() {
			// First create a tag to delete
			tagRequest := &openapi.TagRequest{
				Name:        "Tag to Delete",
				Description: utils.StringPtr("This tag will be deleted"),
				TagType:     openapi.TagTypeExpenditure,
			}

			createResponse, err := s.createTag(tagRequest)
			s.handleErr(
				err,
				"error while creating tag",
			)
			s.Equal(
				http.StatusCreated,
				createResponse.StatusCode,
			)

			var createdTag openapi.Tag
			s.decodeResponse(
				createResponse,
				&createdTag,
			)

			// Now delete the tag
			deleteResponse, err := s.deleteTag(createdTag.Id)
			s.handleErr(
				err,
				"error while deleting tag",
			)

			s.Equal(
				http.StatusNoContent,
				deleteResponse.StatusCode,
			)

			// Verify the tag is actually deleted by trying to get it
			// This would require a getTag function, but we can test the delete response for now
		},
	)
	s.Run(
		"Delete a non-existent tag should return 404",
		func() {
			nonExistentId := "999999"

			deleteResponse, err := s.deleteTag(nonExistentId)
			s.handleErr(
				err,
				"error while making delete request",
			)

			s.assertHttpError(
				deleteResponse,
				http.StatusNotFound,
				domain.ErrTagNotFound.Error(),
			)
		},
	)

	s.Run(
		"Delete same tag twice should return 404 on second attempt",
		func() {
			// Create a tag
			tagRequest := &openapi.TagRequest{
				Name:        "Double Delete Tag",
				Description: utils.StringPtr("This tag will be deleted twice"),
				TagType:     openapi.TagTypeIngress,
			}

			createResponse, err := s.createTag(tagRequest)
			s.handleErr(
				err,
				"error while creating tag",
			)
			s.Equal(
				http.StatusCreated,
				createResponse.StatusCode,
			)

			var createdTag openapi.Tag
			s.decodeResponse(
				createResponse,
				&createdTag,
			)

			// First deletion should succeed
			deleteResponse1, err := s.deleteTag(createdTag.Id)
			s.handleErr(
				err,
				"error while making first delete request",
			)

			s.Equal(
				http.StatusNoContent,
				deleteResponse1.StatusCode,
			)

			// Second deletion should fail with 404
			deleteResponse2, err := s.deleteTag(createdTag.Id)
			s.handleErr(
				err,
				"error while making second delete request",
			)

			s.assertHttpError(
				deleteResponse2,
				http.StatusNotFound,
				domain.ErrTagNotFound.Error(),
			)
		},
	)

	s.Run(
		"Delete tag with different tag types",
		func() {
			tagTypes := []openapi.TagType{
				openapi.TagTypeExpenditure,
				openapi.TagTypeIngress,
				openapi.TagTypeTransfer,
				openapi.TagTypeSavingGoal,
				openapi.TagTypeSavingsWithdrawal,
				openapi.TagTypeSavingsContribution,
			}

			for _, tagType := range tagTypes {
				tagRequest := &openapi.TagRequest{
					Name: fmt.Sprintf(
						"Delete Test %s",
						tagType,
					),
					Description: utils.StringPtr(
						fmt.Sprintf(
							"Testing deletion for %s type",
							tagType,
						),
					),
					TagType: tagType,
				}

				createResponse, err := s.createTag(tagRequest)
				s.handleErr(
					err,
					fmt.Sprintf(
						"error while creating %s tag",
						tagType,
					),
				)
				s.Equal(
					http.StatusCreated,
					createResponse.StatusCode,
				)

				var createdTag openapi.Tag
				s.decodeResponse(
					createResponse,
					&createdTag,
				)

				// Delete the tag
				deleteResponse, err := s.deleteTag(createdTag.Id)
				s.handleErr(
					err,
					fmt.Sprintf(
						"error while deleting %s tag",
						tagType,
					),
				)

				s.Equal(
					http.StatusNoContent,
					deleteResponse.StatusCode,
					fmt.Sprintf(
						"Should successfully delete %s tag",
						tagType,
					),
				)
			}
		},
	)
	s.Run(
		"List all tags after creating several",
		func() {
			// Create multiple tags
			tagRequests := []*openapi.TagRequest{
				{
					Name:        "List Test Tag 1",
					Description: utils.StringPtr("First tag for listing test"),
					TagType:     openapi.TagTypeExpenditure,
					Color:       utils.StringPtr("#FF0000"),
				},
				{
					Name:        "List Test Tag 2",
					Description: utils.StringPtr("Second tag for listing test"),
					TagType:     openapi.TagTypeIngress,
					Color:       utils.StringPtr("#00FF00"),
				},
				{
					Name:    "List Test Tag 3",
					TagType: openapi.TagTypeTransfer,
				},
			}

			var createdTagIds []string
			for i, tagRequest := range tagRequests {
				createResponse, err := s.createTag(tagRequest)
				s.handleErr(
					err,
					fmt.Sprintf(
						"error while creating tag %d",
						i+1,
					),
				)
				s.Equal(
					http.StatusCreated,
					createResponse.StatusCode,
				)

				var createdTag openapi.Tag
				s.decodeResponse(
					createResponse,
					&createdTag,
				)
				createdTagIds = append(
					createdTagIds,
					createdTag.Id,
				)
			}

			// List all tags
			listResponse, err := s.listTags()
			s.handleErr(
				err,
				"error while listing tags",
			)

			s.Equal(
				http.StatusOK,
				listResponse.StatusCode,
			)

			var tagList []openapi.Tag
			s.decodeResponse(
				listResponse,
				&tagList,
			)

			// Should contain at least the tags we created
			s.GreaterOrEqual(
				len(tagList),
				len(tagRequests),
				"Should contain at least the created tags",
			)

			// Verify our created tags are in the list
			createdTagNames := make(map[string]bool)
			for _, tagRequest := range tagRequests {
				createdTagNames[tagRequest.Name] = false
			}

			for _, tag := range tagList {
				if _, exists := createdTagNames[tag.Name]; exists {
					createdTagNames[tag.Name] = true
					// Verify tag structure
					s.NotEmpty(tag.Id)
					s.NotEmpty(tag.Name)
					s.NotNil(tag.CreatedAt)
					s.NotNil(tag.UpdatedAt)
				}
			}

			// All created tags should be found
			for tagName, found := range createdTagNames {
				s.True(
					found,
					fmt.Sprintf(
						"Tag '%s' should be found in the list",
						tagName,
					),
				)
			}

			// Clean up
			for _, tagId := range createdTagIds {
				_, err = s.deleteTag(tagId)
				s.handleErr(
					err,
					fmt.Sprintf(
						"error deleting tag %s",
						tagId,
					),
				)
			}
		},
	)

	s.Run(
		"List tags by type - expenditure",
		func() {
			// Create tags of different types
			expenditureTag := &openapi.TagRequest{
				Name:        "Expenditure Tag for Type Test",
				Description: utils.StringPtr("Testing expenditure type filtering"),
				TagType:     openapi.TagTypeExpenditure,
			}

			ingressTag := &openapi.TagRequest{
				Name:        "Ingress Tag for Type Test",
				Description: utils.StringPtr("Testing ingress type filtering"),
				TagType:     openapi.TagTypeIngress,
			}

			// Create both tags
			expResponse, err := s.createTag(expenditureTag)
			s.handleErr(
				err,
				"error creating expenditure tag",
			)
			s.Equal(
				http.StatusCreated,
				expResponse.StatusCode,
			)

			var createdExpTag openapi.Tag
			s.decodeResponse(
				expResponse,
				&createdExpTag,
			)

			ingResponse, err := s.createTag(ingressTag)
			s.handleErr(
				err,
				"error creating ingress tag",
			)
			s.Equal(
				http.StatusCreated,
				ingResponse.StatusCode,
			)

			var createdIngTag openapi.Tag
			s.decodeResponse(
				ingResponse,
				&createdIngTag,
			)

			// List only expenditure tags
			listResponse, err := s.listTagsByType(openapi.TagTypeExpenditure)
			s.handleErr(
				err,
				"error while listing expenditure tags",
			)

			s.Equal(
				http.StatusOK,
				listResponse.StatusCode,
			)

			var tagList []openapi.Tag
			s.decodeResponse(
				listResponse,
				&tagList,
			)

			// Should contain at least our expenditure tag
			s.GreaterOrEqual(
				len(tagList),
				1,
				"Should contain at least one expenditure tag",
			)

			// All returned tags should be expenditure type
			expenditureTagFound := false
			for _, tag := range tagList {
				s.Equal(
					openapi.TagTypeExpenditure,
					tag.TagType,
					"All tags should be expenditure type",
				)
				if tag.Id == createdExpTag.Id {
					expenditureTagFound = true
				}
			}

			s.True(
				expenditureTagFound,
				"Created expenditure tag should be in the list",
			)

			// Clean up
			_, err = s.deleteTag(createdExpTag.Id)
			s.handleErr(
				err,
				fmt.Sprintf(
					"error deleting tag %s",
					createdExpTag.Id,
				),
			)
			_, err = s.deleteTag(createdIngTag.Id)
			s.handleErr(
				err,
				fmt.Sprintf(
					"error deleting tag %s",
					createdIngTag.Id,
				),
			)
		},
	)

	s.Run(
		"List tags by type - ingress",
		func() {
			// Create an ingress tag
			ingressTag := &openapi.TagRequest{
				Name:        "Ingress Tag for Type Test 2",
				Description: utils.StringPtr("Testing ingress type filtering"),
				TagType:     openapi.TagTypeIngress,
			}

			createResponse, err := s.createTag(ingressTag)
			s.handleErr(
				err,
				"error creating ingress tag",
			)
			s.Equal(
				http.StatusCreated,
				createResponse.StatusCode,
			)

			var createdTag openapi.Tag
			s.decodeResponse(
				createResponse,
				&createdTag,
			)

			// List only ingress tags
			listResponse, err := s.listTagsByType(openapi.TagTypeIngress)
			s.handleErr(
				err,
				"error while listing ingress tags",
			)

			s.Equal(
				http.StatusOK,
				listResponse.StatusCode,
			)

			var tagList []openapi.Tag
			s.decodeResponse(
				listResponse,
				&tagList,
			)

			// Should contain at least our ingress tag
			s.GreaterOrEqual(
				len(tagList),
				1,
				"Should contain at least one ingress tag",
			)

			// All returned tags should be ingress type
			ingressTagFound := false
			for _, tag := range tagList {
				s.Equal(
					openapi.TagTypeIngress,
					tag.TagType,
					"All tags should be ingress type",
				)
				if tag.Id == createdTag.Id {
					ingressTagFound = true
				}
			}

			s.True(
				ingressTagFound,
				"Created ingress tag should be in the list",
			)

			// Clean up
			_, err = s.deleteTag(createdTag.Id)
			s.handleErr(
				err,
				fmt.Sprintf(
					"error deleting tag %s",
					createdTag.Id,
				),
			)
		},
	)

	s.Run(
		"List tags by all supported types",
		func() {
			tagTypes := []openapi.TagType{
				openapi.TagTypeExpenditure,
				openapi.TagTypeIngress,
				// openapi.TagTypeTransfer, // not supported yet
				openapi.TagTypeSavingGoal,
				openapi.TagTypeSavingsWithdrawal,
				openapi.TagTypeSavingsContribution,
			}

			var createdTagIds []string

			// Create one tag for each type
			for _, tagType := range tagTypes {
				tagRequest := &openapi.TagRequest{
					Name: fmt.Sprintf(
						"List Type Test %s",
						tagType,
					),
					Description: utils.StringPtr(
						fmt.Sprintf(
							"Testing list for %s type",
							tagType,
						),
					),
					TagType: tagType,
				}

				createResponse, err := s.createTag(tagRequest)
				s.handleErr(
					err,
					fmt.Sprintf(
						"error creating %s tag",
						tagType,
					),
				)
				s.Equal(
					http.StatusCreated,
					createResponse.StatusCode,
				)

				var createdTag openapi.Tag
				s.decodeResponse(
					createResponse,
					&createdTag,
				)
				createdTagIds = append(
					createdTagIds,
					createdTag.Id,
				)

				// Test listing by this specific type
				listResponse, err := s.listTagsByType(tagType)
				s.handleErr(
					err,
					fmt.Sprintf(
						"error listing %s tags",
						tagType,
					),
				)

				s.Equal(
					http.StatusOK,
					listResponse.StatusCode,
				)

				var tagList []openapi.Tag
				s.decodeResponse(
					listResponse,
					&tagList,
				)

				// Should contain at least our created tag
				s.GreaterOrEqual(
					len(tagList),
					1,
					fmt.Sprintf(
						"Should contain at least one %s tag",
						tagType,
					),
				)

				// All returned tags should be of the requested type
				tagFound := false
				for _, tag := range tagList {
					s.Equal(
						tagType,
						tag.TagType,
						fmt.Sprintf(
							"All tags should be %s type",
							tagType,
						),
					)
					if tag.Id == createdTag.Id {
						tagFound = true
					}
				}

				s.True(
					tagFound,
					fmt.Sprintf(
						"Created %s tag should be in the list",
						tagType,
					),
				)
			}

			// Clean up
			for _, tagId := range createdTagIds {
				_, err := s.deleteTag(tagId)
				s.handleErr(
					err,
					fmt.Sprintf(
						"error deleting tag %s",
						tagId,
					),
				)
			}
		},
	)

}

func (s *Suite) createTag(tagRequest *openapi.TagRequest) (
	*http.Response,
	error,
) {
	body, errBodyPrepare := utils.PrepareRequestBody(tagRequest)

	s.handleErr(
		errBodyPrepare,
		"error while preparing request body",
	)

	req, errReq := http.NewRequestWithContext(
		s.ctx,
		http.MethodPost,
		"http://localhost:9091/tags",
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

func (s *Suite) deleteTag(tagId string) (
	*http.Response,
	error,
) {
	req, errReq := http.NewRequestWithContext(
		s.ctx,
		http.MethodDelete,
		fmt.Sprintf(
			"http://localhost:9091/tags/%s",
			tagId,
		),
		nil,
	)
	s.handleErr(
		errReq,
		"error while creating delete request",
	)

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	client := &http.Client{}
	apiResponse, err := client.Do(req)
	s.handleErr(
		err,
		"error while making delete request",
	)

	return apiResponse, nil
}

func (s *Suite) listTags() (
	*http.Response,
	error,
) {
	req, errReq := http.NewRequestWithContext(
		s.ctx,
		http.MethodGet,
		"http://localhost:9091/tags",
		nil,
	)
	s.handleErr(
		errReq,
		"error while creating list tags request",
	)

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	client := &http.Client{}
	apiResponse, err := client.Do(req)
	s.handleErr(
		err,
		"error while making list tags request",
	)

	return apiResponse, nil
}

func (s *Suite) listTagsByType(tagType openapi.TagType) (
	*http.Response,
	error,
) {
	req, errReq := http.NewRequestWithContext(
		s.ctx,
		http.MethodGet,
		fmt.Sprintf(
			"http://localhost:9091/tags/type/%s",
			tagType,
		),
		nil,
	)
	s.handleErr(
		errReq,
		"error while creating list tags by type request",
	)

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	client := &http.Client{}
	apiResponse, err := client.Do(req)
	s.handleErr(
		err,
		"error while making list tags by type request",
	)

	return apiResponse, nil
}
