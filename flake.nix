{
  inputs.flakelight.url = "github:nix-community/flakelight";
  outputs = {flakelight, ...}:
    flakelight ./. {
      devShell.packages = pkgs: [pkgs.go pkgs.cobra-cli pkgs.mpv pkgs.golangci-lint pkgs.bun pkgs.just];
    };
}
