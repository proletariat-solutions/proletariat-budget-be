package usecase

type UseCases struct {
	Account         *AccountUseCase
	Auth            *AuthUseCase
	HouseholdMember *HouseholdMemberUseCase
	Expenditure     *ExpenditureUseCase
	Category        *CategoryUseCase
	Tags            *TagsUseCase
}
