## 🚀 Examples with YAML

### Preview

<!-- == export: preview-desc / begin == -->

`importer preview` command gives you a quick look at how the file may change when `importer update` and `importer purge` are run against the provided file. This is meant to be useful for testing and debugging.

<!-- == export: preview-desc / end == -->

```console
$ importer preview ./testdata/yaml/demo-before.yaml
---------------------------------------
Content Before:
1:      title: Demo of YAML Importer
2:      # == imptr: description / begin from: ./snippet-description.yaml#[for-demo] ==
3:      dummy: This will be replaced
4:      # == imptr: description / end ==
---------------------------------------

---------------------------------------
Content After Purged:
1:      title: Demo of YAML Importer
2:      # == imptr: description / begin from: ./snippet-description.yaml#[for-demo] ==
3:      # == imptr: description / end ==
---------------------------------------

---------------------------------------
Content After Processed:
1:      title: Demo of YAML Importer
2:      # == imptr: description / begin from: ./snippet-description.yaml#[for-demo] ==
3:      description: |
4:        This demonstrates how importing YAML snippet is made possible, without
5:        changing YAML handling at all.
6:      # == imptr: description / end ==
---------------------------------------

You can replace the file content with either of the commands below:

  importer update ./testdata/yaml/demo-before.yaml     Replace the file content with the Importer processed file.
  importer purge ./testdata/yaml/demo-before.yaml      Replace the file content by removing all data between marker pairs.

You can find more with 'importer help'
```

#### Preview in Action

<!-- == export: preview / begin == -->

https://user-images.githubusercontent.com/23435099/136710122-a0901daf-971b-40bf-9ec0-15f39f0e7958.mp4

<!-- == export: preview / end == -->

### Generate

<!-- == export: generate-desc / begin == -->

`importer generate` imports based on Importer Markers in the given file, and write the result to stdout or file. This can be used for debugging, or create a template file with Importer Markers but keep the file purely for Importer Markers.

<!-- == export: generate-desc / end == -->

#### Generate in Action

<!-- == export: generate / begin == -->

https://user-images.githubusercontent.com/23435099/136703617-9f11e97b-3a87-449a-a5a1-698139392465.mp4

<!-- == export: generate / end == -->

### Purge

<!-- == export: purge-desc / begin == -->

`importer purge` removes any lines between Importer Markers in the given file, and update the file in place. The same operation is executed for `importer update` before importing all the lines, but this "purge" is sometimes useful to see the file without extra data imported.

<!-- == export: purge-desc / end == -->

#### Purge in Action

<!-- == export: purge / begin == -->

https://user-images.githubusercontent.com/23435099/136700548-6c11e599-1cda-4c30-bcfd-840a2c075e37.mp4

<!-- == export: purge / end == -->
