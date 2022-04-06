package postgres

import (
	"2022_1_OnlyGroup_back/app/handlers"
	"2022_1_OnlyGroup_back/app/models"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type postgresPhotosRepository struct {
	connection      *sqlx.DB
	photosTableName string
	avatarTableName string
}

func NewPostgresPhotoRepository(conn *sqlx.DB, photosTableName string, usersTableName string, avatarTableName string) (*postgresPhotosRepository, error) {
	_, err := conn.Exec("CREATE TABLE IF NOT EXISTS " + photosTableName + "(id bigserial primary key, author bigserial not null references " + usersTableName + "(id), left_top_x smallint, left_top_y smallint, right_bottom_x smallint, right_bottom_y smallint, path varchar(32));")
	if err != nil {
		return nil, fmt.Errorf("create table failed: %w", err)
	}
	_, err = conn.Exec("CREATE TABLE IF NOT EXISTS " + avatarTableName + "(user_id bigserial unique references " + usersTableName + "(id), photo_id bigserial unique references " + photosTableName + "(id), left_top_x smallint, left_top_t smallint, right_bottom_x smallint, right_bottom_y smallint);")
	if err != nil {
		return nil, fmt.Errorf("create table failed: %w", err)
	}
	return &postgresPhotosRepository{connection: conn, photosTableName: photosTableName, avatarTableName: avatarTableName}, nil
}

func (repo *postgresPhotosRepository) Create(userId int) (int, error) {
	var photoId int
	err := repo.connection.QueryRow("INSERT INTO "+repo.photosTableName+"(author, path) VALUES($1, NULL) RETURNING id;", userId).Scan(&photoId)
	if err != nil {
		return 0, handlers.ErrBaseApp.Wrap(err, "create photo failed")
	}
	return photoId, err
}

func (repo *postgresPhotosRepository) Save(id int, path string) error {
	res, err := repo.connection.Exec("UPDATE "+repo.photosTableName+" SET path=$1 WHERE id=$2;", path, id)
	if err != nil {
		return handlers.ErrBaseApp.Wrap(err, "save photo failed")
	}
	rowsAffect, err := res.RowsAffected()
	if err != nil {
		return handlers.ErrBaseApp.Wrap(err, "save photo failed")
	}
	if rowsAffect != 1 {
		return handlers.ErrBaseApp.Wrap(nil, "save photo failed: try to save not existed photo")
	}
	return nil
}

func (repo *postgresPhotosRepository) IsSaved(id int) (bool, error) {
	var path *string
	err := repo.connection.QueryRow("SELECT path FROM "+repo.photosTableName+" WHERE id=$1;", id).Scan(&path)
	if errors.Is(err, sql.ErrNoRows) {
		return false, handlers.ErrPhotoNotFound
	}
	if err != nil {
		return false, handlers.ErrBaseApp.Wrap(err, "is saved photo failed")
	}
	return path != nil, nil
}

func (repo *postgresPhotosRepository) GetAuthor(id int) (int, error) {
	var authorId int
	err := repo.connection.QueryRow("SELECT author FROM "+repo.photosTableName+" WHERE id=$1;", id).Scan(&authorId)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, handlers.ErrPhotoNotFound
	}
	if err != nil {
		return 0, handlers.ErrBaseApp.Wrap(err, "get photo author failed")
	}
	return authorId, nil
}

func (repo *postgresPhotosRepository) SetParams(id int, params models.PhotoParams) error {
	res, err := repo.connection.Exec("UPDATE "+repo.photosTableName+" SET left_top_x=$1,left_top_y=$2,right_bottom_x=$3,right_bottom_y=$4 WHERE id=$5;", params.LeftTop[0], params.LeftTop[1], params.RightBottom[0], params.RightBottom[1], id)
	if err != nil {
		return handlers.ErrBaseApp.Wrap(err, "set params photo failed")
	}
	rowsAffect, err := res.RowsAffected()
	if err != nil {
		return handlers.ErrBaseApp.Wrap(err, "set params photo failed")
	}
	if rowsAffect != 1 {
		return handlers.ErrPhotoNotFound
	}
	return nil
}

func (repo *postgresPhotosRepository) GetParams(id int) (models.PhotoParams, error) {
	var leftTopX, leftTopY, rightBottomX, rightBottomY *int
	err := repo.connection.QueryRow("SELECT left_top_x, left_top_y, right_bottom_x, right_bottom_y FROM "+repo.photosTableName+" WHERE id=$1;", id).Scan(&leftTopX, &leftTopY, &rightBottomX, &rightBottomY)
	if errors.Is(err, sql.ErrNoRows) {
		return models.PhotoParams{}, handlers.ErrPhotoNotFound
	}
	if err != nil {
		return models.PhotoParams{}, handlers.ErrBaseApp.Wrap(err, "save photo failed")
	}
	if leftTopX == nil || leftTopY == nil || rightBottomX == nil || rightBottomY == nil {
		return models.PhotoParams{}, handlers.ErrPhotoNotFound
	}
	return models.PhotoParams{LeftTop: [2]int{*leftTopX, *leftTopY}, RightBottom: [2]int{*rightBottomX, *rightBottomY}}, nil
}

func (repo *postgresPhotosRepository) GetPathIfFilled(id int) (string, error) {
	var path *string
	err := repo.connection.QueryRow("SELECT path FROM "+repo.photosTableName+" WHERE id=$1 AND left_top_x IS NOT NULL AND left_top_y IS NOT NULL AND right_bottom_x IS NOT NULL AND right_bottom_y IS NOT NULL;", id).Scan(&path)
	if errors.Is(err, sql.ErrNoRows) {
		return "", handlers.ErrPhotoNotFound
	}
	if err != nil {
		return "", handlers.ErrBaseApp.Wrap(err, "save photo failed")
	}
	if path == nil {
		return "", handlers.ErrPhotoNotFound
	}
	return *path, nil
}

func (repo *postgresPhotosRepository) Delete(id int) error {
	res, err := repo.connection.Exec("DELETE FROM "+repo.photosTableName+" WHERE id=$1;", id)
	if err != nil {
		return handlers.ErrBaseApp.Wrap(err, "delete photo failed")
	}
	rowsAffect, err := res.RowsAffected()
	if err != nil {
		return handlers.ErrBaseApp.Wrap(err, "delete photo failed")
	}
	if rowsAffect != 1 {
		return handlers.ErrPhotoNotFound
	}
	return nil
}

func (repo *postgresPhotosRepository) GetAvatar(userId int) (int, models.PhotoParams, error) {
	model := models.PhotoParams{}
	var photoId int
	err := repo.connection.QueryRow("SELECT photo_id, left_top_x, left_top_y, right_bottom_x, right_bottom_y FROM "+repo.avatarTableName+" WHERE user_id=$1;", userId).Scan(&photoId, &model.LeftTop[0], &model.LeftTop[1], &model.RightBottom[0], &model.RightBottom[1])
	if errors.Is(err, sql.ErrNoRows) {
		return 0, models.PhotoParams{}, handlers.ErrPhotoNotFound
	}
	if err != nil {
		return 0, models.PhotoParams{}, handlers.ErrBaseApp.Wrap(err, "get avatar failed")
	}
	return photoId, model, nil
}

func (repo *postgresPhotosRepository) SetAvatar(id int, params models.PhotoParams, userId int) error {
	var testPhotoFilled int
	err := repo.connection.QueryRow("SELECT photo_id FROM "+repo.avatarTableName+" WHERE user_id=$1;", userId).Scan(&testPhotoFilled)
	if errors.Is(err, sql.ErrNoRows) {
		res, err := repo.connection.Exec("INSERT INTO "+repo.avatarTableName+"(user_id, photo_id, left_top_x, left_top_y, right_bottom_x, right_bottom_y) VALUES ($1, $2, $3, $4, $5, $6", userId, id, params.LeftTop[0], params.LeftTop[1], params.RightBottom[0], params.RightBottom[1])
		if err != nil {
			return handlers.ErrBaseApp.Wrap(err, "set avatar insert failed")
		}
		affected, err := res.RowsAffected()
		if err != nil {
			return handlers.ErrBaseApp.Wrap(err, "set avatar insert failed")
		}
		if affected == 0 {
			return handlers.ErrBaseApp.Wrap(nil, "set avatar insert rows affected = 0")
		}
		return nil
	}
	res, err := repo.connection.Exec("UPDATE "+repo.avatarTableName+" SET photo_id=$1, left_top_x=$2, left_top_y=$3, right_bottom_x=$4, right_bottom_x=$5 WHERE user_id=$6", id, params.LeftTop[0], params.LeftTop[1], params.RightBottom[0], params.RightBottom[1], userId)
	if err != nil {
		return handlers.ErrBaseApp.Wrap(err, "set avatar update failed")
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return handlers.ErrBaseApp.Wrap(err, "set avatar update failed")
	}
	if affected == 0 {
		return handlers.ErrBaseApp.Wrap(nil, "set avatar update rows affected = 0")
	}
	return nil
}

func (repo *postgresPhotosRepository) GetUserPhotos(userId int) (models.UserPhotos, error) {
	var model models.UserPhotos
	err := repo.connection.Select(&model.Photos, "SELECT id FROM "+repo.photosTableName+" WHERE author=$1", userId)
	if err != nil {
		return models.UserPhotos{}, handlers.ErrBaseApp.Wrap(err, "get all user photos failed")
	}
	if model.Photos == nil {
		model.Photos = []int{}
	}
	return model, nil
}
