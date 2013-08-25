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
	"time"
)

type Level struct {
	Map        *system.TiledMap
	Camera     *Camera
	Cast       *Cast
	Player     *Player
	StartX     int
	StartY     int
	Goal       *Tile
	tiles      []Tile
	bombs      []*Bomb
	TileWidth  int
	TileHeight int
}

func LoadLevel(path string, cast *Cast) (out *Level, err error) {
	var (
		tm    *system.TiledMap
		cw    float64
		ch    float64
		tiles []Tile
		bombs []*Bomb
	)
	log.Printf("Loading level from %v\n", path)
	if tm, err = system.LoadMap(path); err != nil {
		return
	}
	cw = float64(tm.Width * tm.Tilewidth)
	ch = float64(tm.Height * tm.Tileheight)
	tiles = make([]Tile, tm.Width*tm.Height)
	bombs = make([]*Bomb, tm.Width*tm.Height)
	out = &Level{
		Map:        tm,
		Cast:       cast,
		Camera:     NewCamera(0, 0, cw, ch),
		TileWidth:  tm.Tilewidth,
		TileHeight: tm.Tileheight,
		tiles:      tiles,
		bombs:      bombs,
	}
	if err = out.parseTiles(); err != nil {
		return
	}
	if err = out.parseObjects(); err != nil {
		return
	}
	out.AddPlayerAtPixel(out.StartX, out.StartY)
	return
}

func (l *Level) Update(diff time.Duration) (err error) {
	var (
		layer *system.TiledLayer
	)
	if layer, err = l.getLayer("tilelayer", "Tiles"); err != nil {
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
	l.Cast.Update(l, diff)
	l.Player.Update(l)
	return
}

func (l *Level) AddPlayerAtPixel(x int, y int) {
	l.Player = NewPlayer(float64(x), float64(y), DOWN|STOPPED, 0)
	l.Cast.AddActor(l.Player)
}

func (l *Level) AddBombAtPixel(x int, y int) (err error) {
	var b *Bomb
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

func (l *Level) TestPixelPassable(x int, y int) bool {
	if t, err := l.getTileAtPixel(x, y); err != nil {
		return false
	} else if b, _ := l.getBombAtPixel(x, y); b != nil {
		// Can't walk over bomb
		return false
	} else {
		return TILES[t.Type].Passable
	}
}

func (l *Level) Explode(b *Bomb) {
	for i, bomb := range l.bombs {
		if b != bomb {
			continue
		}
		l.bombs[i] = nil
		l.Cast.RemoveActor(b)
		break
	}
}

func (l *Level) getLayer(t string, n string) (out *system.TiledLayer, err error) {
	for i, _ := range l.Map.Layers {
		if l.Map.Layers[i].Type != t && l.Map.Layers[i].Name != n {
			continue
		}
		out = &l.Map.Layers[i]
		return
	}
	err = fmt.Errorf("Could not find layer with type %v and name %v", t, n)
	return
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

func (l *Level) getBombAtPixel(x int, y int) (b *Bomb, err error) {
	var i = l.getPixelIndex(x, y)
	if i >= len(l.tiles) || i < 0 {
		err = fmt.Errorf("Pixel at (%v, %v) out of range", x, y)
		return
	}
	b = l.bombs[i]
	return
}

func (l *Level) getTileAtPixel(x int, y int) (t *Tile, err error) {
	var i = l.getPixelIndex(x, y)
	if i >= len(l.tiles) || i < 0 {
		err = fmt.Errorf("No tile at (%v,%v)", x, y)
		return
	}
	t = &l.tiles[i]
	return
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
	if layer, err = l.getLayer("tilelayer", "Tiles"); err != nil {
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
	if layer, err = l.getLayer("objectlayer", "Objects"); err != nil {
		return
	}
	for _, obj := range layer.Objects {
		switch obj.Type {
		case "player":
			l.StartX = obj.X
			l.StartY = obj.Y
		case "goal":
			l.Goal, err = l.getTileAtPixel(obj.X, obj.Y)
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
)

var TILES = map[int]TileType{
	TILE_GRASS: TileType{
		Anim:      system.Anim([]int{1}, 4),
		Passable:  true,
		Breakable: false,
	},
	TILE_STONE: TileType{
		Anim:      system.Anim([]int{2}, 4),
		Passable:  false,
		Breakable: false,
	},
	TILE_BRICK: TileType{
		Anim:      system.Anim([]int{3}, 16),
		Passable:  false,
		Breakable: true,
	},
}

type TileType struct {
	Anim      *system.Animation
	Passable  bool
	Breakable bool
}

type Tile struct {
	X    int
	Y    int
	Type int
}

