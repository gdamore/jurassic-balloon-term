# Jurassic Balloon — Third-Party Notices

Jurassic Balloon is created, designed, and programmed by Garrett D'Amore and
is licensed under the Apache License 2.0. This file records the principal
runtime libraries and the asset attributions that must accompany a release.
The complete audio provenance, source URLs, hashes, and modification notes are
maintained in [`internal/sound/assets/README.md`](internal/sound/assets/README.md).

## Creative Commons attribution assets

### “Soliloquy”

- Creator: Matthew Pablo
- Source: <https://opengameart.org/content/soliloquy>
- License: [CC BY 3.0](https://creativecommons.org/licenses/by/3.0/)
- Changes: the first six seconds were crossfaded with the tail to make a loop;
  the result was resampled and re-encoded.

### “R4_00491_EXP01_Jurassic_world_Alien_planet”

- Creator: kevp888 (Kevin Luce)
- Source: <https://freesound.org/people/kevp888/sounds/797756/>
- License: [CC BY 4.0](https://creativecommons.org/licenses/by/4.0/)
- Changes: a 90-second excerpt beginning at 1:00 was extracted, crossfaded,
  resampled, and re-encoded.

### “Lightning and Thunder Clap”

- Creator: akelley6
- Source: <https://freesound.org/people/akelley6/sounds/551425/>
- License: [CC BY 4.0](https://creativecommons.org/licenses/by/4.0/)
- Changes: the first clap was isolated, faded at its boundaries, resampled,
  and re-encoded.

## Additional audio acknowledgements

Additional audio was sourced from public-domain or CC0 releases by Sorth,
Kistol, congusbongus, kindland, AntumDeluge, rubberduck, Alessandrobraz, and
SketchMan3. Exact works and transformations are listed in the audio provenance
document linked above.

## Runtime libraries

- [Tcell](https://github.com/gdamore/tcell), Apache License 2.0.
- [Oto](https://github.com/ebitengine/oto), Apache License 2.0.
- [oggvorbis](https://github.com/jfreymuth/oggvorbis), MIT License,
  Copyright (c) 2016 Johann Freymuth.

The Apache License 2.0 text is included in this repository as [`LICENSE`](LICENSE).
The MIT notice for oggvorbis follows.

> MIT License
>
> Copyright (c) 2016 Johann Freymuth
>
> Permission is hereby granted, free of charge, to any person obtaining a copy
> of this software and associated documentation files (the “Software”), to deal
> in the Software without restriction, including without limitation the rights
> to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
> copies of the Software, and to permit persons to whom the Software is
> furnished to do so, subject to the following conditions:
>
> The above copyright notice and this permission notice shall be included in all
> copies or substantial portions of the Software.
>
> THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
> IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
> FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
> AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
> LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
> OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
> SOFTWARE.

## Design lineage

The fixed-step loop and sprite-engine concepts were informed by Garrett
D'Amore's [Escape from Proxima 5](https://github.com/gdamore/proxima5), also
released under the Apache License 2.0. Jurassic Balloon reimplements and
reconsiders those ideas rather than copying its source tree.

