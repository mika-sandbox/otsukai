# default target host with user (ubuntu@yuuka.natsuneko.net)
set remote: { host: "hifumi.natsuneko.net", user: "natsuneko" }

task :deploy do
  copy to: :remote, local: "./examples/copy-to-remote/otsukai.rb", remote: "/home/natsuneko/otsukai.rb"
end
