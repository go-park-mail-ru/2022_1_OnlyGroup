package redis

import (
	"2022_1_OnlyGroup_back/app/handlers"
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"math/rand"
	"strconv"
	"strings"
)

const secretRunes = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890"
const sessionSecretSize = 64
const sessionExpires = "86400"
const sessionIdAndSecretSep = "_"

type redisSessionsRepo struct {
	client         *redis.Client
	sessionsPrefix string
}

func CreateRedisSessionRepository(client *redis.Client, sessionsPrefix string) *redisSessionsRepo {
	return &redisSessionsRepo{client: client, sessionsPrefix: sessionsPrefix}
}

func generateSecret(size int) string {
	result := ""
	for i := 0; i < size; i++ {
		result += string(secretRunes[rand.Intn(len(secretRunes))])
	}
	return result
}

func (repo *redisSessionsRepo) AddSession(id int, additionalData string) (string, error) {
	idInString := strconv.Itoa(id)
	secret := strings.Join([]string{idInString, generateSecret(sessionSecretSize)}, sessionIdAndSecretSep)

	key := repo.sessionsPrefix + ":" + idInString
	success, err := repo.client.HSet(context.Background(), key, secret, additionalData).Result()
	//success, err := redis.Int(repo.redisConn.Do("HSET", key, secret, additionalData))
	if err != nil || success != 1 {
		if err != nil {
			err = errors.Wrap(handlers.ErrBaseApp, err.Error())
		}
		err = errors.Wrap(err, "redis session add failed")
		return "", err
	}
	return secret, nil
}

func (repo *redisSessionsRepo) GetIdBySession(secret string) (int, string, error) {
	separated := strings.Split(secret, sessionIdAndSecretSep)
	if len(separated) != 2 {
		return 0, "", handlers.ErrAuthSessionNotFound
	}

	idInString := separated[0]
	id, err := strconv.Atoi(idInString)
	if err != nil {
		err = errors.Wrap(handlers.ErrBaseApp, err.Error())
		err = errors.Wrap(err, "session id atoi failed")
		return 0, "", err
	}

	key := repo.sessionsPrefix + ":" + idInString
	var additionalData string
	err = repo.client.HGet(context.Background(), key, secret).Scan(&additionalData)
	if errors.Is(err, redis.Nil) {
		return 0, "", handlers.ErrAuthSessionNotFound
	}
	if err != nil {
		err = errors.Wrap(handlers.ErrBaseApp, err.Error())
		err = errors.Wrap(err, "redis session get failed")
		return 0, "", err
	}

	return id, additionalData, nil
}

func (repo *redisSessionsRepo) RemoveSession(secret string) error {
	separated := strings.Split(secret, sessionIdAndSecretSep)
	if len(separated) != 2 {
		return handlers.ErrAuthSessionNotFound
	}
	idInString := separated[0]
	key := repo.sessionsPrefix + ":" + idInString
	success, err := repo.client.HDel(context.Background(), key, secret).Result()
	if err != nil {
		err = errors.Wrap(handlers.ErrBaseApp, err.Error())
		err = errors.Wrap(err, "session delete failed")
		return err
	}
	if success != 1 {
		return handlers.ErrAuthSessionNotFound
	}
	return nil
}
