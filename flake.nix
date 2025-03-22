{
  description = "The website hosted on manne.dev";

  inputs.nixpkgs.url = "nixpkgs/nixos-24.05";

  outputs = { self, nixpkgs }:
    let
      # to work with older version of flakes
      lastModifiedDate =
        self.lastModifiedDate or self.lastModified or "19700101";

      # Generate a user-friendly version number.
      version = builtins.substring 0 8 lastModifiedDate;

      # System types to support.
      supportedSystems =
        [ "x86_64-linux" "x86_64-darwin" "aarch64-linux" "aarch64-darwin" ];

      # Helper function to generate an attrset '{ x86_64-linux = f "x86_64-linux"; ... }'.
      forAllSystems = nixpkgs.lib.genAttrs supportedSystems;

      # Nixpkgs instantiated for supported system types.
      nixpkgsFor = forAllSystems (system: import nixpkgs { inherit system; });

    in {

      formatter = forAllSystems (system:
        let pkgs = nixpkgsFor.${system};
        in pkgs.writeShellApplication {
          name = "format";
          runtimeInputs = [ pkgs.nixfmt-classic pkgs.go ];
          text = ''
            find . -name '*.nix' -exec nixfmt {} +
            gofmt -w .
          '';
        });

      packages = forAllSystems (system:
        let pkgs = nixpkgsFor.${system};
        in {
          mannedotdev = pkgs.buildGoModule {
            pname = "mannedotdev";
            inherit version;

            src = ./.;

            postInstall = ''
              mkdir -p $out/data
              cp -r ./static $out/data
              cp ./template.html $out/data/template.html
            '';

            vendorHash = "sha256-7EvBT8JiVkebxOg6+55ZABcTjSipP42X5irOKRvmTRY=";
          };
          default = self.packages.${system}.mannedotdev;
        });

      nixosModules = {
        mannedotdev = import ./nixosModule.nix self.packages;
        default = self.nixosModules.mannedotdev;
      };

      devShells = forAllSystems (system:
        let pkgs = nixpkgsFor.${system};
        in {
          default = pkgs.mkShell {
            buildInputs = with pkgs; [ go gopls gotools go-tools ];
          };
          withpackage =
            pkgs.mkShell { buildInputs = [ self.packages.${system}.default ]; };
        });
    };
}
