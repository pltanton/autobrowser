# Autobrowser

Automatically selects web browsers based on context rules.

## Features

- Fast single binary with minimal dependencies
- Simple rule engine
- Cross-platform (Linux, macOS)

## Configuration

Save as `~/.config/autobrowser/config.toml` or use `-config` flag.

### Example

```toml
# Default when no rules match
default_command = "personal"

# Define commands
[command.work]
cmd = ["firefox", "-p", "work", "{}"]
query_escape = true

[command.personal]
cmd = "firefox {}"

# Define rules
[[rules]]
command = "work"
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

More examples in the `examples` folder.

### Matchers

#### app

Match by source application.

Supported environments: _hyprland_, _gnome_, _sway_, _macos_

```toml
[[rules.matchers]]
type = "app"
class = "Slack"
```

**Linux Properties:**
- `title`: window title (regex)
- `class`: window class

**macOS Properties:**
- `display_name`: app name
- `bundle_id`: App Bundle ID
- `bundle_path`: App Bundle path
- `executable_path`: app executable path

#### url

Match by clicked URL.

```toml
[[rules.matchers]]
type = "url"
host = "github.com"
```

**Properties:**
- `host`: match by host
- `scheme`: match by scheme
- `regex`: match full URL by regex

## Setup

### Linux

#### Gnome

Install [focused-window-dbus extension](https://github.com/flexagoon/focused-window-dbus) to expose the focused window.

#### Installation

**Prebuilt packages:**
Download `.rpm`, `.deb`, `.apk` or `.zst` from releases.

**Manual build:**
```sh
make build-linux
```

Create this `.desktop` file in `~/.local/share/applications/` and set as default browser:

```ini
[Desktop Entry]
Categories=Network;WebBrowser
Exec=/path/to/autobrowser -config ~/.config/autobrowser/config.toml -url %u
Icon=browser
MimeType=x-scheme-handler/http;x-scheme-handler/https
Name=Autobrowser
Terminal=false
Type=Application
```

### Nix home-manager

Works for both Linux and macOS. The flake provides an overlay and a home-manager module.

Example configuration:

```nix
{
  programs.autobrowser = {
    enable = true;
    defaultCommand = "personal";
    
    commands = {
      work = {
        cmd = ["firefox", "ext+container:name=Work&url={}"];
        queryEscape = true;
      };
      personal = {
        cmd = "firefox {}";
      };
    };
    
    rules = [
      {
        command = "work";
        matchers = [
          { 
            type = "app";
            class = "Slack";
          }
          {
            type = "url";
            regex = ".*jira.*";
          }
        ];
      }
    ];
  };
}
```

## Debugging

### macOS

Monitor logs:

```
log stream --predicate 'subsystem == "dev.pltanton.autobrowser"' --style compact --level debug
```

## Acknowledgements

- [b-r-o-w-s-e](https://github.com/BlakeWilliams/b-r-o-w-s-e)
- [Finicky](https://github.com/johnste/finicky)