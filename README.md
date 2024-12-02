# Godot Package Manager (WIP)

Easy package management for Godot.

With this you can eliminate the bloat of the addons folder on your repository and provide a easy way of managing it.

## Build

To build this project you will need `go 1.23`

## The godot-package.json file

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

## Commands

### Init godot-package.json
> **NOTE:** Init command will only add the packages that have a **```godot-package.json```** inside it with the type property as **```addon```** or **```scripts```**. This could be changed in the future.

This will try to list your current dependencies of the ```addons``` folder (if there's any) and add it to the ```godot-package.json```. If there's is no dependencies it will create the ```godot-package.json``` anyway.
Run the following command on the root of your project:
```shell
gpm init
```

Then the CLI will prompt for some information about the project, it will ask for the **name**, **type** and **version** to populate the file.

If the ```godot-package.json``` file exists, it will try to update the dependencies.

---
### Add dependencies
To add a dependency in the ```godot-package.json``` run the following command on the root of your directory:

```shell
gpm add name repository version
```
---
### Installing dependencies
To install the dependencies listed on the ```godot-package.json``` run the following command on the root of your directory:

```shell
gpm install
```

It will install your dependencies on the ```addons``` folder and enable it in the settings.

---
### Removing dependencies
Removing dependencies is as easy as installing, run the following command on the root directory of your godot project:

```shell
gpm remove name
```

This will remove the dependency from the ```addons``` folder and remove it in the settings.


## How does it work

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