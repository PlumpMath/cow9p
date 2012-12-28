# Cow9p - Copy on Write 9p filesystem

## Overview

Simple 9p server that adds a copy on write layer to an existing filesystem.

Invoke with:

	cow9p -s src -d dst

This will start a 9p server (by default listening on tcp port 5640). Any
non-disctructive file operations will try to use the filesystem dst, and will
use src as a fallback if the necessary files do not exist on dst. For
destructive operations (e.g. write) dst will again be used, first, but if the
files are not found, the will be copied to dst from src, rather than operating
on src directly.

## Status

Work in progress. Idea is well defined, but it doesn't work yet. Heck, it
doesn't even build. The interfaces are not considered stable either.

Use at your own risk.

## License

MIT license. see COPYING.

