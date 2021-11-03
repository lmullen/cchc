#!/usr/bin/env Rscript

# Given a set of input texts, identify quotations in them

suppressPackageStartupMessages(library(Matrix))
suppressPackageStartupMessages(library(broom))
suppressPackageStartupMessages(library(dplyr))
suppressPackageStartupMessages(library(fs))
suppressPackageStartupMessages(library(futile.logger))
suppressPackageStartupMessages(library(optparse))
suppressPackageStartupMessages(library(parsnip))
suppressPackageStartupMessages(library(readr))
suppressPackageStartupMessages(library(recipes))
suppressPackageStartupMessages(library(text2vec))
suppressPackageStartupMessages(library(tokenizers))
suppressPackageStartupMessages(library(stringr))

parser <- OptionParser(
  description = "Identify biblical quotations in a batch of texts.",
  usage = "Usage: %prog [options] BATCH --out=OUTPUT",
  epilogue = paste(
    "Input and output files are assumed to be stored as .csv files.",
    "Bible vectorizer and DTM should be a .rda file."
    )) %>%
  add_option(c("-o", "--out"),
             action = "store", type = "character", default = NULL,
             help = "Path to the output file.") %>%
  add_option(c("-b", "--bible"),
             action = "store", type = "character", default = "./bible-payload.rda",
             help = "Path to the Bible vectorizer and document-term model.") %>%
  add_option(c("-m", "--model"),
             action = "store", type = "character", default = "./prediction-payload.rda",
             help = "Path to the prediction model.") %>%
  add_option(c("--tokens"),
             action = "store", type = "integer", default = 2,
             help = "Minimum number of matching tokens (default: 2).") %>%
  add_option(c("--tfidf"),
             action = "store", type = "double", default = 1.0,
             help = "Minimum TF-IDF score to keep a potential match (default: 1.0).") %>%
  add_option(c("-v", "--verbose"),
             action = "store", type = "integer", default = 1,
             help = "Verbosity: 0 = errors and warnings; 1 = information; 2 = debugging.") %>%
  add_option(c("-p", "--potential"),
             action = "store_true", default = FALSE,
             help = "Save the potential matches to a CSV file")
if (!interactive()) {
  # Command line usage
  args <- parse_args(parser, positional_arguments = 1)
} else {
  # For testing
  flog.warn("Using the testing command line arguments since session is interactive.")
  args <- parse_args(parser,
                     args = c("predictor/test/sermons.csv",
                              "--out=predictor/test/sermons-quotations-test.csv",
                              "--bible=predictor/bin/bible-payload.rda",
                              "--model=predictor/bin/prediction-payload.rda",
                              "--verbose=2",
                              "--tokens=2",
                              "--tfidf=1.0"),
                     positional_arguments = 1)
}

# Easier references to outputs
batch_path <- args$args[1]
batch_id <- batch_path %>% path_file() %>% path_ext_remove()
out_path <- args$options$out
bible_path <- args$options$bible
model_file <- args$options$model

# Check validity of inputs and set options
if (args$options$verbose == 0) {
  log_threshold <- flog.threshold(WARN)
} else if (args$options$verbose == 1) {
  log_threshold <- flog.threshold(INFO)
} else if (args$options$verbose == 2) {
  requireNamespace("pryr", quietly = TRUE)
  log_threshold <- flog.threshold(DEBUG)
}
if (!file_exists(batch_path)) {
  flog.fatal("Batch file %s does not exist", batch_path)
  quit(save = "no", status = 1)
}
if (is.null(out_path)) {
  flog.fatal("An output path must be specified.")
  quit(save = "no", status = 1)
}
if (is.null(bible_path) || !file_exists(bible_path)) {
  flog.fatal("Bible payload file %s does not exist or was not specified.", bible_path)
  quit(save = "no", status = 1)
}
if (is.null(model_file) || !file.exists(model_file)) {
  flog.fatal("Model payload file %s does not exist or was not specified.", model_file)
  quit(save = "no", status = 1)
}
if (!dir_exists(path_dir(out_path))) {
  flog.fatal("The output directory must exist.")
  quit(save = "no", status = 1)
}
if (file_exists(out_path)) {
  flog.debug("The output file already exists. It will be overwritten.")
}
if (args$options$tokens < 0 || args$options$tfidf <0) {
  flog.fatal("The number of tokens and TF-IDF score must be positive.")
  quit(save = "no", status = 1)
}

flog.info("Beginning processing: %s.", batch_id)

flog.debug("Loading the Bible payload.")
bible <- new.env()
load(bible_path, envir = bible)

flog.debug("Loading the prediction model payload.")
load(model_file)

flog.debug("Reading batch of texts: %s.", batch_path)
texts <- read_csv(batch_path,
                  col_types = "ccc",
                  col_names = c("job_id", "doc_id", "text"))
flog.debug("Number of texts: %s.", nrow(texts))

flog.debug("Creating n-gram tokens from the texts.")
texts <- texts %>%
  mutate(tokens_ngrams = bible$bible_tokenizer(text, type = "ngrams"))

# Don't store the text once we don't need it any longer
texts <- texts %>% select(-text)

flog.debug("Creating the document-term matrix for the batch.")
token_it <- itoken(texts$tokens_ngrams,
                   ids = texts$doc_id,
                   progressbar = FALSE, n_chunks = 20)
docs_dtm <- create_dtm(token_it, bible$bible_vectorizer)
texts <- texts %>% select(-tokens_ngrams) # Don't store the n-gram tokens any more

flog.debug("Getting the count of matching tokens.")
token_count_m <- tcrossprod(bible$bible_dtm, docs_dtm)
# tidy() is deprecated for sparse matrices, but suppress that warning
suppressWarnings(
token_count <- token_count_m %>%
  tidy() %>%
  rename(verse_id = row, doc_id = column, tokens = value)
)

flog.debug("Computing the TF-IDF matrix for the Bible DTM.")
tfidf = TfIdf$new()
bible$bible_tfidf <- tfidf$fit_transform(bible$bible_dtm)

flog.debug("Getting the TF-IDF scores.")
# tidy() is deprecated for sparse matrices, but suppress that warning
suppressWarnings(
tfidf_score <- tcrossprod(bible$bible_tfidf, docs_dtm) %>%
  tidy() %>%
  rename(verse_id = row, doc_id = column, tfidf = value)
)

flog.debug("Getting the proportion of the matched verses.")
proportion_m <- (1 / rowSums(bible$bible_dtm)) * token_count_m
# tidy() is deprecated for sparse matrices, but suppress that warning
suppressWarnings(
proportion <- proportion_m %>%
  tidy() %>%
  rename(verse_id = row, doc_id = column, proportion = value)
)

flog.debug("Creating the potential matches data frame.")
potential_matches <- token_count %>%
  left_join(tfidf_score, by = c("verse_id", "doc_id")) %>%
  left_join(proportion, by = c("verse_id", "doc_id")) %>%
  as_tibble()

if (args$options$potential) {
  potential_dir <- "/tmp/potential/"
  fs::dir_create(potential_dir)
  potential_file <- str_c(potential_dir, fs::path_file(fs::file_temp()), ".csv")
  flog.info("Saving the potential matches to ", potential_file)
  write_csv(potential_matches, potential_file)
}

pnum <- function(x) { prettyNum(x, big.mark = ",") }
n_potential <- nrow(potential_matches)
potential_matches <- potential_matches %>%
  filter(tokens >= args$options$tokens | tfidf >= args$options$tfidf)
n_keepers <- nrow(potential_matches)
prop_keepers <- n_keepers / n_potential
flog.debug("Kept %s potential matches out of %s total (%s%%).",
          pnum(n_keepers), pnum(n_potential), round(prop_keepers * 100, 1))

# If we didn't get any potential matches, then end early
if (n_keepers < 1) {
  flog.debug("Quitting early since there were no potential matches")
  file_create(out_path)
  quit(save = "no", status = 0) 
}

# Center and scale the measurements as we did the training data
measurements <- bake(data_recipe, new_data = potential_matches %>% select(-verse_id, -doc_id))

# Do the predictions
probs <- predict(model$fit, measurements, type = "response")
names(probs) <- NULL
potential_matches$probability <- probs
predictions <- potential_matches %>%
  select(verse_id, doc_id, probability)

# Keep only one version per verse/document, and give the KJV a slight boost
quotations <- predictions %>%
  filter(probability >= 0.57) %>%
  mutate(reference_id = str_remove(verse_id, " \\(.+\\)")) %>%
  mutate(prob_adj = if_else(str_detect(verse_id, "(KJV)"), probability + 0.05, probability)) %>%
  group_by(doc_id, reference_id) %>%
  filter(prob_adj == max(prob_adj)) %>%
  slice(1) %>%
  ungroup() %>%
  select(doc_id, reference_id, verse_id, probability)

quotations <- texts %>%
  left_join(quotations, by = "doc_id") %>%
  filter(!is.na(probability))

write_csv(quotations, out_path, col_names = FALSE)
flog.info("Successfully identified %s quotations.", pnum(nrow(quotations)))

