package repositories

import (
	"database/sql"
	"strings"

	. "github.com/agl/music_library/internal/domain/entities"
	. "github.com/agl/music_library/internal/logger"
)

func GetIdSongsByPage(filter Song, limit int, offset int, conn *sql.DB) ([]Song, error) {
	Log.Info("Fetching songs by page")

	rows, err := conn.Query(
		`SELECT artist, name 
         FROM Library 
         WHERE 
             ($1 = '' OR artist = $1) AND ($2 = '' OR name = $2) 
         LIMIT $3 OFFSET $4`,
		filter.Group, filter.Name, limit, offset)

	if err != nil {
		Log.Info("Failed to query songs by filter", "error", err)
		return nil, err
	}

	defer rows.Close()

	var songs []Song
	for rows.Next() {
		var song Song

		if err := rows.Scan(&song.Group, &song.Name); err != nil {
			Log.Info("Failed to scan song", "error", err)

			return nil, err
		}

		songs = append(songs, song)
	}

	Log.Debug("Successfully fetched songs")

	return songs, nil
}

func GetSongTextPaginated(song Song, page, limit int, conn *sql.DB) ([]string, error) {
	Log.Info("Fetching song text with pagination", "group", song.Group, "title", song.Name, "page", page, "limit", limit)

	var songText string
	err := conn.QueryRow(
		`SELECT text FROM Song_details 
         WHERE id = (SELECT id FROM Library WHERE artist = $1 AND name = $2)`,
		song.Group, song.Name,
	).Scan(&songText)

	if err != nil {
		Log.Error("Failed to fetch song text", "group", song.Group, "title", song.Name, "error", err)
		return nil, err
	}

	verses := strings.Split(songText, "\n\n")

	start := (page - 1) * limit
	end := start + limit

	if start >= len(verses) {
		return []string{}, nil
	}
	if end > len(verses) {
		end = len(verses)
	}

	return verses[start:end], nil
}

func DeleteSongByDetails(group, song string, conn *sql.DB) (bool, error) {
	Log.Debug("Deleting song by details", "group", group, "song", song)

	_, err := conn.Exec(
		`DELETE FROM Library WHERE artist = $1 AND name = $2`,
		group, song,
	)

	if err != nil {
		Log.Info("Failed to delete song", "group", group, "song", song, "error", err)

		return false, err
	}

	Log.Debug("Successfully deleted song", "group", group, "song", song)

	return true, nil
}

func UpdateSongByDetails(id int, group, song string, conn *sql.DB) (bool, error) {
	Log.Debug("Updating song by details", "id", id, "group", group, "song", song)

	_, err := conn.Exec(
		`UPDATE Library SET artist = $1, name = $2 WHERE id = $3`,
		group, song, id,
	)

	if err != nil {
		Log.Info("Failed to update song", "id", id, "group", group, "song", song, "error", err)

		return false, err
	}

	Log.Debug("Successfully updated song", "id", id, "group", group, "song", song)

	return true, nil
}

func AddNewSong(song Song, details SongDetails, conn *sql.DB) (bool, error) {
	Log.Debug("Adding new song", "group", song.Group, "song", song.Name)

	var id int
	err := conn.QueryRow(
		`INSERT INTO Library ("artist", "name") VALUES ($1, $2) RETURNING ID`,
		song.Group, song.Name,
	).Scan(&id)

	if err != nil {
		Log.Info("Failed to add new song", "group", song.Group, "song", song.Name, "error", err)

		return false, err
	}

	_, err = conn.Exec(
		`INSERT INTO Song_details ("id", "releaseDate", "text", "song_link") VALUES ($1, $2, $3, $4)`,
		id, details.ReleaseDate, details.Text, details.Link)

	if err != nil {
		Log.Info("Failed to add new song details", "group", song.Group, "song", song.Name, "error", err)

		return false, err
	}

	Log.Debug("Successfully added new song", "group", song.Group, "song", song.Name)

	return true, nil
}
