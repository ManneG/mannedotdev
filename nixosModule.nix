packages:
{ config, pkgs, lib, ... }:
let
  webserver = packages.${pkgs.system}.default;
  inherit (lib) mkIf mkOption mkEnableOption;
  cfg = config.mannedotdev;
in /* assert cfg.enable -> cfg.acme.credentials != null;
assert cfg.enable -> cfg.acme.email != null;  */{
  options.mannedotdev = {
    enable = mkEnableOption "Enable the webserver";
    acme.credentials = mkOption {
      type = lib.types.path;
      default = null;
      example = /run/secrets/credentials.ini;
      description = "Path to CloudFlare API token";
    };
    acme.email = mkOption {
      type = lib.types.str;
      default = null;
      example = "foo@example.com";
      description = "ACME Email adress";
    };
  };

  config = mkIf cfg.enable {
    services.nginx = {
      enable = true;
      recommendedTlsSettings = true;
      recommendedOptimisation = true;
      recommendedGzipSettings = true;
      recommendedProxySettings = true;
      virtualHosts."manne.dev" = {
        onlySSL = true;
        useACMEHost = "manne.dev"; # Depends on security.acme.certs."manne.dev"

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

    security.acme = {
      certs."manne.dev" = {
        credentialFiles."CLOUDFLARE_DNS_API_TOKEN_FILE" = cfg.acme.credentials;
        dnsProvider = "cloudflare";
        email = cfg.acme.email;
        extraDomainNames = [ ];
        group = "nginx";
        reloadServices = [ "nginx.service" ];
      };

      acceptTerms = true;
    };

    networking.firewall.allowedTCPPorts = [ 443 ];
  };
}
