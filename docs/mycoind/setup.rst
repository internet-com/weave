--------------------
Prepare Requirements
--------------------

Before you can run this code, you need to have a number
of programs set up on your machine. In particular, you
will need a bash shell (or similar), and development tooling
for both go and node.

**WARNING**

This is only tested under Linux and OSX.
It will most likely not work under Windows, even with Cygwin.
If you use windows, please make a PR with any adjustments so
this tutorial works for Windows as well.

Install Go
==========

You will need to have the go tooling installed, version 1.9+.
If you do not already have it, please
`download <https://golang.org/dl/>`_ and
`follow the instructions <https://golang.org/doc/install>`__
from the official golang homepage. Make sure to read down
to `Test Your Installation <https://golang.org/doc/install#testing>`__.

We assume a standard setup in the Makefiles, especially to
build tendermint nicely. That means you must set up `GOPATH`,
you must check out all source code under `$GOPATH/src`,
and you must add the default install directory to your `PATH`.

.. code-block:: console

    # these two lines should be in .bashrc or similar
    export GOPATH=$HOME/go
    export PATH=$PATH:$GOPATH/bin
    # this will properly place the code in $HOME/go/src/github.com/confio/weave
    go get github.com/confio/weave


Go related tools
----------------

You must also make sure to have a few other developer tools
installed. If you are a developer in any language, they are
probably there. Just double check.

* git
* make
* curl


Install Node
============

You will need node 8+ to run the example client. Unless you know what you
are doing, stick to even numbered versions (6, 8, 10, ...), the odd numbers
are unstable and get deprecated every few weeks it seems. For ease
of updating later, I advise you to install `nvm <https://github.com/creationix/nvm#installation>`__ and then add the most recent stable version

.. code-block:: console

    nvm ls-remote | grep -v v9 | tail
    nvm install v8.11.1 # or other recent version

    # test it out
    node --version
    node
    > let {x, y} = {x: 10, y:10}


Node related tools
------------------

Yarn is a faster alternative to npm for installing modules, so
we use that as default.

.. code-block:: console

    npm install -g yarn
