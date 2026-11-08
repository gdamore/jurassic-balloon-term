# Jurassic Balloon — Game Design Handoff

Status: working design handoff for Codex
Source: the Jurassic Balloon portion of the ChatGPT conversation titled
“App Store Fees Comparison” (`6aef8d4f-a134-83ea-bd80-8ccd68ea2fd1`)
Source transcript: [JURASSIC_BALLOON_SOURCE_TRANSCRIPT.md](JURASSIC_BALLOON_SOURCE_TRANSCRIPT.md)

This document distinguishes decisions Garrett explicitly made from proposals
that remain open. It should guide implementation without turning exploratory
ChatGPT suggestions into requirements.

## One-sentence pitch

Jurassic Balloon is a lighthearted, real-time, forward-scrolling arcade game
about surviving a prehistoric landscape in a leaky hot-air balloon improvised
from a parachute, using a simple burner control to manage delayed vertical
motion while avoiding aerial, terrain, and altitude hazards.

## Design priorities

1. **Simple input, sophisticated situations.** The central decision is when to
   activate and release the burner.
2. **Predictable, inertial flight.** The balloon should reward anticipation,
   not twitchy correction or arbitrary randomness.
3. **Telegraphed hazards.** Weather and environmental forces need visible
   signals so failure feels earned.
4. **Fast, funny failure.** In the terminal game, any hazard collision pops the
   balloon. The pop should have personality and lead quickly to another try.
5. **Personality over graphical complexity.** Compact, readable sprites and a
   light tone matter more than visual detail.
6. **Terminal and mobile may diverge.** Share ideas where useful; do not force
   either version into an abstraction designed for the other.

## Confirmed decisions

### Premise and tone

- The working title is **Jurassic Balloon**.
- The player crashed in a prehistoric setting and made a balloon from the
  parachute.
- The balloon is crude and leaky.
- The game should be lighthearted and accessible to a casual audience.
- A visible character in the basket is desirable, within the graphical limits
  of the target platform.
- More serious or sophisticated spinoffs, including a possible first-person VR
  interpretation, are separate future ideas and must not expand this project.

### Core structure

- The game is a casual side-scroller with automatic, forward-only travel.
- It takes inspiration from the immediacy of an endless runner, but can contain
  authored levels, checkpoints/save points, distinct challenges, obstacles,
  enemies, and possibly power-ups.
- Normal horizontal speed is constant.
- Environmental conditions can change speed; storm levels are the clearest
  example and are expected to increase speed and difficulty.
- The terminal version is unequivocally **real time**, not turn based.
- During prototype playtesting, number keys `1` through `5` immediately start
  the corresponding campaign stage from either the title screen or gameplay.

### Initial campaign structure

- The first campaign has five stages totaling about fifteen minutes for a
  clean run. Deaths will make a first completion longer.
- Progress is called **distance traveled** even while constant world speed
  makes it directly proportional to elapsed play time.
- The initial target stages are Fern Valley (2:00 / 384 units), Red-Eye Run
  (2:30 / 480), Thunderhead Pass (3:00 / 756), Long-Neck Grove (3:30 / 672),
  and Carnivore Country (4:00 / 768).
- A visible completion flag marks each stage save/unlock boundary. Crossing an
  intermediate flag continues directly into the next stage without pausing or
  returning the balloon to the ground; altitude, heat, burner state, and
  vertical velocity carry across. The crossed flag and visible outgoing
  scenery, clouds, and enemies keep moving left instead of disappearing;
  carried enemies remain animated and dangerous until they leave the screen.
  Replacement stage clouds enter from beyond the right edge. Only the final
  campaign flag stops play.
  Checkpoints within longer stages are authored section boundaries, initially
  spaced roughly 60–75 seconds apart. Death should eventually offer restart at
  the latest unlocked stage, and unlocked stages remain directly selectable.
- Fuel is not required for this first campaign iteration. Reconsider it only
  after flight, hazards, stage pacing, and restarts are working together.
- The current Fern Valley prototype introduces a deterministic grove of tall,
  dangerous ordinary trees in its second half. This breaks up the otherwise
  safe bottom-third flight lane without raising the floor across the whole
  stage.
- Its first half uses harmless, low fern scenery and three pairs of compact
  three-horn grazers to establish the prehistoric setting before the dangerous
  grove. The grazers are static, fatal ground hazards with raised horns, a
  downward beak, a distinct head, and separate frill, body, belly, and leg
  shading. The grove mixes broadleaf, palm, and cycad silhouettes while
  retaining comparable collision reach.
- Red-Eye Run introduces a green, red-eyed pterosaur alongside the ordinary
  variant. It uses a larger 18×6 silhouette, cruises at twice the ordinary
  speed, and has a stronger but deliberately capped lunge rather than simply
  doubling the ordinary boost. Both pterosaur variants visibly open their
  mouths for the lunge frame that emits the cry. Fern Valley has one authored
  pair; Red-Eye Run has one staggered three-member group containing a single
  red-eye. These flocks are limited encounters, not the default spawn pattern.
  Ferns, a mixed line of
  dangerous trees, and three more grazer pairs keep its landscape populated.
  Grazer feet overlap the terrain surface slightly so uneven ground never
  makes them appear to hover.
- Thunderhead Pass uses persistent screen-wide animated rain, four
  independently phased thunderheads, sparse fern scenery, six authored tree
  hazards, and one grazer pair around distance 318 to break up a former lull.
  Ordinary clouds and both pterosaur variants are absent. A dark, four-row
  collidable thundercloud ceiling uses rounded crown clusters, rolling
  half-block undersides, and highlight/body/shadow color variation so it reads
  as weather rather than a slab. Its texture drifts laterally at 1.9
  cells/second, between the far-terrain and foreground rates, while retaining a
  continuous cloud mass without sky-colored holes. It makes the top rows dangerous and telegraphs
  deterministic pseudo-random lightning strikes that extend five cells below
  it. Two much larger thunderheads occupy the lower half of the
  playfield. Each briefly draws a harmless bright-white preview of the eventual bolt
  before lightning reaches the terrain and remains live while crossing the
  balloon's lane, forcing the player to climb rather than settle at one
  altitude. Warning and strike masks end one cell inside the live foreground
  surface rather than continuing behind it. The brilliant yellow strike drops
  rapidly through pale cyan to electric blue, disappears shortly after
  entering the balloon's lane, and terminates in a brief ground-impact burst
  so the terrain contact reads clearly. Ordinary thunderhead and ceiling bolts
  use the same quick yellow-to-blue decay. Only the bright-yellow portion is
  damaging; cyan and blue afterglow frames are visual feedback only. The whole foreground scrolls at
  4.2 cells/second rather than
  the ordinary 3.2 to convey strong wind, while far terrain retains its slower
  parallax ratio. Entering the stage switches to storm ambience and threat
  music; leaving it restores the ordinary daytime soundscape. Both boundaries
  use a roughly two-second moving weather front: existing terrain and visible
  scenery continue scrolling, rain and the thundercloud ceiling sweep across
  the screen, and wind speed eases between the ordinary and storm values.

A post-campaign endless run that unlocks after completing all five stages is a
strong next-mode idea, not yet a committed implementation. It would measure
survival time/distance and steadily increase encounter difficulty rather than
use a completion flag.

Later-stage concepts to explore, not current requirements:

- Give Carnivore Country its own more threatening score rather than reusing
  the clear-weather orchestral track.
- Add an over-water stage with marine predators and strongly telegraphed
  waterspouts as vertical hazards.
- Add Lava Fields with pterosaurs, visible volcanic ejecta, and visually
  explicit updrafts; avoid surprise forces and pyroclastic-eruption clutter.

### Flight and controls

- Gravity matters.
- The player controls vertical motion through the burner.
- The burner has a slightly delayed effect. Holding it continues the effect,
  with an upper limit on vertical response.
- The desired feel is balloon-like: inertial, predictable, and simple enough
  that the player can also concentrate on hazards.
- The normal control vocabulary should remain extremely small, ideally burner
  on/off as the primary gameplay input.
- A start screen precedes the run and offers independent music and sound-effect
  toggles plus room for a short premise blurb.
- Pterosaur lunges use two restrained CC0 designed-dinosaur calls rather than
  a synthetic squawk or shrill bird scream. Calls are limited to one every 2.2
  seconds to prevent each flap from becoming an audio interruption; this is
  characterful game audio, not a claim about reconstructed pterosaur sound.
- Ordinary daytime ambience uses a continuous CC0 field recording of real
  birds rather than randomized processed creature calls.
- Clear-weather stages use a low-mixed, gently looped version of Matthew
  Pablo's slower orchestral "Soliloquy." Fern Valley retains music, but it
  should sit behind the subdued jungle ambience and gameplay cues rather than
  compete for attention.
- The balloon begins resting on the terrain. Stage time, distance, and hazards
  remain stopped until sustained burner heat produces actual lift, so the
  player is not dropped into an immediate fall.
- The terminal balloon uses a compact three-frame, true-color burner flame.
  Terminal emoji are avoided because their width and presentation vary by
  terminal and font.
- Modern terminal support for key-release events must be recognized. Kitty
  keyboard protocol and Win32 input mode can provide genuine press/release
  behavior. Capability detection and a legacy fallback are appropriate.

The exact equations, heat/leak state, acceleration cap, velocity cap, and tuning
constants are **not yet decisions**. They must be found through playtesting.

One concrete candidate proposed in the conversation is a normalized heat/lift
state `h` that rises gradually while the burner is held and decays after
release:

```text
burner on:   dh/dt = (1 - h) / tauHeat
burner off:  dh/dt = -h / tauCool
vertical:    ay = gravity - maxLift*h + environmentAY
```

The assistant suggested testing a roughly 0.3–0.5 second burner response and a
0.7–1.2 second cooling/leak response, with moderate capped ascent and somewhat
faster capped descent. These are prototype starting values only, not agreed
physics.

### Hazards and environmental communication

- Dangerous high-altitude conditions keep the player below an upper region.
- Ground-level danger includes carnivorous dinosaurs and terrain hazards.
- Carnivore Country should include squat, frilled Dilophosaurus-inspired
  spitters. At fixed intervals they telegraph and launch damaging projectiles
  upward, creating a predictable timing hazard for low and middle flight. The
  venom-spitting behavior is a playful cinematic invention, not a claim about
  the real animal. Exact cadence, projectile path, and warning animation remain
  playtesting decisions.
- Aerial threats include pterosaurs.
- Environmental forces must be signaled rather than appearing arbitrarily.
- Rain is the chosen storm signal.
- Volcanoes may create lift while their ejecta creates a simultaneous damage
  hazard. Both the volcano and its effects should be visually understandable.
- Distant terrain may include decorative cinder cones with narrow, gently
  pulsing, branching lava flows. Sparse variations in cone height, crater
  profile, channel direction, and branching establish volcanic character
  across stages without adding collision or pyroclastic-eruption gameplay.
- Predictable motion is important because the balloon reacts with delay.

### Terminal collision rules

- The terminal game does not need physical collision resolution.
- There should be no bounce, reversal, or “push out of overlap” behavior.
- Collisions are damage events only.
- **Any damaging collision pops the balloon.** There is no health bar or
  escalating damage model in the terminal core.

Likely implementation choices, but not separately confirmed design decisions:

- transparent sprite cells do not collide;
- clouds remain non-damaging scenery; and
- gameplay enters a pop/death state on first contact, preventing repeated
  effects from a persistent overlap.

Dynamic, animation-dependent collision masks are useful but risky. Proxima 5
had an alien that could become stuck at an edge because wing animation changed
its colliding cells. Damage-only collisions remove the resolution feedback loop
that caused that class of bug, but animation-aware masks still require tests.
Whether every hazard uses its animated mask or a stable forgiving mask remains
a tuning decision.

### Terminal sprites

The terminal version needs compact, readable, colored sprites for:

- the balloon, with clearly different burner-on and burner-off states;
- pterosaurs;
- clouds;
- trees and possibly sauropods or other ground hazards; and
- a balloon-pop sequence inspired by Proxima 5’s explosion graphics.

The balloon, pterosaurs, and ground hazards should animate. Clouds can remain
static while scrolling. Exact sprite dimensions and art are not settled. A
small sprite-gallery command is an appropriate way to tune them outside the
game.

Candidate terminal sizes suggested during the discussion—not commitments—were:

| Sprite | Candidate footprint | Candidate animation |
| --- | ---: | --- |
| Balloon and basket | 7–9 × 5–6 cells | Burner flame and slight sway |
| Pterosaur | 5–7 × 2–3 cells | Two or three wing positions |
| Cloud | 8–14 × 2–3 cells | Static |
| Tree | 3–5 × 3–5 cells | Optional sway |
| Sauropod | 9–13 × 4–5 cells | Head/neck movement |
| Balloon pop | 9–13 × 5–7 cells | Four to six short frames |

The mobile discussion also proposed separately animated envelope, basket,
pilot, burner flame, and (at that stage) weapon components. That decomposition
is a possibility for mobile art, not a terminal requirement.

### Technology tracks

#### Terminal

- Go and Tcell.
- The preferred playfield is 120×40 cells. At startup the game asks Tcell to
  resize the terminal window to that size; terminals that cannot or will not
  resize retain their existing dimensions.
- Proxima 5 is a useful reference for sprites, animation, collision detection,
  layering, and level ideas.
- Proxima 5 was Garrett’s first substantial game (apart from a number-guessing
  game written as a child around 1980). It was created mostly in a weekend after
  realizing that Tcell, originally made for a business application, could be
  used for a terminal game.
- Proxima 5 may contain poor or expedient design choices. Critique freely; copy
  behavior only after deciding it fits Jurassic Balloon.
- The terminal game should be a real game in its own right, not merely a
  disposable graphical prototype or a reduced mobile renderer.

#### Mobile

- A commercial mobile version for both iOS and Android is of genuine interest.
- Go is strongly attractive because Garrett knows it deeply.
- Ebitengine is the leading candidate discussed for mobile, but it has not been
  adopted by an implementation decision in this repository.
- The mobile version may have smoother physics, richer art and audio, touch
  input, and mechanics that do not appear in the terminal game.
- Do not require shared code between terminal and mobile. Reuse concepts and
  vocabulary first; extract shared code later only when a demonstrated benefit
  exceeds the coupling cost.

## Strong possibilities, not yet commitments

### Burner fuel

Fuel is the strongest post-core mechanic discussed, but Garrett introduced it
as “just an idea.” Treat it as a planned experiment, not a first-milestone
requirement.

Possible behavior:

- Burner use consumes a finite fuel supply.
- Running out of fuel disables the burner and lets the balloon descend; it need
  not be a separate instant-death condition.
- Certain visually distinctive treetops replenish fuel when the balloon grazes
  a collection region near the canopy top.
- Fuel trees should be distinguishable by shape or symbols as well as color so
  they remain readable across terminal capabilities.
- Pickup placement near hazards can create controlled difficulty and
  risk/reward decisions.
- Required pickups must be reachable given inertia and delayed burner response,
  not merely geometrically unobstructed.
- Early levels might have unlimited or generous fuel, with scarcity introduced
  later.

This mechanic is appealing because it creates depth without adding an input.
It should be added only after basic flight and avoidance are already fun.

### Additional level ideas

- Storm stages with rain signaling increased forward speed.
- Volcano stages with thermal lift and visible ballistic ejecta.
- Checkpoints and save points.
- Predictable pterosaur flight patterns.
- Sauropods raising or lowering their heads.
- Other carnivores threatening low flight, possibly including leaps.
- Handcrafted sequences with optional procedural variation.
- Landing and short on-foot phases.
- Power-ups other than weapons.
- Scores, achievements, or endless survival modes.

All except the basic recognizable hazards are possibilities. On-foot play in
particular is a scope expansion and should remain deferred.

## Superseded or explicitly deprioritized ideas

- **Variable ordinary forward speed:** superseded by constant normal speed;
  exceptional environments such as storms may change it.
- **Unsignaled updrafts/turbulence:** rejected in spirit. Environmental forces
  must be telegraphed. A visible volcano is a valid source of lift.
- **Turn-based terminal play:** explicitly rejected; the game is real time.
- **Terminals cannot report key release:** factually superseded. Modern Kitty
  and Win32 input paths can report it; legacy terminals need a fallback.
- **Physical collision resolution:** rejected. Jurassic Balloon needs detection
  and consequences, not bounce/reversal/separation.
- **Partial damage or worsening leaks after a hit:** superseded for the terminal
  core by “any collision pops the balloon.”
- **Lance/spear in the terminal version:** removed to keep that version simpler.
- **Weapons as a planned mobile mechanic:** initially discussed (lance aimed at
  a pterosaur mouth, limited gun or slingshot shots), later downgraded to an
  optional experiment. The current bias is no weapons unless playtesting shows
  a clear need.
- **A single shared terminal/mobile engine:** not a goal. The versions may use
  different mechanics and architectures.
- **Godot as the assumed mobile engine:** displaced by interest in a Go-based
  Ebitengine track, though no mobile engine is formally selected.
- **Ads or recurring monetization:** inconsistent with the stated product goal.

## Monetization and distribution

The commercial discussion applies to the possible mobile version, not the
terminal repository.

### Direction with clear support

- No ad-supported design.
- A free demo followed by a one-time full-game purchase is the preferred model
  being explored.

The assistant recommended one application containing both demo and full game,
using a non-consumable unlock so progress carries forward and purchases can be
restored. It also recommended avoiding subscriptions and manipulative virtual
currency. Those are sensible proposals consistent with Garrett’s stated
preferences, but only the rejection of ads and interest in a one-time unlock
were explicit user decisions.

### Not decided

- Final price. The conversation mentioned roughly USD $1.99–$4.99 and leaned
  toward $2.99–$3.99, but those were assistant suggestions, not a decision.
- Exact demo boundary.
- Whether mobile is definitely shipping or merely worth exploring.
- Store implementation, entitlement restoration, tax, and compliance details.
- Whether the terminal version is free/open-source, although that is the
  natural current direction.

## Practical first implementation brief — terminal track

The current repository already has a fixed-step Tcell loop and a small
terminal-independent sprite/scene package. The next implementation should prove
the core game, not broaden the feature list.

### Goal

Produce a repeatable, playable one-control survival slice in which a balloon
flies forward in real time, responds to delayed burner input, collides with one
or two clearly readable hazards, pops, and restarts quickly.

### Milestone 1: flight laboratory

1. Replace direct arrow-key positioning with abstract `BurnerActive` input.
2. Implement fixed-timestep vertical simulation with tunable parameters:
   gravity, burner response delay, lift/heat rise and decay, maximum ascent and
   descent speed, and leak/cooling behavior.
3. Keep normal horizontal world speed constant; represent forward travel by
   scrolling the world past a mostly fixed horizontal balloon position.
4. Add a developer HUD or flags exposing altitude, vertical velocity, burner
   state, and tuning constants.
5. Add deterministic unit tests for burner response, caps, pause, and restart.

Acceptance test: with no hazards, maintaining altitude for one to two minutes
is understandable and enjoyable rather than exhausting or random.

### Milestone 2: sprite gallery and input capability

1. Add a `cmd/sprite-gallery` program showing candidate sizes and animations
   for balloon on/off, pterosaur wing states, clouds, tree, sauropod, and pop.
2. Use a small palette with strong silhouettes and fallbacks for limited-color
   or Unicode-poor terminals.
3. Add a terminal input adapter that exposes key-down/key-up when available,
   clears held state on focus loss, and offers a documented legacy fallback.
4. Keep the simulation unaware of keyboard protocol details.

Acceptance test: burner state is immediately legible, sprites do not merge into
one another, and the balloon leaves enough forward visibility at 80×24.

### Milestone 3: first survival slice

1. Add a scrolling ground heightmap and dangerous upper boundary.
2. Add one predictable pterosaur pattern.
3. Detect damage-mask entry, including a swept or conservative check so high
   speed cannot tunnel through a one-cell hazard.
4. On first collision, freeze the live simulation, run a short pop animation,
   then allow immediate restart.
5. Add score/distance only if it helps evaluate the loop.

Acceptance test: every failure is visually understandable; no collision causes
bounce, reversal, trapping, or repeated death events.

### Deferred until the slice is fun

- Fuel consumption and fuel-bearing treetops.
- Storm and volcano levels.
- Checkpoint persistence.
- Additional enemy families.
- On-foot phases.
- Weapons of any kind.
- Mobile/Ebitengine code.
- Store purchases or entitlement handling.
- Final art, audio, procedural generation, and monetization UI.

## Open questions for playtesting

- Is the best model delayed impulse, a heat/lift state, or a combination?
- Which values should be capped: lift, acceleration, vertical velocity, or more
  than one?
- Does deliberate fast descent need an input, or are gravity and burner release
  sufficient?
- How much forward visibility is needed at 80×24 for inertial avoidance?
- Should the balloon use a stable forgiving collision mask while hazards use
  animated occupied-cell masks?
- Are hazards fun before fuel is introduced?
- Do fuel-bearing trees improve risk/reward or make burner management too busy?
- Does mobile need any mechanic beyond burner timing and environmental
  avoidance?
- Should terminal levels be authored, generated, or assembled from authored
  segments?

## Scope guardrail

The first success criterion is not “feature complete.” It is: **holding and
releasing one control makes a leaky balloon enjoyable to fly through a small,
fair, readable prehistoric obstacle course.** Everything else is earned by that
prototype.
