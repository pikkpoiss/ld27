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
	"sort"
	"time"
)

type Cast struct {
	Texture     *system.Texture
	Actors      []system.Drawable
	Width       int
	Height      int
	OffsetX     int
	OffsetY     int
	TextureCols int
}

func LoadCast(path string, width int, height int, th int, tw int) (c *Cast, err error) {
	var t *system.Texture
	if t, err = system.LoadTexture(path, system.IntNearest, width, height); err != nil {
		return
	}
	c = &Cast{
		Texture:     t,
		Width:       width,
		Height:      height,
		OffsetX:     width - tw,
		OffsetY:     height - th,
		TextureCols: t.Width / width,
		Actors:      []system.Drawable{},
	}
	return
}

func (c *Cast) Overlaps(a *Actor, b *Actor) bool {
	var (
		w    = float64(c.Width - c.OffsetX)
		h    = float64(c.Height - c.OffsetY)
		xmin = a.x >= b.x && a.x <= b.x+w
		xmax = a.x+w >= b.x && a.x+w <= b.x+w
		ymin = a.y >= b.y && a.y <= b.y+h
		ymax = a.y+h >= b.y && a.y+h <= b.y+h
	)
	return (xmin || xmax) && (ymin || ymax)
}

func (c *Cast) AddActor(a system.Drawable) {
	if a == nil {
		return
	}
	c.Actors = append(c.Actors, a)
	sort.Sort(system.ByY{c.Actors})
	return
}

func (c *Cast) RemoveActor(a system.Drawable) {
	for i, actor := range c.Actors {
		if a != actor {
			continue
		}
		c.Actors = append(c.Actors[:i], c.Actors[i+1:]...)
		break
	}
}

func (c *Cast) Update(level *Level, diff time.Duration) {
	sort.Sort(system.ByY{c.Actors})
	for _, a := range ACTOR_ANIMATIONS {
		a.Next()
	}
}

type Actor struct {
	x          float64
	y          float64
	State      int
	textureRow int
	flipX      bool
	Rate       float64
	Padding    int
	Bomb       *Bomb
}

func NewActor(x float64, y float64, state int, textureRow int) *Actor {
	return &Actor{
		x:          x,
		y:          y,
		State:      state,
		textureRow: textureRow,
	}
}

func (a *Actor) X() float64 {
	return a.x
}

func (a *Actor) Y() float64 {
	return a.y
}

func (a *Actor) SetX(x float64) {
	a.x = x
}

func (a *Actor) SetY(y float64) {
	a.y = y
}

func (a *Actor) FlipX() bool {
	return a.flipX
}

func (a *Actor) TextureRow() int {
	return a.textureRow
}

func (a *Actor) Update(level *Level) bool {
	if !a.TestState(WALKING) {
		return true
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
	return true
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
	return anim.Curr()
}

func (a *Actor) TestState(state int) bool {
	return a.State&state == state
}

const UNSET_MASK = 1<<10 - 1

func (a *Actor) UnsetState(mask int) {
	a.State &= UNSET_MASK ^ mask
}

func (a *Actor) SetState(mask int) {
	a.State |= mask
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
	x = a.getClamped(a.x, l.TileWidth)
	y = int(a.y + a.Rate)
	if l.TestPixelPassable(a, x+a.Padding, y+l.TileHeight) &&
		l.TestPixelPassable(a, x+l.TileWidth-a.Padding, y+l.TileHeight) {
		if x == int(a.x) {
			// Only move once we've clamped.
			a.y = float64(y)
		}
		a.x = float64(x)
	}
}

func (a *Actor) moveUp(l *Level) {
	var (
		x int
		y int
	)
	x = a.getClamped(a.x, l.TileWidth)
	y = int(a.y - a.Rate)
	if l.TestPixelPassable(a, x+a.Padding, y) &&
		l.TestPixelPassable(a, x+l.TileWidth-a.Padding, y) {
		if x == int(a.x) {
			// Only move once we've clamped.
			a.y = float64(y)
		}
		a.x = float64(x)
	}
}

func (a *Actor) moveRight(l *Level) {
	var (
		x int
		y int
	)
	y = a.getClamped(a.y, l.TileHeight)
	x = int(a.x + a.Rate)
	if l.TestPixelPassable(a, x+l.TileWidth, y+a.Padding) &&
		l.TestPixelPassable(a, x+l.TileWidth, y+l.TileHeight-a.Padding) {
		if y == int(a.y) {
			// Only move once we've clamped.
			a.x = float64(x)
		}
		a.y = float64(y)
	}
}

func (a *Actor) moveLeft(l *Level) {
	var (
		x int
		y int
	)
	y = a.getClamped(a.y, l.TileHeight)
	x = int(a.x - a.Rate)
	if l.TestPixelPassable(a, x, y+a.Padding) &&
		l.TestPixelPassable(a, x, y+l.TileHeight-a.Padding) {
		if y == int(a.y) {
			// Only move once we've clamped.
			a.x = float64(x)
		}
		a.y = float64(y)
	}
}

type Player struct {
	*Actor
}

func (p *Player) SetDirection(dir int) {
	p.UnsetState(LEFT | RIGHT | UP | DOWN)
	if dir == RIGHT {
		p.flipX = true
	} else {
		p.flipX = false
	}
	p.SetState(dir)
}

func (p *Player) SetMovement(mov int) {
	p.UnsetState(WALKING | STOPPED)
	p.SetState(mov)
}

func NewPlayer(x float64, y float64, state int, offset int) (p *Player) {
	return &Player{
		Actor: &Actor{
			x:          x,
			y:          y,
			State:      state,
			textureRow: offset,
			Rate:       2.0,
			Padding:    12,
		},
	}
}

type Bomb struct {
	*Actor
	Elapsed time.Duration
	Expires time.Duration
	Radius  int
}

func NewBomb(x float64, y float64) (b *Bomb) {
	return &Bomb{
		Actor: &Actor{
			x:          x,
			y:          y,
			State:      BOMB,
			textureRow: 1,
			Padding:    12,
		},
		Elapsed: 0,
		Expires: time.Duration(5) * time.Second,
		Radius:  2,
	}
}

func (b *Bomb) Update(level *Level) bool {
	if b.Elapsed >= b.Expires {
		level.Explode(b)
		return false
	}
	return true
}

func (b *Bomb) AddTime(diff time.Duration) {
	b.Elapsed += diff
}

type Fire struct {
	*Actor
	Elapsed time.Duration
	Expires time.Duration
}

func NewFire(x float64, y float64) (f *Fire) {
	return &Fire{
		Actor: &Actor{
			x:          x,
			y:          y,
			State:      FLAME,
			textureRow: 2,
		},
		Elapsed: 0,
		Expires: time.Duration(500) * time.Millisecond,
	}
}

func (f *Fire) AddTime(diff time.Duration) {
	f.Elapsed += diff
}

func (f *Fire) Update(level *Level) bool {
	if f.Elapsed >= f.Expires {
		level.Extinguish(f)
		return false
	}
	return true
}

const (
	LEFT    = 1 << iota
	RIGHT   = 1 << iota
	UP      = 1 << iota
	DOWN    = 1 << iota
	WALKING = 1 << iota
	STOPPED = 1 << iota
	BOMB    = 1 << iota
	FLAME   = 1 << iota
	GOAL    = 1 << iota
)

var ACTOR_ANIMATIONS = map[int]*system.Animation{
	LEFT | STOPPED:                   system.Anim([]int{6}, 4),
	RIGHT | STOPPED:                  system.Anim([]int{6}, 4),
	UP | STOPPED:                     system.Anim([]int{3}, 4),
	DOWN | STOPPED:                   system.Anim([]int{0}, 4),
	LEFT | WALKING:                   system.Anim([]int{6, 7, 6, 8}, 4),
	RIGHT | WALKING:                  system.Anim([]int{6, 7, 6, 8}, 4),
	UP | WALKING:                     system.Anim([]int{3, 4, 3, 5}, 4),
	DOWN | WALKING:                   system.Anim([]int{0, 1, 0, 2}, 4),
	BOMB:                             system.Anim([]int{0, 1}, 4),
	GOAL:                             system.Anim([]int{2}, 4),
	FLAME:                            system.Anim([]int{0}, 4),
	FLAME | UP:                       system.Anim([]int{1}, 4),
	FLAME | DOWN:                     system.Anim([]int{2}, 4),
	FLAME | LEFT:                     system.Anim([]int{3}, 4),
	FLAME | RIGHT:                    system.Anim([]int{4}, 4),
	FLAME | UP | LEFT:                system.Anim([]int{5}, 4),
	FLAME | UP | RIGHT:               system.Anim([]int{6}, 4),
	FLAME | DOWN | RIGHT:             system.Anim([]int{7}, 4),
	FLAME | DOWN | LEFT:              system.Anim([]int{8}, 4),
	FLAME | DOWN | LEFT | RIGHT:      system.Anim([]int{9}, 4),
	FLAME | UP | DOWN | LEFT:         system.Anim([]int{10}, 4),
	FLAME | UP | LEFT | RIGHT:        system.Anim([]int{11}, 4),
	FLAME | UP | RIGHT | DOWN:        system.Anim([]int{12}, 4),
	FLAME | UP | RIGHT | DOWN | LEFT: system.Anim([]int{13}, 4),
	FLAME | RIGHT | LEFT:             system.Anim([]int{14}, 4),
	FLAME | UP | DOWN:                system.Anim([]int{15}, 4),
}
