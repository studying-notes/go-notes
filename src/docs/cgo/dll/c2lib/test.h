#ifndef TEST_H
#define TEST_H

#ifdef TEST_DLL_EXPORT
#define TEST_API __declspec(dllexport)
#else
#define TEST_API __declspec(dllimport)
#endif

TEST_API void greet(char *name);

TEST_API char *name();

#endif