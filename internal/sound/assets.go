// Copyright 2026 Garrett D'Amore
// Licensed under the Apache License, Version 2.0.

package sound

import (
	"bytes"
	"embed"
	"fmt"
	"sync"

	"github.com/gdamore/jurassic-balloon-term/internal/audio"
	"github.com/jfreymuth/oggvorbis"
)

//go:embed assets/*.ogg
var assetFiles embed.FS

var decodedAssets = sync.OnceValue(func() map[string]audio.Clip {
	names := []string{
		"music_clear.ogg",
		"music_threat.ogg",
		"music_failure.ogg",
		"music_victory.ogg",
		"ambience_jungle.ogg",
		"ambience_rain.ogg",
		"effect_balloon_pop.ogg",
		"effect_lightning_crack.ogg",
		"effect_ground_strike.ogg",
		"effect_thunder.ogg",
		"effect_basket_1.ogg",
		"effect_basket_2.ogg",
		"effect_pterosaur_1.ogg",
		"effect_pterosaur_2.ogg",
		"effect_dinosaur_1.ogg",
		"effect_dinosaur_2.ogg",
		"effect_sauropod.ogg",
		"effect_carnivore.ogg",
	}
	clips := make(map[string]audio.Clip, len(names))
	for _, name := range names {
		data, err := assetFiles.ReadFile("assets/" + name)
		if err != nil {
			panic(fmt.Sprintf("sound: read embedded asset %s: %v", name, err))
		}
		pcm, format, err := oggvorbis.ReadAll(bytes.NewReader(data))
		if err != nil {
			panic(fmt.Sprintf("sound: decode embedded asset %s: %v", name, err))
		}
		if format.SampleRate != audio.DefaultSampleRate || format.Channels != audio.DefaultChannels {
			panic(fmt.Sprintf("sound: asset %s is %d Hz/%d channels, want %d Hz/%d channels",
				name, format.SampleRate, format.Channels, audio.DefaultSampleRate, audio.DefaultChannels))
		}
		clips[name] = audio.Clip{PCM: pcm}
	}
	return clips
})

func recorded(name string) audio.Clip {
	clip, ok := decodedAssets()[name]
	if !ok {
		panic("sound: unknown embedded asset " + name)
	}
	return clip
}
