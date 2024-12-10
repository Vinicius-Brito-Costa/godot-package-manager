# Godot Package Manager<!-- omit from toc -->

Easy package management for Godot.

With this you can eliminate the bloat of the addons folder on your repository and provide a easy way of managing it.
> Do not run `gpm` with Godot opened, close Godot run your commands then start it again.

## Index <!-- omit from toc -->

- [Build and install](#build-and-install)
- [The godot-package.json file](#the-godot-packagejson-file)
- [Commands](#commands)
  - [Init godot-package.json](#init-godot-packagejson)
  - [Add dependencies](#add-dependencies)
  - [Installing dependencies](#installing-dependencies)
  - [Removing dependencies](#removing-dependencies)
- [How does it work](#how-does-it-work)
- [FAQ](#faq)
  - [Where's the logging folder located?](#wheres-the-logging-folder-located)

## <a name="build-and-install"></a>Build and install

> To build this project you will need `go 1.23` or higher.

Just run the following command:

```shell
go build # for current os
env GOOS=windows GOARCH=amd64 go build  # for windows
env GOOS=windows GOARCH=386 go build  # for windows
env GOOS=linux GOARCH=arm go build  # for linux
env GOOS=darwin GOARCH=arm64 go build  # for mac
```

Now it should have a `gpm` executable in the folder, you can add it to you path and use it anywhere!

**OR**

You could run the following command to install it into you `go/bin` (that should already be on the path):
```shell
go install
```

**OR**

Just download the executable for your current platform on the [releases](https://github.com/Vinicius-Brito-Costa/godot-package-manager-backup/releases) tab and add it to your path.

## <a name="the-godot-packagejson-file"></a>The godot-package.json file

This file is used to store and manage your dependencies. We need to keep it inside the root directory of the project.

This file has the following structure:
```json
{
  "project": {
    "name": "Project",
    "type": "game",
    "version": "1.0.0",
    "repository": "github",
    "description": "Soft Description",
    "godotVersion": "4.4"
  },
  "plugins": [
    {
      "repository": "github",
      "name": "ramokz/phantom-camera",
      "version": "v0.8"
    },
    {
        "repository": "github",
        "name": "Burloe/GoLogger",
        "version": "1.2"
    }
  ]
}
```

## <a name="commands"></a>Commands

### <a name="init-godot-packagejson"></a>Init godot-package.json
> **NOTE:** Init command will only add the packages that have a **```godot-package.json```** inside it with the type property as **```addon```** or **```scripts```**. This could be changed in the future.

This will try to list your current dependencies of the ```addons``` folder (if there's any) and add it to the ```godot-package.json```. If there's is no dependencies it will create the ```godot-package.json``` anyway.
Run the following command on the root of your project:
```shell
gpm init
```

Then the CLI will prompt for some information about the project, it will ask for the **name**, **type** and **version** to populate the file.

If the ```godot-package.json``` file exists, it will try to update the dependencies.

---
### <a name="add-dependencies"></a>Add dependencies
To add a dependency in the ```godot-package.json``` run the following command on the root of your directory:

```shell
gpm add name repository version
```
---
### <a name="installing-dependencies"></a>Installing dependencies
To install the dependencies listed on the ```godot-package.json``` run the following command on the root of your directory:

```shell
gpm install
```

It will install your dependencies on the ```addons``` folder and enable it in the settings.

---
### <a name="removing-dependencies"></a>Removing dependencies
Removing dependencies is as easy as installing, run the following command on the root directory of your godot project:

```shell
gpm remove name
```

This will remove the dependency from the ```addons``` folder and remove it in the settings.


## <a name="how-dows-it-work"></a>How does it work

The ```gpm``` will try to access the repository with the given name and search for a release with the version.

Example with a github repository:

**repository:** github  
**name:** ramokz/phantom-camera  
**version:** v0.8

```
GET https://github.com/ramokz/phantom-camera/archive/refs/tags/v0.8.zip
```

Next it will unzip it and try to locate an ```addons``` folder, if it cannot find it, it will search for ```plugin.cfg``` to locate the folder in which the addon is.

Then it will add it to the project ```addons``` folder and enable it in the Settings.

## <a name="faq"></a>FAQ

### <a name="wheres-the-logging-folder-located"></a>Where's the logging folder located?
You can find the log files in `username/.godot-package-manager/`