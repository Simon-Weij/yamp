{
  inputs = {
    nixpkgs.url = "github:Simon-Weij/nixpkgs/wails3";
    flakelight.url = "github:nix-community/flakelight";
  };
  outputs = {flakelight, ...} @ inputs:
    flakelight ./. {
      inherit inputs;
      devShell = pkgs: {
        packages = with pkgs; [wails3 go pnpm nodejs-slim pkg-config webkitgtk_6_0 golangci-lint just glib-networking];
        env.GIO_EXTRA_MODULES = "${pkgs.glib-networking}/lib/gio/modules";
      };
    };
}
