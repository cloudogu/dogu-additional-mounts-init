# Usage

This application support currently following commands:

- copy

## Copy

Copy copies all files from given source paths to destination paths.
The files have to be regular files and already existing files in the target will be overwritten.

An error during execution does not stop the whole process and does not remove previous copied files!


### Example

`target/dogu-data-seeder copy -source=./cmd -target=./cmdCopy -source=./build -target=./buildCopy`

Where the source n will be copied to the target n.
