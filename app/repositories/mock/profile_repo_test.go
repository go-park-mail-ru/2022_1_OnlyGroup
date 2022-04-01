package mock

import (
	"2022_1_OnlyGroup_back/app/handlers"
	"2022_1_OnlyGroup_back/app/models"
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
			err := ProfileMockTest.Add(test)
			assert.Equal(nil, err)
			res, err := ProfileMockTest.Get(idx)
			assert.Equal(res, test)
			assert.Equal(nil, err)
		})
	}
	_, err := ProfileMockTest.Get(TestNum + 1)
	assert.Equal(err, handlers.ErrProfileNotFound)

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
			err := ProfileMockTest.Add(test)
			assert.Equal(nil, err)
			res, err := ProfileMockTest.GetShort(idx)
			expectRes := models.ShortProfile{FirstName: test.FirstName, LastName: test.LastName, City: test.City}
			assert.Equal(expectRes, res)
			assert.Equal(nil, err)
		})
	}
	_, err := ProfileMockTest.Get(TestNum + 1)
	assert.Equal(err, handlers.ErrProfileNotFound)
}

func TestAddGetEmptyProfile(t *testing.T) {
	assert := assert.New(t)
	var ProfileMockTest ProfileMock

	for idx := 0; idx < TestNum; idx++ {
		t.Run("", func(t *testing.T) {
			//t.Parallel()
			err := ProfileMockTest.AddEmpty(idx)
			assert.Equal(nil, err)

			res, err := ProfileMockTest.Get(idx)
			assert.Equal(models.Profile{Interests: []string{}, UserId: idx}, res)
			assert.Equal(nil, err)
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
	_, err := ProfileMockTest.FindCandidate(TestNum + 1)
	assert.Equal(err, handlers.ErrProfileNotFound)

	for _, test := range profileRepoTest {
		t.Run("", func(t *testing.T) {
			//t.Parallel()
			err := ProfileMockTest.Add(test)
			assert.Equal(nil, err)
		})
	}

	for i := 0; i < TestNum; i++ {
		candidateProfiles, err := ProfileMockTest.FindCandidate(i)
		assert.Equal(nil, err)
		if len(candidateProfiles.Candidates) != 3 {
			t.Error("len candidate < 3")
		}
		for _, value := range candidateProfiles.Candidates {
			_, err := ProfileMockTest.Get(value)
			assert.Equal(nil, err)
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
			err := ProfileMockTest.Add(test)
			assert.Equal(nil, err)
		})
	}
	err := ProfileMockTest.AddEmpty(TestNum + 1)
	assert.Equal(nil, err)
	err = ProfileMockTest.CheckFiled(TestNum + 1)
	assert.Equal(err, handlers.ErrProfileNotFiled)

	for idx := 0; idx < TestNum; idx++ {
		t.Run("", func(t *testing.T) {
			//t.Parallel()
			err := ProfileMockTest.AddEmpty(idx)
			assert.Equal(nil, err)

			err = ProfileMockTest.CheckFiled(idx)
			assert.Equal(nil, err)
		})
	}
	err = ProfileMockTest.CheckFiled(TestNum + 2)
	assert.Equal(err, handlers.ErrProfileNotFound)
}

func TestNewChangeProfile(t *testing.T) {
	assert := assert.New(t)

	ProfileMockTest := NewProfileMock()
	profileRepoTest := models.Profile{FirstName: faker.FirstName(), LastName: faker.LastName(), Birthday: faker.Date(), City: "Moscow", Interests: []string{faker.Word(), faker.Word()}, AboutUser: faker.Sentence(), UserId: 0, Gender: faker.Gender()}

	err := ProfileMockTest.Change(0, profileRepoTest)
	assert.Equal(nil, err)

	err = ProfileMockTest.Change(8, profileRepoTest)
	assert.Equal(err, handlers.ErrProfileNotFound)

}

func TestDeleteChangeProfile(t *testing.T) {
	assert := assert.New(t)

	ProfileMockTest := NewProfileMock()

	err := ProfileMockTest.Delete(5)
	assert.Equal(nil, err)

	err = ProfileMockTest.Delete(8)
	assert.Equal(err, handlers.ErrProfileNotFound)

}
