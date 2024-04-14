# About

Automatically choosing web-browser depends on environment context rules.

## Features

- suckless solution with no redundant dependencies
- fast single binary
- simple rule engine with custom interpretable DSL
- extendible (TODO)
- crossplatform (TODO)

# Configuration

## Example

```
firefox -p job {}:url.regex='.*jira.*';app.class=Slack # Open all jira links from slack with job firefox profile
firefox 'ext+container:name=Isolated&url={}':app.class=org.telegram.desktop # Open all links from the telegram app using Isolated firefox container

# Default fallback
firefox {}:fallback
```

## Configuration syntax

The application just evaluates configuration rules 1 by 1 and applies url to a first matched command. Syntax can be described as: 

```
<browser_command>:<matcher_knd>.<prop_name>=<prop_value>[;<matcher_knd>.<prop_name>=<prop_value>]
```

Browser command is a sequence of words, divided by spaces. The first word is an executable name and the others are arguments. `{}` char sequence will be replaced with a clicked URL.

You can escape spaces or other _non-word characters_ can be escaped by a single-quote string.

## Matchers

### fallback

This matcher always succeeds. Use it at the end of a configuration to specify the default browser. 

### app

Matches by source application.

Currently supported desktop environments: _hyprland_, _gnome_ (TODO), _macos_ (TODO).

Properties:

- *title*: match by source window title with regex
- *class*: match by window class

### url

Match by a clicked URL.

Properties:

- *host*: match URL by host
- *scheme*: match URL by scheme
- *regex*: match full URL by regex

# Setup

## Linux

### Linux manual

```sh
go install github.com/pltanton/autobrowser/cmd/autobrowser-linux@latest
```

Create config at `~/.config/autobrowser.config`.
Then add the following .desktop file to `~/.local/share/applications/` and set it as the default browser. 
Change paths for your setup if needed.

```ini
[Desktop Entry]
Categories=Network;WebBrowser
Exec=~/go/bin/autobrowser-linux -config ~/.config/autobrowser.config -url %u
Icon=browser
MimeType=x-scheme-handler/http;x-scheme-handler/https
Name=Autobrowser: select browser by contextual rules
Terminal=false
Type=Application
```

### Nixos flakes with Home manager

In your `flake.nix`:

```nix
{
  autobrowser.url = "github:pltanton/autobrowser";
  autobrowser.inputs.nixpkgs.follows = "nixpkgs";

  outputs = { self, nixpkgs, home-manager, ddcsync }: {
    modules = [
      ({ pkgs, ... }: {
        nixpkgs.overlays = [ autobrowser.overlays.default ]; # To use programm as package
      })
    ];

    # To use with home-manager
    homeConfigurations."USER@HOSTNAME"= home-manager.lib.homeManagerConfiguration {
      pkgs = nixpkgs.legacyPackages.x86_64-linux;
      modules = [
        autobrowser.homeManagerModules.default
        { 
            programs.autobrowser = {
                enable = true; 
                rules = [
                    "firefox 'ext+container:name=Work&url={}':app.class=Slack"
                ];
                default = "firefox {}";
            };
        }
        # ...
      ];
    };
  };
```
