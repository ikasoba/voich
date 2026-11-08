{
  description = "Go Template";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    { nixpkgs, flake-utils, ... }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
        nativeBuildInputs = with pkgs; [
          go
          gopls
        ];
        buildInputs =
          with pkgs;
          [
            pkg-config
            portaudio
            pipewire

            glfw
            # libglvnd
            # libxcursor
            # libx11
            # libxext
            # libxi
            # libxinerama
            # libxrandr
            # libxxf86vm
          ];
      in
      {
        devShells.default = pkgs.mkShell {
          inherit nativeBuildInputs buildInputs;
          shellHook = ''
            export LD_LIBRARY_PATH="$LD_LIBRARY_PATH:$NIX_LD_LIBRARY_PATH"
          '';
        };

        packages.default = pkgs.buildGoModule rec {
          name = "template";
          src = ./.;

          inherit buildInputs;

          vendorHash = null;
        };
      }
    );
}
