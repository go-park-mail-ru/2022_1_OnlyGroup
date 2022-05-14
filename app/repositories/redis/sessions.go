package redis

import (
	"2022_1_OnlyGroup_back/app/handlers/http"
	"2022_1_OnlyGroup_back/pkg/randomGenerator"
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"strconv"
	"strings"
)

const sessionSecretSize = 64
const sessionIdAndSecretSep = "_"
const sessionSplitNumber = 2

type redisSessionsRepo struct {
	client          *redis.Client
	sessionsPrefix  string
	secretGenerator randomGenerator.RandomGenerator
}

func NewRedisSessionRepository(client *redis.Client, sessionsPrefix string, generator randomGenerator.RandomGenerator) *redisSessionsRepo {
	return &redisSessionsRepo{client: client, sessionsPrefix: sessionsPrefix, secretGenerator: generator}
}

func (repo *redisSessionsRepo) addSessionInternal(id int, additionalData string, generatedSecret string) (string, error) {
	idInString := strconv.Itoa(id)
	secret := strings.Join([]string{idInString, generatedSecret}, sessionIdAndSecretSep)

	key := strings.Join([]string{repo.sessionsPrefix, idInString}, sessionIdAndSecretSep)
	_, err := repo.client.HSet(context.Background(), key, secret, additionalData).Result()
	if err != nil {
		return "", http.ErrBaseApp.Wrap(err, "redis session add failed")
	}
	return secret, nil
}

func (repo *redisSessionsRepo) Add(id int, additionalData string) (string, error) {
	generatedSecret, err := repo.secretGenerator.String(sessionSecretSize)
	if err != nil {
		return "", err
	}
	return repo.addSessionInternal(id, additionalData, generatedSecret)
}

func (repo *redisSessionsRepo) Get(secret string) (int, string, error) {
	separated := strings.Split(secret, sessionIdAndSecretSep)
	if len(separated) != sessionSplitNumber {
		return 0, "", http.ErrAuthSessionNotFound
	}

	idInString := separated[0]
	id, err := strconv.Atoi(idInString)
	if err != nil {
		return 0, "", http.ErrAuthSessionNotFound.Wrap(err, "session id atoi failed")
	}

	key := strings.Join([]string{repo.sessionsPrefix, idInString}, sessionIdAndSecretSep)
	var additionalData string
	err = repo.client.HGet(context.Background(), key, secret).Scan(&additionalData)
	if errors.Is(err, redis.Nil) {
		return 0, "", http.ErrAuthSessionNotFound
	}
	if err != nil {
		return 0, "", http.ErrBaseApp.Wrap(err, "redis session get failed")
	}

	return id, additionalData, nil
}

func (repo *redisSessionsRepo) Remove(secret string) error {
	separated := strings.Split(secret, sessionIdAndSecretSep)
	if len(separated) != sessionSplitNumber {
		return http.ErrAuthSessionNotFound
	}
	idInString := separated[0]
	key := strings.Join([]string{repo.sessionsPrefix, idInString}, sessionIdAndSecretSep)
	success, err := repo.client.HDel(context.Background(), key, secret).Result()
	if err != nil {
		return http.ErrBaseApp.Wrap(err, "session delete failed")
	}
	if success != 1 {
		return http.ErrAuthSessionNotFound
	}
	return nil
}
