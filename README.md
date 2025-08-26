# About

Automatically choosing web-browser depends on environment context rules.

## Features

- suckless solution with no redundant dependencies
- fast single binary
- simple rule engine with custom interpretable DSL
- cross-platform

# Configuration

Autobrowser uses TOML for configuration, providing a structured and maintainable approach.

## Configuration Example

```toml
# Default command to use when no rules match
default_command = "personal"

# Define commands that can be executed
[command.work]
cmd = ["firefox", "-p", "work", "{}"]
query_escape = true

[command.personal]
cmd = "firefox {}"

# Define rules that determine which command to use
[[rules]]
command = "work"
# Multiple matchers in a rule work as AND conditions
[[rules.matchers]]
type = "app"
class = "Slack"
[[rules.matchers]]
type = "url"
regex = ".*jira.*"

[[rules]]
command = "work"
[[rules.matchers]]
type = "app"
class = "org.telegram.desktop"
```

Other examples can be found in the `examples` folder.

## Configuration Syntax

The configuration consists of three main sections:

1. **Commands** - Define reusable browser commands with the `command.<name>` structure
2. **Rules** - Define conditions for browser selection using `[[rules]]` arrays
3. **Default Command** - Define the default browser when no rules match with `default_command`

Save your configuration as `~/.config/autobrowser/config.toml` or pass a custom path to Autobrowser using the `-config` flag.

## Matchers

Matchers define the conditions for rules to be applied.

### fallback

The fallback matcher always succeeds. Instead of using a matcher, you can simply set the `default_command` property:

```toml
default_command = "firefox"
```

### app

Matches by source application.

Currently supported desktop environments: _hyprland_, _gnome_, _sway_, _macos_. When using Home Manager, Linux and macOS app matcher properties are properly typed and validated.

```toml
[[rules.matchers]]
type = "app"
class = "Slack"
```

Hyprland/Sway/Gnome Properties:

- _title_: match by source window title with regex
- _class_: match by window class

MacOS Properties:

- _display_name_ - match by app display name (ex: `Slack`)
- _bundle_id_ - match by App Bundle ID (ex: `com.tinyspeck.slackmacgap`)
- _bundle_path_ - match by App Bundle path (ex: `/Applications/Slack.app`)
- _executable_path_ - match by app executable path (ex: `/Applications/Slack.app/Contents/MacOS/Slack`)

When using these properties in Nix Home Manager configuration, use camelCase format: `displayName`, `bundleId`, `bundlePath`, `executablePath`.

Note: When using Home Manager, you only need to specify the properties you want to use, and they are strictly typed.

### url

Match by a clicked URL.

```toml
[[rules.matchers]]
type = "url"
host = "github.com"
```

Properties:

- _host_: match URL by host
- _scheme_: match URL by scheme
- _regex_: match full URL by regex

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

Create config at `~/.config/autobrowser/config.toml`.
Then add the following .desktop file to `~/.local/share/applications/` and set it as the default browser.
Change paths for your setup if needed.

```ini
[Desktop Entry]
Categories=Network;WebBrowser
Exec=~/go/bin/autobrowser -config ~/.config/autobrowser/config.toml -url %u
Icon=browser
MimeType=x-scheme-handler/http;x-scheme-handler/https
Name=Autobrowser: select browser by contextual rules
Terminal=false
Type=Application
```

## Nix home-manager

This setup works both for Linux and macOS (Darwin) environments.

The flake provides an overlay (`overlays.default`) and a module for home-manager (`homeModules.default`).

The home-manager module provides strictly typed matcher configurations for improved type safety and validation. Each matcher type (`url`, `app`, `fallback`) has its own specific set of properties that are properly validated.

Note: Nix configuration uses camelCase for properties (e.g., `bundleId`), while the generated TOML uses snake_case (e.g., `bundle_id`).

Example of home-manager module configuration:

```nix
{
  pkgs,
  inputs,
  ...
}: {
  programs.autobrowser = {
    enable = true;
    package = pkgs.autobrowser;
    defaultCommand = "personal";
    
    commands = {
      work = {
        cmd = ["firefox", "ext+container:name=Work&url={}"];
        placeholder = "{}";
        queryEscape = true;  # Will generate query_escape = true in TOML
      };
      personal = {
        cmd = "firefox {}";
      };
      youtube = {
        cmd = ["mpv", "{}"];
      };
    };
    
    rules = [
      # Open Jira links from Slack in work container
      {
        command = "work";
        matchers = [
          # Strictly typed app matcher for Linux
          { 
            type = "app";    # Required field
            class = "Slack"; # Window class (Linux)
            # title = ".*";  # Optional window title regex
          }
          # Strictly typed URL matcher
          {
            type = "url";      # Required field
            regex = ".*jira.*"; # Match URL by regex pattern
          }
        ];
      },
      # Open GitHub links in personal browser
      {
        command = "personal";
        matchers = [
          {
            type = "url";
            host = "github.com";
          }
        ];
      },
      # Open YouTube videos in mpv
      {
        command = "youtube";
        matchers = [
          {
            type = "url";
            regex = ".*youtube\\.com/watch.*|.*youtu\\.be/.*";
          }
        ];
      },
      # macOS specific example
      {
        command = "personal";
        matchers = [
          {
            type = "app";
            # macOS specific properties (use camelCase in Nix)
            displayName = "Safari";      # Match by app display name
            # bundleId = "com.apple.Safari";  # Optional
            # bundlePath = "/Applications/Safari.app";  # Optional
            # executablePath = "/Applications/Safari.app/Contents/MacOS/Safari";  # Optional
          }
        ];
      },
      # Simple fallback example
      {
        command = "personal";
        matchers = [
          { type = "fallback"; } # Always matches
        ];
      }
    ];
  };
}
```

# Acknowledgements

- [b-r-o-w-s-e](https://github.com/BlakeWilliams/b-r-o-w-s-e) project and [related article](https://blakewilliams.me/posts/handling-macos-url-schemes-with-go): great example of handling URLs with Golang on macOS
- [Finicky](https://github.com/johnste/finicky) project: inspiration for Autobrowser, good example of handling more complex URL events
