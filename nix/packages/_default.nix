{
  lib,
  buildGoModule,
  stdenv,
  darwin,
}:
buildGoModule {
  pname = "autobrowser";
  version = "1.0.2";
  vendorHash =
    if stdenv.isDarwin
    then "sha256-/8llw+85SbNKxlAfwZBJmHNYZunCZeXiMmoGzZ4eMYs="
    else "sha256-05D0rsPh/QLCL5i5c/xNTBozdRkPmtRQa5KU/Y0Y4pA=";

  src = import ../src.nix {inherit lib;};
  modRoot =
    if stdenv.isDarwin
    then "macos"
    else "linux";

  postInstall = lib.optionalString stdenv.isDarwin ''
    # Create macOS app bundle
    mkdir -p $out/Applications/Autobrowser.app/Contents/{MacOS,Resources}

    # Copy the binary to the app bundle
    mv $out/bin/autobrowser $out/Applications/Autobrowser.app/Contents/MacOS/

    # Copy app bundle assets
    cp $src/macos/assets/Info.plist $out/Applications/Autobrowser.app/Contents/
    cp $src/macos/assets/AppIcon.icns $out/Applications/Autobrowser.app/Contents/Resources/
  '';

  meta = with lib; {
    homepage = "https://github.com/pltanton/autobrowser";
    description = "Automatically determine browser depending on provided rules";
    license = licenses.gpl3Only;
    platforms = platforms.linux ++ platforms.darwin;
    mainProgram = "autobrowser";
  };
}
