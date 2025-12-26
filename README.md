
# GaMaTeT

[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)

GaMaTeT is a minimalistic, and yet feature rich, reimagining of Tetris.

The game features:
* Single player games.
* Local or LAN multiplayer games.
* Battle mode: Players playing against each other.
* Co-operative mode: Players playing on the same field.
* Variable field size: horizontally and vertically.
* Different piece types: rotating or mirroring.
* Different piece sizes: trominoes, tetrominoes or pentominoes.
* Various special pieces and blocks.

From technical point, the game features:
* Small size: The entire game is just a *single 7MB executable* file.
* Zero assets: All textures and 3D models are procedurally generated.
* Minimal memory footprint.
* Multiplayer communication using UDP protocol.
* 2D view or 3D view: Toggle with F12 key.

## About the Project

GaMaTeT is a product of passion. No AI tools were used for the code development in any way. The goal was not to make the best game or the best looking game, but simply to go through the entire process and enjoy every step of it. This is why the project imports almost no libraries - there's no game engine library, no networking library, no utility libraries. It was just OpenGL and me.

This game does not have any assets. All 3D models and textures are procedurally generated. It’s not a major achievement - the 3D models are all cube-like objects and the textures are just Perlin noise - but it does help set the tone for how I want this project to evolve.

The whole project is an attempt to escape strict corporate routine of my daily work and to remind myself how programming can be interesting, dynamic and enjoyable. How programs can be small, fast and efficient. The way it used to be in the good old days.

The project has also been my playground. A place where I could try out different programming styles and patterns. A place where I can try something different.

### Project History

* 2020: The development started in November 2020. The first things implemented are the *field* and the *piece* packages.
* 2021: I experimented with OpenGL. I wasn’t taking the project very seriously at the time, so progress was slow.
* 2022: Commits are chaotic, with changes scattered all over the place. I occasionally paused work for several months at a time. Later, I decided to squash all those early commits into a single one.
* 2023: I stopped the work on the game itself and focused solely on the network engine. This code was later moved into its own repository: [udpstar](https://github.com/marko-gacesa/udpstar).
* 2024: I continued with the network engine. Once it was finished, I went back to the game and completed text rendering and implemented the menu engine.
* 2025: Development intensified. The two big sides - the network engine and the game engine - finally came together: The game became playable over the network. I decided to open source the whole thing.

### Wish list

* There's no sound. I don't want music, because people can (and should) play their own favourite music. But I would like to have sound effects. And I would like them, in the spirit of the project, to be procedurally generated if possible.
* The only input device currently supported is keyboard. I don't think mouse input is necessary for a game like this, but I would like to support game controllers.
* Terminal rendering.
* Different game types.

## Support Development

If you enjoy this project, consider supporting it with a small donation:

[![GitHub](https://img.shields.io/badge/github_sponsors-EA4AAA?logo=githubsponsors&logoColor=white&style=for-the-badge)](https://github.com/marko-gacesa)
[![Buy Me A Coffee](https://img.shields.io/badge/buy_me_a_coffee-FFDD00?logo=buy-me-a-coffee&logoColor=black&style=for-the-badge)](https://www.buymeacoffee.com/marko.gacesa)
[![PayPal](https://img.shields.io/badge/paypal-002991?logo=paypal&logoColor=white&style=for-the-badge)](https://paypal.me/markogacesa77)

Thank you!

## License

This project is licensed under the terms of the [GNU General Public License v3.0](./LICENSE).

You are free to use, modify, and distribute this software under those terms.
See the LICENSE file for full details.
