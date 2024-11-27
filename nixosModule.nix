packages:
{ config, pkgs, lib, ... }:
let
  webserver = packages.${pkgs.system}.default;
  inherit (lib) mkIf mkEnableOption;
  cfg = config.mannedotdev;
in {
  options.mannedotdev = { enable = mkEnableOption "Enable the webserver"; };

  config = mkIf cfg.enable {
    services.nginx = {
      enable = true;
      recommendedTlsSettings = true;
      recommendedOptimisation = true;
      recommendedGzipSettings = true;
      recommendedProxySettings = true;
      virtualHosts."manne.dev" = {
        addSSL = false;
        enableACME = false;
        listen = [{
          addr = "95.217.0.108";
          port = 80;
        }];

        locations."/" = {
          proxyPass = "http://127.0.0.1:8080";
          proxyWebsockets = false;
        };
      };
    };

    systemd.services.mannedotdev = {
      enable = true;
      after = [ "network.target" ];
      wantedBy = [ "multi-user.target" ];
      description = "Webserver hosting manne.dev";
      serviceConfig = {
        Type = "simple";
        ExecStart = "${webserver}/bin/mannedotdev";
        WorkingDirectory = "${webserver}/data";
      };
    };

    networking.firewall.allowedTCPPorts = [ 80 ];
  };
}
