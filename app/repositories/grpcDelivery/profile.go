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

func (repo *ProfileGrpc) Get(profileId int) (profile models.Profile, err error) {
	deliveryProfile, err := repo.delivery.Get(context.Background(), &proto.UserID{Id: int64(profileId)})
	if err != nil {
		return profile, http.ErrBaseApp.Wrap(err, "failed get profile Grpc")
	}
	profile = models.GRPCToModelProfile(deliveryProfile)
	return profile, nil
}

func (repo *ProfileGrpc) GetShort(profileId int) (shortProfile models.ShortProfile, err error) {
	deliveryModel, err := repo.delivery.GetShort(context.Background(), &proto.UserID{Id: int64(profileId)})
	if err != nil {
		return shortProfile, http.ErrBaseApp.Wrap(err, "failed get profile Grpc")
	}
	shortProfile = models.GRPCToModelShortProfile(deliveryModel)
	return shortProfile, nil
}

func (repo *ProfileGrpc) Change(profileId int, profile models.Profile) (err error) {
	nothing, err := repo.delivery.Change(context.Background(), models.ModelProfileToGRPC(&profile))
	if err != nil {
		return http.ErrBaseApp.Wrap(err, "failed get profile Grpc")
	}
	return
}

func (repo *ProfileGrpc) Delete(profileId int) (err error) {
	return
}

func (repo *ProfileGrpc) Add(profile models.Profile) (err error) {

	return
}

func (repo *ProfileGrpc) AddEmpty(profileId int) (err error) {
	return
}

func (repo *ProfileGrpc) FindCandidate(profileId int) (candidateProfiles models.VectorCandidate, err error) {

	return
}

func (repo *ProfileGrpc) CheckFiled(profileId int) (err error) {
	return
}

func (repo *ProfileGrpc) GetFilters(userId int) (models.Filters, error) {
	var filters models.Filters

	return filters, nil
}

func (repo *ProfileGrpc) ChangeFilters(userId int, filters models.Filters) error {
	return nil
}

func (repo *ProfileGrpc) GetInterests() ([]models.Interest, error) {
	var interests []models.Interest

	return interests, nil
}

func (repo *ProfileGrpc) GetDynamicInterest(interest string) ([]models.Interest, error) {
	var interests []models.Interest
	return interests, nil
}

func (repo *ProfileGrpc) CheckInterests(interests []models.Interest) error {

	return nil
}

func (repo *ProfileGrpc) SetAction(profileId int, likes models.Likes) (err error) {

	return
}

func (repo *ProfileGrpc) GetMatched(profileId int) (likesVector models.LikesMatched, err error) {
	return
}
