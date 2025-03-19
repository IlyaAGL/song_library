package entities

// SongDetails model
// @Description Extended song metadata
type SongDetails struct {
  ReleaseDate string `json:"releaseDate" example:"2003-09-15"`
  Text        string `json:"text" example:"It's bugging me, grating me"`
  Link        string `json:"link" example:"https://www.youtube.com/watch?v=3dm_5qWWDV8"`
}
