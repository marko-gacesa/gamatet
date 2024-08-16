module gamatet

go 1.21

replace github.com/marko-gacesa/udpstar => ../udpstar

require (
	github.com/go-gl/gl v0.0.0-20231021071112-07e5d0ea2e71
	github.com/go-gl/glfw/v3.3/glfw v0.0.0-20231223183121-56fa3ac82ce7
	github.com/go-gl/mathgl v1.1.0
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0
	github.com/google/go-cmp v0.5.7
	github.com/marko-gacesa/appctx v0.0.0-20220908091727-0bf8492596d2
	github.com/marko-gacesa/udpstar v0.0.0-20231130162919-86d574d6bdcb
	golang.org/x/image v0.14.0
	golang.org/x/text v0.14.0
)

require golang.org/x/xerrors v0.0.0-20191204190536-9bdfabe68543 // indirect
