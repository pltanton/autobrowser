{ lib, stdenv, buildGoModule, makeWrapper, ddcutil }:
buildGoModule {
  pname = "autobrowser";
  version = "0";
  vendorHash = "sha256-GWrflGa6evgehcHQujac67llQJWnRIQbBFna26DYizk=";
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
