#define TEST_DLL_EXPORT
#include "test.h"
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

void greet(char *name)
{
    printf("Hello, %s!\n", name);
}

char *name()
{
    char buf[] = "Gopher";
    char *n = malloc(strlen(buf) + 1);
    strcpy(n, buf);
    n[strlen(buf)] = '\0';
    return n;
}

// gcc test.c -shared -o test.dll