# default target host with user (ubuntu@yuuka.natsuneko.net)
set remote: { host: "hifumi.natsuneko.net", user: "natsuneko" }

task :deploy do
  copy to: :local, local: "./examples/copy-to-local/k3s-install.sh", remote: "/home/natsuneko/k3s-install.sh"
end
