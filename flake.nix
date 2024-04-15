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
          packages.default = pkgs.callPackage ./nix/default.nix { };

          devShells.default = with pkgs; mkShell {
            buildInputs = [ 
              go 
              clang 
              gnustep.make 
              gnustep.base 
              gnustep.gui 
              gnustep.libobjc 
              gnustep.stdenv
              gnustep.gworkspace
            ];
            CC = "clang";
            shellHook = ''
              export CC="clang"
              export CGO="-fobjc-nonfragile-abi"
            '';
          };
        }) //
    {
      overlays.default = final: prev: rec {
        autobrowser = final.pkgs.callPackage ./nix/default.nix { };
      };


      homeManagerModules.default = import ./nix/hm-module.nix;
    };
}
