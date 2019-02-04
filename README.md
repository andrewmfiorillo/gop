# `gop` - simple governing of your Python versions

`gop` is simple to setup, simple to use, and powerful to boot.

There are no tricky settings, options, or crazy dependencies. `gop` is just a helpful tool that gets the job done, and is the Golang cousin of [`p`, a bash script for managing Python versions](https://github.com/Raphx/p).

`gop` is also great for getting started using Python development versions. Use `gop latest` to get up and running with the latest development version of Python!

## Installation

To build from source, install [Go](https://golang.org/) and then run

```shell
go get -u github.com/andrewmfiorillo/gop/cmd/gop
```


The two settings you may want to change are `P_PREFIX` (where are installed versions are kept) and your `PATH` (which makes them easily accessible)

In Linux or OSX, this would look like

```shell
# In ~/.bash_profile, or the equivalent
export P_PREFIX=$HOME
export PATH="$P_PREFIX/p/versions/bin:$PATH"
```

Or in Windows, you set the environment variables with [Rapid Environment Editor](https://www.rapidee.com/en/about) or your favorite tool

```batch
set P_PREFIX=%USERPROFILE%
set PATH=%P_PREFIX%;%PATH%
```

## Usage

```
Commands:
    gop <version>                  Activate to Python <version>
    gop ls, list                   Output the versions of Python available
        gop ls installed           Output the installed versions of Python
        gop ls latest              Output the latest Python version available
        gop ls stable              Output the latest stable Python version available
    gop latest                     Activate to the latest Python release
    gop stable                     Activate to the latest stable Python release
    gop status                     Output current status
    gop install <version> --force  Install Python <version> but do NOT activate
    gop use <version> [args ...]   Execute Python <version> with [args ...]
    gop bin <version>              Output bin path for <version>
    gop rm <version ...>           Remove the given version(s)
    gop default, disable           Use default (system) Python installation
    gop help, h [command]          Shows a list of commands or help for one command

Options:
  --verbose
  --help, -h     show help
  --version, -v  print the version

```

<!-- ### `gop`

Executing `gop` without any arguments displays a list of installed Python versions, and the current activated version.

```
$ p

    2.7.14
    3.3.4
  ο 3.6.5
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

### `p [<version>|latest|stable]`

Activate, or install the specified Python version if not already installed. `latest` and `stable` can be used to quickly install the latest or latest stable version respectively.

```
$ p 3.3.4

     install : Python-3.3.4
       fetch : https://www.python.org/ftp/python/3.3.4/Python-3.3.4.tgz
   configure : 3.3.4
     compile : 3.3.4

  Success: Installed Python 3.3.4!
```

### `p status`

Show the version, bin path, and status of current activated Python version.

```
$ p status
     version : 3.6.5
         bin : /home/raphx/.python/p/versions/python/3.6.5/bin/python
      latest : no
      stable : yes
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
/home/raphx/.python/p/versions/python/3.6.5/bin/python
```

### `p rm`

Remove an installed Python version. Multiple versions can also be provided using space as delimiter. The `rm` subcommand also accepts `stable` and `latest` identifier.

```
$ p rm 2.7.14
      remove : 2.7.14

  Success: Removed Python 2.7.14!
```

### `p default`

Remove the Python symlink created by `p`, thereby reverting to use the default or system installed Python, if there are any. For instance, I have a system Python with version 3.6.5:

```
$ p default
    activate : default

  Success: Now using default system Python!
``` -->

## How does `gop` work?

`gop` stores each Python version installed under the directory `$P_PREFIX/p/versions/python`. When a Python version is activated, `p` creates symbolic links in `$P_PREFIX/p/versions`, pointing to the:

 - `bin`
 - `include`
 - `lib`
 - `share`

directories of the activated Python version.

For example, Python version 3.6.5 is installed, and it will be placed under the directory:

```
$P_PREFIX/p/versions/python/3.6.5
```

Activating version 3.6.5 will create symlinks that points to directories under the activated Python installation:

```
$P_PREFIX/p/versions/bin     -> $P_PREFIX/p/versions/python/3.6.5/bin
$P_PREFIX/p/versions/include -> $P_PREFIX/p/versions/python/3.6.5/include
$P_PREFIX/p/versions/lib     -> $P_PREFIX/p/versions/python/3.6.5/lib
$P_PREFIX/p/versions/share   -> $P_PREFIX/p/versions/python/3.6.5/share
```

`$P_PREFIX` allows you to customize where python versions are installed, and defaults to `$HOME` (`%USERPROFILE%` on Windows) if unspecified. To use the Python that `gop` installs, you must either call its full path (given with `gop bin`) or add `$P_PREFIX/p/versions/bin` to your `$PATH`.

When installing Python 3, the symlink `python` and `pip` are also created for `python3` and `pip3` executables respectively, for the sake of convenience.

## FAQs

**What about `pip`?**

`pip` is installed by default for Python 3. For Python 2.7.9 or above, you can run the following command to install `pip`:

```
python -m ensurepip
```

For Python version less than 2.7.9, you will have to install `pip` manually.

**Can I use `gop` to manage Python project dependencies?**

You could, though it is recommended to use [`pipenv`](https://docs.pipenv.org/) to do so. `gop` can be used to install the Python version for your project, and `pipenv` can later be used to setup the project virtual environment using the Python version installed.

## Attribution

This is a fork from [p](https://github.com/Raphx/p) by [Raphx](https://github.com/Raphx).

## License

[MIT](LICENSE)
