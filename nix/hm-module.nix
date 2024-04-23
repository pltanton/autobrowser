{ config, lib, pkgs, ... }:
with lib;
let
  cfg = config.programs.autobrowser;
  configFile = pkgs.writeText "autobrowser.config"
    (builtins.concatStringsSep "\n" (
      (lib.mapAttrsToList (k: v: "${k}:=${v}") cfg.variables) ++
      cfg.rules ++
      [ "${cfg.default}:fallback" ]
    ));
in
{
  options.programs.autobrowser = {
    enable = lib.mkEnableOption "whenever to enable autobrowser as default browser";
    package = mkPackageOption pkgs "autobrowser" { };
    variables = mkOption {
      type = with lib.types; attrsOf str;
      description = "Attribute set of variables";
      default = { };
    };
    rules = mkOption {
      type = with lib.types; listOf str;
      example = [ "firefox {}:app.class=telegram" "firefox -p work {}:url.regex='.*atlassian.org.*'" ];
      description = "List of rules";
    };
    default = mkOption {
      type = lib.types.str;
      description = "Default browser command";
      default = "";
      example = "firefox {}";
    };
  };
  config = mkIf cfg.enable {
    home.packages = [
      (pkgs.writeTextDir "share/applications/autobrowser.desktop"
        (lib.generators.toINI { } {
          "Desktop Entry" = {
            Type = "Application";
            Exec = "${cfg.package}/bin/autobrowser -config ${configFile} -url %u";
            Terminal = false;
            Name = "Autobrowser: select browser by contextual rules";
            Icon = "browser";
            Categories = "Network;WebBrowser";
            MimeType = "x-scheme-handler/http;x-scheme-handler/https";
          };
        }))
    ];

    xdg.mimeApps = {
      defaultApplications = {
        "x-scheme-handler/http" = "autobrowser.desktop";
        "x-scheme-handler/https" = "autobrowser.desktop";
        "x-scheme-handler/about" = "autobrowser.desktop";
      };
    };
  };
}

