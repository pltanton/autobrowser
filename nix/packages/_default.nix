{
  lib,
  buildGoModule,
}:
buildGoModule {
  pname = "autobrowser";
  version = "0";
  vendorHash = "sha256-dvu80fFm3vIBjhk9k9Z5h9J5qTbvl3Tq1MCMQVJ+ru8=";
  meta = with lib; {
    homepage = "https://github.com/pltanton/autobrowser";
    description = "Automatically determine browser depending on provided rules";
    license = licenses.gpl3Only;
    platforms = platforms.linux;
    mainProgram = "autobrowser";
  };
  src = import ../src.nix {inherit lib;};

  modRoot = "linux";
}
