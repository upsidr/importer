# Details of Importer

## Background

Importer aims to achieve one simple goal: break long files into small files.

In most of programming languages, you can define variables and reuse code, and some even try to emphasise on DRY - "Don't Repeat Yourself". Importer's goal is similar in idea, but it tries to be as "dumb" as possible. It doesn't guarantee the code reuse is done all the time, there is essentially no compilation, and nothing complicated. This means, if Importer runs only once, and the file content gets updated, Importer doesn't care. Importer is a simple enough tool, which can be embedded as a part of some other automation. If you want to ensure the Importer's `generate` command is always run against a given file, you should be putting some CI job using Importer - but Importer shouldn't know anything about CI itself.

Also, Importer's idea is not about DRY - instead, it's about "write once and reuse".

For example, in any programming language, if you update a single variable, it could have huge impact as the variable could be referred by many other codes. You will need some smart tools such as IDE or CI jobs to fully understand the impact of the change. You may miss some change that wasn't intended to take place because of the code reuse.

Also, if you were to simply copy and paste in multiple places you want to use the same value, the impact would be very easy to see when reviewing the change, but when it comes to updating them, you would have to update many places.

That's why Importer was created. You can stick to simple format such as Markdown, but when it becomes harder to maintain due to the size and complexity, you can "Import" some content from other files, while keeping Markdown as a pure Markdown.

## Implementation Details

Importer is a simple regular expression file reader. It looks for special Importer Annotation comment, which has no meaning in that given language, but Importer can parse such comment and wires up other file content.

In language like Markdown and YAML, a single file is made to be as is, meaning you cannot import some other files. This makes them really simple and easy to get started, but when you try to do something more involved, it becomes more difficult to maintain very quickly.

Because Importer tries to be "dumb", it doesn't actually know much about the given file syntax. Importer looks for Importer Annotation comments, parses them, and generates the updated version of that file, with specified lines imported into it.

Because the goal of Importer is very simple, the implementation is based on simple regular expressions. It is not made to be performant, nor capable of handling complex scenarios. But it works for most cases, such as Markdown and YAML. Other file typse may benefit from this approach. If there is any other file types that could benefit from this, we will look to expand our support in the future.
