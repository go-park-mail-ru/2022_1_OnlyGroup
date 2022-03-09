package mock

import (
	"2022_1_OnlyGroup_back/app/models"
	"2022_1_OnlyGroup_back/pkg/errors"
	"github.com/bxcodec/faker/v3"
	"github.com/stretchr/testify/assert"
	"testing"
)

const TestNum = 6

func TestAddGetProfile(t *testing.T) {
	assert := assert.New(t)
	var ProfileMockTest ProfileMock
	var profileRepoTest []models.Profile
	for i := 0; i < TestNum; i++ {
		profileRepoTest = append(profileRepoTest, models.Profile{FirstName: faker.FirstName(), LastName: faker.LastName(), Birthday: faker.Date(), City: "Moscow", Interests: []string{faker.Word(), faker.Word()}, AboutUser: faker.Sentence(), UserId: i, Gender: faker.Gender()})
	}

	for idx, test := range profileRepoTest {
		t.Run("", func(t *testing.T) {
			//t.Parallel()
			err := ProfileMockTest.AddProfile(test)
			assert.Equal(err, nil)
			res, err := ProfileMockTest.GetProfile(idx)
			assert.Equal(res, test)
			assert.Equal(err, nil)
		})
	}
	_, err := ProfileMockTest.GetProfile(TestNum + 1)
	assert.Equal(err, errors.ErrProfileNotFound)

}

func TestAddGetShortProfile(t *testing.T) {
	assert := assert.New(t)
	var ProfileMockTest ProfileMock
	var profileRepoTest []models.Profile
	for i := 0; i < TestNum; i++ {
		profileRepoTest = append(profileRepoTest, models.Profile{FirstName: faker.FirstName(), LastName: faker.LastName(), Birthday: faker.Date(), City: "Moscow", Interests: []string{faker.Word(), faker.Word()}, AboutUser: faker.Sentence(), UserId: i, Gender: faker.Gender()})
	}

	for idx, test := range profileRepoTest {
		t.Run("", func(t *testing.T) {
			//t.Parallel()
			err := ProfileMockTest.AddProfile(test)
			assert.Equal(err, nil)
			res, err := ProfileMockTest.GetShortProfile(idx)
			assert.Equal(res.City, test.City)
			assert.Equal(res.LastName, test.LastName)
			assert.Equal(res.FirstName, test.FirstName)
			assert.Equal(err, nil)
		})
	}
	_, err := ProfileMockTest.GetProfile(TestNum + 1)
	assert.Equal(err, errors.ErrProfileNotFound)
}

func TestAddGetEmptyProfile(t *testing.T) {
	assert := assert.New(t)
	var ProfileMockTest ProfileMock

	for idx := 0; idx < TestNum; idx++ {
		t.Run("", func(t *testing.T) {
			//t.Parallel()
			err := ProfileMockTest.AddEmptyProfile(idx)
			assert.Equal(err, nil)

			res, err := ProfileMockTest.GetProfile(idx)
			assert.Equal(res.City, "")
			assert.Equal(res.LastName, "")
			assert.Equal(res.FirstName, "")
			assert.Equal(res.Birthday, "")
			assert.Equal(res.AboutUser, "")
			assert.Equal(res.UserId, idx)
			assert.Equal(res.Gender, "")
			assert.Equal(res.Interests, []string{})
			assert.Equal(err, nil)
		})
	}
}

func TestFindCandidate(t *testing.T) {
	assert := assert.New(t)
	var ProfileMockTest ProfileMock
	var profileRepoTest []models.Profile
	for i := 0; i < TestNum; i++ {
		profileRepoTest = append(profileRepoTest, models.Profile{FirstName: faker.FirstName(), LastName: faker.LastName(), Birthday: faker.Date(), City: "Moscow", Interests: []string{faker.Word(), faker.Word()}, AboutUser: faker.Sentence(), UserId: i, Gender: faker.Gender()})
	}
	_, err := ProfileMockTest.FindCandidateProfile(TestNum + 1)
	assert.Equal(err, errors.ErrProfileNotFound)

	for _, test := range profileRepoTest {
		t.Run("", func(t *testing.T) {
			//t.Parallel()
			err := ProfileMockTest.AddProfile(test)
			assert.Equal(err, nil)
		})
	}

	for i := 0; i < TestNum; i++ {
		candidateProfiles, err := ProfileMockTest.FindCandidateProfile(i)
		assert.Equal(err, nil)
		if len(candidateProfiles.Candidates) != 3 {
			t.Error("len candidate < 3")
		}
		for _, value := range candidateProfiles.Candidates {
			_, err := ProfileMockTest.GetProfile(value)
			assert.Equal(err, nil)
		}
	}
}

func TestCheckEmptyProfile(t *testing.T) {
	assert := assert.New(t)
	var ProfileMockTest ProfileMock
	var profileRepoTest []models.Profile

	for i := 0; i < TestNum; i++ {
		profileRepoTest = append(profileRepoTest, models.Profile{FirstName: faker.FirstName(), LastName: faker.LastName(), Birthday: faker.Date(), City: "Moscow", Interests: []string{faker.Word(), faker.Word()}, AboutUser: faker.Sentence(), UserId: i, Gender: faker.Gender()})
	}
	for _, test := range profileRepoTest {
		t.Run("", func(t *testing.T) {
			//t.Parallel()
			err := ProfileMockTest.AddProfile(test)
			assert.Equal(err, nil)
		})
	}
	err := ProfileMockTest.AddEmptyProfile(TestNum + 1)
	assert.Equal(err, nil)
	err = ProfileMockTest.CheckProfileFiled(TestNum + 1)
	assert.Equal(err, errors.ErrProfileNotFiled)

	for idx := 0; idx < TestNum; idx++ {
		t.Run("", func(t *testing.T) {
			//t.Parallel()
			err := ProfileMockTest.AddEmptyProfile(idx)
			assert.Equal(err, nil)

			err = ProfileMockTest.CheckProfileFiled(idx)
			assert.Equal(err, nil)
		})
	}
	err = ProfileMockTest.CheckProfileFiled(TestNum + 2)
	assert.Equal(err, errors.ErrProfileNotFound)
}

func TestNewChangeProfile(t *testing.T) {
	assert := assert.New(t)

	ProfileMockTest := NewProfileMock()
	profileRepoTest := models.Profile{FirstName: faker.FirstName(), LastName: faker.LastName(), Birthday: faker.Date(), City: "Moscow", Interests: []string{faker.Word(), faker.Word()}, AboutUser: faker.Sentence(), UserId: 0, Gender: faker.Gender()}

	err := ProfileMockTest.ChangeProfile(0, profileRepoTest)
	assert.Equal(err, nil)

	err = ProfileMockTest.ChangeProfile(8, profileRepoTest)
	assert.Equal(err, errors.ErrProfileNotFound)

}

func TestDeleteChangeProfile(t *testing.T) {
	assert := assert.New(t)

	ProfileMockTest := NewProfileMock()

	err := ProfileMockTest.DeleteProfile(5)
	assert.Equal(err, nil)

	err = ProfileMockTest.DeleteProfile(8)
	assert.Equal(err, errors.ErrProfileNotFound)

}
