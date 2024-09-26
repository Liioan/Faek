## Settings steps:

-   default output (list)
    -   terminal
    -   file
-   output file name (text)
    -   if blank use default "faekOutput.ts"

## steps:

-   array name (text) - default "arr"
-   fields (text)
    -   input empty
        -   fields.len > 0 -> next step
        -   fields.len == 0 -> continue
    -   input invalid (no type, wrong type) -> continue
    -   input contains fields with options
        -   set input type to list
        -   options stored in `map[string]string`
        -   if custom option
            -   set input type to text
            -   date -> `new Date(input)`
            -   img -> `https://.../input` in form 100x100
