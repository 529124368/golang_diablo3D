package engine

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/g3n/engine/audio"
	"github.com/g3n/engine/util/logger"
)

type Muscic struct {
	// BGM
	musicGame *audio.Player
	skill     *audio.Player
}

func init() {
	go func() {
		_, check := os.Stat("town1.wav")
		if os.IsNotExist(check) {
			datas, err := asset.ReadFile("asset/music/town1.wav")
			if err != nil {
				panic(err)
			}
			ioutil.WriteFile("town1.wav", datas, 0o600)
		} else {
			fmt.Println("has")
		}
	}()
	go func() {
		_, check := os.Stat("ba_skill.wav")
		if os.IsNotExist(check) {
			datas, err := asset.ReadFile("asset/music/ba_skill.wav")
			if err != nil {
				panic(err)
			}
			ioutil.WriteFile("ba_skill.wav", datas, 0o600)
		}

	}()

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
	a.musicGame = createPlayer("town1.wav")
	a.musicGame.SetLooping(true)
	a.skill = createPlayer("ba_skill.wav")
	return a
}
