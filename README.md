# go3dprint

<p align="center">
  <img src="img/frontend.png" width="70%">
</p>

`go3dprint` is a needlessly distributed phallic object generator.

It's an educational demo project made to showcase the basics of creating and rendering 3D mesh with Go. It uses the [`sdfx`](https://github.com/deadsy/sdfx) and [`fauxgl`](https://github.com/fogleman/fauxgl) libraries.

This project was presented at [dotGo](https://www.dotgo.eu/) 2019, in Paris.

<p align="center">
  <a href="https://www.youtube.com/watch?v=ZACOc-NwV0c" target="_new"><img src="https://img.youtube.com/vi/ZACOc-NwV0c/0.jpg" width="40%"></a>
</p>

For the live demo a simpler version was used. This is the "full version," which has the following key differences:

- Functionality is split into loosely coupled microservices
- Running as lightweight containers
- That communicate via API calls and websockets.

And thanks to [Garden](https://garden.io/), the workflow is configured so that it:

- Re-builds and re-deploys on every code change
- Can use hot reload, so containers can be update without restarting
- Uses the same tooling for all environments—local, CI, remote.

The differences and the conversion process are explained in detail in the article [_Needlessly Distributed Phallic Object Generator_](https://medium.com/garden-io/the-needlessly-distributed-phallic-object-generator-2da47672be6f).

## How It Works

This project is split into three microservices: `mesh`, `render`, and `web`.

<p align="center">
  <img src="img/graph.png" width="35%">
</p>

The way it works is:

- `web` constantly polls `mesh` for 2D and 3D objects.
- If it receives a 2D object, it displays it on the browser.
- If it receives a 3D object, it POSTs it to the `render` service.
- The `render` service, in turn, returns an image of the rendered mesh, which is then displayed on the browser.

To enable live feedback as one explores different 2D/3D forms, this project uses Garden to re-build and re-deploy services whenever the source code changes.

Simply install Garden, clone this repository, and `garden dev`.

## Instructions

1. [Install Garden](https://docs.garden.io/getting-started/1-installation).
2. Clone this repo.
3. Run `garden dev --hot=mesh`, and leave it running.
4. You should see the ingress endpoint for the `web` service in the output. Open that in your browser.
5. Now go and mess around! Open `mesh/main.go` and look for the commented out lines in the `magic()` function. Try uncommenting the sections one by one and observe the results in the browser (it takes a few seconds to rebuild each time).

<p align="center">
  <img src="img/dashboard.gif">
</p>
