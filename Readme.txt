To run the program:
'./Main -start=1 -end=123345 -timeout=22 -threads=6'

To run the tests
'go test *.go -v'

Q/A:
1) What would you do to make the program better had you more time?
-   Create a way get task completed per thread before timeout.
-   Implement a way to use keyboard interrupt along with the timeout

2) How would you prove that your solution is correct?
-   Verify same data set againts different thread numbers
-   Use the same thread number against a variety of data sets to see efficiency
-   Use robust unit tests to verify the calculations 
