#!/usr/bin/env Rscript
# DataViz.R — downstream visualisation for metagenomics pipeline
# Colours points by SEASON, shapes by LOCATION derived from sample names
# Sample naming convention: Season_Location_Number (e.g. Summer_Braddock_1)

suppressPackageStartupMessages({
  library(ggplot2)
  library(dplyr)
  library(tidyr)
  library(stringr)
  library(tibble)
  library(readr)
  library(glue)
  library(viridis)
  library(ape)
})

args <- commandArgs(trailingOnly = TRUE)
year <- if (length(args) > 0) args[1] else "2019"

if (!dir.exists("Plots")) dir.create("Plots")

alpha_plot <- function(df, value_col, metric) {
  df %>%
    mutate(Location = str_split_fixed(Sample, "_", 3)[,2]) %>%   # 2nd token
    ggplot(aes(x = Location, y = .data[[value_col]], fill = Location)) +
    geom_boxplot() +
    labs(title = glue("{metric} by location"), x = NULL, y = metric) +
    theme_bw() +
    theme(legend.position = "none")
}

read_distance <- function(path) {
  df <- read_csv(path, show_col_types = FALSE)
  df <- column_to_rownames(df, var = names(df)[1])
  df <- df %>% select(where(~ !all(is.na(.x))))
  samples <- intersect(rownames(df), colnames(df))
  df <- df[samples, samples, drop = FALSE]
  as.matrix(df)
}

plot_metric <- function(metric) {
  m <- read_distance(glue("Matrices/{metric}BetaDiversityMatrix_{year}.csv"))
  
  # ------ heat map ------
  m_long <- as_tibble(m, rownames = "row") %>% 
    pivot_longer(-row, names_to = "col", values_to = "dist")
  heat <- ggplot(m_long, aes(row, col, fill = dist)) +
    geom_tile() +
    coord_equal() +
    scale_fill_viridis(limits = c(0, 1), name = "Distance") +
    theme_minimal(base_size = 9) +
    theme(axis.text.x = element_text(angle = 90, vjust = 0.5, hjust = 1),
          axis.title = element_blank())
  ggsave(glue("Plots/{metric}HeatMap_{year}.png"), heat, width = 8, height = 6, dpi = 300)
  
  # ------ PCoA ------
  vec <- pcoa(as.dist(m))$vectors %>% as_tibble()
  sample_names <- colnames(m)
  parts <- str_split_fixed(sample_names, "_", 3)
  vec$Season   <- parts[,1]
  vec$Site     <- parts[,2]
  
  p <- ggplot(vec, aes(Axis.1, Axis.2, colour = Season, shape = Site)) +
    geom_point(size = 3, stroke = 0.9) +
    labs(title = glue("{metric} PCoA"), x = "PCoA 1", y = "PCoA 2") +
    theme_bw()
  
  ggsave(glue("Plots/{metric}PCoA_{year}.png"), p, width = 8, height = 6, dpi = 300)
}

# ---------- alpha diversity ----------
richness <- read_csv(glue("Matrices/RichnessMatrix_{year}.csv"), show_col_types = FALSE)
simpson  <- read_csv(glue("Matrices/SimpsonMatrix_{year}.csv"),  show_col_types = FALSE)

ggsave(glue("Plots/RichnessBoxPlots_{year}.png"),
       alpha_plot(richness, "Richness", "Richness"), width = 8, height = 6, dpi = 300)

ggsave(glue("Plots/SimpsonsBoxPlots_{year}.png"),
       alpha_plot(simpson, "SimpsonsIndex", "Simpson's Index"), width = 8, height = 6, dpi = 300)

# ---------- beta diversity ----------
for (metric in c("BrayCurtis", "Jaccard")) plot_metric(metric)

message("Plots written to Plots/ for year ", year)
