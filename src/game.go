// Copyright 2013 Arne Roomann-Kurrik
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"./system"
	"log"
	"time"
)

const (
	UPDATE_HZ int = 60
	PAINT_HZ  int = 60
	BG_R      int = 0
	BG_G      int = 255
	BG_B      int = 0
	BG_A      int = 0
)

type Game struct {
	Controller *system.Controller
	exit       chan bool
}

func NewGame(ctrl *system.Controller) (game *Game, err error) {
	game = &Game{
		Controller: ctrl,
		exit:       make(chan bool, 1),
	}
	game.Controller.SetClearColor(BG_R, BG_G, BG_B, BG_A)
	game.handleKeys()
	game.handleClose()
	return
}

func (g *Game) handleClose() {
	g.Controller.SetCloseCallback(func() int {
		g.exit <- true
		return 0
	})
}

func (g *Game) checkKeys() {
	switch {
	case g.Controller.Key(system.KeySpace) == 1:
		log.Printf("Space\n")
		g.exit <- true
	}
}

func (g *Game) handleKeys() {
	g.Controller.SetKeyCallback(func(key int, state int) {
		log.Printf("handleKeys: %v %v\n", key, state)
	})
}

func (g *Game) Run() (err error) {
	go func() {
		update := time.NewTicker(time.Second / time.Duration(UPDATE_HZ))
		for true {
			<-update.C
			g.checkKeys()
		}
	}()
	running := true
	paint := time.NewTicker(time.Second / time.Duration(PAINT_HZ))
	for running == true {
		<-paint.C
		g.Controller.Paint()
		select {
		case <-g.exit:
			paint.Stop()
			running = false
		default:
		}
	}
	return
}
