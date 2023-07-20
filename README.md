# cardboard
A '📦-inspired' Interpreted Language Implemented In Go

# Example Script

``cardboard`` follows a simple ``C``-like syntax structure, without being a strictly typed language. 

An example (proposed) code script can be seen below:

```
< Function Declaration >
put add = box(a, b) {
    put y = a + b;
    unbox y;
}

< Variable Declaration >
put x = 10;
put y = 20 - 15;

< Printing To Output >
show(add(x, y));
```

The syntax of cardboard is liable to change as I develop it, but the design focus for ``cardboard`` will always be simplicity and ease of use. 

# Development Plans
Developing programming languages is an interest I wish to experiment with, and therefore ``cardboard`` will always be in development! 

The following is a subset of the ideas I wish to implement in the language:

- [ ] Variable Declarations
- [ ] Arithmetic Operations
- [ ] Function Declarations
- [ ] Printing Functionality
- [ ] Comments
- [ ] Arrays
- [ ] Constants
- [ ] Structs / Classes

# Contribution
It would be really cool if you could fork this Repo and work on new features! I'll be happy to merge them right away 😀

# Resources
I'm working through Thorsten Ball's amazing book called 'Writing An Interpreter In Go'. Really recommend it 😀