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
	"strings"
)

type Menu interface {
	SelectNext()
	SelectPrev()
	Select(int)
	Choose()
	GetMap() *system.TiledMap
	Draw()
}

type MenuHandler func(selection int)

type BasicMenu struct {
	Map      *system.TiledMap
	Camera   *Camera
	Handler  MenuHandler
	buttons  []*Button
	selected int
}

func LoadMenu(path string, handler MenuHandler) (out *BasicMenu, err error) {
	var (
		tm *system.TiledMap
		cw float64
		ch float64
	)
	log.Printf("Loading menu from %v\n", path)
	if tm, err = system.LoadMap(path); err != nil {
		return
	}
	cw = float64(tm.Width * tm.Tilewidth)
	ch = float64(tm.Height * tm.Tileheight)
	out = &BasicMenu{
		Map:     tm,
		Camera:  NewCamera(0, 0, cw, ch),
		Handler: handler,
	}
	if err = out.parseButtons(); err != nil {
		return
	}
	out.Select(0)
	return
}

func (m *BasicMenu) parseButtons() (err error) {
	var (
		layer *system.TiledLayer
		gid   int
		v     int
		i     int
	)
	if layer, err = m.Map.GetLayer("tilelayer", "Buttons"); err != nil {
		return
	}
	for i, gid = range layer.Data {
		if gid == 0 {
			continue
		}
		if v, err = m.Map.GetTilesetOffset(gid); err != nil {
			return
		}
		log.Printf("Got button: %v %v\n", i, v)
		m.buttons = append(m.buttons, &Button{
			Type:  v,
			index: i,
			gid:   gid,
		})
	}
	return
}

func (m *BasicMenu) Draw() {
}

func (m *BasicMenu) GetMap() *system.TiledMap {
	return m.Map
}

func (m *BasicMenu) Select(i int) {
	var (
		layer *system.TiledLayer
		b     *Button
		err   error
	)
	if layer, err = m.Map.GetLayer("tilelayer", "Buttons"); err != nil {
		return
	}
	if len(m.buttons) == 0 {
		return
	}
	b = m.buttons[m.selected]
	layer.Data[b.index] = b.gid
	i = i % len(m.buttons)
	m.selected = i
	b = m.buttons[i]
	layer.Data[b.index] = b.gid + 1
}

func (m *BasicMenu) SelectNext() {
	m.Select(m.selected + 1)
}

func (m *BasicMenu) SelectPrev() {
	m.Select(m.selected - 1 + len(m.buttons))
}

func (m *BasicMenu) Choose() {
	m.Handler(m.buttons[m.selected].Type)
}

type OverlayMenu struct {
	*BasicMenu
	Text  []string
	Curr  int
	Font  *system.Font
	TextX float64
	TextY float64
}

func LoadOverlayMenu(path string, handler MenuHandler, font *system.Font) (out *OverlayMenu, err error) {
	var menu *BasicMenu
	if menu, err = LoadMenu(path, handler); err != nil {
		return
	}
	out = &OverlayMenu{
		BasicMenu: menu,
		Text:      []string{},
		Curr:      0,
		Font:      font,
	}
	out.parseObjects()
	return
}

func (m *OverlayMenu) parseObjects() {
	var (
		layer *system.TiledLayer
		err   error
	)
	if layer, err = m.Map.GetLayer("objectgroup", "Text"); err != nil {
		return
	}
	for _, obj := range layer.Objects {
		switch obj.Type {
		case "text":
			m.TextX = float64(obj.X)
			m.TextY = float64(obj.Y)
		}
	}
}

func (m *OverlayMenu) SetText(text []string) {
	m.Text = text
	m.Curr = 0
}

func (m *OverlayMenu) SelectNext() {
	m.advance()
}

func (m *OverlayMenu) SelectPrev() {
	m.advance()
}

func (m *OverlayMenu) Choose() {
	m.advance()
}

func (m *OverlayMenu) Draw() {
	if len(m.Text) > m.Curr {
		var y = m.TextY * 2
		var lines = strings.Split(m.Text[m.Curr], "\n")
		for _, line := range lines {
			// Scaling is a hack, since we're pixel doubling
			m.Font.Printf(m.TextX*2, y, "%v", line)
			y += 32
		}
	}
}

func (m *OverlayMenu) advance() {
	m.Curr += 1
	if m.Curr >= len(m.Text) {
		m.Handler(BUTTON_START)
	}
}

type Button struct {
	Type  int
	gid   int
	index int
}

const (
	BUTTON_START = 0
	BUTTON_EXIT  = 2
)
