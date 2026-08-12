#ifndef GO_FASTTEXT_WRAPPER_H
#define GO_FASTTEXT_WRAPPER_H

#ifdef __cplusplus
extern "C" {
#endif

/* Returns 0 on success; on failure writes a NUL-terminated message into err. */
int ft_load_model(const char* path, char* err, int err_len);

void ft_free_model(void);

/*
 * Top-1 language prediction.
 * Writes language code without "__label__" prefix into lang, and probability
 * into *prob. Returns 0 on success.
 */
int ft_predict(
    const char* text,
    char* lang,
    int lang_len,
    float* prob,
    char* err,
    int err_len);

#ifdef __cplusplus
}
#endif

#endif /* GO_FASTTEXT_WRAPPER_H */
