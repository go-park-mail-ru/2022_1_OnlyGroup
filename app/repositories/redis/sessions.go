package redis

import (
	"2022_1_OnlyGroup_back/app/handlers"
	"2022_1_OnlyGroup_back/pkg/sessionGenerator"
	"context"
	"fmt"
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
	secretGenerator sessionGenerator.SessionGenerator
}

func NewRedisSessionRepository(client *redis.Client, sessionsPrefix string, generator sessionGenerator.SessionGenerator) *redisSessionsRepo {
	return &redisSessionsRepo{client: client, sessionsPrefix: sessionsPrefix, secretGenerator: generator}
}

func (repo *redisSessionsRepo) addSessionInternal(id int, additionalData string, generatedSecret string) (string, error) {
	idInString := strconv.Itoa(id)
	secret := strings.Join([]string{idInString, generatedSecret}, sessionIdAndSecretSep)

	key := strings.Join([]string{repo.sessionsPrefix, idInString}, sessionIdAndSecretSep)
	success, err := repo.client.HSet(context.Background(), key, secret, additionalData).Result()
	if err != nil || success != 1 {
		if err != nil {
			return "", handlers.ErrBaseApp.Wrap(err, "redis session add failed")
		}
		return "", handlers.ErrBaseApp.Wrap(fmt.Errorf("hset returned not 1 sucsess result"), "redis session add failed")
	}
	return secret, nil
}

func (repo *redisSessionsRepo) Add(id int, additionalData string) (string, error) {
	generatedSecret := repo.secretGenerator.String(sessionSecretSize)
	return repo.addSessionInternal(id, additionalData, generatedSecret)
}

func (repo *redisSessionsRepo) Get(secret string) (int, string, error) {
	separated := strings.Split(secret, sessionIdAndSecretSep)
	if len(separated) != sessionSplitNumber {
		return 0, "", handlers.ErrAuthSessionNotFound
	}

	idInString := separated[0]
	id, err := strconv.Atoi(idInString)
	if err != nil {
		return 0, "", handlers.ErrAuthSessionNotFound.Wrap(err, "session id atoi failed")
	}

	key := strings.Join([]string{repo.sessionsPrefix, idInString}, sessionIdAndSecretSep)
	var additionalData string
	err = repo.client.HGet(context.Background(), key, secret).Scan(&additionalData)
	if errors.Is(err, redis.Nil) {
		return 0, "", handlers.ErrAuthSessionNotFound
	}
	if err != nil {
		return 0, "", handlers.ErrBaseApp.Wrap(err, "redis session get failed")
	}

	return id, additionalData, nil
}

func (repo *redisSessionsRepo) Remove(secret string) error {
	separated := strings.Split(secret, sessionIdAndSecretSep)
	if len(separated) != sessionSplitNumber {
		return handlers.ErrAuthSessionNotFound
	}
	idInString := separated[0]
	key := strings.Join([]string{repo.sessionsPrefix, idInString}, sessionIdAndSecretSep)
	success, err := repo.client.HDel(context.Background(), key, secret).Result()
	if err != nil {
		return handlers.ErrBaseApp.Wrap(err, "session delete failed")
	}
	if success != 1 {
		return handlers.ErrAuthSessionNotFound
	}
	return nil
}
