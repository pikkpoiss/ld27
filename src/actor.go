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
	if t, err = system.LoadTexture(path, system.IntNearest, width); err != nil {
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
		X:      x,
		Y:      y,
		State:  state,
		Offset: offset,
	}
	c.Actors = append(c.Actors, a)
	return
}

func (c *Cast) Update() {
	for _, a := range ACTOR_ANIMATIONS {
		a.Next()
	}
}

type Actor struct {
	X      float64
	Y      float64
	State  int
	Offset int
	FlipX  bool
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
	return a.State & state == state
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
