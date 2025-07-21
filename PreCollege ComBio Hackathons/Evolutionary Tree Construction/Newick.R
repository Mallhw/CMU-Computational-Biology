# Ensure the 'ape' package is loaded for phylogenetic tree operations
if (!requireNamespace("ape", quietly = TRUE)) {
  install.packages("ape")
  library(ape)
} else {
  library(ape)
}

library(readxl)

# Under the "Session" menu, set the working directory to the source file location.

# First, we write a function completing a simple demo of how to draw a Newick tree
NewickDemo <- function() {
  # Define the tree from a Newick format string
  myTree <- ape::read.tree(
    text = "(A:0.1,B:0.2,(C:0.3,D:0.4)E:0.5);"
  )
  
  # Plot the tree with tip labels
  plot(myTree, show.tip.label = TRUE, show.node.label = TRUE)
  
  # Add edge labels showing rounded branch lengths
  edgelabels(round(myTree$edge.length, 1), cex=1)  # Adjust cex for size of text
}

# calling our function
NewickDemo()

# Next, we write a function that will generate the Newick format of a Hemoglobin subunit alpha tree
HBATree <- function(filename = "Output/HBA1/HBA1.png", width = 4000, height = 3000) {
  
  # Start PNG device
  png(filename, width = width, height = height)
  
  # Read the tree from a pre-specified file
  tree <- read.tree("Output/HBA1/hba1.tre")
  
  # Plot the tree
  plot(tree, edge.color = "black", show.node.label = FALSE, show.tip.label = TRUE, cex=2)
  
  # Add edge labels showing rounded branch lengths
  edgelabels(round(tree$edge.length, 1), cex=1)  # Adjust cex for size of text
  
  # Close the PNG device
  dev.off()
}

# Calling our function
HBATree()


# We produce a plot of SARS-CoV-2 genomes from the UK, colored by year
COVIDByYear <- function(filename = "Output/UK-Genomes/COVIDByYear.png", width = 6000, height = 4500) {

  # Start PNG device
  png(filename, width = width, height = height)
  
  # Read the tree from a pre-specified file
  tree <- read.tree("Output/UK-Genomes/sars-cov-2.tre")

  # Define a color palette for different years
  color_palette <- c("2020" = "green", "2021" = "blue", "2022" = "purple", "2023" = "red", "2024" = "orange")
  
  # Extract years from the tip labels of the tree
  years <- substr(tree$tip.label, 1, 4)
  
  # Assign colors based on the extracted years
  node_colors <- color_palette[years]
  
  # Plot the tree
  plot(tree, edge.color = "black", show.node.label = TRUE, show.tip.label = FALSE)
  
  # Add colored node labels
  tiplabels(pch = 19, col = node_colors, cex = 2, srt = 45, adj = c(1, 0.5), font = 2)
  
  # Add a legend to the plot
  legend("topleft",               # Position of the legend
         legend = names(color_palette),  # The text in the legend, pulling from the names of the color_palette vector
         col = color_palette,      # The colors for the symbols in the legend
         pch = 19,                 # Type of point to use, same as in nodelabels
         title = "Year",           # Title of the legend
         cex = 8)                # Font size of the text in the legend
  
  # Close the PNG device
  dev.off()
}

#calling our function
COVIDByYear()

Process16SSequences <- function(
    year,
    width  = 4000,
    height = 3000,
    res    = NA_real_  # keep default device resolution unless you set explicitly
) {
  stopifnot(is.numeric(year), length(year) == 1, !is.na(year))
  
  # Build year-specific paths
  yrdir    <- sprintf("16S_%d", year)
  out_dir  <- file.path("Output", yrdir)
  tree_file <- file.path(out_dir, "colonized_bacteria.tre")
  out_png   <- file.path(out_dir, sprintf("16STree_%d.png", year))
  
  # Ensure output directory exists
  if (!dir.exists(out_dir)) dir.create(out_dir, recursive = TRUE, showWarnings = FALSE)
  
  # Read the tree
  tree <- read.tree(tree_file)
  
  # Plot
  if (is.na(res)) {
    png(out_png, width = width, height = height)
  } else {
    png(out_png, width = width, height = height, res = res)
  }
  plot(tree, edge.color = "black", show.node.label = FALSE, show.tip.label = TRUE, cex = 3)
  edgelabels(round(tree$edge.length, 1), cex = 1)
  dev.off()
  
  invisible(tree)
}

Process16SSequences(2024)
Process16SSequences(2025)

## =====================================================================
## 16S Tree Relabeling: Student Initials + Genus species
## =====================================================================
## Dependencies --------------------------------------------------------
if (!requireNamespace("ape", quietly = TRUE)) install.packages("ape")
if (!requireNamespace("readxl", quietly = TRUE)) install.packages("readxl")
suppressPackageStartupMessages({
  library(ape)
  library(readxl)
  library(stringr)
})

## ---------------------------------------------------------------------
## Utility: Safe column resolver (case-insensitive, punctuation-insensitive)
## ---------------------------------------------------------------------
.find_col <- function(df_names, candidates) {
  # normalize names: lower and strip non-alphanum
  norm <- function(x) gsub("[^a-z0-9]", "", tolower(x))
  nms_norm <- norm(df_names)
  cand_norm <- norm(candidates)
  for (c in seq_along(candidates)) {
    hit <- which(nms_norm == cand_norm[c])
    if (length(hit)) return(df_names[hit[1]])
  }
  # fallback: try partial contains
  for (c in seq_along(candidates)) {
    hit <- grep(cand_norm[c], nms_norm, fixed = TRUE)
    if (length(hit)) return(df_names[hit[1]])
  }
  return(NA_character_)
}

## ---------------------------------------------------------------------
## Utility: Auto-detect header row in taxonomy Excel
## ---------------------------------------------------------------------
## Some spreadsheets have 1–N title rows above the real header. This
## function scans the top 'peek_n' rows and finds the first row that
## contains *all* required column names (case-insensitive).
.ReadTaxonomySmart <- function(path,
                               required_cols = c("Sample Name", "Genus", "species"),
                               peek_n = 20) {
  # Read top peek_n rows with no col names
  peek <- suppressMessages(read_excel(path, col_names = FALSE, n_max = peek_n))
  found_row <- NA_integer_
  req_low <- tolower(required_cols)
  for (i in seq_len(nrow(peek))) {
    row_vals <- tolower(trimws(as.character(unlist(peek[i, ], use.names = FALSE))))
    if (all(req_low %in% row_vals)) {
      found_row <- i
      break
    }
  }
  if (is.na(found_row)) {
    # Could not auto-detect; fall back to default read (first row = header)
    warning("Could not auto-detect taxonomy header row; using first row as header.")
    return(read_excel(path))
  } else {
    skip <- found_row - 1
    return(read_excel(path, skip = skip))
  }
}

## ---------------------------------------------------------------------
## Utility: Canonicalize sample IDs for matching tree labels to taxonomy
## ---------------------------------------------------------------------
.clean_sample_id <- function(x) {
  x <- trimws(as.character(x))
  
  # Remove obvious run metadata, but KEEP 16S block so we can parse it
  x <- sub("[-_]SANGER.*$", "", x, ignore.case = TRUE)
  x <- sub("[-_]R[0-9]+.*$", "", x, ignore.case = TRUE)
  x <- sub("[-_]?seq[FR]?.*$", "", x, ignore.case = TRUE)  # -seqF / -seqR / etc.
  
  # If 16S-block exists, grab just the <SiteLetter><Digits><Initials> chunk after it
  # Ex: Kangas0427-16S-G21DL-TRS -> G21DL
  #     Kangas0427_16s-T14MY     -> T14MY
  has16 <- grepl("(?i)16s[-_]*[A-Za-z][0-9]+[A-Za-z]{1,3}", x, perl = TRUE)
  if (any(has16)) {
    x_sub <- sub("(?i).*16s[-_]*([A-Za-z][0-9]+[A-Za-z]{1,3}).*", "\\1", x[has16], perl = TRUE)
    x[has16] <- x_sub
  }
  
  # Drop everything non-alphanumeric, uppercase
  x <- gsub("[^A-Za-z0-9]", "", x)
  toupper(x)
}

## ---------------------------------------------------------------------
## Utility: Extract student initials from raw tree tip labels
## ---------------------------------------------------------------------
## Looks for final underscore + 1–3 letters BEFORE any sequencing suffix.
## e.g. "Kangas7753_A2_AD-SANGER_R1-16S-rRNA" -> "AD"
## If none found, returns "".
ExtractInitials <- function(labels) {
  # Remove downstream run metadata; keep 16S block
  base <- gsub("[-_]SANGER.*$", "", labels, ignore.case = TRUE)
  base <- gsub("[-_]R[0-9]+.*$", "", base,   ignore.case = TRUE)
  base <- gsub("[-_]?seq[FR]?.*$", "", base, ignore.case = TRUE)  # -seqF/-seqR
  
  # Normalize dashes to underscores (simplify regex)
  base_us <- gsub("-", "_", base)
  
  out <- character(length(base_us))
  out[] <- ""
  
  ## ---- Pattern A: 16S_<Site><Digits><Initials>[delim or end] (2025 style) ----
  mA <- regexpr("(?i)16s[_-]*[A-Za-z][0-9]+([A-Za-z]{1,3})(?:[_-].*)?$", base_us, perl = TRUE)
  hitA <- mA > 0
  if (any(hitA)) {
    out[hitA] <- toupper(sub("(?i).*16s[_-]*[A-Za-z][0-9]+([A-Za-z]{1,3})(?:[_-].*)?$",
                             "\\1", base_us[hitA], perl = TRUE))
  }
  
  ## ---- Pattern B: underscore + initials at end (2024 style) -------------------
  todoB <- which(!hitA)
  if (length(todoB)) {
    mB <- regexpr("_[A-Za-z]{1,3}$", base_us[todoB], perl = TRUE)
    hitB <- mB > 0
    if (any(hitB)) {
      out[todoB[hitB]] <- toupper(sub("^.*_([A-Za-z]{1,3})$", "\\1",
                                      base_us[todoB][hitB], perl = TRUE))
    }
  }
  
  ## ---- Pattern C: trailing letters after a digit (fallback: G23LH etc.) -------
  todoC <- which(out == "")
  if (length(todoC)) {
    mC <- regexpr("[0-9]([A-Za-z]{1,3})$", base_us[todoC], perl = TRUE)
    hitC <- mC > 0
    if (any(hitC)) {
      out[todoC[hitC]] <- toupper(sub(".*[0-9]([A-Za-z]{1,3})$", "\\1",
                                      base_us[todoC][hitC], perl = TRUE))
    }
  }
  
  out
}

## ---------------------------------------------------------------------
## Core relabeling function
## ---------------------------------------------------------------------
RelabelTreeWithGenusSpecies <- function(
    tree_file,
    taxonomy_file,
    sample_col_candidates  = c("Sample Name","Sample","SampleName","Sample_Name"),
    genus_col_candidates   = c("Genus","genus"),
    species_col_candidates = c("species","Species","specific epithet","Specific epithet"),
    fallback_species_text = "sp.",
    add_initials = TRUE,
    keep_unmatched_original = TRUE,
    verbose = TRUE,
    auto_header = TRUE
) {
  # --- read data ---
  tree <- read.tree(tree_file)
  tax  <- if (auto_header) {
    .ReadTaxonomySmart(taxonomy_file,
                       required_cols = c(sample_col_candidates[1],
                                         genus_col_candidates[1],
                                         species_col_candidates[1]))
  } else {
    read_excel(taxonomy_file)
  }
  
  # --- resolve column names dynamically ---
  s_col <- .find_col(names(tax), sample_col_candidates)
  g_col <- .find_col(names(tax), genus_col_candidates)
  p_col <- .find_col(names(tax), species_col_candidates)
  missing <- c()
  if (is.na(s_col)) missing <- c(missing, "Sample")
  if (is.na(g_col)) missing <- c(missing, "Genus")
  if (is.na(p_col)) missing <- c(missing, "species")
  if (length(missing)) {
    stop("Could not find required column(s) in taxonomy file: ",
         paste(missing, collapse=", "))
  }
  
  # --- extract vectors ---
  SampleName <- trimws(as.character(tax[[s_col]]))
  GenusX     <- trimws(as.character(tax[[g_col]]))
  SpeciesX   <- trimws(as.character(tax[[p_col]]))
  SpeciesX[is.na(SpeciesX) | SpeciesX == "" |
             tolower(SpeciesX) %in% c("na","nan")] <- fallback_species_text
  FullSpecies <- str_trim(paste(GenusX, SpeciesX))
  
  # --- canonical keys ---
  tax_canon <- .clean_sample_id(SampleName)
  
  # --- match tree tips ---
  tip_orig  <- tree$tip.label
  tip_canon <- .clean_sample_id(tip_orig)
  idx       <- match(tip_canon, tax_canon)
  
  # --- species labels for matched entries ---
  core_species <- FullSpecies[idx]
  
  # --- initials ---
  initials <- if (add_initials) ExtractInitials(tip_orig) else rep("", length(tip_orig))
  has_init <- nzchar(initials)
  
  matched_labels <- ifelse(
    is.na(core_species),
    NA_character_,
    ifelse(has_init, paste(initials, core_species), core_species)
  )
  
  final_labels <- matched_labels
  if (keep_unmatched_original) {
    final_labels[is.na(final_labels)] <- tip_orig[is.na(final_labels)]
  }
  
  tree$tip.label <- final_labels
  
  # --- reporting ---
  matched <- !is.na(core_species)
  n_matched   <- sum(matched)
  n_unmatched <- length(matched) - n_matched
  if (verbose) {
    message(sprintf("Matched %d / %d tips to taxonomy.", n_matched, length(matched)))
    if (n_unmatched > 0) {
      message(sprintf("Unmatched tips: %d. Use $mapping to inspect.", n_unmatched))
    }
  }
  
  # --- mapping data.frame ---
  mapping <- data.frame(
    TreeLabel_orig  = tip_orig,
    TreeLabel_canon = tip_canon,
    Initials        = initials,
    Matched         = matched,
    FullSpecies     = core_species,
    FinalLabel      = final_labels,
    Tax_SampleName  = SampleName[idx],
    stringsAsFactors = FALSE
  )
  
  return(list(
    tree = tree,
    mapping = mapping,
    n_matched = n_matched,
    n_unmatched = n_unmatched
  ))
}

## ---------------------------------------------------------------------
## Wrapper: relabel + plot to PNG
## ---------------------------------------------------------------------
Process16SSequencesSpeciesLabels <- function(
    year,
    width        = 4000,
    height       = 3000,
    res          = 300,
    add_initials = TRUE,
    auto_header  = TRUE,
    verbose      = TRUE,
    taxonomyFile = NULL   # <-- NEW: user can pass explicit path
) {
  stopifnot(is.numeric(year), length(year) == 1, !is.na(year))
  
  # Year-specific dirs/files
  yrdir    <- sprintf("16S_%d", year)
  out_dir  <- file.path("Output", yrdir)
  data_dir <- file.path("Data",   yrdir)
  
  # tree produced by your Go pipeline
  tree_file <- file.path(out_dir, "colonized_bacteria.tre")
  
  # ------------------------------------------------------------------
  # Resolve taxonomy file
  # ------------------------------------------------------------------
  if (is.null(taxonomyFile)) {
    # Look in data_dir for anything *Taxonomy*.xls[x]
    cand <- list.files(
      data_dir,
      pattern = "[Tt]axonomy\\.(xlsx|xls)$",
      full.names = TRUE
    )
    if (length(cand) == 0) {
      stop("No taxonomy file found for year ", year,
           " in ", data_dir,
           ". Expected something like '*Taxonomy.xlsx' or '.xls'.")
    }
    if (length(cand) > 1 && verbose) {
      message("Multiple taxonomy files found in ", data_dir,
              "; using first: ", basename(cand[1]),
              " (pass taxonomyFile= to override).")
    }
    taxonomyFile <- cand[1]
  } else {
    # If user passed a path that is not absolute, resolve relative to data_dir
    if (!file.exists(taxonomyFile)) {
      candidate <- file.path(data_dir, taxonomyFile)
      if (file.exists(candidate)) {
        taxonomyFile <- candidate
      } else {
        stop("Taxonomy file not found: ", taxonomyFile,
             " (also looked in ", candidate, ").")
      }
    }
  }
  
  # Name output PNG with year
  out_png <- file.path(out_dir, sprintf("colonized_bacteria_species_tree_%d.png", year))
  
  # Ensure output directory exists
  if (!dir.exists(out_dir)) dir.create(out_dir, recursive = TRUE, showWarnings = FALSE)
  
  # Relabel + plot
  rel <- RelabelTreeWithGenusSpecies(
    tree_file     = tree_file,
    taxonomy_file = taxonomyFile,
    add_initials  = add_initials,
    auto_header   = auto_header,
    verbose       = verbose
  )
  
  png(out_png, width = width, height = height, res = res)
  plot(rel$tree, edge.color = "black", show.node.label = FALSE,
       show.tip.label = TRUE, cex = 0.7)
  dev.off()
  
  invisible(rel)
}

 
Process16SSequencesSpeciesLabels(2024)
Process16SSequencesSpeciesLabels(2025)

