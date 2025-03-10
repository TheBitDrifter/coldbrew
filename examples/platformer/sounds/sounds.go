package sounds

import blueprintclient "github.com/TheBitDrifter/blueprint/client"

// More audio players since steps happen quite rapidly
var Run = blueprintclient.SoundConfig{
	Path:             "run.wav",
	AudioPlayerCount: 6,
}

// Two audio players for two players
var Jump = blueprintclient.SoundConfig{
	Path:             "jump.wav",
	AudioPlayerCount: 2,
}

// Two audio players for two players
var Land = blueprintclient.SoundConfig{
	Path:             "land.wav",
	AudioPlayerCount: 2,
}
