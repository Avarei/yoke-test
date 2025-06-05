{
  description = "A basic flake with a shell";
  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
  inputs.systems.url = "github:nix-systems/default";
  inputs.flake-utils = {
    url = "github:numtide/flake-utils";
    inputs.systems.follows = "systems";
  };

  outputs =
    { nixpkgs, flake-utils, ... }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
        yoke = pkgs.buildGo124Module rec {
              pname = "yoke";
              version = "v0.13.4";
              subPackages = [ "cmd/yoke" ];
              src = pkgs.fetchFromGitHub {
                owner = "yokecd";
                repo = "yoke";
                rev = "v0.13.4";
                sha256 = "sha256-3WTXz+ZDPUJ2+v4xVqQCXDHhdi3rJLowsDAzvovJbLQ=";
              };
              vendorHash = "sha256-Z3hkYD6QnKS1kEkuF0aLfypaq+J/8ECApkU1UYVukU4=";
              doCheck = false;
            };
      in
      {
        devShells.default = pkgs.mkShell {
          packages = [
            pkgs.go
            yoke
            pkgs.cobra-cli
            pkgs.mage
          ];
        };
      }
    );
}
