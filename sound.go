package coldbrew

import (
	"bytes"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

// Sound represents an audio resource that can be played multiple times simultaneously
type Sound struct {
	name          string
	rawData       []byte
	audioCtx      *audio.Context
	players       []*audio.Player
	currentVolume float64
	isPaused      bool
}

// newSound creates a new Sound with multiple audio players for concurrent playback
func newSound(name string, data []byte, audioCtx *audio.Context, audioPlayerCount int) (Sound, error) {
	snd := Sound{
		name:          name,
		rawData:       data,
		audioCtx:      audioCtx,
		currentVolume: 1.0,
		players:       make([]*audio.Player, audioPlayerCount),
	}
	for i := range snd.players {
		// Create a new reader for each player
		audioReader := bytes.NewReader(data)
		// Decode the audio data for each player separately
		audioStream, err := wav.DecodeWithSampleRate(audioCtx.SampleRate(), audioReader)
		if err != nil {
			return Sound{}, err
		}
		// Create a new player with its own stream
		player, err := audioCtx.NewPlayer(audioStream)
		if err != nil {
			return Sound{}, err
		}
		snd.players[i] = player
	}
	return snd, nil
}

// GetPlayer returns the audio player at the specified index
func (s Sound) GetPlayer(i int) *audio.Player {
	return s.players[i]
}

// GetAnyAvailable returns an available player that is not currently playing,
// or the first player if all are in use
func (s Sound) GetAnyAvailable() *audio.Player {
	// First, try to find any player that's not currently playing
	for _, player := range s.players {
		if !player.IsPlaying() {
			return player
		}
	}
	// If i is out of bounds, return the first player as a last resort
	return s.players[0]
}
