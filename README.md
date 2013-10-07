# What is this?

A little Go program to use [cote-de-porc][cdp]'s JSON API.

# Usage:

    $ ./cdporc -h
    Usage of ./cdporc:
      -d="ID": Delete quote (shorthand)
      -del="ID": Delete quote
      -l=false: List quotes (shorthand)
      -list=false: List quotes
      -p="ID": Publish quote (shorthand)
      -publish="ID": Publish quote
      -r=false: Get a random quote (shorthand)
      -rand=false: Get a random quote
      -t="all": List quotes of this type (shortand)
      -type="all": List quotes of this type

The program expect the `CDP_SERVER` environment variable to point to your
installation of the [cote-de-porc][cdp] software, example:

    $ CDP_SERVER="http://user:password@cote-de-porc-server.com" ./cdporc -l

# Building:

You need Go 1.1.

    $ git clone http://github.com/oz/cdporc
    $ cd cdporc
    $ go build -o cdporc

# Will this program change the world?

Of course!

[cdp]: https://github.com/spk/cotedeporc
