## 🚀 Examples with YAML

### Preview

<!-- == export: preview / begin == -->

`importer preview` command gives you a quick look at how the file may change when `importer update` and `importer purge` are run against the provided file. This is meant to be useful for testing and debugging.

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

<!-- == export: preview / end == -->

### Generate

<!-- == export: generate / begin == -->

![Generate in Action][generate-in-action]

[generate-in-action]: /assets/images/importer-generate-yaml-demo.gif "Generate in Action"

<!-- == export: generate / end == -->

### Purge

<!-- == export: purge / begin == -->

https://github.com/upsidr/importer/blob/adjust-readme/assets/images/importer-purge-yaml-demo.mp4?raw=true

<!-- == export: purge / end == -->
