{
  description = "Open specified browser depends on contextual rules";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
  inputs.flake-utils.url = "github:numtide/flake-utils";

  outputs = { nixpkgs, flake-utils, ... }:
    flake-utils.lib.eachDefaultSystem
      (system:
        let
          pkgs = import nixpkgs { inherit system; };
        in
        {
          packages.autobrowser = pkgs.callPackage ./nix/linux.nix { };
          packages.default = pkgs.callPackage ./nix/linux.nix { };

          checks.autobrowser-common = pkgs.callPackage ./nix/common.nix { };

          devShells. default = with pkgs; mkShell {
            buildInputs = [ go ];
          };
        }) //
    {
      overlays.default = final: prev: rec {
        autobrowser = final.pkgs.callPackage ./nix/linux.nix { };
      };

      homeManagerModules.default = import ./nix/hm-module.nix;
    };
}
