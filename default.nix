{ buildGoModule, lib, fetchFromGitHub, stdenv, ... }:

buildGoModule rec {
  pname = "skuld";
  name = "skuld-${version}";
  version = "0.6.4";

  src = fetchFromGitHub {
    owner = "DEEP-IMPACT-AG";
    repo = "skuld";
    sha256 = "nE8/h+poDVRkSLWnAaqRDOxNsXqhU/OyG8UHHc5I53c=";
    rev = "v${version}";
  };

  vendorSha256 = "LxSeLzkseBpV4Pqs+rnlBKMNOIpkpa5uZiImkVVF1SI=";

  modSha256 = stdenv.lib.fakeSha256;

  subPackages = ["."];

  runVend = false;

  meta = with lib; {
    description = "CLI utility to help developers use 2FA with AWS";
    homepage = "https://github.com/DEEP-IMPACT-AG/skuld";
    license = licenses.asl20;
    platforms = platforms.linux ++ platforms.darwin;
  };
}
