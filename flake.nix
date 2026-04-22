{
  inputs.flakelight.url = "github:nix-community/flakelight";
  outputs = { flakelight, ... }:
    flakelight ./. {
      devShell.packages = pkgs: with pkgs; [
        python313
        python313Packages.pyside6
        qt6.qtdeclarative
      ];
      devShell.env = pkgs: {
        LD_LIBRARY_PATH = pkgs.lib.makeLibraryPath [
          pkgs.stdenv.cc.cc.lib
          pkgs.qt6.qtbase
        ];
        QML2_IMPORT_PATH = pkgs.lib.makeSearchPath "lib/qt-6/qml" [
          pkgs.qt6.qtdeclarative
        ];
        QML_IMPORT_PATH = pkgs.lib.makeSearchPath "lib/qt-6/qml" [
          pkgs.qt6.qtdeclarative
          pkgs.python313Packages.pyside6
        ];
      };
    };
}
