# Shorthand [![Build Status](https://travis-ci.org/koki/shorthand.svg?branch=master)](https://travis-ci.org/koki/shorthand)

Encode your domain-specific knowledge and automatically generate a user-friendly API.

Applied to Kubernetes, it allows us to provide an ergonomic spec format (and tools!) that's simpler to learn and easier to get right.

## Development - Working on Shorthand

Dependencies are managed using `glide`, which puts them in the `vendor` folder.
If something (e.g. `cobra`) breaks, try `glide cache-clear` and then `glide up`.

If you have Go set up, `go install` builds faster than the Docker-based scripts.
Then you can run the `shorthand` tool from the command line.
