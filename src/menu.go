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

type MenuHandler func(selection int)

type Menu struct {
	Map      *system.TiledMap
	Camera   *Camera
	Handler  MenuHandler
	buttons  []*Button
	selected int
}

func LoadMenu(path string, handler MenuHandler) (out *Menu, err error) {
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
	out = &Menu{
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

func (m *Menu) parseButtons() (err error) {
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

func (m *Menu) Select(i int) {
	var (
		layer *system.TiledLayer
		b     *Button
		err   error
	)
	if layer, err = m.Map.GetLayer("tilelayer", "Buttons"); err != nil {
		return
	}
	b = m.buttons[m.selected]
	layer.Data[b.index] = b.gid
	i = i % len(m.buttons)
	m.selected = i
	b = m.buttons[i]
	layer.Data[b.index] = b.gid + 1
}

func (m *Menu) SelectNext() {
	m.Select(m.selected + 1)
}

func (m *Menu) SelectPrev() {
	m.Select(m.selected - 1 + len(m.buttons))
}

func (m *Menu) Choose() {
	m.Handler(m.buttons[m.selected].Type)
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
