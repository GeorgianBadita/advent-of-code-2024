# advent-of-code-2024

This time I'm using a little script that fetches the problem input for me, as well as
copies a template that I'm going to use for each problem. Check below if you'd like to use it as well

## Install

1. Make sure you have [rust and cargo](https://www.rust-lang.org/tools/install) installed
2. `cd aoc-problem-fetch && cargo build --release`
3. `mv ./target/release/aoc-problem-fetch /usr/local/bin/aoc`
4. `aoc --help` should display the following:

```
Downloads the input for a particular day/year of advent of code.

It will create a folder named day_{d}_{year}, with a file named input.txt.

If a template dir path is specified, it will copy all files within the folder to day_{d}_{year}.

Usage: aoc [OPTIONS] -d <DAY> -y <YEAR>

Options:
  -d <DAY>
          The day of the problem

  -y <YEAR>
          The year of the problem

  -s <SEESION_COOKIE>
          Session cookie to be used for downloading the input. If this argument is provided, the cookie
          will be saved in a .session.lock file and future usages of this tool will look for that file.
          If the .session.lock file is missing, this argument MUST be provided

  -t <TEMPLATE_PATH>
          Path to a template folder containing the boilerplate code that's being used for all problems

  -h, --help
          Print help (see a summary with '-h')
```
