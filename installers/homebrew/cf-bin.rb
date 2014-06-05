require 'formula'

class CfBin < Formula
  homepage 'https://github.com/starkandwayne/cf-cli'
  url 'https://github.com/starkandwayne/cf-cli/releases/download/v6.0.1/cf-darwin-amd64.tgz'
  sha1 '548a83996ade1fb4c4334e4ebcfd558434c01daf'

  def install
    system 'curl -O https://raw.github.com/starkandwayne/cf-cli/v6.0.1/LICENSE'
    bin.install 'cf-darwin-amd64' => 'cf'
    doc.install 'LICENSE'
  end

  test do
    system "#{bin}/cf"
  end
end
