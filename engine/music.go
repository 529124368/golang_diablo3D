package engine

import (
	"github.com/g3n/engine/audio"
	"github.com/g3n/engine/util/logger"
)

type Muscic struct {
	// BGM
	musicGame *audio.Player
	skill     *audio.Player
}

func NewMusic() *Muscic {
	a := new(Muscic)
	// Helper function to create player and handle errors
	createPlayer := func(fileName string) *audio.Player {
		p, err := audio.NewPlayer(fileName)
		if err != nil {
			logger.Error("Failed to create sound player: %v", err)
		}
		return p
	}
	// Music
	a.musicGame = createPlayer("asset/music/town1.wav")
	a.musicGame.SetLooping(true)
	a.skill = createPlayer("asset/music/ba_skill.wav")
	return a
}
