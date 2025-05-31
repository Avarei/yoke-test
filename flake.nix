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
              version = "v0.12.9";
              subPackages = [ "cmd/yoke" ];
              src = pkgs.fetchFromGitHub {
                owner = "yokecd";
                repo = "yoke";
                rev = "v0.12.9";
                sha256 = "sha256-4n4hzwOuzS3bZ2vAr0fn+3urlT7ihC+cStRddmtqKPg";
              };
              vendorHash = "sha256-Lqzi7oRmnhINZY+Tbkh42qhNaKtExUc3kUBjufxCyLw=";
              doCheck = false;
            };
      in
      {
        devShells.default = pkgs.mkShell {
          packages = [
            pkgs.go
            yoke
            pkgs.cobra-cli
          ];
        };
      }
    );
}
