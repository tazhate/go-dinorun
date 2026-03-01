# go-dinorun-tuned

A souped-up fork of [go-dinorun](https://github.com/ahmad-alkadri/go-dinorun) — the Chrome Dinosaur Game in your terminal, now with god mode, speed control, a full-width gradient progress bar, and instant restart.

## What's new

| Feature | Flag | Default |
|---|---|---|
| Speed control | `--speed 1..10` | `5` |
| Score goal / progress bar | `--goal N` | `1000` |
| God mode (no collisions) | `--immortal` | off |
| No enemies mode | `--no-enemies` | off |
| **Restart on Space** after death | — | always on |

### Progress bar

A truecolor gradient bar spans the full terminal width at the top of the screen. It fills as your score climbs toward `--goal`:

- Colors sweep **blue → cyan → green → yellow → red** based on position
- Sub-character precision using Unicode block elements (`▏▎▍▌▋▊▉█`)
- Glow effect at the leading edge
- Set `--goal 0` to hide the bar entirely

## Installation

```sh
go install github.com/tazhate/go-dinorun@latest
```

Or clone and build locally:

```sh
git clone https://github.com/tazhate/go-dinorun.git
cd go-dinorun
go install
```

Requires a truecolor-capable terminal (most modern terminals: kitty, alacritty, iTerm2, GNOME Terminal, Windows Terminal).

## Usage

```sh
go-dinorun [flags]
```

```
  --speed N      Game speed: 1 (slow) to 10 (insane). Default: 5
  --goal N       Score goal for the progress bar. Default: 1000. Set 0 to disable.
  --immortal     God mode — collisions are ignored, score keeps climbing
  --no-enemies   Run without cactuses or pteranodons
```

### Examples

```sh
# Normal game
go-dinorun

# Fast game with a longer goal
go-dinorun --speed 7 --goal 5000

# Zen mode — just vibe
go-dinorun --immortal --no-enemies --speed 3

# Practice collisions at slow speed
go-dinorun --speed 2
```

## Controls

| Key | Action |
|---|---|
| `Space` | Jump |
| `Space` (after death) | Restart immediately |
| Any other key (after death) | Exit and show final score |
| `Esc` / `Ctrl+C` | Quit |

## Uninstall

```sh
rm -f $(which go-dinorun)
```

---

Based on [ahmad-alkadri/go-dinorun](https://github.com/ahmad-alkadri/go-dinorun). Original game logic and sprites by Ahmad Al-Kadri.
