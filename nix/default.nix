{ lib, stdenv, buildGoModule, makeWrapper, ddcutil }:
buildGoModule {
  pname = "autobrowser";
  version = "0";

  vendorHash = "sha256-y8Q4P4k3KyMHkbTU/usD82XIOw4hO0Uj8AbCClsZgwc=";

  src = lib.cleanSourceWith {
    filter = name: type:
      let baseName = baseNameOf (toString name);
      in !(lib.hasSuffix ".nix" baseName);
    src = lib.cleanSource ../.;
  };

  allowGoReference = true;

  nativeBuildInputs = [ makeWrapper ];

  postFixup = ''
    mv $out/bin/autobrowser-linux $out/bin/autobrowser
  '';

  meta = with lib; {
    homepage = "https://github.com/pltanton/autobrowser";
    description = "Automatically determine browser depending on provided rules";
    license = licenses.bsd3;
    platforms = platforms.linux;
    mainProgram = "autobrowser";
  };
}