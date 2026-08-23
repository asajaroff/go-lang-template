// Command server is the single compiled binary for this service: `serve`
// runs the HTTP server (Factor VII: it binds its own port, nothing external
// injects it into a host process), `migrate` runs schema migrations as a
// one-off admin process (Factor XII) using the exact same binary/image.
package main

import "os"

func main() {
	if err := newRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}
