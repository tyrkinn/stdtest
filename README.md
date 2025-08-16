# `stdtest` Simplest possible CLI-testing CLI

Now it's under active development and not production ready yet

## Usage

1. Create `.stdtest` file in your project root

For example you have `add` cli to sum two numbers and output result in stdout.
Here's how your `.stdtest` may look like.

(Formatting like that is not required btw)

```text
add 1  2  ->  3 [Simple addition]
add 1  0  ->  1 [Addition with zero]
add 1 -1  ->  0 [Addition with negative number]
```

2. Just run stdtest CLI in this directory

! Note, that cli you gonna test should be in PATH or in directory you're running from !

You'll get results for passed and failed tests in format:

```text
Test Results for `add`: 
PASSED: 2, FAILED: 1

FAILED "Addition with negative number":
    Call         : add 1 -1
    Expected     : 0
    Got          : 1
```

## Roadmap

- [ ] Implement Property-Based Tests feature

    Possible implementation:

    ```text
    x := prop(int)
    y := prop(int)

    add x  y -> add y x [Commutativity]
    add x  0 -> x       [Additivity]
    add x -x -> 0       [Opposites addition]
    ```

    And then run it with 

    `stdtest --prop=100` where `prop` flag indicates property based testing usage and takes number of random tests executed

- [ ] Specify executable relative or global location 

    Possible implementation:

    ```text
    #cli add = "./build/debug/version/add"

    add 1  2  ->  3 [Simple addition]
    add 1  0  ->  1 [Addition with zero]
    add 1 -1  ->  0 [Addition with negative number]
    ```

- [ ] Add exit code and stderr message testing possibilities

    Possible implementation (here we have divide cli):

    ```text
    divide 5 5  ->  1                                [Divide by itself]
    divide 5 0  =$= 1 <2> "ERROR: Division by zero"  [Divizion by zero]
    ```

    ! Not sure about this, but in my mind it's ok to use `=$=` as exit code checker and `<DESCRIPTOR>` as descriptor value checker.
    But maybe it should be special syntax for error case !
