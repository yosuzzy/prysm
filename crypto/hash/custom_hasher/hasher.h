#ifndef __CUSTOM_HASHER__
#define __CUSTOM_HASHER__

#include <stdint.h>
extern void sha256_4_avx(unsigned char* output, const unsigned char* input, uint64_t blocks);
extern void sha256_8_avx2(unsigned char* output, const unsigned char* input, uint64_t blocks);
extern void sha256_shani(unsigned char* output, const unsigned char* input, uint64_t blocks);
extern void sha256_1_avx(unsigned char* output, const unsigned char* input);
void sha256_armv8_neon_x1(unsigned char* output, const unsigned char* input, uint64_t blocks);
void sha256_armv8_neon_x4(unsigned char* output, const unsigned char* input, uint64_t blocks);
#endif
