# godepth
Count the maxdepth of go functions. It's helpful to see if a function
needs to be splitted into several smaller functions, for readability
purpose.

To install, run

    $ go get github.com/arthurgustin/godepth

# examples:

    $ godepth .
    $ godepth -over 3
    $ godepth -avg
