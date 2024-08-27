# default target host with user (ubuntu@yuuka.natsuneko.net)
set remote: { host: "hifumi.natsuneko.net", user: "natsuneko" }
set timeout: 10

task :deploy do
  run local: "invalid-command", stdout: true
end

hook before: :deploy do
  run local: "echo 'Running Before Task'", stdout: true
end

hook after: :deploy do
  run local: "echo 'Running After Task'", stdout: true

  if task_success
    # this sample does not run the following command
    run local: "echo 'Running When Task is Successful'", stdout: true
  end
end