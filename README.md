### Personal Search Engine

- Many times we would prefer to have our own search engine for looking up any term or things we might have seen in the
  past. This project aims to facilitate that.
- It is a simple search engine backed by sqlite's FTS.
- Uses can be something like if we want to save some link or want to index all the links in the blog post and later on
  search in our db

### How to run

- `make init` to set up the database
- `make run` to start the server

### Things to Do

- What if unable to find elements / wrong elements. or in other words add validators
- Display the indexing as a progress bar as submitting the link blocks the call.
- show how many links found
- How to swap terminal with file for logging?
- How to build and containerize the app and deploy on AWS?