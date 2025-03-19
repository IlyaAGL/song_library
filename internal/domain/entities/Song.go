package entities

// @Schema(name="Song")
// @Title Song
// @Description Musical composition entity
type Song struct {
	Group string `json:"group" example:"Muse"`
	Name  string `json:"song" example:"Hysteria"`
}
