# Otsukai

Otsukai: Simple, and Easy to use deployment application for your servers.

## Installation

### From Source

```bash
$ git clone https://github.com/mika-f/otsukai
$ cd otsukai
$ go install
```

### From Binary

Download the binary from the [releases page](https://github.com/mika-f/otsukai/releases) and place it in your `$PATH`.

## Usage

At the first time, you need to create a configuration file.

> **NOTE**
> The deployment recipe has subset of Ruby, but not run as Ruby.

```ruby
# default target host, with user
set target: { host: "yuuka.natsuneko.net", user: "ubuntu" }

task :deploy do
  if changed("/path/to/docker-compose.yml", from: :last - commit)
    # run with sudo
    run_with :sudo do
      # run docker compose down on remote
      run remote: "docker compose down -f /remote/path/to/docker-compose.yml"

      # copy file/directory from local (/path/to/docker-compose.yml) to remote (/home/ubuntu/docker-compose.yml)
      copy to: :remote, local: "/path/to/docker-compose.yml", remote: "/home/ubuntu/docker-compose.yml"

      # run docker compose on remote
      run remote: "docker compose up -d -f /remote/path/to/docker-compose.yml"
    end
  end
end
```

Then, you can deploy your application by running the following command.

```bash
# check syntax before deploy
$ otsukai test --recipe examples/docker-compose/otsukai.rb

# deploy
$ otsukai deploy --recipe examples/docker-compose/otsukai.rb

# deploy (dry-run)
$ otsukai deploy --recipe examples/docker-compose/otsukai.rb --dry-run
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
