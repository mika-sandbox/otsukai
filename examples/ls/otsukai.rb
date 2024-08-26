# default target host with user (ubuntu@yuuka.natsuneko.net)
set remote: { host: "hifumi.natsuneko.net", user: "natsuneko" }
set timeout: 10

task :deploy do
  run remote: "ls -al", stdout: true
end

hook after: :deploy do
  if task_success
    echo "Deploy Successful"
  end
end