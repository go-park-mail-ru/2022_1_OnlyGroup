package grpcDelivery

import (
	"2022_1_OnlyGroup_back/app/handlers/http"
	"2022_1_OnlyGroup_back/app/models"
	"2022_1_OnlyGroup_back/microservices/profile"
	"2022_1_OnlyGroup_back/microservices/profile/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProfileHandler struct {
	useCase profile.ProfileGRPCUseCases
}

func NewProfileHandler(useCase profile.ProfileGRPCUseCases) *ProfileHandler {
	return &ProfileHandler{useCase: useCase}
}

func (handler *ProfileHandler) Get(ctx context.Context, userId *proto.UserID) (*proto.Profile, error) {
	profileModel, err := handler.useCase.Get(int(userId.Id))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Profile service error: Get")
	}
	return models.ModelProfileToGRPC(&profileModel), err
}
func (handler *ProfileHandler) GetShort(ctx context.Context, userId *proto.UserID) (*proto.ShortProfile, error) {
	shortProfileModel, err := handler.useCase.GetShort(int(userId.Id))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Profile service error: GetShort")
	}
	return models.ModelShortProfileToGRPC(&shortProfileModel), nil
}
func (handler *ProfileHandler) Change(ctx context.Context, model *proto.Profile) (*proto.Nothing, error) {
	err := handler.useCase.Change(int(model.UserId), models.GRPCToModelProfile(model))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Profile service error: Change")
	}
	return &proto.Nothing{}, nil
}
func (handler *ProfileHandler) Delete(ctx context.Context, model *proto.UserID) (*proto.Nothing, error) {
	err := handler.useCase.Delete(int(model.Id))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Profile service error: change")
	}
	return &proto.Nothing{}, nil
}
func (handler *ProfileHandler) Add(context.Context, *proto.Profile) (*proto.Nothing, error) {
	return &proto.Nothing{}, status.Errorf(codes.Unimplemented, "method Add not implemented")
}
func (handler *ProfileHandler) CheckFiled(context.Context, *proto.UserID) (*proto.Nothing, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckFiled not implemented")
}
func (handler *ProfileHandler) AddEmpty(ctx context.Context, model *proto.UserID) (*proto.Nothing, error) {
	err := handler.useCase.AddEmpty(int(model.Id))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Profile service error: AddEmpty")
	}
	return &proto.Nothing{}, nil
}
func (handler *ProfileHandler) FindCandidate(ctx context.Context, model *proto.UserID) (*proto.VectorCandidate, error) {
	modelCandidates, err := handler.useCase.GetCandidates(int(model.Id))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Profile service error: FindCandidate")
	}
	return models.ModelVectorCandidateToGRPC(&modelCandidates), nil
}
func (handler *ProfileHandler) GetFilters(ctx context.Context, model *proto.UserID) (*proto.Filters, error) {
	modelFilters, err := handler.useCase.GetFilters(int(model.Id))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Profile service error: GetFilters")
	}
	return models.ModelFiltersToGRPC(int(model.Id), &modelFilters), nil
}
func (handler *ProfileHandler) ChangeFilters(ctx context.Context, model *proto.Filters) (*proto.Nothing, error) {
	err := handler.useCase.ChangeFilters(int(model.Id), models.GRPCToModelFilters(model))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Profile service error: ChangeFilters")
	}
	return &proto.Nothing{}, nil
}
func (handler *ProfileHandler) GetInterests(ctx context.Context, model *proto.Nothing) (*proto.Interests, error) {
	modelInterest, err := handler.useCase.GetInterest()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Profile service error: GetInterests")
	}
	return models.ModelInterestsToGRPC(modelInterest), nil
}
func (handler *ProfileHandler) GetDynamicInterest(ctx context.Context, model *proto.StrInterest) (*proto.Interests, error) {
	modelInterest, err := handler.useCase.GetDynamicInterests(model.StrInterest)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Profile service error: GetDynamicInterest")
	}
	return models.ModelInterestsToGRPC(modelInterest), nil
}
func (handler *ProfileHandler) SetAction(ctx context.Context, model *proto.Likes) (*proto.Nothing, error) {
	err := handler.useCase.SetAction(int(model.WhoId), models.GRPCToModelLikes(model))
	if err == http.ErrBadRequest {
		return nil, status.Errorf(codes.InvalidArgument, "Bad request")
	}
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Profile service error: SetAction")
	}
	return &proto.Nothing{}, nil
}
func (handler *ProfileHandler) GetMatched(ctx context.Context, model *proto.UserID) (*proto.LikesMatched, error) {
	likesVector, err := handler.useCase.GetMatched(int(model.Id))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Profile service error: GetMatched")
	}
	return models.ModelLikesMatchedToGRPC(likesVector), nil
}

func (handler *ProfileHandler) CheckInterests(ctx context.Context, model *proto.Interests) (*proto.Nothing, error) {
	err := handler.useCase.CheckInterests(models.GRPCToModelInterests(model))
	if err == http.ErrBadRequest {
		return nil, status.Errorf(codes.InvalidArgument, "Bad request")
	}
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Profile service error: CheckInterest")
	}

	return &proto.Nothing{}, nil
}
