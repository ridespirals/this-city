// Package fonts embeds UI typefaces for the renderer.
package fonts

import _ "embed"

// Space Mono by Colophon Foundry, licensed under the SIL Open Font License 1.1.
// See OFL.txt in this directory.

//go:embed SpaceMono-Regular.ttf
var SpaceMonoRegular []byte

//go:embed SpaceMono-Bold.ttf
var SpaceMonoBold []byte
