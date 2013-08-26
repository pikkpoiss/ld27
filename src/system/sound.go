// Copyright 2013 Arne Roomann-Kurrik + Wes Goodman
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

package system

import (
	"fmt"
	"github.com/banthar/Go-SDL/mixer"
	"github.com/banthar/Go-SDL/sdl"
)

type Sound struct {
}

// Creates / Initializes the sound system.
func NewSound() (s *Sound, err error) {
	sdl.Init(sdl.INIT_AUDIO)
	if mixer.OpenAudio(mixer.DEFAULT_FREQUENCY, mixer.DEFAULT_FORMAT,
		mixer.DEFAULT_CHANNELS, 4096) != 0 {
		err = fmt.Errorf(sdl.GetError())
		return
	}
	s = &Sound{}
	return
}

// Plays the supplied sound file on repeat.
func (s *Sound) PlayMusic(path string) (err error) {
	var m = mixer.LoadMUS(path)
	if m == nil {
		err = fmt.Errorf(sdl.GetError())
		return
	}
	m.PlayMusic(-1)
	return
}

func (s *Sound) GetEffect(path string) (c *mixer.Chunk, err error) {
	c = mixer.LoadWAV(path)
	if c == nil {
		err = fmt.Errorf(sdl.GetError())
		return
	}
	return
}

func (s *Sound) GetMusic(path string) (m *mixer.Music, err error) {
	m = mixer.LoadMUS(path)
	if m == nil {
		err = fmt.Errorf(sdl.GetError())
		return
	}
	return
}

// Call to clean up after you're done.
func (s *Sound) Terminate() {
	sdl.Quit()
}
