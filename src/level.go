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
	"fmt"
	"log"
	"strings"
	"time"
)

type Level struct {
	Map        *system.TiledMap
	Camera     *Camera
	Cast       *Cast
	Player     *Player
	Goal       *Actor
	tiles      []Tile
	bombs      []*Bomb
	fire       []*Fire
	TileWidth  int
	TileHeight int
	Won        bool
	Died       bool
}

func LoadLevel(path string, cast *Cast) (out *Level, err error) {
	var (
		tm    *system.TiledMap
		cw    float64
		ch    float64
		count int
	)
	log.Printf("Loading level from %v\n", path)
	if tm, err = system.LoadMap(path); err != nil {
		return
	}
	cw = float64(tm.Width * tm.Tilewidth)
	ch = float64(tm.Height * tm.Tileheight)
	count = tm.Width * tm.Height
	out = &Level{
		Map:        tm,
		Cast:       cast,
		Camera:     NewCamera(0, 0, cw, ch),
		TileWidth:  tm.Tilewidth,
		TileHeight: tm.Tileheight,
		Won:        false,
		Died:       false,
		tiles:      make([]Tile, count),
		bombs:      make([]*Bomb, count),
		fire:       make([]*Fire, count),
	}
	if err = out.parseTiles(); err != nil {
		return
	}
	if err = out.parseObjects(); err != nil {
		return
	}
	return
}

func (l *Level) Update(diff time.Duration) (err error) {
	var (
		layer *system.TiledLayer
	)
	if layer, err = l.Map.GetLayer("tilelayer", "Tiles"); err != nil {
		return
	}
	for _, t := range TILES {
		t.Anim.Next()
	}
	for i, t := range l.tiles {
		layer.Data[i] = TILES[t.Type].Anim.Curr()
	}
	for _, b := range l.bombs {
		if b != nil {
			b.AddTime(diff)
			b.Update(l)
		}
	}
	for i, f := range l.fire {
		if f != nil {
			f.AddTime(diff)
			f.Update(l)
			l.setFireDirection(i)
		}
	}
	l.Cast.Update(l, diff)
	l.Player.Update(l)
	if l.checkActorBurned(l.Player.Actor) {
		l.Died = true
	}
	if l.Cast.Overlaps(l.Player.Actor, l.Goal) {
		l.Won = true
	}
	return
}

func (l *Level) AddBombFromActor(a *Actor) {
	var (
		x   = int(a.X() + float64(l.TileWidth)/2.0)
		y   = int(a.Y() + float64(l.TileHeight)/2.0)
		err error
		b   *Bomb
	)
	if b, err = l.addBombAtPixel(x, y); err != nil {
		return
	}
	a.Bomb = b
}

func (l *Level) checkActorBurned(a *Actor) bool {
	var (
		x = int(a.X() + float64(l.TileWidth)/2.0)
		y = int(a.Y() + float64(l.TileHeight)/2.0)
		i int
	)
	i = l.getPixelIndex(x, y)
	if l.fire[i] != nil {
		return true
	}
	return false
}

func (l *Level) addBombAtPixel(x int, y int) (b *Bomb, err error) {
	if b, err = l.getBombAtPixel(x, y); err != nil || b != nil {
		return
	}
	var i = l.getPixelIndex(x, y)
	x, y = l.getPixelFromIndex(i)
	b = NewBomb(float64(x), float64(y))
	l.bombs[i] = b
	l.Cast.AddActor(b)
	return
}

func (l *Level) TestPixelPassable(a *Actor, x int, y int) bool {
	var (
		t   *Tile
		err error
	)
	if t, err = l.getTileAtPixel(x, y); err != nil {
		return false
	} else if b, _ := l.getBombAtPixel(x, y); b != nil {
		return b == a.Bomb
	} else {
		a.Bomb = nil
	}
	return TILES[t.Type].Passable
}

func (l *Level) Explode(b *Bomb) {
	for i, bomb := range l.bombs {
		if b != bomb {
			continue
		}
		var (
			x int
			y int
		)
		x = l.iToX(i)
		y = l.iToY(i)
		l.bombs[i] = nil
		l.Cast.RemoveActor(b)
		if l.addFire(x, y) {
			l.addFireColumn(x, y, b.Radius, 1, 0)
			l.addFireColumn(x, y, b.Radius, -1, 0)
			l.addFireColumn(x, y, b.Radius, 0, 1)
			l.addFireColumn(x, y, b.Radius, 0, -1)
		}
	}
}

func (l *Level) Extinguish(f *Fire) {
	for i, fire := range l.fire {
		if f != fire {
			continue
		}
		l.Cast.RemoveActor(f)
		l.fire[i] = nil
		break
	}
}

func (l *Level) GetDescription() (text []string) {
	var (
		ok  bool
		raw string
	)
	if raw, ok = l.Map.Properties["text"]; !ok {
		return
	}
	raw = strings.Replace(raw, "[BR]", "\n", -1)
	text = strings.Split(raw, "|")
	return
}

func (l *Level) setFireDirection(i int) {
	var (
		x   int
		y   int
		f1  *Fire
		f2  *Fire
		err error
	)
	if f1, err = l.getFire(i); err != nil || f1 == nil {
		return
	}
	x = l.iToX(i)
	y = l.iToY(i)
	if f2, err = l.getFire(l.xyToI(x, y-1)); err == nil && f2 != nil {
		f1.SetState(UP)
		f2.SetState(DOWN)
	} else {
		f1.UnsetState(UP)
	}
	if f2, err = l.getFire(l.xyToI(x, y+1)); err == nil && f2 != nil {
		f1.SetState(DOWN)
		f2.SetState(UP)
	} else {
		f1.UnsetState(DOWN)
	}
	if f2, err = l.getFire(l.xyToI(x-1, y)); err == nil && f2 != nil {
		f1.SetState(LEFT)
		f2.SetState(RIGHT)
	} else {
		f1.UnsetState(LEFT)
	}
	if f2, err = l.getFire(l.xyToI(x+1, y)); err == nil && f2 != nil {
		f1.SetState(RIGHT)
		f2.SetState(LEFT)
	} else {
		f1.UnsetState(RIGHT)
	}
}

func (l *Level) addFireColumn(x int, y int, r int, stepx int, stepy int) {
	for i := 1; i <= r; i++ {
		if !l.addFire(x+(stepx*i), y+(stepy*i)) {
			break
		}
	}
}

func (l *Level) addFire(x int, y int) bool {
	if x < 0 || y < 0 || x >= l.Map.Width || y >= l.Map.Height {
		return false
	}
	var (
		i         = l.xyToI(x, y)
		t         *Tile
		f         *Fire
		b         *Bomb
		err       error
		ttype     TileType
		continues bool
		px        int
		py        int
	)
	continues = true
	if f, err = l.getFire(i); err != nil || f != nil {
		return continues
	}
	if b, err = l.getBomb(i); err != nil || b != nil {
		b.AddTime(time.Duration(1000) * time.Hour)
		return continues
	}
	if t, err = l.getTile(i); err != nil {
		return continues
	}
	ttype = TILES[t.Type]
	if ttype.StopsFire {
		continues = false
		if ttype.Breakable {
			t.Type = ttype.NextState
		} else {
			return continues
		}
	}
	px, py = l.getPixelFromIndex(i)
	f = NewFire(float64(px), float64(py))
	l.fire[i] = f
	l.Cast.AddActor(f)
	l.setFireDirection(i)
	return continues
}

func (l *Level) getPixelIndex(x int, y int) (i int) {
	x = x / l.Map.Tilewidth
	y = y / l.Map.Tileheight
	i = l.Map.Width*y + x
	return i
}

func (l *Level) getPixelFromIndex(i int) (x int, y int) {
	x = l.iToX(i) * l.TileWidth
	y = l.iToY(i) * l.TileHeight
	return
}

func (l *Level) getFire(i int) (f *Fire, err error) {
	if i >= len(l.fire) || i < 0 {
		err = fmt.Errorf("Pixel at (%v) out of range", i)
		return
	}
	f = l.fire[i]
	return
}

func (l *Level) getBomb(i int) (b *Bomb, err error) {
	if i >= len(l.bombs) || i < 0 {
		err = fmt.Errorf("Pixel at (%v) out of range", i)
		return
	}
	b = l.bombs[i]
	return
}

func (l *Level) getBombAtPixel(x int, y int) (*Bomb, error) {
	var i = l.getPixelIndex(x, y)
	return l.getBomb(i)
}

func (l *Level) getTile(i int) (t *Tile, err error) {
	if i >= len(l.tiles) || i < 0 {
		err = fmt.Errorf("Pixel at (%v) out of range", i)
		return
	}
	t = &l.tiles[i]
	return
}

func (l *Level) getTileAtPixel(x int, y int) (*Tile, error) {
	var i = l.getPixelIndex(x, y)
	return l.getTile(i)
}

func (l *Level) xyToI(x int, y int) int {
	return y*l.Map.Width + x
}

func (l *Level) iToX(i int) int {
	return i % l.Map.Width
}

func (l *Level) iToY(i int) int {
	return i / l.Map.Width
}

func (l *Level) parseTiles() (err error) {
	var (
		layer *system.TiledLayer
	)
	if layer, err = l.Map.GetLayer("tilelayer", "Tiles"); err != nil {
		return
	}
	for i, v := range layer.Data {
		l.tiles[i] = Tile{
			X:    l.iToX(i),
			Y:    l.iToY(i),
			Type: v,
		}
	}
	return
}

func (l *Level) parseObjects() (err error) {
	var (
		layer *system.TiledLayer
	)
	if layer, err = l.Map.GetLayer("objectgroup", "Objects"); err != nil {
		return
	}
	for _, obj := range layer.Objects {
		switch obj.Type {
		case "player":
			l.Player = NewPlayer(float64(obj.X), float64(obj.Y), DOWN|STOPPED, 0)
			l.Cast.AddActor(l.Player)
		case "goal":
			l.Goal = NewActor(float64(obj.X), float64(obj.Y), GOAL, 1)
			l.Cast.AddActor(l.Goal)
		}
		if err != nil {
			return
		}
	}
	return
}

const (
	TILE_GRASS = 1 + iota
	TILE_STONE
	TILE_BRICK
	TILE_BREAKABLE_STONE_1
	TILE_BREAKABLE_STONE_2
)

var TILES = map[int]TileType{
	TILE_GRASS: TileType{
		Anim:      system.Anim([]int{1}, 4),
		Passable:  true,
		Breakable: false,
		StopsFire: false,
	},
	TILE_STONE: TileType{
		Anim:      system.Anim([]int{2}, 4),
		Passable:  false,
		Breakable: false,
		StopsFire: true,
	},
	TILE_BREAKABLE_STONE_1: TileType{
		Anim:      system.Anim([]int{4}, 4),
		Passable:  false,
		Breakable: true,
		StopsFire: true,
		NextState: TILE_BREAKABLE_STONE_2,
	},
	TILE_BREAKABLE_STONE_2: TileType{
		Anim:      system.Anim([]int{5}, 4),
		Passable:  false,
		Breakable: true,
		StopsFire: true,
		NextState: TILE_GRASS,
	},
	TILE_BRICK: TileType{
		Anim:      system.Anim([]int{3}, 16),
		Passable:  false,
		Breakable: true,
		StopsFire: true,
		NextState: TILE_GRASS,
	},
}

type TileType struct {
	Anim      *system.Animation
	Passable  bool
	Breakable bool
	StopsFire bool
	NextState int
}

type Tile struct {
	X    int
	Y    int
	Type int
}
