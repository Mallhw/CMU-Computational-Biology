# This R script takes a distance matrix between genomes as input.
# It produces a PCoA plot of the distances between genomes.

# Import needed libraries. Please install these in R beforehand using install.packages("package_name")

library(ggcorrplot)
library(reshape)
library(stringr)
library(ggplot2)
library(ape) #ape library to compute PCoA of our matrix

# Now set working directory. This should be wherever the files are stored and is the only line that the user needs to edit.

# Use "Session" --> "Set Working Directory" --> "To Source File Location"

#PLOT 1: Generating a PCoA plot of the data
# Read in the file and process the table.
table <- read.csv(file="Matrices/JaccardBetaDiversityMatrix.csv")

#trim the first column out because it only contains names
table <- table[-c(1)]
table <- table[, -c(3454)] # trim out weird extra column at end of the matrix file

matrix <- as.matrix(table)

pcoa_data <- pcoa(matrix, correction="none", rn=NULL) #This step may take a minute or two
pcoa_vectors <- data.frame(pcoa_data$vectors)
# columns contains a vector for each point after PCoA tries to assign data points to vectors to preserve distances between points.

colnames(table)


# Now, plot the data
ggplot(pcoa_vectors, aes(x=Axis.1, y=Axis.2)) + geom_point()
ggsave("Plots/JaccardPCoA.png")


#PLOT 2: Using the matrix labelled using our variant classifier to make a colored PCoA plot
table2 <- read.csv(file="Matrices/JaccardBetaDiversityMatrixLabelled.csv")

#trim the first column out because it only contains names
table2 <- table2[-c(1)]
table2 <- table2[, -c(3454)] # trim out weird extra column at end of the matrix file

matrix2 <- as.matrix(table2)

pcoa_data <- pcoa(matrix2, correction="none", rn=NULL)
pcoa_vectors <- data.frame(pcoa_data$vectors)
# columns contains a vector for each point after PCoA tries to assign data points to vectors to preserve distances between points.

colnames(table2)


Variant <- sub(".*_.*_.*_.*_", "", colnames(table2))
cbind(pcoa_vectors, Variant) # adding column

# Now, plot the data, colored by variant
ggplot(pcoa_vectors, aes(x=Axis.1, y=Axis.2, color=Variant)) + geom_point()
ggsave("Plots/JaccardPCoALabelled.png")

#PLOT 3: Now, sort the matrix by variant and create a heatmap
table3 <- read.csv(file="Matrices/JaccardBetaDiversityMatrixLabelled.csv")
table3$X <- sub(".*_.*_.*_.*_", "", table3$X )
table3 <- table3[order(table3$X ),]
table3  <-table3[-c(1)]
colnames(table3)<- Variant
table3 <- table3[,order(colnames(table3))]
matrix <- as.matrix(table3)
# the following code is just all the technical stuff to build a heatmap out of the distance matrix.
co=melt(matrix)
head(co)
ggplot(co, aes(X1, X2)) + # x and y axes => Var1 and Var2
  geom_tile(aes(fill = value)) + # background colours are mapped according to the value column
  scale_fill_gradient2(low = "#6D9EC1",
                       mid = "white",
                       high = "#E46726",
                       midpoint = 0.1, limit= c(0,0.2)) +
  theme(panel.grid.major.x=element_blank(), #no gridlines
        panel.grid.minor.x=element_blank(),
        panel.grid.major.y=element_blank(),
        panel.grid.minor.y=element_blank(),
        panel.background=element_rect(fill="white"), # background=white
        axis.text.x = element_text(angle=90, hjust = 1,vjust=1,size = 8,face = "bold"),
        plot.title = element_text(size=20,face="bold"),
        axis.text.y = element_text(size = 8,face = "bold")) +
  ggtitle("Jaccard Heat Map") +
  theme(legend.title=element_text(face="bold", size=14)) +
  scale_x_discrete(name="") +
  scale_y_discrete(name="")
ggsave("Plots/JaccardHeatMap.png")
