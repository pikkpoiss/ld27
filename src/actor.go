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
	"math"
)

type Cast struct {
	Texture *system.Texture
	Actors  []*Actor
	Width   int
	Height  int
	OffsetX int
	OffsetY int
}

func LoadCast(path string, width int, height int, th int, tw int) (c *Cast, err error) {
	var t *system.Texture
	if t, err = system.LoadTexture(path, system.IntNearest, width, height); err != nil {
		return
	}
	c = &Cast{
		Texture: t,
		Width:   width,
		Height:  height,
		OffsetX: width - tw,
		OffsetY: height - th,
	}
	return
}

func (c *Cast) AddActor(x float64, y float64, state int, offset int) (a *Actor) {
	a = &Actor{
		X:       x,
		Y:       y,
		State:   state,
		Offset:  offset,
		Rate:    2.0,
		Padding: 12,
	}
	c.Actors = append(c.Actors, a)
	return
}

func (c *Cast) Update(level *Level) {
	for _, a := range c.Actors {
		if !a.TestState(WALKING) {
			continue
		}
		switch {
		case a.TestState(DOWN):
			a.moveDown(level)
		case a.TestState(UP):
			a.moveUp(level)
		case a.TestState(RIGHT):
			a.moveRight(level)
		case a.TestState(LEFT):
			a.moveLeft(level)
		}
	}
	for _, a := range ACTOR_ANIMATIONS {
		a.Next()
	}
}

type Actor struct {
	X       float64
	Y       float64
	State   int
	Offset  int
	FlipX   bool
	Rate    float64
	Padding int
}

// Attempts to round X or Y values to tile boundaries if they're within Padding
func (a *Actor) getClamped(v float64, size int) int {
	var (
		clamped = math.Floor(v/float64(size)+0.5) * float64(size)
		diff    = math.Abs(clamped - v)
	)
	if int(diff) <= a.Padding {
		return int(clamped)
	}
	return int(v)
}

func (a *Actor) moveDown(l *Level) {
	var (
		x int
		y int
	)
	x = a.getClamped(a.X, l.TileWidth)
	y = int(a.Y + a.Rate)
	if l.TestPixelPassable(x+a.Padding, y+l.TileHeight) &&
		l.TestPixelPassable(x+l.TileWidth-a.Padding, y+l.TileHeight) {
		if x == int(a.X) {
			// Only move once we've clamped.
			a.Y = float64(y)
		}
		a.X = float64(x)
	}
}

func (a *Actor) moveUp(l *Level) {
	var (
		x int
		y int
	)
	x = a.getClamped(a.X, l.TileWidth)
	y = int(a.Y - a.Rate)
	if l.TestPixelPassable(x+a.Padding, y) &&
		l.TestPixelPassable(x+l.TileWidth-a.Padding, y) {
		if x == int(a.X) {
			// Only move once we've clamped.
			a.Y = float64(y)
		}
		a.X = float64(x)
	}
}

func (a *Actor) moveRight(l *Level) {
	var (
		x int
		y int
	)
	y = a.getClamped(a.Y, l.TileHeight)
	x = int(a.X + a.Rate)
	if l.TestPixelPassable(x+l.TileWidth, y+a.Padding) &&
		l.TestPixelPassable(x+l.TileWidth, y+l.TileHeight-a.Padding) {
		if y == int(a.Y) {
			// Only move once we've clamped.
			a.X = float64(x)
		}
		a.Y = float64(y)
	}
}

func (a *Actor) moveLeft(l *Level) {
	var (
		x int
		y int
	)
	y = a.getClamped(a.Y, l.TileHeight)
	x = int(a.X - a.Rate)
	if l.TestPixelPassable(x, y+a.Padding) &&
		l.TestPixelPassable(x, y+l.TileHeight-a.Padding) {
		if y == int(a.Y) {
			// Only move once we've clamped.
			a.X = float64(x)
		}
		a.Y = float64(y)
	}
}

const UNSET_MASK = 1<<10 - 1

func (a *Actor) unsetState(mask int) {
	a.State &= UNSET_MASK ^ mask
}

func (a *Actor) setState(mask int) {
	a.State |= mask
}

func (a *Actor) SetDirection(dir int) {
	a.unsetState(LEFT | RIGHT | UP | DOWN)
	if dir == RIGHT {
		a.FlipX = true
	} else {
		a.FlipX = false
	}
	a.setState(dir)
}

func (a *Actor) SetMovement(mov int) {
	a.unsetState(WALKING | STOPPED)
	a.setState(mov)
}

func (a *Actor) TestState(state int) bool {
	return a.State&state == state
}

func (a *Actor) GetFrame() int {
	var (
		anim *system.Animation
		ok   bool
	)
	if anim, ok = ACTOR_ANIMATIONS[a.State]; !ok {
		log.Printf("No animation for state %v", a.State)
		anim = ACTOR_ANIMATIONS[LEFT|STOPPED]
	}
	return anim.Curr() + a.Offset
}

const (
	LEFT    = 1 << iota
	RIGHT   = 1 << iota
	UP      = 1 << iota
	DOWN    = 1 << iota
	WALKING = 1 << iota
	STOPPED = 1 << iota
)

var ACTOR_ANIMATIONS = map[int]*system.Animation{
	LEFT | STOPPED:  system.Anim([]int{6}, 4),
	RIGHT | STOPPED: system.Anim([]int{6}, 4),
	UP | STOPPED:    system.Anim([]int{3}, 4),
	DOWN | STOPPED:  system.Anim([]int{0}, 4),
	LEFT | WALKING:  system.Anim([]int{6, 7, 6, 8}, 4),
	RIGHT | WALKING: system.Anim([]int{6, 7, 6, 8}, 4),
	UP | WALKING:    system.Anim([]int{3, 4, 3, 5}, 4),
	DOWN | WALKING:  system.Anim([]int{0, 1, 0, 2}, 4),
}
