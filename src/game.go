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
	Maps       []string
	Level      *Level
	Camera     *Camera
	exit       chan bool
}

func NewGame(ctrl *system.Controller) (game *Game, err error) {
	game = &Game{
		Controller: ctrl,
		Maps: []string{
			"data/level01.json",
		},
		exit: make(chan bool, 1),
	}
	game.Controller.SetClearColor(BG_R, BG_G, BG_B, BG_A)
	game.handleKeys()
	game.handleClose()
	err = game.SetLevel(0)
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

func (g *Game) SetLevel(i int) (err error) {
	var (
		index = (i + len(g.Maps)) % len(g.Maps)
		path  = g.Maps[index]
	)
	if g.Level, err = LoadLevel(path); err != nil {
		return
	}
	return
}

func (g *Game) Run() (err error) {
	go func() {
		update := time.NewTicker(time.Second / time.Duration(UPDATE_HZ))
		for true {
			<-update.C
			g.checkKeys()
			g.Level.Update()
		}
	}()
	running := true
	paint := time.NewTicker(time.Second / time.Duration(PAINT_HZ))
	for running == true {
		<-paint.C
		g.Level.Camera.SetProjection()
		BeginPaint()
		PaintMap(g.Controller, g.Level.Map)
		EndPaint()
		select {
		case <-g.exit:
			paint.Stop()
			running = false
		default:
		}
	}
	return
}
