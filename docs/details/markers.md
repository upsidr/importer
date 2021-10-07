<!-- == importer-skip-update == -->

# Markers

## 🌗 Overview

<!-- == export: basic-marker / begin == -->

**Markers** are a simple comment with special syntax Importer understands. Importer is a simple CLI tool, and these markers are the key to make all the import and export to happen. There are several types of markers.

| Name                 | Description                                                    |
| -------------------- | -------------------------------------------------------------- |
| Importer Marker      | Main marker, used to import data from other file.              |
| Exporter Marker      | Supplemental marker used to define line range in target files. |
| Skip Importer Update | Special marker to suppress `importer update`.                  |
| Auto Generated Note  | Special marker for `importer generate` information.            |

<!-- == export: basic-marker / end == -->

## ☕️ Basic Syntax

The markers always follow the pattern of `== some-importer-marker-input ==`.

With YAML, this would be `# == some-importer-marker-input ==`.\
With Markdown, this would be `<!-- == some-importer-marker-input == -->`.

The main markers **Importer Markers** and **Exporter Markers** are both made up of pairs, `begin` and `end`.

Other markers are used to update Importer behaviours.

## 🖋 Importer Marker

![Importer Marker Syntax](/assets/images/importer-marker-syntax.png)

> NOTE: The above example is from [`/testdata/markdown/simple-before.md`](/testdata/markdown/simple-before.md).

### 1️⃣ Importer Marker Type

- Tell Importer to start the import setup.
- This can be represented with `importer`, `import`, `imptr` or `i`.
- Do not forget to add `:` at the end.

### 2️⃣ Importer Marker Name

- Any name of your choice, with no whitespace character.
- The same name cannot be used in a single file.

### 📃 Separate with `/`

- Add separator using `/`. The spaces around the `/` are required as of now.

### 3️⃣ Either `begin` or `end`

- Each Importer Marker must be a pair to operate.

### 4️⃣ Importer Marker Details

- `from: FILENAME#OPTION`: Define where to import from.
  - `FILENAME`: Specify the location of target file, which can be a URL or relative path from the source file.
  - `OPTION`: Define which line(s) to import.
    - `NUM1~NUM2`: Import line range from `NUM1` to `NUM2`.\
      Leaving `NUM1` empty means from the beginning of the file.\
      Leaving `NUM2` empty means to the end of the file.
    - `NUM1,NUM2`: Import each lines specified (e.g. `NUM1`, `NUM2`) one by one.
    - `[Exporter-Marker]`: Import lines based on Exporter Markers defined in the target file.
- `indent: [align|absolute NUM|extra NUM|keep]`: Update indentation for the imported data.
  - `align`: Align to the indentation of Importer Marker.
  - `absolute NUM` (e.g. `absolute 2`): Update indentation to `NUM` spaces. This ignores the original indentation from the imported data, but keeps the tree structure.
  - `extra NUM` (e.g. `extra 4`): Add extra indentation of `NUM` spaces.
  - `keep` (default): Keep the indentation from the imported data.

### Examples

#### With `/testdata/markdown/simple-before.md`

<details>
<summary>Preview Importer CLI in action</summary>

```console
$ importer preview ./testdata/markdown/simple-before.md
---------------------------------------
Content Before:
1:      # Simple Markdown Test
2:
3:      <!-- == imptr: lorem / begin from: ./snippet-lorem.md#5~12 == -->
4:
5:      Any content here will be removed by Importer.
6:
7:      <!-- == imptr: lorem / end == -->
8:
9:      Content after marker is left untouched.
---------------------------------------

---------------------------------------
Content After Purged:
1:      # Simple Markdown Test
2:
3:      <!-- == imptr: lorem / begin from: ./snippet-lorem.md#5~12 == -->
4:      <!-- == imptr: lorem / end == -->
5:
6:      Content after marker is left untouched.
---------------------------------------

---------------------------------------
Content After Processed:
1:      # Simple Markdown Test
2:
3:      <!-- == imptr: lorem / begin from: ./snippet-lorem.md#5~12 == -->
4:      "Lorem ipsum dolor sit amet,
5:      consectetur adipiscing elit,
6:      sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.
7:      Ut enim ad minim veniam,
8:      quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.
9:      Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur.
10:     Excepteur sint occaecat cupidatat non proident,
11:     sunt in culpa qui officia deserunt mollit anim id est laborum."
12:     <!-- == imptr: lorem / end == -->
13:
14:     Content after marker is left untouched.
---------------------------------------

You can replace the file content with either of the commands below:

  importer update ./testdata/markdown/simple-before.md     Replace the file content with the Importer processed file.
  importer purge ./testdata/markdown/simple-before.md      Replace the file content by removing all data between marker pairs.

You can find more with 'importer help'
```

</details>

#### With `/testdata/yaml/snippet-description.yaml` using Exporter Marker

<details>
<summary>Preview Importer CLI in action</summary>

```console
$ cat ./testdata/yaml/snippet-description.yaml
# == export: for-demo / begin ==
description: |
  This demonstrates how importing YAML snippet is made possible, without
  changing YAML handling at all.
# == export: for-demo / end ==

$ importer preview ./testdata/yaml/demo-before.yaml
---------------------------------------
Content Before:
1:      title: Demo of YAML Importer
2:      # == import: description / begin from: ./snippet-description.yaml#[for-demo] ==
3:      dummy: This will be replaced
4:      # == import: description / end ==
---------------------------------------

---------------------------------------
Content After Purged:
1:      title: Demo of YAML Importer
2:      # == import: description / begin from: ./snippet-description.yaml#[for-demo] ==
3:      # == import: description / end ==
---------------------------------------

---------------------------------------
Content After Processed:
1:      title: Demo of YAML Importer
2:      # == import: description / begin from: ./snippet-description.yaml#[for-demo] ==
3:      description: |
4:        This demonstrates how importing YAML snippet is made possible, without
5:        changing YAML handling at all.
6:      # == import: description / end ==
---------------------------------------

You can replace the file content with either of the commands below:

  importer update ./testdata/yaml/demo-before.yaml     Replace the file content with the Importer processed file.
  importer purge ./testdata/yaml/demo-before.yaml      Replace the file content by removing all data between marker pairs.

You can find more with 'importer help'
```

</details>

## 📝 Exporter Marker

![Exporter Marker Syntax](/assets/images/exporter-marker-syntax.png)

> NOTE: The above example is from [`/testdata/yaml/snippet-k8s-resource.yaml`](/testdata/yaml/snippet-k8s-resource.yaml).

### 1️⃣ Exporter Marker Type

- Tell Importer to start the export setup.
- This can be represented with `exporter`, `export`, `exptr` or `e`.
- Do not forget to add `:` at the end.

### 2️⃣ Importer Marker Name

- Any name of your choice, with no whitespace character.
- The same name cannot be used in a single file.

### 📃 Separate with `/`

- Add separator using `/`. The spaces around the `/` are required as of now.

### 3️⃣ Either `begin` or `end`

- Each Exporter Marker must be a pair to operate.

### Examples

#### Import with `/testdata/yaml/snippet-k8s-resource.yaml`

<details>
<summary>Preview Importer CLI in action</summary>

```console
$ importer preview ./testdata/yaml/k8s-color-svc-before.yaml
---------------------------------------
Content Before:
1:      ---
2:      apiVersion: v1
3:      kind: Service
4:      metadata:
5:        name: color-svc-only-green
6:        labels:
7:          app.kubernetes.io/name: color-svc-only-green
8:      spec:
9:        ports:
10:         - name: http
11:           port: 8800
12:           targetPort: 8800
13:       selector:
14:         app.kubernetes.io/name: color-svc-only-green
15:     ---
16:     apiVersion: apps/v1
17:     kind: Deployment
18:     metadata:
19:       name: color-svc-only-green
20:     spec:
21:       replicas: 1
22:       selector:
23:         matchLabels:
24:           app.kubernetes.io/name: color-svc-only-green
25:           app.kubernetes.io/version: v1
26:       template:
27:         metadata:
28:           labels:
29:             app.kubernetes.io/name: color-svc-only-green
30:             app.kubernetes.io/version: v1
31:         spec:
32:           serviceAccountName: color-svc
33:           containers:
34:             # == i: latest-color-svc / begin from: ./snippet-k8s-color-svc.yaml#[latest-svc] indent: align ==
35:             - image: docker.io/rytswd/color-svc:latest
36:               name: color-svc
37:               command:
38:                 - color-svc
39:               ports:
40:                 - containerPort: 8800
41:               # == i: latest-color-svc / end ==
42:
43:               env:
44:                 # == i: color-svc-default-envs / begin from: ./snippet-k8s-color-svc.yaml#[basic-envs] indent: align ==
45:                 # == i: color-svc-default-envs / end ==
46:                 - name: DISABLE_RED
47:                   value: "true"
48:                 - name: DISABLE_GREEN
49:                   value: "false" # The same as default
50:                 - name: DISABLE_BLUE
51:                   value: "true"
52:                 - name: DISABLE_YELLOW
53:                   value: "true"
54:
55:               # == i: resource-footprint / begin from: ./snippet-k8s-resource.yaml#[min-resource] indent: align ==
56:               data: |
57:                 this will be purged
58:               # == i: resource-footprint / end ==
---------------------------------------

---------------------------------------
Content After Purged:
1:      ---
2:      apiVersion: v1
3:      kind: Service
4:      metadata:
5:        name: color-svc-only-green
6:        labels:
7:          app.kubernetes.io/name: color-svc-only-green
8:      spec:
9:        ports:
10:         - name: http
11:           port: 8800
12:           targetPort: 8800
13:       selector:
14:         app.kubernetes.io/name: color-svc-only-green
15:     ---
16:     apiVersion: apps/v1
17:     kind: Deployment
18:     metadata:
19:       name: color-svc-only-green
20:     spec:
21:       replicas: 1
22:       selector:
23:         matchLabels:
24:           app.kubernetes.io/name: color-svc-only-green
25:           app.kubernetes.io/version: v1
26:       template:
27:         metadata:
28:           labels:
29:             app.kubernetes.io/name: color-svc-only-green
30:             app.kubernetes.io/version: v1
31:         spec:
32:           serviceAccountName: color-svc
33:           containers:
34:             # == i: latest-color-svc / begin from: ./snippet-k8s-color-svc.yaml#[latest-svc] indent: align ==
35:               # == i: latest-color-svc / end ==
36:
37:               env:
38:                 # == i: color-svc-default-envs / begin from: ./snippet-k8s-color-svc.yaml#[basic-envs] indent: align ==
39:                 # == i: color-svc-default-envs / end ==
40:                 - name: DISABLE_RED
41:                   value: "true"
42:                 - name: DISABLE_GREEN
43:                   value: "false" # The same as default
44:                 - name: DISABLE_BLUE
45:                   value: "true"
46:                 - name: DISABLE_YELLOW
47:                   value: "true"
48:
49:               # == i: resource-footprint / begin from: ./snippet-k8s-resource.yaml#[min-resource] indent: align ==
50:               # == i: resource-footprint / end ==
---------------------------------------

---------------------------------------
Content After Processed:
1:      ---
2:      apiVersion: v1
3:      kind: Service
4:      metadata:
5:        name: color-svc-only-green
6:        labels:
7:          app.kubernetes.io/name: color-svc-only-green
8:      spec:
9:        ports:
10:         - name: http
11:           port: 8800
12:           targetPort: 8800
13:       selector:
14:         app.kubernetes.io/name: color-svc-only-green
15:     ---
16:     apiVersion: apps/v1
17:     kind: Deployment
18:     metadata:
19:       name: color-svc-only-green
20:     spec:
21:       replicas: 1
22:       selector:
23:         matchLabels:
24:           app.kubernetes.io/name: color-svc-only-green
25:           app.kubernetes.io/version: v1
26:       template:
27:         metadata:
28:           labels:
29:             app.kubernetes.io/name: color-svc-only-green
30:             app.kubernetes.io/version: v1
31:         spec:
32:           serviceAccountName: color-svc
33:           containers:
34:             # == i: latest-color-svc / begin from: ./snippet-k8s-color-svc.yaml#[latest-svc] indent: align ==
35:             - image: docker.io/rytswd/color-svc:latest
36:               name: color-svc
37:               command:
38:                 - color-svc
39:               ports:
40:                 - containerPort: 8800
41:               # == i: latest-color-svc / end ==
42:
43:               env:
44:                 # == i: color-svc-default-envs / begin from: ./snippet-k8s-color-svc.yaml#[basic-envs] indent: align ==
45:                 - name: ENABLE_DELAY
46:                   value: "true"
47:                 - name: DELAY_DURATION_MILLISECOND
48:                   value: "500"
49:                 - name: ENABLE_CORS
50:                   value: "true"
51:                 # == i: color-svc-default-envs / end ==
52:                 - name: DISABLE_RED
53:                   value: "true"
54:                 - name: DISABLE_GREEN
55:                   value: "false" # The same as default
56:                 - name: DISABLE_BLUE
57:                   value: "true"
58:                 - name: DISABLE_YELLOW
59:                   value: "true"
60:
61:               # == i: resource-footprint / begin from: ./snippet-k8s-resource.yaml#[min-resource] indent: align ==
62:               resources:
63:                 requests:
64:                   cpu: 10m
65:                   memory: 10Mi
66:                 limits:
67:                   cpu: 30m
68:                   memory: 30Mi
69:               # == i: resource-footprint / end ==
---------------------------------------

You can replace the file content with either of the commands below:

  importer update ./testdata/yaml/k8s-color-svc-before.yaml     Replace the file content with the Importer processed file.
  importer purge ./testdata/yaml/k8s-color-svc-before.yaml      Replace the file content by removing all data between marker pairs.

You can find more with 'importer help'
```

</details>

## 🖍 Other Markers

There are some special markers that are used to update Importer behaviour.

### Skip Importer Update: `== importer-skip-update ==`

- Mark the file not to be updated by `importer update` command.

### Auto Generated Note: `== improter-generated-from: FILENAME ==`

- This is auto-generated by `importer generate FILENAME --out TARGET_FILE`.
- This tells how the file was generated by using `FILENAME` as input.

## 🔬 Other Details

### Importer File Option Details

| Name                       | Example         | Description                                                                                                                                                                                                                                                                                                 |
| -------------------------- | --------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Target Path                | `from: xyz.md`  | Defines where to import from. This is a relative path from the file containing the marker.<br /><br /> **Known Limitations**: Path cannot contain whitespace characters.                                                                                                                                    |
| Separator                  | `#`             | This is to separate Target Path and Target Detail. It can have as many preceding whispace characters.                                                                                                                                                                                                       |
| Target Detail - Line Range | `[1~33]`        | Imports only provided line ranges. You can omit before or after `~` to indicate the range starts from the beginning of the file, or ends at the end of the file.                                                                                                                                            |
| Target Detail - Line List  | `[1,2,5]`       | Imports only provided lines. The lines are comma separated, and you can also use line range in the same target detail. <br /><br /> **Known Limitations**: The order of lines is not persisted, and thus if you define `[3,2,1]`, you would actually see lines imported as line#1, line#2, and then line#3. |
| Target Detail - Marker     | `[some-marker]` | Searches for the matching Export Marker in the target file. More about Export Marke below. <br /><br /> **Known Limitations**: You can only provide single marker.                                                                                                                                          |
