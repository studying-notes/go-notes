/*
base64.h - c header for a base64 decoding algorithm

This is part of the libb64 project, and has been placed in the public domain.
For details, see http://sourceforge.net/projects/libb64
*/

#ifndef BASE64_BASE64_H
#define BASE64_BASE64_H

char *base64_encode(const char *input);

char *base64_decode(const char *input);

#endif //BASE64_BASE64_H
