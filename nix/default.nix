{ lib, stdenv, buildGoModule, makeWrapper, ddcutil }:
let
  sources = lib.cleanSourceWith {
    filter = name: type:
      let baseName = baseNameOf (toString name);
      in !(lib.hasSuffix ".nix" baseName);
    src = lib.cleanSource ../.;
  };
in
buildGoModule {
  pname = "autobrowser";
  version = "0";

  vendorHash = "sha256-4vLAS5eQyvE5bsQ35q0PYdu1zUxYT34Y0gC/6nSfPI8=";

  src = sources;

  modRoot = "linux";

  meta = with lib; {
    homepage = "https://github.com/pltanton/autobrowser";
    description = "Automatically determine browser depending on provided rules";
    license = licenses.gpl3Only;
    platforms = platforms.linux;
    mainProgram = "autobrowser";
  };
}
