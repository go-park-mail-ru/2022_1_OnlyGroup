package models

import "2022_1_OnlyGroup_back/microservices/profile/proto"

type Filters struct {
	AgeFilter    [2]int `json:",omitempty" validate:"ageFilter"`
	GenderFilter int    `json:",omitempty" validate:"min=0, max=1"`
	HeightFilter [2]int `json:",omitempty" validate:"heightFilter"`
}

func ModelFiltersToGRPC(userId int, model *Filters) *proto.Filters {
	var ageFilter []int64
	for _, val := range model.AgeFilter {
		ageFilter = append(ageFilter, int64(val))
	}
	var heightFilter []int64
	for _, val := range model.HeightFilter {
		heightFilter = append(heightFilter, int64(val))
	}

	return &proto.Filters{AgeFilter: ageFilter, HeightFilter: heightFilter, GenderFilter: int64(model.GenderFilter), Id: int64(userId)}
}

func GRPCToModelFilters(grpcModel *proto.Filters) Filters {
	var ageFilter [2]int
	for idx, val := range grpcModel.AgeFilter {
		ageFilter[idx] = int(val)
	}
	var heightFilter [2]int
	for idx, val := range grpcModel.HeightFilter {
		heightFilter[idx] = int(val)
	}

	return Filters{GenderFilter: int(grpcModel.GenderFilter), HeightFilter: heightFilter, AgeFilter: ageFilter}
}
