#include <stdio.h>
#include <stdarg.h>

extern char* gcmalloc(long size);



// the _GN2io5print function wrapper.
void _GN2io5print(char *fmt, ...) {
	va_list args;
	va_start(args, fmt);
	vprintf(fmt, args);
	va_end(args);
}


// the _GN2io6format function wrapper.
char* _GN2io6format(char *fmt, ...) {
	va_list checkArgs;
	va_start(checkArgs, fmt);
	long size = vsnprintf(NULL, 0, fmt, checkArgs) + 1;
	va_end(checkArgs);
	
	// Allocate memory for the string
	char* buffer = gcmalloc(size);
	
	// Reparse the args... There is no way around this, sadly
	va_list args;
	va_start(args, fmt);
	vsnprintf(buffer, size, fmt, args) + 1;
	va_end(args);
	return buffer;
}

// Since geode doesn't have any way of using c structs at the time being,
// I represent FILE* as void* and just trust the user (tm) and cast.
char* __openfile(char* path, char* mode) {
	FILE* f = fopen(path, mode);
	return (char*)f;
}

char __readchar(char* a) {
	FILE* f = (FILE*)a;
	return (char)fgetc(f);
}

int __fileeof(char* a) {
	return feof((FILE*)a);
}

int __filewritestring(char* a, char* data) {
	return fputs(data, (FILE*)a);
}