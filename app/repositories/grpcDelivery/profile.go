package grpcDelivery

import (
	"2022_1_OnlyGroup_back/app/handlers/http"
	"2022_1_OnlyGroup_back/app/models"
	"2022_1_OnlyGroup_back/microservices/profile/proto"
	"context"
)

type ProfileGrpc struct {
	delivery proto.ProfileRepositoryClient
}

func NewProfileGrpc(delivery proto.ProfileRepositoryClient) *ProfileGrpc {
	return &ProfileGrpc{delivery: delivery}
}

func (repo *ProfileGrpc) Get(profileId int) (models.Profile, error) {
	msg, err := repo.delivery.Get(context.Background(), models.ModelUserIdToGRPC(profileId))
	if err != nil {
		return models.Profile{}, http.AppErrorFromGRPC(err)
	}
	return models.GRPCToModelProfile(msg), nil
}

func (repo *ProfileGrpc) GetShort(profileId int) (models.ShortProfile, error) {
	msg, err := repo.delivery.GetShort(context.Background(), models.ModelUserIdToGRPC(profileId))
	if err != nil {
		return models.ShortProfile{}, http.AppErrorFromGRPC(err)
	}
	return models.GRPCToModelShortProfile(msg), nil
}

func (repo *ProfileGrpc) Change(profileId int, profile models.Profile) error {
	_, err := repo.delivery.Change(context.Background(), models.ModelProfileToGRPC(&profile))
	if err != nil {
		return http.AppErrorFromGRPC(err)
	}
	return nil
}

func (repo *ProfileGrpc) Delete(profileId int) error {
	_, err := repo.delivery.Delete(context.Background(), models.ModelUserIdToGRPC(profileId))
	if err != nil {
		return http.AppErrorFromGRPC(err)
	}
	return nil
}

func (repo *ProfileGrpc) Add(profile models.Profile) error {
	_, err := repo.delivery.Add(context.Background(), models.ModelProfileToGRPC(&profile))
	if err != nil {
		return http.AppErrorFromGRPC(err)
	}
	return nil
}

func (repo *ProfileGrpc) AddEmpty(profileId int) error {
	_, err := repo.delivery.AddEmpty(context.Background(), models.ModelUserIdToGRPC(profileId))
	if err != nil {
		return http.AppErrorFromGRPC(err)
	}
	return nil
}

func (repo *ProfileGrpc) FindCandidate(profileId int) (models.VectorCandidate, error) {
	msg, err := repo.delivery.FindCandidate(context.Background(), models.ModelUserIdToGRPC(profileId))
	if err != nil {
		return models.VectorCandidate{}, http.AppErrorFromGRPC(err)
	}
	return models.GRPCToModelVectorCandidate(msg), nil
}

func (repo *ProfileGrpc) CheckFiled(profileId int) error {
	_, err := repo.delivery.CheckFiled(context.Background(), models.ModelUserIdToGRPC(profileId))
	if err != nil {
		return http.AppErrorFromGRPC(err)
	}
	return nil
}

func (repo *ProfileGrpc) GetFilters(userId int) (models.Filters, error) {
	filtersModel, err := repo.delivery.GetFilters(context.Background(), models.ModelUserIdToGRPC(userId))
	if err != nil {
		return models.Filters{}, http.AppErrorFromGRPC(err)
	}
	return models.GRPCToModelFilters(filtersModel), nil
}

func (repo *ProfileGrpc) ChangeFilters(userId int, filters models.Filters) error {
	_, err := repo.delivery.ChangeFilters(context.Background(), models.ModelFiltersToGRPC(userId, &filters))
	if err != nil {
		return http.AppErrorFromGRPC(err)
	}
	return nil
}

func (repo *ProfileGrpc) GetInterests() ([]models.Interest, error) {
	msg, err := repo.delivery.GetInterests(context.Background(), &proto.Nothing{})
	if err != nil {
		return nil, http.AppErrorFromGRPC(err)
	}
	return models.GRPCToModelInterests(msg), nil
}

func (repo *ProfileGrpc) GetDynamicInterest(interest string) ([]models.Interest, error) {
	msg, err := repo.delivery.GetDynamicInterest(context.Background(), models.ModelStrToGRPC(interest))
	if err != nil {
		return nil, http.AppErrorFromGRPC(err)
	}
	return models.GRPCToModelInterests(msg), nil
}

func (repo *ProfileGrpc) CheckInterests(interests []models.Interest) error {
	_, err := repo.delivery.CheckInterests(context.Background(), models.ModelInterestsToGRPC(interests))
	if err != nil {
		return http.AppErrorFromGRPC(err)
	}
	return nil
}

func (repo *ProfileGrpc) SetAction(profileId int, likes models.Likes) error {
	_, err := repo.delivery.SetAction(context.Background(), models.ModelLikesToGRPC(profileId, likes))
	if err != nil {
		return http.AppErrorFromGRPC(err)
	}
	return err
}

func (repo *ProfileGrpc) GetMatched(profileId int) (models.LikesMatched, error) {
	msg, err := repo.delivery.GetMatched(context.Background(), models.ModelUserIdToGRPC(profileId))
	if err != nil {
		return models.LikesMatched{}, http.AppErrorFromGRPC(err)
	}
	return models.GRPCToModelLikesMatched(msg), nil
}
