set default: :deploy

# default target host with user (ubuntu@yuuka.natsuneko.net)
set target: { host: "yuuka.natsuneko.net", user: "ubuntu" }

task :deploy do
  if changed("/path/to/docker-compose.yml", from: :last_commit)
    # run as sudo
    run_as :sudo do
      # run docker compose down on remote
      run remote: "docker compose down -f /remote/path/to/docker-compose.yml"
      # copy file/directory from local (/path/to/docker-compose.yml) to remote (/home/ubuntu/docker-compose.yml)
      copy to: :remote, local: "/path/to/docker-compose.yml", remote: "/home/ubuntu/docker-compose.yml"
      # run docker compose on remote
      run remote: "docker compose up -d -f /remote/path/to/docker-compose.yml"
    end
    #   # run as natsuneko
    #   run_as :natsuneko do
    #   end
  end
end

hook after: :deploy do
  if task_success()
    echo "Deploy Successful"
  end
end