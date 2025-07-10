{
  lib,
  buildGoModule,
  stdenv,
  darwin,
}:
buildGoModule {
  pname = "autobrowser";
  version = "0";
  vendorHash =
    if stdenv.isDarwin
    then "sha256-w/MoA6uOgbQVPFzApJEDLeEFviKbvKjpdaIltyZ3he0="
    else "sha256-dvu80fFm3vIBjhk9k9Z5h9J5qTbvl3Tq1MCMQVJ+ru8=";

  src = import ../src.nix {inherit lib;};
  modRoot =
    if stdenv.isDarwin
    then "macos"
    else "linux";

  buildInputs = lib.optionals stdenv.isDarwin [
    darwin.apple_sdk.frameworks.Cocoa
    darwin.apple_sdk.frameworks.Foundation
  ];

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
