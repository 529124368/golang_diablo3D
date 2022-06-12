package main

import (
	"test/engine"
)

func main() {
	//创建游戏
	game := engine.New()
	//启动游戏
	game.Run()
}
