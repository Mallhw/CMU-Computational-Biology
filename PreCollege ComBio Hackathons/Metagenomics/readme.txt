Copy the Metagenomics folder into your "go/src" directory.

Note that we provide helper functions in helperFunctions.go from the Chapter 1 exercises. You're welcome to use your own functions, too.

After you have passed all of these tests, navigate into the parent Metagenomics directory (using "cd .."). We will fill in main.go, after which you can call "go build" and then execute the resulting executable file.

The "Data" folder contains samples for our Three Rivers Metagenomics project. Our code will use the functions that you write and then write matrices to the "Matrices" folder.

The "DataViz.R" file can then be opened in R. You will need to set your working directory in the appropriate line in the R file, and make sure that you have the packages installed that it imports. You can then run this file, and it will generate plots based on your data in the "Matrices" folder that will be placed into the "Plots" directory. We will then visualize and interpret these plots!