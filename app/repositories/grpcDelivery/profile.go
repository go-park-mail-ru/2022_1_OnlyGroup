package grpcDelivery

import (
	"2022_1_OnlyGroup_back/app/handlers/http"
	"2022_1_OnlyGroup_back/app/models"
	"2022_1_OnlyGroup_back/microservices/profile/proto"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProfileGrpc struct {
	delivery proto.ProfileRepositoryClient
}

func NewProfileGrpc(delivery proto.ProfileRepositoryClient) *ProfileGrpc {
	return &ProfileGrpc{delivery: delivery}
}

func (repo *ProfileGrpc) Get(profileId int) (models.Profile, error) {
	msg, err := repo.delivery.Get(context.Background(), models.ModelUserIdToGRPC(profileId))
	if status.Code(err) == codes.Unavailable {
		return models.Profile{}, http.ErrServiceUnavailable.Wrap(err, "failed get profile Grpc")
	}
	if err != nil {
		return models.Profile{}, http.ErrBaseApp.Wrap(err, "failed get profile Grpc")
	}
	return models.GRPCToModelProfile(msg), nil
}

func (repo *ProfileGrpc) GetShort(profileId int) (models.ShortProfile, error) {
	msg, err := repo.delivery.GetShort(context.Background(), models.ModelUserIdToGRPC(profileId))
	if status.Code(err) == codes.Unavailable {
		return models.ShortProfile{}, http.ErrServiceUnavailable.Wrap(err, "failed get shortProfile Grpc")
	}
	if err != nil {
		return models.ShortProfile{}, http.ErrBaseApp.Wrap(err, "failed get shortProfile Grpc")
	}
	return models.GRPCToModelShortProfile(msg), nil
}

func (repo *ProfileGrpc) Change(profileId int, profile models.Profile) error {
	_, err := repo.delivery.Change(context.Background(), models.ModelProfileToGRPC(&profile))
	if status.Code(err) == codes.Unavailable {
		return http.ErrServiceUnavailable.Wrap(err, "failed change profile Grpc")
	}
	if err != nil {
		return http.ErrBaseApp.Wrap(err, "failed change profile Grpc")
	}
	return nil
}

func (repo *ProfileGrpc) Delete(profileId int) error {
	_, err := repo.delivery.Delete(context.Background(), models.ModelUserIdToGRPC(profileId))
	if status.Code(err) == codes.Unavailable {
		return http.ErrServiceUnavailable.Wrap(err, "failed delete profile Grpc")
	}
	if err != nil {
		return http.ErrBaseApp.Wrap(err, "failed delete profile Grpc")
	}
	return nil
}

func (repo *ProfileGrpc) Add(profile models.Profile) error {
	_, err := repo.delivery.Add(context.Background(), models.ModelProfileToGRPC(&profile))
	if status.Code(err) == codes.Unavailable {
		return http.ErrServiceUnavailable.Wrap(err, "failed add profile Grpc")
	}
	if err != nil {
		return http.ErrBaseApp.Wrap(err, "failed add profile Grpc")
	}
	return nil
}

func (repo *ProfileGrpc) AddEmpty(profileId int) error {
	_, err := repo.delivery.AddEmpty(context.Background(), models.ModelUserIdToGRPC(profileId))
	if status.Code(err) == codes.Unavailable {
		return http.ErrServiceUnavailable.Wrap(err, "failed add emptyProfile Grpc")
	}
	if err != nil {
		return http.ErrBaseApp.Wrap(err, "failed add emptyProfile Grpc")
	}
	return nil
}

func (repo *ProfileGrpc) FindCandidate(profileId int) (models.VectorCandidate, error) {
	msg, err := repo.delivery.FindCandidate(context.Background(), models.ModelUserIdToGRPC(profileId))
	if status.Code(err) == codes.Unavailable {
		return models.VectorCandidate{}, http.ErrServiceUnavailable.Wrap(err, "failed findCandidate Grpc")
	}
	if err != nil {
		return models.VectorCandidate{}, http.ErrBaseApp.Wrap(err, "failed findCandidate Grpc")
	}
	return models.GRPCToModelVectorCandidate(msg), nil
}

func (repo *ProfileGrpc) CheckFiled(profileId int) error {
	_, err := repo.delivery.CheckFiled(context.Background(), models.ModelUserIdToGRPC(profileId))
	if status.Code(err) == codes.Unavailable {
		return http.ErrServiceUnavailable.Wrap(err, "failed checkFiled Grpc")
	}
	if err != nil {
		return http.ErrBaseApp.Wrap(err, "failed checkFiled Grpc")
	}
	return nil
}

func (repo *ProfileGrpc) GetFilters(userId int) (models.Filters, error) {
	filtersModel, err := repo.delivery.GetFilters(context.Background(), models.ModelUserIdToGRPC(userId))
	if status.Code(err) == codes.Unavailable {
		return models.Filters{}, http.ErrServiceUnavailable.Wrap(err, "failed getFilters Grpc")
	}
	if err != nil {
		return models.Filters{}, http.ErrBaseApp.Wrap(err, "failed getFilters Grpc")
	}
	return models.GRPCToModelFilters(filtersModel), nil
}

func (repo *ProfileGrpc) ChangeFilters(userId int, filters models.Filters) error {
	_, err := repo.delivery.ChangeFilters(context.Background(), models.ModelFiltersToGRPC(userId, &filters))
	if status.Code(err) == codes.Unavailable {
		return http.ErrServiceUnavailable.Wrap(err, "failed changeFilters Grpc")
	}
	if err != nil {
		return http.ErrBaseApp.Wrap(err, "failed changeFilters Grpc")
	}
	return nil
}

func (repo *ProfileGrpc) GetInterests() ([]models.Interest, error) {
	msg, err := repo.delivery.GetInterests(context.Background(), &proto.Nothing{})
	if status.Code(err) == codes.Unavailable {
		return nil, http.ErrServiceUnavailable.Wrap(err, "failed getInterests Grpc")
	}
	if err != nil {
		return nil, http.ErrBaseApp.Wrap(err, "failed getInterests Grpc")
	}
	return models.GRPCToModelInterests(msg), nil
}

func (repo *ProfileGrpc) GetDynamicInterest(interest string) ([]models.Interest, error) {
	msg, err := repo.delivery.GetDynamicInterest(context.Background(), models.ModelStrToGRPC(interest))
	if status.Code(err) == codes.Unavailable {
		return nil, http.ErrServiceUnavailable.Wrap(err, "failed getDynamicInterest Grpc")
	}
	if err != nil {
		return nil, http.ErrBaseApp.Wrap(err, "failed getDynamicInterest Grpc")
	}
	return models.GRPCToModelInterests(msg), nil
}

func (repo *ProfileGrpc) CheckInterests(interests []models.Interest) error {
	_, err := repo.delivery.CheckInterests(context.Background(), models.ModelInterestsToGRPC(interests))
	if status.Code(err) == codes.Unavailable {
		return http.ErrServiceUnavailable.Wrap(err, "failed checkInterests Grpc")
	}
	if status.Code(err) == codes.Internal {
		return http.ErrBaseApp.Wrap(err, "failed checkInterests Grpc")
	}
	if status.Code(err) == codes.InvalidArgument {
		return http.ErrBadRequest
	}
	return nil
}

func (repo *ProfileGrpc) SetAction(profileId int, likes models.Likes) error {
	_, err := repo.delivery.SetAction(context.Background(), models.ModelLikesToGRPC(profileId, likes))
	if status.Code(err) == codes.Unavailable {
		return http.ErrServiceUnavailable.Wrap(err, "failed setAction Grpc")
	}
	if status.Code(err) == codes.Internal {
		return http.ErrBaseApp.Wrap(err, "failed setAction Grpc")
	}
	if status.Code(err) == codes.InvalidArgument {
		return http.ErrBadRequest
	}
	return err
}

func (repo *ProfileGrpc) GetMatched(profileId int) (models.LikesMatched, error) {
	msg, err := repo.delivery.GetMatched(context.Background(), models.ModelUserIdToGRPC(profileId))
	if status.Code(err) == codes.Unavailable {
		return models.LikesMatched{}, http.ErrServiceUnavailable.Wrap(err, "failed getMatched Grpc")
	}
	if err != nil {
		return models.LikesMatched{}, http.ErrBaseApp.Wrap(err, "failed getMatched Grpc")
	}
	return models.GRPCToModelLikesMatched(msg), nil
}
