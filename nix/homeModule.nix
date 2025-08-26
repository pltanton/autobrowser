{
  config,
  lib,
  pkgs,
  ...
}:
with lib; let
  cfg = config.programs.autobrowser;
  configText = ''
    default_command = "${cfg.defaultCommand}"

    ${lib.concatStringsSep "\n\n" (lib.mapAttrsToList (name: command: ''
      [command.${name}]
      cmd = ${builtins.toJSON (if builtins.isList command.cmd then command.cmd else command.cmd)}
      ${lib.optionalString (command.placeholder != null) "placeholder = ${builtins.toJSON command.placeholder}"}
      ${lib.optionalString command.queryEscape "query_escape = true"}
    '') cfg.commands)}

    ${lib.concatStringsSep "\n\n" (map (rule: ''
      [[rules]]
      command = "${rule.command}"
      ${lib.concatStringsSep "\n" (map (matcher: ''
        [[rules.matchers]]
        type = "${matcher.type}"
        ${lib.concatStringsSep "\n" (lib.mapAttrsToList (k: v:
          if k != "type" && v != null then "${(lib.replaceStrings ["displayName" "bundleId" "bundlePath" "executablePath"] ["display_name" "bundle_id" "bundle_path" "executable_path"] k)} = ${builtins.toJSON v}" else ""
        ) (removeAttrs (lib.filterAttrs (n: v: v != null) matcher) ["type"]))}
      '') rule.matchers)}
    '') cfg.rules)}
  '';
in {
  options.programs.autobrowser = {
    enable = lib.mkEnableOption "whenever to enable autobrowser as default browser";
    package = mkPackageOption pkgs "autobrowser" {};
    commands = mkOption {
      type = with lib.types; attrsOf (submodule {
        options = {
          cmd = mkOption {
            type = either str (listOf str);
            description = "Command to execute (string or list of strings)";
            example = ''["firefox", "--new-tab", "{}"]'';
          };
          placeholder = mkOption {
            type = nullOr str;
            description = "Placeholder to replace with URL (default is {})";
            default = null;
            example = "{{url}}";
          };
          queryEscape = mkOption {
            type = bool;
            description = "Whether to apply URL query escaping";
            default = false;
          };
        };
      });
      description = "Commands that can be executed";
      default = {};
      example = {
        personal = {
          cmd = "firefox {}";
        };
        work = {
          cmd = ["firefox" "--private-window" "{}"];
          queryEscape = true;
        };
      };
    };

    rules = mkOption {
      type = with lib.types; listOf (submodule {
        options = {
          command = mkOption {
            type = str;
            description = "Command to use when rule matches";
            example = "personal";
          };
          matchers = mkOption {
            type = with lib.types; listOf (submodule {
              options = {
                type = mkOption {
                  type = enum [ "url" "app" "fallback" ];
                  description = "Matcher type";
                  example = "url";
                };

                # URL matcher options
                regex = mkOption {
                  type = nullOr str;
                  description = "Match URL by regex pattern (for URL matchers)";
                  default = null;
                  example = ".*github\\.com.*";
                };
                host = mkOption {
                  type = nullOr str;
                  description = "Match URL by host (for URL matchers)";
                  default = null;
                  example = "github.com";
                };
                scheme = mkOption {
                  type = nullOr str;
                  description = "Match URL by scheme (for URL matchers)";
                  default = null;
                  example = "https";
                };

                # App matcher options - Linux
                class = mkOption {
                  type = nullOr str;
                  description = "Match by window class (Linux, for app matchers)";
                  default = null;
                  example = "Firefox";
                };
                title = mkOption {
                  type = nullOr str;
                  description = "Match window title by regex pattern (Linux, for app matchers)";
                  default = null;
                  example = ".*GitHub.*";
                };

                # App matcher options - macOS
                displayName = mkOption {
                  type = nullOr str;
                  description = "Match by app display name (macOS, for app matchers)";
                  default = null;
                  example = "Safari";
                };
                bundleId = mkOption {
                  type = nullOr str;
                  description = "Match by App Bundle ID (macOS, for app matchers)";
                  default = null;
                  example = "com.apple.Safari";
                };
                bundlePath = mkOption {
                  type = nullOr str;
                  description = "Match by App Bundle path (macOS, for app matchers)";
                  default = null;
                  example = "/Applications/Safari.app";
                };
                executablePath = mkOption {
                  type = nullOr str;
                  description = "Match by app executable path (macOS, for app matchers)";
                  default = null;
                  example = "/Applications/Safari.app/Contents/MacOS/Safari";
                };
              };
            });
            description = "List of matchers (all must match for rule to apply)";
            example = [
              { type = "url"; regex = ".*github.com.*"; }
              { type = "app"; class = "Terminal"; }
            ];
          };
        };
      });
      description = "List of rules to match";
      default = [];
      example = [
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

    defaultCommand = mkOption {
      type = lib.types.str;
      description = "Default command to use when no rules match";
      default = "";
      example = "personal";
    };

    desktop = mkOption {
      type = lib.types.package;
      description = "Desktop entry for autobrowser";
      default = pkgs.writeTextDir "share/applications/autobrowser.desktop"
        (lib.generators.toINI {} {
          "Desktop Entry" = {
            Type = "Application";
            Exec = "${cfg.package}/bin/autobrowser -url %u";
            Terminal = false;
            Name = "Autobrowser: select browser by contextual rules";
            Icon = "browser";
            Categories = "Network;WebBrowser";
            MimeType = "x-scheme-handler/http;x-scheme-handler/https";
          };
        });
    };
  };

  config = lib.mkMerge [
    # Type validation through assertions
    {
      assertions = flatten (map (rule:
        map (matcher: [
          {
            # For URL matchers, only URL fields should be set
            assertion = matcher.type == "url" -> (
              matcher.class == null &&
              matcher.title == null &&
              matcher.displayName == null &&
              matcher.bundleId == null &&
              matcher.bundlePath == null &&
              matcher.executablePath == null
            );
            message = "URL matcher should only use URL-specific fields (regex, host, scheme)";
          }
          {
            # For app matchers, only app fields should be set
            assertion = matcher.type == "app" -> (
              matcher.regex == null &&
              matcher.host == null &&
              matcher.scheme == null
            );
            message = "App matcher should only use app-specific fields";
          }
          {
            # For fallback matchers, no extra fields should be set
            assertion = matcher.type == "fallback" -> (
              matcher.regex == null &&
              matcher.host == null &&
              matcher.scheme == null &&
              matcher.class == null &&
              matcher.title == null &&
              matcher.displayName == null &&
              matcher.bundleId == null &&
              matcher.bundlePath == null &&
              matcher.executablePath == null
            );
            message = "Fallback matcher should not define any additional fields";
          }
        ]) rule.matchers
      ) (cfg.rules or []));
    }

    (mkIf cfg.enable {
      xdg.configFile."autobrowser/config.toml".text = configText;

      home.packages =
        [cfg.package]
        ++ (
          if pkgs.stdenv.isLinux
          then [cfg.desktop]
          else []
        );

      xdg.mimeApps = mkIf pkgs.stdenv.isLinux {
        defaultApplications = {
          "x-scheme-handler/http" = "autobrowser.desktop";
          "x-scheme-handler/https" = "autobrowser.desktop";
          "x-scheme-handler/about" = "autobrowser.desktop";
        };
      };
    })
  ];
}
