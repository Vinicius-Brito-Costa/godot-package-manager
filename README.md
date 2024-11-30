# Godot Package Manager (WIP)

Easy package management for Godot.

With this you can eliminate the bloat of the addons folder on your repository and provide a easy way of managing it.

## The godot-package.json file

This file is used to store and manage your dependencies. We need to keep it inside the root directory of the project.

This file has the following structure:
```json

```

## Commands

### Init godot-package.json

This will try to list your current dependencies of the ```addons``` folder (if there's any) and add it to the ```godot-package.json```.
Run the following command on the root of your project:
```shell
gpm init
```

### Add dependencies
To add a dependency in the ```godot-package.json``` run the following command on the root of your directory:

```shell
gpm add repository name version
```

### Installing dependencies
To install the dependencies listed on the ```godot-package.json``` run the following command on the root of your directory:

```shell
gpm install
```

It will install your dependencies on the ```addons``` folder and enable it in the settings.

### Removing dependencies
Removing dependencies is as easy as installing, run the following command on the root directory of your godot project:

```shell
gpm remove name
```

This will remove the dependency from the ```addons``` folder and remove it in the settings.

## How does it work

The ```gpm``` will try to access the repository and search for a release with the version.

Example with a github repository:

**repository:** https://github.com/ramokz/phantom-camera  
**version:** v0.8

```
GET https://github.com/ramokz/phantom-camera/archive/refs/tags/v0.8.zip
```

Next it will unzip it and try to locate an ```addons``` folder, if it cannot find it, it will search for ```plugin.cfg``` to locate the folder in which the addon is.

Then it will add it to the project ```addons``` folder and enable it in the Settings.