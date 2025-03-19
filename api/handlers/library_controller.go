package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	. "github.com/agl/music_library/internal/domain/entities"
	. "github.com/agl/music_library/internal/logger"
	"github.com/agl/music_library/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type SongHandler struct {
	BASE_URL string
	conn *sql.DB
}

func NewSongHandler(conn *sql.DB) *SongHandler {
	godotenv.Load("../../.env")
	BASE_URL := os.Getenv("BASE_URL")

	return &SongHandler{BASE_URL: BASE_URL, conn: conn}
}

// @title Music Library API
// @version 1.0.0
// @description API for managing songs and their metadata
// @host localhost:6060
// @BasePath /
func SetAPI(conn *sql.DB) {
	r := gin.Default()
	handler := NewSongHandler(conn)

	r.GET("/songs", handler.GetSongs)
	r.GET("/song-text", handler.GetSongText)
	r.DELETE("/song", handler.DeleteSong)
	r.PUT("/updated-song/:id", handler.UpdateSong)
	r.POST("/new-song", handler.AddSong)

	r.Run(":6060")
}

// GetSords request for receiving list of songs by filter
// @Summary Get multiple songs
// @Description Get songs with pagination and filtering
// @Tags Songs
// @Produce json
// @Param page query int false "Page number (starts from 1)"
// @Param limit query int false "Items per page"
// @Param group query string false "Filter by group name"
// @Param title query string false "Filter by song title"
// @Success 200 {array} Song
// @Failure 500 {object} map[string]string
// @Router /songs [get]
func (h *SongHandler) GetSongs(ctx *gin.Context) {
	page := ctx.Query("page")
	limit := ctx.Query("limit")

	group := ctx.DefaultQuery("group", "")
	title := ctx.DefaultQuery("title", "")

	song := &Song{Group: group, Name: title}

	songs, err := services.GetMultipleSongs(*song, page, limit, h.conn)

	if err != nil {
		Log.Info("Error has occurred while trying to receive the songs")

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Something wrong happened",
		})

		return
	} else {
		Log.Info("Songs were successfully received")
	}

	ctx.JSON(http.StatusOK, songs)
}

// GetSongText обрабатывает запрос на получение текста песни
// @Summary Get song text
// @Description Retrieve the text of a song based on group and song name
// @Tags Songs
// @Produce json
// @Param group query string true "Group name"
// @Param song query string true "Song title"
// @Success 200 {string} string
// @Failure 500 {object} map[string]string
// @Router /song-text [get]
func (handler *SongHandler) GetSongText(ctx *gin.Context) {
	page := ctx.Query("page")
	limit := ctx.Query("limit")

	group := ctx.Query("group")
	songName := ctx.Query("song")

	song := &Song{Group: group, Name: songName}

	text, err := services.GetSongText(*song, page, limit, handler.conn)

	if err != nil {
		Log.Info("Error has occurred while trying to receive the song text")

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Something wrong happened",
		})

		return
	} else {
		Log.Info("Song text was successfully received")
	}

	ctx.JSON(http.StatusOK, text)
}

// DeleteSong обрабатывает запрос на удаление песни
// @Summary Delete a song
// @Description Remove a song from the library
// @Tags Songs
// @Accept json
// @Produce json
// @Param song body Song true "Song to delete"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /song [delete]
func (handler *SongHandler) DeleteSong(ctx *gin.Context) {
	var song Song
 
	if err := ctx.ShouldBindJSON(&song); err != nil {
		Log.Error("Invalid request body")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect body"})

		return
	}

	status, err := services.DeleteSong(song.Group, song.Name, handler.conn)

	if err != nil {
		Log.Info("Delete operation didn't complete :(")

		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": status,
		})

		return
	}

	Log.Info("Song was deleted successfully")

	ctx.JSON(http.StatusOK, gin.H{"message": status})
}

// UpdateSong обрабатывает запрос на обновление песни
// @Summary Update a song
// @Description Update an existing song by ID
// @Tags Songs
// @Accept json
// @Produce json
// @Param id path string true "Song ID"
// @Param song body Song true "Updated song data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /updated-song/{id} [put]
func (handler *SongHandler) UpdateSong(ctx *gin.Context) {
	var song Song
 
	if err := ctx.ShouldBindJSON(&song); err != nil {
		Log.Error("Invalid request body")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect body"})

		return
	}

	id := ctx.Param("id")
	status, err := services.UpdateSong(id, song.Group, song.Name, handler.conn)

	if err != nil {
		Log.Info("Could not update the song")

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": status})
		
		return
	}
	Log.Info("Song was updated successfully")

	ctx.JSON(http.StatusOK, gin.H{"message": status})
}

// AddSong обрабатывает запрос на добавление новой песни
// @Summary Add a new song
// @Description Create a new song and fetch its details
// @Tags Songs
// @Accept json
// @Produce json
// @Param song body Song true "New song data"
// @Success 200 {object} SongDetails
// @Failure 400 {object} map[string]string
// @Router /new-song [post]
func (handler *SongHandler) AddSong(ctx *gin.Context) {
	var song Song
 
	if err := ctx.ShouldBindJSON(&song); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})

		return
	}

	params := url.Values{}
	params.Add("group", song.Group)
	params.Add("song", song.Name)

	fullURL := fmt.Sprintf("%s?%s", handler.BASE_URL, params.Encode())
	resp, err := http.Get(fullURL)

	if err != nil {
		Log.Error("Could not receive song details data")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error :("})

		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		Log.Error("Could not read the response body")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid body received"})

		return
	}

	Log.Info("Song details successfully received")

	var songDetails SongDetails
	if err := json.Unmarshal(body, &songDetails); err != nil {
		Log.Error("Could not parse the response body")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not parse the response body"})

		return
	}

	status, err := services.AddSong(song, songDetails, handler.conn)

	if err != nil {
		Log.Info("Could not add the song :(")

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": status,
		})

		return
	}

	Log.Info("Song details successfully parsed and everything is added")

	ctx.JSON(http.StatusOK, songDetails)
}