# Test Markdown

<!-- == imptr: some_importer / begin from: ./somefile#1~2 == -->

some data between an annotation pair, which gets purged.

This annotation for "another_importer" gets ignored as it is within another annotation pair.

<!-- == imptr: another_importer / begin from: ./another_file#1~2 == -->
<!-- == imptr: another_importer / end == -->

<!-- == imptr: some_importer / end == -->
