module gamatet

go 1.21

replace github.com/marko-gacesa/udpstar => ../udpstar

require (
	github.com/go-gl/gl v0.0.0-20190320180904-bf2b1f2f34d7
	github.com/go-gl/glfw/v3.3/glfw v0.0.0-20201108214237-06ea97f0c265
	github.com/go-gl/mathgl v1.0.0
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0
	github.com/google/go-cmp v0.5.7
	github.com/marko-gacesa/udpstar v0.0.0-20231130162919-86d574d6bdcb
	golang.org/x/image v0.0.0-20190321063152-3fc05d484e9f
)

require github.com/marko-gacesa/appctx v0.0.0-20220908091727-0bf8492596d2 // indirect
