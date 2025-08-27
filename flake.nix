{
  description = "Open specified browser depends on contextual rules";

  inputs.flakelight.url = "github:nix-community/flakelight";

  outputs = {flakelight, ...}:
    flakelight ./. {
      systems = ["x86_64-linux" "aarch64-linux" "x86_64-darwin" "aarch64-darwin"];
      devShell.packages = pkgs: with pkgs; [go alejandra dprint];
      formatters = {
        "*.yml" = "dprint fmt";
        "*.md" = "dprint fmt";
        "*.nix" = "alejandra";
      };
    };
}
