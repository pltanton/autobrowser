{ lib, stdenv, buildGoModule, makeWrapper, ddcutil }:
buildGoModule {
  pname = "autobrowser";
  version = "0";
  vendorHash = "sha256-4vLAS5eQyvE5bsQ35q0PYdu1zUxYT34Y0gC/6nSfPI8=";
  meta = with lib; {
    homepage = "https://github.com/pltanton/autobrowser";
    description = "Automatically determine browser depending on provided rules";
    license = licenses.gpl3Only;
    platforms = platforms.linux;
    mainProgram = "autobrowser";
  };
  src = import ./src.nix { inherit lib; };

  modRoot = "linux";

}
