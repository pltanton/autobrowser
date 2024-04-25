# About

Automatically choosing web-browser depends on environment context rules.

## Features

- suckless solution with no redundant dependencies
- fast single binary
- simple rule engine with custom interpretable DSL
- extendable (TODO)
- cross-platform

# Configuration

## Example

```
work:=firefox -p job {}:url.regex='.*jira.*'

work:app.class=Slack # Open all jira links from slack with job firefox profile
work:app.class=org.telegram.desktop # Open all links from the telegram app using Isolated firefox container

# Default fallback
firefox {}:fallback
```

Other examples can be found in `examples` folder

## Configuration syntax

The application just evaluates configuration rules 1 by 1 and applies url to a first matched command. Syntax can be described as: 

```
<browser_command>:<matcher_knd>.<prop_name>=<prop_value>[;<matcher_knd>.<prop_name>=<prop_value>]
```

Browser command is a sequence of words, divided by spaces. The first word is an executable name and the others are arguments. `{}` char sequence will be replaced with a clicked URL.

You can escape spaces or other _non-word characters_ can be escaped by a single-quote string.

To avoid repeating of same browser command you can user assignment syntax `command_name:=your command {}` for further use.

## Matchers

### fallback

This matcher always succeeds. Use it at the end of a configuration to specify the default browser. 

### app

Matches by source application.

Currently supported desktop environments: _hyprland_, _gnome_, _sway_, _macos_.

Hyprland Properties:

- *title*: match by source window title with regex
- *class*: match by window class

MacOS Properties:

- *display_name* - match by app display name (ex: `Slack`)
- *bundle_id* - match by App Bundle ID (ex: `com.tinyspeck.slackmacgap`)
- *bundle_path* - match by App Bundle path (ex: `/Applications/Slack.app`)
- *executable_path* - match by app executable path (ex: `/Applications/Slack.app/Contents/MacOS/Slack`)

### url

Match by a clicked URL.

Properties:

- *host*: match URL by host
- *scheme*: match URL by scheme
- *regex*: match full URL by regex

# Setup

## Linux

### Gnome

Due to stupidity of Gonme shell interface there is no legal way to recieve focused winow for Gnome with wayland: https://www.reddit.com/r/gnome/comments/pneza1/gdbus_call_for_moving_windows_not_working_in/

To be able to use the `app` matcher, please [install the extenions](https://extensions.gnome.org/extension/5592/focused-window-d-bus/), that exposes currently focused window via dbus interface: https://github.com/flexagoon/focused-window-dbus 

### Prebuilt packages

You can find `.rpm`, `.deb`, `.apk` and `.zst` packages on the release page.

### Linux manual

```sh
go install github.com/pltanton/autobrowser/linux/cmd/autobrowser@latest
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


# Acknowledgements

* [b-r-o-w-s-e](https://github.com/BlakeWilliams/b-r-o-w-s-e) project and [related article](https://blakewilliams.me/posts/handling-macos-url-schemes-with-go): great example of handling URLs with Golang on macOS
* [Finicky](https://github.com/johnste/finicky) project: inspiration for Autobrowser, good example of handling more complex URL events