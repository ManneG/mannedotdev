packages: { config, pkgs, lib, ...}:
let
  webserver = packages.${pkgs.system}.default;
  inherit (lib) mkIf mkEnableOption;
  cfg = config.mannedotdev;
in {
  options.mannedotdev = {
    enable = mkEnableOption "Enable the webserver";
  };

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
        listen = [
          {
            addr = "";
            port = 80;
          }
        ];

        locations."/" = {
          proxyPass = "127.0.0.1:8080";
          proxyWebsockets = false;
          extraConfig = ''
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            
          '';
        };
      };
    };

    systemd.user.services.mannedotdev-service = {
      enable = true;
      after = [ "network.target" ];
      wantedBy = [ "multi-user.target" ];
      description = "Webserver hosting manne.dev";
      serviceConfig = {
          Type = "simple";
          ExecStart = webserver;
      };
    };

    networking.firewall.allowedTCPPorts = [ 80 8080 ];
  };
}