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
      ${if builtins.isList command.cmd then "cmd = ${builtins.toJSON command.cmd}" else "cmd = ${builtins.toJSON command.cmd}"}
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
          if k != "type" then "${k} = ${builtins.toJSON v}" else ""
        ) matcher)}
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
            type = listOf attrs;
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
            { type = "app"; class = "Slack"; }
            { type = "url"; regex = ".*jira.*"; }
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
  config = mkIf cfg.enable {
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
  };
}
