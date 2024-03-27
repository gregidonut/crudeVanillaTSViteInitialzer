# crudVanillaTSViteInitializer


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
  using the apps first argument as the name of 
  the project, root html title and h1 text content

I really wanted to craft a better 
implementation of this as a portfolio project 
that demonstrates my knowledge in viper and 
cobra-cli but it seems I've been needing to 
make this a lot recently and its just getting 
in my nerves how I can't just set this up in 
webstorm 