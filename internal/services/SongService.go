package services

import (
	"database/sql"
	"fmt"
	"strconv"

	. "github.com/agl/music_library/internal/domain/entities"
	. "github.com/agl/music_library/internal/logger"
	"github.com/agl/music_library/internal/repositories"
)

func GetMultipleSongs(song Song, page, limit string, conn *sql.DB) ([]Song, error) {
	Log.Debug("Fetching songs", "page", page)
	limit_int, err := strconv.Atoi(limit)

	if err != nil || limit_int < 0 {
		Log.Info("Invalid limit format", "limit", limit)

		return nil, err
	}

	page_int, err := strconv.Atoi(page)

	if err != nil || page_int < 0 {
		Log.Info("Invalid page format", "page", page)
		return nil, err
	}

	offset := (page_int - 1) * limit_int

	songs, err := repositories.GetIdSongsByPage(song, limit_int, offset, conn)

	if err != nil {
		Log.Info("Could not receive songs data")

		return nil, err
	}

	return songs, nil
}

func GetSongText(song Song, page, limit string, conn *sql.DB) ([]string, error) {
	limit_int, err := strconv.Atoi(limit)

	if err != nil || limit_int < 0 {
		Log.Info("Invalid limit format", "limit", limit)

		return nil, err
	}

	page_int, err := strconv.Atoi(page)

	if err != nil || page_int < 0 {
		Log.Info("Invalid page format", "page", page)
		return nil, err
	}

	text, err := repositories.GetSongTextPaginated(song, page_int, limit_int, conn)

	if err != nil {
		Log.Info("Something went wrong during fetching song text")

		return nil, err
	}

	Log.Info("Successfully fetched song text")

	return text, nil
}

func DeleteSong(group, song string, conn *sql.DB) (string, error) {
	Log.Debug("Deleting song", "group", group, "name", song)

	isDeleted, err := repositories.DeleteSongByDetails(group, song, conn)

	if err != nil {
		Log.Info("Something went wrong during deleting the song")

    return "Could not delete the song :(", err
	}

	if isDeleted {
		return "Successfully deleted", nil
	}

	return "There is no song with that parametres", fmt.Errorf("There is no song with that parametres")
}

func UpdateSong(id string, group, song string, conn *sql.DB) (string, error) {
	numberId, err := strconv.Atoi(id)

	if err != nil || numberId < 0 {
		Log.Info("Invalid id")

		return "Error", err
	}

	isUpdated, err := repositories.UpdateSongByDetails(numberId, group, song, conn)

	if err != nil {
		Log.Info("Something went wrong during updating the song details")

		return "Error", err
	}

	if isUpdated {
		return "Successfully updated", nil
	}

	return "Error", fmt.Errorf("Something wrong happened")

}

func AddSong(song Song, details SongDetails, conn *sql.DB) (string, error) {
	Log.Debug("Adding song", "group", song.Group, "name", song.Name)

	isAdded, err := repositories.AddNewSong(song, details, conn)

	if err != nil {
		Log.Info("Something went wrong during adding the song details")

		return "Error", err
	}

	if isAdded {
		return "Successfully added new song", nil
	}

	return "Error", fmt.Errorf("Something wrong happened")
}
