# crudeVanillaTSViteInitializer

Is a simple easy to break script in golang to
initialize vanilla-ts vite project with the
following:

- eslint
- prettier
- git
- a crude implementation of a file walker to
  imitate a static file server using
  directories as page names if they have an
  index.html file in them
- a hello world equivalent root index.html
  using the app's working directory as the
  name of the project, root html title and h1 text
  content
- also added some styles with colors  that's 
  sort of inspired by the [ dracula theme ][1]

I really wanted to craft a better
implementation of this as a portfolio project
that demonstrates my knowledge in viper and
cobra-cli but it seems I've been needing to
make this a lot recently and its just getting
in my nerves how I can't just set this up in
webstorm 

## Limitations

- biggest one right now is the project name it 
  can only do lowercase letters or the create 
  vite command will ask for a prompt which I 
  do not want to deal with yet. would be nice 
  if the app can do camelcase dir names as the 
  project name or maybe even spaces?? nah.
- I'd also want to put an flag argument where 
  the html h1 and title can be changed or 
  prefixed or suffixed by something like if i 
  the project name was 'bar' I could do a 
  --pref="foo" and then the title and h1 would 
  be 'foo bar' or something like that

## Installation 

1. compile the binary after cloning the repo
2. copy or move the referenceapp dir to 
   somewhere suitable

## Usage

no prefix
```shell
env VITEINIT_REFERENCE_PATH=/path/to/referenceapp  ./path/to/binary
```

with prefix
```shell
env VITEINIT_REFERENCE_PATH=/path/to/referenceapp  ./path/to/binary -pref <string>
```

On my computer I sym-linked the 
`path/to/the/repo/cmd/viteinit/refernceapp` to a 
`~/.config` directory named 'viteinit' I just 
thought it was a good place for it, compiled the 
binary to `~/bin` as `viteinit` and exported 
the environment variable in my `~/.zshrc` so i 
can just go:

```shell
viteinit
```

[1]: https://draculatheme.com/contribute