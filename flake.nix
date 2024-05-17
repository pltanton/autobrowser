{
  description = "Open specified browser depends on contextual rules";

  inputs.flakelight.url = "github:nix-community/flakelight";

  outputs = inputs@{ flakelight, ... }: flakelight ./. {
    devShell.packages = pkgs: with pkgs; [ go ];
    homeModule = ./nix/hm-module.nix;
  };
}
