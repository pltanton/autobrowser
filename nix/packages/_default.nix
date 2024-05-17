{ lib, buildGoModule }:
buildGoModule {
  pname = "autobrowser";
  version = "0";
  # vendorHash = "sha256-W4WltWSZ1hJPUusaH3iRHo4HD5xmz3+/kOGxEspVF30=";
  vendorHash = "sha256-GWrflGa6evgehcHQujac67llQJWnRIQbBFna26DYizk=";
  meta = with lib; {
    homepage = "https://github.com/pltanton/autobrowser";
    description = "Automatically determine browser depending on provided rules";
    license = licenses.gpl3Only;
    platforms = platforms.linux;
    mainProgram = "autobrowser";
  };
  src = import ../src.nix { inherit lib; };

  modRoot = "linux";
}
