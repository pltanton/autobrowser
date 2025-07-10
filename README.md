# About

Automatically choosing web-browser depends on environment context rules.

## Features

- suckless solution with no redundant dependencies
- fast single binary
- simple rule engine with custom interpretable DSL
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

Hyprland/Sway/Gnome Properties:

- _title_: match by source window title with regex
- _class_: match by window class

MacOS Properties:

- _display_name_ - match by app display name (ex: `Slack`)
- _bundle_id_ - match by App Bundle ID (ex: `com.tinyspeck.slackmacgap`)
- _bundle_path_ - match by App Bundle path (ex: `/Applications/Slack.app`)
- _executable_path_ - match by app executable path (ex: `/Applications/Slack.app/Contents/MacOS/Slack`)

### url

Match by a clicked URL.

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

Create config at `~/.config/autobrowser.config`.
Then add the following .desktop file to `~/.local/share/applications/` and set it as the default browser.
Change paths for your setup if needed.

```ini
[Desktop Entry]
Categories=Network;WebBrowser
Exec=~/go/bin/autobrowser -config ~/.config/autobrowser.config -url %u
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

# Acknowledgements

- [b-r-o-w-s-e](https://github.com/BlakeWilliams/b-r-o-w-s-e) project and [related article](https://blakewilliams.me/posts/handling-macos-url-schemes-with-go): great example of handling URLs with Golang on macOS
- [Finicky](https://github.com/johnste/finicky) project: inspiration for Autobrowser, good example of handling more complex URL events
