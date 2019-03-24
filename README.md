# go3dprint

A totally overengineered thingy to make 3D models and visualize them in Go.

To use:
- Install [Garden](https://docs.garden.io/basics/installation).
- Clone this repo.
- Go into it and `garden dev`.
- Open the `web` ingress on your browser.

Structure:
- The `mesh` module generates the 2D and 3D stuff.
- The `render` module makes the 3D stuff visible.
- The `web` module fetches things from the above and shows you.

Commit history might contain nightmares—you've been warned.