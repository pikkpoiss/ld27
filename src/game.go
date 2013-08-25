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
	LevelIndex int
	Render     bool
	Camera     *Camera
	exit       chan bool
}

func NewGame(ctrl *system.Controller) (game *Game, err error) {
	game = &Game{
		Controller: ctrl,
		Maps: []string{
			"data/level01.json",
			"data/level02.json",
		},
		LevelIndex: 0,
		Render:     false,
		exit:       make(chan bool, 1),
	}
	game.Controller.SetClearColor(BG_R, BG_G, BG_B, BG_A)
	game.handleKeys()
	game.handleClose()
	if err = game.setLevel(); err != nil {
		return
	}
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
	case g.Controller.Key(system.KeyEsc) == 1:
		g.exit <- true
	}
}

func (g *Game) handleKeys() {
	g.Controller.SetKeyCallback(func(key int, state int) {
		switch {
		case state == 1 && key == system.KeyUp:
			g.Level.Player.SetDirection(UP)
			g.Level.Player.SetMovement(WALKING)
		case state == 1 && key == system.KeyDown:
			g.Level.Player.SetDirection(DOWN)
			g.Level.Player.SetMovement(WALKING)
		case state == 1 && key == system.KeyLeft:
			g.Level.Player.SetDirection(LEFT)
			g.Level.Player.SetMovement(WALKING)
		case state == 1 && key == system.KeyRight:
			g.Level.Player.SetDirection(RIGHT)
			g.Level.Player.SetMovement(WALKING)
		case state == 1 && key == 87: //w
			log.Printf("Autowin\n")
			g.Level.Won = true
		case state == 0:
			switch {
			case g.Level.Player.TestState(UP) && key == system.KeyUp ||
				g.Level.Player.TestState(DOWN) && key == system.KeyDown ||
				g.Level.Player.TestState(LEFT) && key == system.KeyLeft ||
				g.Level.Player.TestState(RIGHT) && key == system.KeyRight:
				g.Level.Player.SetMovement(STOPPED)
			case key == system.KeySpace:
				g.Level.AddBombFromActor(g.Level.Player.Actor)
			}
		default:
			log.Printf("handleKeys: %v %v\n", key, state)
		}
	})
}

func (g *Game) setLevel() (err error) {
	var (
		index = (g.LevelIndex + len(g.Maps)) % len(g.Maps)
		path  = g.Maps[index]
		cast  *Cast
	)
	if cast, err = g.getCast("data/actors.png", 32, 64); err != nil {
		return
	}
	if g.Level, err = LoadLevel(path, cast); err != nil {
		return
	}
	g.Render = true
	return
}

func (g *Game) getCast(path string, width int, height int) (cast *Cast, err error) {
	return LoadCast(path, width, height, 32, 32)
}

func (g *Game) Run() (err error) {
	go func() {
		var (
			now    time.Time
			last   time.Time
			update *time.Ticker
			diff   time.Duration
		)
		update = time.NewTicker(time.Second / time.Duration(UPDATE_HZ))
		last = time.Now()
		for true {
			<-update.C
			now = time.Now()
			diff = now.Sub(last)
			g.checkKeys()
			g.Level.Update(diff)
			last = now
		}
	}()
	running := true
	paint := time.NewTicker(time.Second / time.Duration(PAINT_HZ))
	for running == true {
		<-paint.C
		if g.Render {
			g.Level.Camera.SetProjection()
			BeginPaint()
			PaintMap(g.Controller, g.Level.Map)
			PaintCast(g.Controller, g.Level.Cast)
			EndPaint()
		}
		select {
		case <-g.exit:
			paint.Stop()
			running = false
		default:
		}
		if g.Level.Won {
			g.Render = false
			g.LevelIndex += 1
			g.setLevel()
		}
	}
	return
}
