{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };

  outputs = {nixpkgs, ...}: let
    systems = ["x86_64-linux" "aarch64-linux" "x86_64-darwin" "aarch64-darwin"];
    forAllSystems = f: nixpkgs.lib.genAttrs systems (system: f nixpkgs.legacyPackages.${system});
  in {
    devShells = forAllSystems (pkgs: {
      default = pkgs.mkShell {
        packages = let
          wails3 = pkgs.callPackage ./nix/wails3.nix {};
        in
          with pkgs; [wails3 go pnpm nodejs-slim pkg-config webkitgtk_6_0 golangci-lint just glib-networking clang clang-tools];

        env.GIO_EXTRA_MODULES = "${pkgs.glib-networking}/lib/gio/modules";

        buildInputs = with pkgs; [
          mpv
        ];
      };
    });
  };
}
