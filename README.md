Wordle
---
A project to learn about golang channels and DynamoDB with a helpful side effect of being useful for solving a daily Wordle.
Eventually there will be a REPL for generating guesses based on the Words stored in DynamoDB. 

Dictionaries / Submodules
---
The project relies on several sources to compile a list of candidate words for the Wordle.
* [Github Wordset Dictionary](https://github.com/wordset/wordset-dictionary.git)
* [Github English Words Dictionary](https://github.com/dwyl/english-words.git)

To run the DynamoDB population scripts both repositories (submodules) need to be cloned
```shell
git submodule update --recursive --remote
```