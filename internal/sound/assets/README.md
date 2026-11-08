# Audio asset provenance

Most files in this directory are derived from sources published under the
[CC0 1.0 Universal public-domain dedication](https://creativecommons.org/publicdomain/zero/1.0/).
The jungle ambience and ground-strike recordings are licensed under
[CC BY 4.0](https://creativecommons.org/licenses/by/4.0/) and are attributed
below. Creators and source pages are preserved for every file to make the
license history auditable.

All derived files were converted with FFmpeg 8.1 to 48 kHz, stereo Ogg/Vorbis
for a single portable runtime format. Music and ordinary effects were otherwise
left unchanged. Creature calls received the pitch and filtering treatments
described below. `SHA256SUMS` records the exact derived files in this directory.

| Derived file | Source work and creator | Upstream file SHA-256 | Changes |
| --- | --- | --- | --- |
| `music_clear.ogg` | [Soliloquy](https://opengameart.org/content/soliloquy) by Matthew Pablo, CC BY 3.0 | `480aa2ae1211e217cf7506704a73d223747942ed1dfaa4e296f2c2ef592830a4` | The first six seconds were crossfaded with the tail to form a gentle loop, then the result was resampled and re-encoded; runtime gain is deliberately low |
| `music_threat.ogg` | [Foreboding Pixels](https://opengameart.org/content/foreboding-pixels), Sorth | `f54726be447aa69040171eccdcc7b13d46135d4caa635e7e255c990ffecb4770` | Loop version resampled and transcoded from MP3 |
| `music_failure.ogg` | [Game Over](https://opengameart.org/content/game-over), Kistol | `eb16aeb17fed5b13575702e082b54db00b296f6f4f2e9ee89c63aca895bc085e` | Resampled and re-encoded |
| `music_victory.ogg` | [Glorious victory fanfare NES](https://opengameart.org/content/glorious-victory-fanfare-nes), congusbongus | `7be73d492c038bce957bbbfbe479c7210fd58b50ae52f3ddf58b5514b0fd894e` | Resampled, converted to stereo, and re-encoded |
| `ambience_wind.ogg` | [wind whoosh loop](https://opengameart.org/content/wind-whoosh-loop), SketchMan3 | `0cfdbd3f21ed449689a9024264edf267d0037a6495813a7010827672ec191dac` | Re-encoded |
| `ambience_jungle.ogg` | [R4_00491_EXP01_Jurassic_world_Alien_planet](https://freesound.org/people/kevp888/sounds/797756/) by kevp888 (Kevin Luce), CC BY 4.0 | `00ee9e50858ed01f2db4d10c1226fa863e1aa7d2d847ecdc254cf61fdae87c74` | A 90-second section beginning at `1:00` was extracted and given a two-second tail-to-head crossfade for seamless looping; runtime gain is deliberately reduced so it remains behind gameplay cues |
| `ambience_rain.ogg` | [rain and thunders](https://opengameart.org/content/rain-and-thunders), kindland, `Dark_Rainy_Night(ambience).ogg` | `115c233f9f8bc793f12723030aaae6cfda923e565c32a8e54bb274b2a8199ee1` | A thunder-rich 75-second excerpt, resampled, given short boundary fades for clean looping, and re-encoded |
| `effect_balloon_pop.ogg` | [Balloon Sounds](https://opengameart.org/content/balloon-sounds), AntumDeluge | `a1f5dc1bfc466e232960888e482d116c2d9a7bab4ed1aedfbfea32d828ba5ad7` | Resampled, converted to stereo, and re-encoded |
| `effect_lightning_crack.ogg` | [25 CC0 bang / firework SFX](https://opengameart.org/content/25-cc0-bang-firework-sfx), rubberduck, `bang_04.ogg` | `fb3b76187c36e185a331d2250268132dbf2eb490a6c0414d69e3570c72b4f2e6` | Low boom reduced, upper crack emphasized, gain adjusted, and re-encoded |
| `effect_ground_strike.ogg` | [Lightning and Thunder Clap](https://freesound.org/people/akelley6/sounds/551425/) by akelley6, CC BY 4.0 | `48bfad24c75c0e46ce8ebaa78c171d4b33c08f3288bc3b364d8b85ce17ae35d4` | The first clap from approximately `0:01.55–0:04.05` was isolated, given a short boundary fade, resampled, and re-encoded |
| `effect_thunder.ogg` | [100 CC0 SFX #2](https://opengameart.org/content/100-cc0-sfx-2), rubberduck, `sfx100v2_thunder_01.ogg` | `8a4f1dee999546b363f989eff92174375eeae7c241dc9beaa0892fb131ef12a9` | Unchanged |
| `effect_basket_1.ogg` | [100 CC0 metal and wood SFX](https://opengameart.org/content/100-cc0-metal-and-wood-sfx), rubberduck, `wood_squeak_01.ogg` | `c6c90e02993b2e044937f344351b417aeff6841aa4cd683f208d85ececf23392` | Re-encoded |
| `effect_basket_2.ogg` | Same pack, `wood_squeak_02.ogg` | `841793c92034e419afec66538f7cf08f0f36f307aa3a1a5aaf7b664a0a178dee` | Re-encoded |
| `effect_pterosaur_1.ogg` | [Dinosaur call](https://freesound.org/people/Alessandrobraz/sounds/708641/) by Alessandrobraz, CC0 | `9196c7b7a3fe002f9a815ba9c87b682cb6ab187ab0cbb4c9bb93f99844b8f070` | First version isolated, reduced 4 dB, lightly band-limited, faded, resampled, and re-encoded |
| `effect_pterosaur_2.ogg` | Same source | `9196c7b7a3fe002f9a815ba9c87b682cb6ab187ab0cbb4c9bb93f99844b8f070` | Second version isolated, reduced 7 dB to match its sharper transient to the first call, lightly band-limited, faded, resampled, and re-encoded |
| `effect_dinosaur_1.ogg` | [80 CC0 creature SFX](https://opengameart.org/content/80-cc0-creature-sfx), rubberduck, `roar_01.ogg` | `e9d44326860a1c3aba917684faccfd967a4fefa284d7df773e05ae890ec242e0` | Slowed/pitched down and low-pass filtered |
| `effect_dinosaur_2.ogg` | Same pack, `roar_02.ogg` | `ceb532ee0fe775d6f86a821099b75f1790ab8022b963cb837c8c1fa4fffee5c2` | Slowed/pitched down and low-pass filtered |
| `effect_sauropod.ogg` | Same pack, `monster_01.ogg` | `e8bb5c9d3d25172666c60ed302d354144b67e66170bdf039ec88f4b6b0cd879f` | Substantially slowed/pitched down and low-pass filtered |
| `effect_carnivore.ogg` | Same pack, `roar_03.ogg` | `983396c808eb658ed3bd223ffaf9ea7f756dba7bc45ecc3fde117083026a4556` | Slowed/pitched down and low-pass filtered |

These are aesthetic selections, not claims about scientifically accurate
dinosaur vocalization. The processed bird and creature recordings are meant
to give the lighthearted game distinct, readable sound families.

## Audition candidates

These sources are noted for evaluation only and are not included in the game.

| Intended use | Candidate | Region of interest | License note |
| --- | --- | --- | --- |
| Alternate ground lightning strike | [Thunder and Lightning Strike](https://freesound.org/people/be-steele/sounds/712962/) by be-steele | Initial close strike around `0:18.6–0:20.1`; audition a longer cut if the following impulse is needed | CC BY 4.0; usable with attribution and a modification notice |
| Pterosaur call | [Pterosaur Sounds](https://freesound.org/people/The_Great_Q/sounds/733871/) by The_Great_Q | Approximately `0:09–0:18` | CC BY-NC 4.0; unsuitable for a potentially commercial release unless the creator grants separate permission |
