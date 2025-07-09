package port

type Ports struct {
	Account          *AccountRepo
	Auth             *AuthRepo
	Category         *CategoryRepo
	Expenditure      *ExpenditureRepo
	HouseholdMembers *HouseholdMembersRepo
	Ingress          *IngressRepo
	SavingGoal       *SavingsGoalRepo
	Tags             *TagsRepo
	Transaction      *TransactionRepo
}
