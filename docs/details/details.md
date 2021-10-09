# Details of Importer

## Background

Importer aims to achieve one simple goal: break long files into small files.

In most of programming languages, you can define variables and reuse code, and some even try to emphasise on DRY - "Don't Repeat Yourself". Importer's goal is similar in idea, but it tries to be as "dumb" as possible. It doesn't guarantee the code reuse is done all the time, there is essentially no compilation, and no complicated logic at all. This means, if the file gets manually updated after Importer update, Importer doesn't care. It is a simple enough tool, which can be embedded as a part of some other automation. If you want to ensure the Importer's `update` command is always run against a given file, you should be putting some CI job using Importer - but Importer shouldn't know anything about CI itself.

Also, Importer's idea is not about DRY - instead, it's about "reuse for clarity".

For example, in any programming language, if you update a single variable, it could have a huge impact as the variable could be referred in many places. You will need some smart tools such as IDE, CI jobs, or test caess to fully understand the impact. You may miss some change that wasn't intended to take place because of the hidden dependencies.

However, you probably won't simply copy and paste the same code in multiple places. It may be very easy to see even for reviewer, but when it comes to updating them, you would have to update in many places. This is very error prone, as you can forget to update one file, and it is often difficult to spot those.

That's why Importer was created. You can stick to simple file format, but when it becomes harder to maintain due to the size and complexity, you can split up the file into multiple files and "Import" as necessary. This keeps the actual file as is, meaning Markdown will still be a pure Markdown.

## Implementation Details

Importer is a simple regular expression file reader. It looks for special Importer Marker comments, which have no meaning in that given language, but Importer can parse such a comment and wires up other file content. [You can find more about Importer Marker and other markers here.](/docs/details/markers.md)

In language like Markdown and YAML, a single file is made to be as is, meaning you cannot import other files. This makes them really simple and easy to get started, but when you try to do something a bit more involved, it becomes difficult to maintain very quickly.

Because Importer tries to be "dumb", it doesn't actually know much about the given file syntax. Importer looks for Importer Marker comments, parses them, and generates the updated version of that file, with specified lines imported into it.

Because the goal of Importer is very simple, the implementation is based on simple regular expressions. It is not made to be performant, nor capable of handling complex scenarios. But it works for most cases, such as Markdown and YAML. Other file typse may benefit from this approach. If there is any other file types that could benefit from this, we will look to expand our support in the future.
