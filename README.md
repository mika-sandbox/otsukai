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
  if changed(path: "/path/to/docker-compose.yml", from: :last_commit)
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

## Functions

### `set` (local func)

Set a variable value. Example:

```ruby
set target: { host: "yuuka.natsuneko.net", user: "ubuntu" }
set default: :deploy
set app_root: "/usr/local/"
```

### `task` (local func)

Define a task with name. Example:

```ruby
# define `deploy` task
task :deploy do 
  # ...
end

# define `rollback` task
task :rollback do 
  # ...
end
```

### `changed` (local func)

Check the specified path has changed from specified refs.

```ruby
changed(path: "/path/to/file", commit_from: :last_commit, commit_to: :head) # returns bool
```

the `commit_from` supports the following args:

- `:last_commit`  : the specified file is changed in last commit
- `:fetch_commit` : the specified file is changed in remote fetched commit (ref: [`git-rev-parse#FETCH_HEAD`](https://git-scm.com/docs/git-rev-parse#Documentation/git-rev-parse.txt-codeFETCHHEADcode))  
- `:before_merge` : the specified file is changed in before merged commit (ref: [`git-rev-parse#ORIG_HEAD`](https://git-scm.com/docs/git-rev-parse#Documentation/git-rev-parse.txt-codeORIGHEADcode))
- `:after_merge`  : the specified file is changed in merged commit(s) (ref: [`git-rev-parse#MERGE_HEAD`](https://git-scm.com/docs/git-rev-parse#Documentation/git-rev-parse.txt-codeMERGEHEADcode))

the `commit_to` supports the following args:

- Not Yet Implemented

### `copy` (local / remote func)

Copy file/directory between from local/remote to remote/local.

```ruby
copy(to: :remote, local: "/path/to/file", remote: "/path/to/dest")
```

the `to` supports the following args:

- `:remote` : copy from local to remote
- `:local`  : copy from remote to local

the local and remote is path of the file or directory.
if the directory is specified, copy recursively.

### `run` (local / remote func)

Run commands in local/remote.

```ruby
run(remote: "echo 'Hello, World'")
run(local:  "echo 'Hello, World'")
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
