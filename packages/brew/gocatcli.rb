class Gocatcli < Formula
  desc "Command-line catalog tool for your offline data"
  homepage "https://github.com/deadc0de6/gocatcli"
  url "https://github.com/deadc0de6/gocatcli/archive/refs/tags/v1.0.3.tar.gz"
  sha256 "0011a3a65ab6e894b12de7e55cb33d413a590a3e7c422af30e83ba94d28476b4"
  license "GPL-3.0"

  depends_on "go" => :build

  def install
    system "make", "build"
    bin.install "bin/gocatcli"
  end

  test do
    assert_match "gocatcli version #{version}", shell_output("#{bin}/gocatcli --version")
  end
end