package resthttp

import (
	"ghorkov32/proletariat-budget-be/internal/core/domain"
	"ghorkov32/proletariat-budget-be/openapi"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"time"
)

func FromOAPIAccount(a *openapi.Account) *domain.Account {
	return &domain.Account{
		ID:                 &a.Id,
		Name:               a.Name,
		Type:               domain.AccountType(a.Type),
		Currency:           a.Currency,
		InitialBalance:     a.InitialBalance,
		CurrentBalance:     a.CurrentBalance,
		Description:        a.Description,
		Institution:        a.Institution,
		AccountNumber:      a.AccountNumber,
		AccountInformation: a.AccountInformation,
		OwnerID:            &a.Owner.Id,
		Active:             *a.Active,
		CreatedAt:          a.CreatedAt,
		UpdatedAt:          a.UpdatedAt,
	}
}

// FromOAPIAccountRequest converts an OpenAPI AccountRequest to domain Account
func FromOAPIAccountRequest(a *openapi.AccountRequest) *domain.Account {
	var ownerID *string
	if a.Owner != nil {
		ownerID = &a.Owner.Id
	}

	return &domain.Account{
		Name:               a.Name,
		Type:               domain.AccountType(a.Type),
		Currency:           a.Currency,
		InitialBalance:     a.InitialBalance,
		CurrentBalance:     a.InitialBalance, // Set current balance to initial balance for new accounts
		Description:        a.Description,
		Institution:        a.Institution,
		AccountNumber:      a.AccountNumber,
		AccountInformation: a.AccountInformation,
		OwnerID:            ownerID,
		Active:             *a.Active,
	}
}

// ToOAPIAccount converts a domain Account to OpenAPI Account
func ToOAPIAccount(account domain.Account) *openapi.Account {
	var id string
	if account.ID != nil {
		id = *account.ID
	}

	return &openapi.Account{
		Id:                 id,
		Name:               account.Name,
		Type:               openapi.AccountType(account.Type),
		Currency:           account.Currency,
		InitialBalance:     account.InitialBalance,
		CurrentBalance:     account.CurrentBalance,
		Description:        account.Description,
		Institution:        account.Institution,
		AccountNumber:      account.AccountNumber,
		AccountInformation: account.AccountInformation,
		Owner:              ToOAPIHouseholdMember(account.Owner),
		Active:             &account.Active,
		CreatedAt:          account.CreatedAt,
		UpdatedAt:          account.UpdatedAt,
	}
}

// ToOAPIAccountList converts domain AccountList to OpenAPI AccountList
func ToOAPIAccountList(al *domain.AccountList) *openapi.AccountList {
	oapiAccounts := make(
		[]openapi.Account,
		0,
		len(al.Accounts),
	)
	for _, account := range al.Accounts {
		oapiAccounts = append(
			oapiAccounts,
			*ToOAPIAccount(account),
		)
	}
	return &openapi.AccountList{
		Accounts: &oapiAccounts,
		Metadata: &openapi.ListMetadata{
			Total:  al.Metadata.Total,
			Limit:  al.Metadata.Limit,
			Offset: al.Metadata.Offset,
		},
	}
}

func FromOAPIAccountListParams(params *openapi.ListAccountsParams) *domain.AccountListParams {
	return &domain.AccountListParams{
		Type:     params.Type,
		Currency: params.Currency,
		Active:   params.Active,
		Limit:    params.Limit,
		Offset:   params.Offset,
	}
}

func ToOAPIHouseholdMember(member *domain.HouseholdMember) *openapi.HouseholdMember {
	return &openapi.HouseholdMember{
		Id:        member.ID,
		FirstName: member.FirstName,
		LastName:  member.LastName,
		Nickname:  member.Nickname,
		Role:      member.Role,
		Active:    &member.Active,
		CreatedAt: member.CreatedAt,
		UpdatedAt: member.UpdatedAt,
	}
}

func FromOAPIHouseholdMember(member *openapi.HouseholdMember) *domain.HouseholdMember {
	return &domain.HouseholdMember{
		ID:        member.Id,
		FirstName: member.FirstName,
		LastName:  member.LastName,
		Nickname:  member.Nickname,
		Role:      member.Role,
		Active:    *member.Active,
		CreatedAt: member.CreatedAt,
		UpdatedAt: member.UpdatedAt,
	}
}

func FromOAPIHouseholdMemberRequest(member *openapi.HouseholdMemberRequest) *domain.HouseholdMember {
	return &domain.HouseholdMember{
		FirstName: member.FirstName,
		LastName:  member.LastName,
		Nickname:  member.Nickname,
		Role:      member.Role,
		Active:    *member.Active,
	}
}

func ToOAPIHouseholdMemberList(members *domain.HouseholdMemberList) *openapi.HouseholdMemberList {
	oapiMembers := make(
		[]openapi.HouseholdMember,
		0,
		len(members.HouseholdMembers),
	)
	for _, member := range members.HouseholdMembers {
		oapiMembers = append(
			oapiMembers,
			*ToOAPIHouseholdMember(&member),
		)
	}
	return &openapi.HouseholdMemberList{
		Members: &oapiMembers,
	}
}

func FromOAPIHouseholdMemberListParams(params *openapi.ListHouseholdMembersParams) *domain.HouseholdMemberListParams {
	return &domain.HouseholdMemberListParams{
		Role:   params.Role,
		Active: params.Active,
	}
}

func FromOAPIExpenditure(e *openapi.ExpenditureRequest) *domain.Transaction {
	return &domain.Transaction{
		AccountID:       e.AccountId,
		Amount:          e.Amount,
		Currency:        e.Currency,
		TransactionDate: e.Date.Time,
		Description:     e.Description,
		TransactionType: domain.TransactionTypeExpenditure,
	}
}

func FromOAPIExpenditureRequestTransaction(e *openapi.ExpenditureRequest) *domain.Transaction {
	if e.Date.Time.IsZero() {
		e.Date.Time = time.Now() // Set transaction date to current time if not provided in the request
	}
	return &domain.Transaction{
		AccountID:       e.AccountId,
		Amount:          e.Amount,
		Currency:        e.Currency,
		TransactionDate: e.Date.Time,
		Description:     e.Description,
		TransactionType: domain.TransactionTypeExpenditure,
	}
}

func FromOAPIIngressTransaction(i *openapi.Ingress) *domain.Transaction {
	return &domain.Transaction{
		AccountID:       i.AccountId,
		Amount:          i.Amount,
		Currency:        i.Currency,
		TransactionDate: i.Date.Time,
		Description:     *i.Description,
		TransactionType: domain.TransactionTypeIngress,
		CreatedAt:       *i.CreatedAt,
	}
}

func FromOAPIIngressRequestTransaction(i *openapi.IngressRequest) *domain.Transaction {
	return &domain.Transaction{
		AccountID:       i.AccountId,
		Amount:          i.Amount,
		Currency:        i.Currency,
		TransactionDate: i.Date.Time,
		Description:     *i.Description,
		TransactionType: domain.TransactionTypeIngress,
	}
}

func FromOAPITransferDebit(
	t *openapi.Transfer,
	sourceAccountCurrency string,
) *domain.Transaction {
	return &domain.Transaction{
		AccountID:       t.SourceAccountId,
		Amount:          *t.SourceAmount,
		Currency:        sourceAccountCurrency,
		TransactionDate: t.Date.Time,
		Description:     *t.Description,
		TransactionType: domain.TransactionTypeTransfer,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
}

func FromOAPITransferCredit(
	t *openapi.Transfer,
	destinationAccountCurrency string,
) *domain.Transaction {
	return &domain.Transaction{
		AccountID:       t.DestinationAccountId,
		Amount:          *t.DestinationAmount,
		Currency:        destinationAccountCurrency,
		TransactionDate: t.Date.Time,
		Description:     *t.Description,
		TransactionType: domain.TransactionTypeTransfer,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
}

func ToOAPICategory(category *domain.Category) *openapi.Category {
	return &openapi.Category{
		Id:              category.ID,
		Name:            category.Name,
		Description:     *category.Description,
		Active:          &category.Active,
		Color:           &category.Color,
		BackgroundColor: &category.BackgroundColor,
		CategoryType:    ToOAPICategoryCategoryType(&category.CategoryType),
	}
}

func ToOAPICategoryCategoryType(categoryType *domain.CategoryType) *openapi.CategoryType {
	var categoryTypeStr openapi.CategoryType
	switch *categoryType {
	case domain.CategoryTypeIngress:
		categoryTypeStr = openapi.CategoryTypeIngress
		break
	case domain.CategoryTypeExpenditure:
		categoryTypeStr = openapi.CategoryTypeExpenditure
		break
	case domain.CategoryTypeTransfer:
		categoryTypeStr = openapi.CategoryTypeTransfer
		break
	case domain.CategoryTypeSavingGoal:
		categoryTypeStr = openapi.CategoryTypeSavingGoal
		break
	default:
		break
	}
	return &categoryTypeStr
}

func ToOAPICategoryList(categories *[]domain.Category) *[]openapi.Category {
	if categories == nil {
		return nil
	}
	oapiCategories := make(
		[]openapi.Category,
		0,
		len(*categories),
	)
	for _, category := range *categories {
		oapiCategories = append(
			oapiCategories,
			*ToOAPICategory(&category),
		)
	}
	return &oapiCategories
}

func FromOAPITag(tag openapi.Tag) *domain.Tag {
	return &domain.Tag{
		ID:              tag.Id,
		Name:            tag.Name,
		Description:     tag.Description,
		Color:           tag.Color,
		BackgroundColor: tag.BackgroundColor,
		TagType:         *FromOAPITagType(&tag.TagType),
	}
}

func FromOAPITagType(tagType *openapi.TagType) *domain.TagType {
	if tagType == nil {
		return nil
	}
	var tagTypeDomain domain.TagType

	switch *tagType {
	case openapi.TagTypeIngress:
		tagTypeDomain = domain.TagTypeIngress
		break
	case openapi.TagTypeExpenditure:
		tagTypeDomain = domain.TagTypeExpenditure
		break
	case openapi.TagTypeTransfer:
		tagTypeDomain = domain.TagTypeTransfer
		break
	case openapi.TagTypeSavingGoal:
		tagTypeDomain = domain.TagTypeSavingsGoal
		break
	case openapi.TagTypeSavingsContribution:
		tagTypeDomain = domain.TagTypeSavingsContribution
		break
	case openapi.TagTypeSavingsWithdrawal:
		tagTypeDomain = domain.TagTypeSavingsWithdrawal
		break
	case openapi.TagTypeTransaction:
		break
	default:
		break
	}
	return &tagTypeDomain
}

func ToOAPITagType(tagType *domain.TagType) *openapi.TagType {
	if tagType == nil {
		return nil
	}
	var tagTypeStr openapi.TagType
	switch *tagType {
	case domain.TagTypeIngress:
		tagTypeStr = openapi.TagTypeIngress
		break
	case domain.TagTypeExpenditure:
		tagTypeStr = openapi.TagTypeExpenditure
		break
	case domain.TagTypeTransfer:
		tagTypeStr = openapi.TagTypeTransfer
		break
	case domain.TagTypeSavingsContribution:
		tagTypeStr = openapi.TagTypeSavingsContribution
		break
	case domain.TagTypeSavingsWithdrawal:
		tagTypeStr = openapi.TagTypeSavingsWithdrawal
		break
	case domain.TagTypeSavingsGoal:
		tagTypeStr = openapi.TagTypeSavingGoal
		break
	default:
		break
	}
	return &tagTypeStr
}

func FromOAPITagRequest(tagRequest openapi.TagRequest) *domain.Tag {
	return &domain.Tag{
		Name:            tagRequest.Name,
		Description:     tagRequest.Description,
		Color:           tagRequest.Color,
		BackgroundColor: tagRequest.BackgroundColor,
		TagType:         domain.TagType(tagRequest.TagType),
	}
}

func ToOAPITag(tag *domain.Tag) *openapi.Tag {
	return &openapi.Tag{
		Id:              tag.ID,
		Name:            tag.Name,
		Description:     tag.Description,
		Color:           tag.Color,
		BackgroundColor: tag.BackgroundColor,
		TagType:         *ToOAPITagType(&tag.TagType),
	}
}

func FromOAPIExpenditureRequest(e *openapi.ExpenditureRequest) *domain.Expenditure {
	var tagList []*domain.Tag
	if e.Tags != nil {
		tagList = make(
			[]*domain.Tag,
			0,
			len(*e.Tags),
		)
		for _, tag := range *e.Tags {
			tagList = append(
				tagList,
				FromOAPITag(tag),
			)
		}
	}

	var date time.Time
	if e.Date.Time.IsZero() {
		date = time.Now()
	} else {
		date = e.Date.Time
	}
	var declared bool
	if e.Declared == nil {
		declared = false
	} else {
		declared = *e.Declared
	}

	var planned bool
	if e.Planned == nil {
		planned = false
	} else {
		planned = *e.Planned
	}

	return &domain.Expenditure{
		Category:    FromOAPICategory(&e.Category),
		Declared:    declared,
		Planned:     planned,
		Transaction: FromOAPIExpenditureRequestTransaction(e),
		Tags:        &tagList,
		Date:        date,
	}
}

func ToOAPIExpenditure(e *domain.Expenditure) *openapi.Expenditure {
	var tagList []openapi.Tag
	if e.Tags != nil {
		tagList = make(
			[]openapi.Tag,
			0,
			len(*e.Tags),
		)
		for _, tag := range *e.Tags {
			tagList = append(
				tagList,
				*ToOAPITag(tag),
			)
		}
	}

	return &openapi.Expenditure{
		AccountId:   e.Transaction.AccountID,
		Amount:      e.Transaction.Amount,
		Category:    *ToOAPICategory(e.Category),
		CreatedAt:   e.Transaction.CreatedAt,
		Currency:    e.Transaction.Currency,
		Date:        openapi_types.Date{Time: e.Date},
		Declared:    &e.Declared,
		Description: e.Transaction.Description,
		Id:          e.ID,
		Planned:     &e.Planned,
		Tags:        &tagList,
		UpdatedAt:   e.Transaction.UpdatedAt,
	}
}

func FromOAPIExpenditureListParams(p *openapi.ListExpendituresParams) *domain.ExpenditureListParams {
	return &domain.ExpenditureListParams{
		CategoryID:  p.CategoryId,
		StartDate:   &p.StartDate.Time,
		EndDate:     &p.EndDate.Time,
		Declared:    p.Declared,
		Planned:     p.Planned,
		Currency:    p.Currency,
		Description: p.Description,
		AccountID:   p.AccountId,
		Tags:        p.Tags,
		Limit:       p.Limit,
		Offset:      p.Offset,
	}
}

func ToOAPIExpenditureList(p *domain.ExpenditureList) *openapi.ExpenditureList {
	var expenditures []openapi.Expenditure
	for _, e := range p.Expenditures {
		expenditures = append(
			expenditures,
			*ToOAPIExpenditure(&e),
		)
	}

	return &openapi.ExpenditureList{
		Metadata: &openapi.ListMetadata{
			Total:  p.Metadata.Total,
			Offset: p.Metadata.Offset,
			Limit:  p.Metadata.Limit,
		},
		Expenditures: &expenditures,
	}
}

func FromOAPICategoryType(c *openapi.CategoryType) *domain.CategoryType {
	if c == nil {
		return nil
	}
	var categoryType domain.CategoryType
	switch *c {
	case openapi.CategoryTypeExpenditure:
		categoryType = domain.CategoryTypeExpenditure
		return &categoryType
	case openapi.CategoryTypeIngress:
		categoryType = domain.CategoryTypeIngress
		return &categoryType
	case openapi.CategoryTypeTransfer:
		categoryType = domain.CategoryTypeTransfer
		return &categoryType
	case openapi.CategoryTypeSavingGoal:
		categoryType = domain.CategoryTypeSavingGoal
		return &categoryType
	default:
		return &categoryType
	}
}

func ToOAPICategoryType(c *domain.CategoryType) openapi.CategoryType {
	switch *c {
	case domain.CategoryTypeExpenditure:
		return openapi.CategoryTypeExpenditure
	case domain.CategoryTypeIngress:
		return openapi.CategoryTypeIngress
	case domain.CategoryTypeTransfer:
		return openapi.CategoryTypeTransfer
	case domain.CategoryTypeSavingGoal:
		return openapi.CategoryTypeSavingGoal
	}
	return openapi.CategoryTypeExpenditure
}

func FromOAPICategoryTypeType(c openapi.CategoryType) *domain.CategoryType {
	var categoryType domain.CategoryType
	switch c {
	case openapi.CategoryTypeExpenditure:
		categoryType = domain.CategoryTypeExpenditure
		break
	case openapi.CategoryTypeIngress:
		categoryType = domain.CategoryTypeIngress
		break
	case openapi.CategoryTypeTransfer:
		categoryType = domain.CategoryTypeTransfer
		break
	case openapi.CategoryTypeSavingGoal:
		categoryType = domain.CategoryTypeSavingGoal
		break
	}
	return &categoryType
}

func FromOAPICategoryRequest(
	c *openapi.CategoryRequest,
	id *string,
) *domain.Category {
	categoryID := ""
	if id != nil {
		categoryID = *id
	}
	return &domain.Category{
		ID:              categoryID,
		Name:            c.Name,
		Description:     &c.Description,
		Active:          true,
		Color:           *c.Color,
		BackgroundColor: *c.BackgroundColor,
		CategoryType:    domain.CategoryType(*c.CategoryType),
	}
}

func FromOAPICategory(c *openapi.Category) *domain.Category {
	category := FromOAPICategoryTypeType(*c.CategoryType)
	return &domain.Category{
		ID:              c.Id,
		Name:            c.Name,
		Description:     &c.Description,
		Active:          *c.Active,
		Color:           *c.Color,
		BackgroundColor: *c.BackgroundColor,
		CategoryType:    *category,
	}
}
