package usecase

import (
	"context"
	"errors"
	"ghorkov32/proletariat-budget-be/internal/core/domain"
	"ghorkov32/proletariat-budget-be/internal/core/port"
)

type TagsUseCase struct {
	tagsRepo *port.TagsRepo
}

func NewTagsUseCase(tagsRepo *port.TagsRepo) *TagsUseCase {
	return &TagsUseCase{tagsRepo: tagsRepo}
}

func (u *TagsUseCase) ListTags(
	ctx context.Context,
	tagType *domain.TagType,
) (
	[]*domain.Tag,
	error,
) {
	if tagType == nil {
		tags, err := (*u.tagsRepo).List(ctx)
		if err != nil {
			return nil, err
		}
		return *tags, nil
	}
	tags, err := (*u.tagsRepo).ListByType(
		ctx,
		*tagType,
		nil,
	)
	if err != nil {
		return nil, err
	}
	return *tags, nil
}

func (u *TagsUseCase) CreateTag(
	ctx context.Context,
	tag *domain.Tag,
) (
	*domain.Tag,
	error,
) {
	if err := tag.Validate(); err != nil {
		return nil, err
	}
	_, err := (*u.tagsRepo).GetByNameAndType(
		ctx,
		(*tag).Name,
		(*tag).TagType,
	)
	if err == nil {
		return nil, domain.ErrTagAlreadyExists
	}
	id, err := (*u.tagsRepo).Create(
		ctx,
		*tag,
	)
	if err != nil {
		return nil, err
	}
	tag.ID = id
	return tag, nil
}

func (u *TagsUseCase) DeleteTag(
	ctx context.Context,
	id string,
) error {
	err := (*u.tagsRepo).Delete(
		ctx,
		id,
	)
	if err != nil {
		if errors.Is(
			err,
			port.ErrRecordNotFound,
		) {
			return domain.ErrTagNotFound
		}
	}
	return nil
}
