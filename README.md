# About

Automatically choosing web-browser depends on environment context rules.

## Features

- suckless solution with no redundant dependencies
- fast single binary
- simple rule engine with custom interpretable DSL
- cross-platform

# Configuration

## Example

```toml
# Variables - Define browser commands that can be reused in rules
[variables]
work = "firefox 'ext+container:name=Work&url={}'"
personal = "firefox {}"

# Rules - Define which browser to use based on context
[[rules]]
command = "work"
[rules.matchers]
app_class = "Slack"

[[rules]]
command = "work"
[rules.matchers]
url_regex = ".*jira.*"

# Default browser to use if no rules match
default = "personal"
```

Other examples can be found in the `examples` folder

## Configuration syntax

Autobrowser uses TOML for configuration. The application evaluates rules in order and applies the URL to the first matched command.

The configuration consists of three main sections:

1. **Variables**: Define reusable browser commands
2. **Rules**: Define matchers and corresponding browser commands
3. **Default**: Specify the fallback browser if no rules match

In browser commands, the `{}` placeholder will be replaced with the clicked URL.

## Matchers

The following matchers are available in the TOML configuration:

### App Matchers

Match by source application.

Currently supported desktop environments: _hyprland_, _gnome_, _sway_, _macos_.

Hyprland/Sway/Gnome Properties:

- `app_title`: Match by source window title with regex
- `app_class`: Match by window class

MacOS Properties:

- `app_display_name`: Match by app display name (ex: `Slack`)
- `app_bundle_id`: Match by App Bundle ID (ex: `com.tinyspeck.slackmacgap`)
- `app_bundle_path`: Match by App Bundle path (ex: `/Applications/Slack.app`)
- `app_executable_path`: Match by app executable path (ex: `/Applications/Slack.app/Contents/MacOS/Slack`)

### URL Matchers

Match by a clicked URL.

Properties:

- `url_host`: Match URL by host
- `url_scheme`: Match URL by scheme
- `url_regex`: Match full URL by regex

### Fallback

Set `fallback = true` in a rule's matchers or use the `default` setting at the root of the config to specify a default browser.

# Setup

## Linux

### Gnome

Due to stupidity of Gonme shell interface there is no legal way to recieve focused winow for Gnome with wayland: https://www.reddit.com/r/gnome/comments/pneza1/gdbus_call_for_moving_windows_not_working_in/

To be able to use the `app` matcher, please [install the extenions](https://extensions.gnome.org/extension/5592/focused-window-d-bus/), that exposes currently focused window via dbus interface: https://github.com/flexagoon/focused-window-dbus

### Prebuilt packages

You can find `.rpm`, `.deb`, `.apk` and `.zst` packages on the release page.

### Linux manual

Clone the repository and run, you can find a result binary in the `out` directory.

```sh
make build-linux
```

Create a TOML configuration file at `~/.config/autobrowser.toml`.
Then add the following .desktop file to `~/.local/share/applications/` and set it as the default browser.
Change paths for your setup if needed.

```ini
[Desktop Entry]
Categories=Network;WebBrowser
Exec=~/go/bin/autobrowser -config ~/.config/autobrowser.toml -url %u
Icon=browser
MimeType=x-scheme-handler/http;x-scheme-handler/https
Name=Autobrowser: select browser by contextual rules
Terminal=false
Type=Application
```

## Nix home-manager

This setup works booth for linux and darwin environments.

Actual flakes provides overlay (`overlays.default`) and module for home-manager (`autobrowser.homeModules.default`).

Example of home-manager module configuration:

```nix
{
  inputs,
  ...
}: {
  programs.autobrowser = {
    package = inputs.autobrowser.packages.x86_64-linux.default;
    enable = true;
    variables = {
      work = "firefox 'ext+container:name=Work&url={}'";
      home = "firefox {}";

      # Example for darwin (MacOS) configuration
      work-darwin = "open -a 'Zen' 'ext+container:name=Work&url={}'";
    };
    # Your configuration will be automatically converted to TOML format internally
    rules = [
      "work:app.class=Slack"
      "work:app.class=org.telegram.desktop;app.title='.*Work related group name.*'"

      # Example for darwin (MacOS) configuration
      "work-darwin:app.bundle_id='com.tinyspeck.slackmacgap'"
    ];
    default = "home";
  };
}
```

# Migration from Old Format to TOML

Autobrowser now uses TOML for configuration. Here's how to migrate your existing configuration:

## Variable Definitions

Old format:
```
work:=firefox 'ext+container:name=Work&url={}'
```

New TOML format:
```toml
[variables]
work = "firefox 'ext+container:name=Work&url={}'"
```

## Rules

Old format:
```
work:app.class=Slack
firefox {}:url.regex='.*github\.com.*'
firefox {}:fallback
```

New TOML format:
```toml
[[rules]]
command = "work"
[rules.matchers]
app_class = "Slack"

[[rules]]
command = "firefox {}"
[rules.matchers]
url_regex = ".*github\\.com.*"

# For fallback, either use:
[[rules]]
command = "firefox {}"
[rules.matchers]
fallback = true

# Or more simply:
default = "firefox {}"
```

Note: The order of rules in the TOML file matters, just as it did in the old format.

# Acknowledgements

- [b-r-o-w-s-e](https://github.com/BlakeWilliams/b-r-o-w-s-e) project and [related article](https://blakewilliams.me/posts/handling-macos-url-schemes-with-go): great example of handling URLs with Golang on macOS
- [Finicky](https://github.com/johnste/finicky) project: inspiration for Autobrowser, good example of handling more complex URL events
