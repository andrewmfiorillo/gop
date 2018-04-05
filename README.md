# `p`

## Python Version Management, Simplified.

![introduction](https://cloud.githubusercontent.com/assets/1139621/7488032/37f37308-f389-11e4-8995-89f7cba5ad8b.gif)

`p` is powerful and feature-packed, yet simple, both in setup and use. There are no tricky settings, options, or crazy dependencies. `p` is just a helpful ~600 line Bash script that gets the job done.

**`p` let's you quickly switch between Python versions whenever you need to, removing the barrier between Python 2.x.x and 3.x.x.**

`p` was heavily inspired by [`n`, a version manager for Node.js](https://github.com/tj/n).

`p` is also great for getting started using Python development versions. Use `p latest` to get up and running with the latest development version of Python!

## Installation

After downloading the Bash script, simply copy it over to your `$PATH` and `p` will take care of the rest.

```
curl -sSLo p https://raw.githubusercontent.com/Raphx/p/master/bin/p
chmod +x p
mv p /usr/local/bin
```

So far, `p` has only been tested in Bash.

## Usage

```
Usage: p [COMMAND] [args]

Commands:

p                              Output versions installed
p status                       Output current status
p <version>                    Activate to Python <version>
  p latest                     Activate to the latest Python release
  p stable                     Activate to the latest stable Python release
p use <version> [args ...]     Execute Python <version> with [args ...]
p bin <version>                Output bin path for <version>
p rm <version ...>             Remove the given version(s)
p ls                           Output the versions of Python available
  p ls latest                  Output the latest Python version available
  p ls stable                  Output the latest stable Python version available

Options:

-V, --version   Output current version of p
-h, --help      Display help information
```

### `p`

Executing `p` without any arguments displays a list of installed Python versions, and the current activated version.

```
$ p

    2.7.14
  ο 3.6.5
    3.7.0
```

### `p status`

Shows the version, bin path, and status of current activated Python version.

```
$ p status
     version : 3.6.5
         bin : /home/raphx/.python/p/versions/python/3.6.5/python
      latest : no
      stable : yes
```

### `p [<version>|latest|stable]`

Activate, or install the specified Python version if not already installed. `latest` and `stable` can be used to quickly install the latest or stable version respectively.

```
$ p 3.3.4

     install : Python-3.3.4
      create : /home/raphx/.python/p/versions/python/3.3.4
       fetch : https://www.python.org/ftp/python/3.3.4/Python-3.3.4.tgz
   configure : 3.3.4
     compile : 3.3.4

  Success: Installed Python 3.3.4!
```

### `p use`

Quickly use the specified version to execute a one-off command, even when the version is not activated.

```
$ p use 2.7.14 -c "import sys; print sys.version"
2.7.14 (default, Apr  5 2018, 22:47:52)
[GCC 7.3.1 20180312]
```

### `p bin`

Output the bin path for the current activated Python version.

```
$ p bin
/home/raphx/.python/p/versions/python/3.6.5/python
```

### `p rm`

Remove an installed Python version.

```
$ p rm 2.7.14
      remove : 2.7.14

  Success: Removed Python 2.7.14!
```

### `p ls [latest|stable]`

List available Python versions. If `latest` or `stable` is supplied, show only the corresponding version.

```
$ p ls

# --snip--

    2.7.12
    2.7.13
    2.7.14
    3.0.1
    3.1.1
    3.1.2
    3.1.3
    3.1.4
    3.1.5
    3.2.1
    3.2.2
    3.2.3
    3.2.4
    3.2.5
    3.2.6
    3.3.0

# --snip--
```

## FAQs

**How does `p` work?**

`p` stores each Python version installed in `$P_PREFIX/p/versions/python`. When a Python version is activated, `p` creates a symbolic link to the Python binary located at `$P_PREFIX/p/versions/python/python`. `$P_PREFIX` allows you to customize where python versions are installed, and defaults to `/usr/local` if unspecified.

**How do I revert back to my default Python version?**

Simply run `p default` and `p` will remove the symbolic link described above; therefore reverting back to your default Python version.

**Does `p` download the source each time I activate or install a version?**

Nope. `p` stores the source for each of the versions installed, allowing for quick activations between already-installed versions.

**How do I get this working on Windows?**

No Windows support planned at the moment, PRs are always welcomed!

## TODO

* also manage pip
* also manage PyPy

## Attribution

This is a fork from the original [p](https://github.com/qw3rtman/p) by [Nimit Klara](https://github.com/qw3rtman).

## License

[MIT](LICENSE)
