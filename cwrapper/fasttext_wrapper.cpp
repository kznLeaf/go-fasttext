#include "fasttext_wrapper.h"

#include <cstdio>
#include <cstring>
#include <memory>
#include <sstream>
#include <string>
#include <utility>
#include <vector>

#include "fasttext.h"

namespace {

std::unique_ptr<fasttext::FastText> g_model;

void set_error(char* err, int err_len, const char* msg) {
  if (err == nullptr || err_len <= 0) {
    return;
  }
  if (msg == nullptr) {
    msg = "unknown error";
  }
  std::snprintf(err, static_cast<size_t>(err_len), "%s", msg);
}

std::string strip_label_prefix(const std::string& label) {
  static const char kPrefix[] = "__label__";
  constexpr size_t kPrefixLen = sizeof(kPrefix) - 1;
  if (label.size() >= kPrefixLen &&
      label.compare(0, kPrefixLen, kPrefix) == 0) {
    return label.substr(kPrefixLen);
  }
  return label;
}

} // namespace

extern "C" {

int ft_load_model(const char* path, char* err, int err_len) {
  if (path == nullptr || path[0] == '\0') {
    set_error(err, err_len, "model path is empty");
    return 1;
  }

  try {
    auto model = std::make_unique<fasttext::FastText>();
    model->loadModel(std::string(path));
    g_model = std::move(model);
    if (err != nullptr && err_len > 0) {
      err[0] = '\0';
    }
    return 0;
  } catch (const std::exception& e) {
    set_error(err, err_len, e.what());
    return 1;
  } catch (...) {
    set_error(err, err_len, "failed to load model");
    return 1;
  }
}

void ft_free_model(void) {
  g_model.reset();
}

int ft_predict(
    const char* text,
    char* lang,
    int lang_len,
    float* prob,
    char* err,
    int err_len) {
  if (g_model == nullptr) {
    set_error(err, err_len, "model is not loaded");
    return 1;
  }
  if (text == nullptr) {
    set_error(err, err_len, "text is null");
    return 1;
  }
  if (lang == nullptr || lang_len <= 0 || prob == nullptr) {
    set_error(err, err_len, "invalid output buffers");
    return 1;
  }

  try {
    std::stringstream ioss;
    ioss << text;
    if (text[0] == '\0' || text[std::strlen(text) - 1] != '\n') {
      ioss << '\n';
    }

    std::vector<std::pair<fasttext::real, std::string>> predictions;
    if (!g_model->predictLine(ioss, predictions, /*k=*/1, /*threshold=*/0.0f) ||
        predictions.empty()) {
      set_error(err, err_len, "no prediction");
      return 1;
    }

    const std::string code = strip_label_prefix(predictions[0].second);
    std::snprintf(lang, static_cast<size_t>(lang_len), "%s", code.c_str());
    *prob = predictions[0].first;
    if (err != nullptr && err_len > 0) {
      err[0] = '\0';
    }
    return 0;
  } catch (const std::exception& e) {
    set_error(err, err_len, e.what());
    return 1;
  } catch (...) {
    set_error(err, err_len, "prediction failed");
    return 1;
  }
}

} // extern "C"
