# Exporter

For Importer to be efficient, you can provide "Exporter" marker within the file.

```markdown
<!-- == export: simple_instruction / begin == -->

Some instruction
This will be exported as "simple_instruction".

<!-- == export: simple_instruction / end == -->
```

With the above, you can simplify Importer annotation. You can find more in `using-exporter-before.md` and `using-exporter-after.md`

## Exporter data

<!-- == export: test_exporter / begin == -->

✨✨✨✨✨✨✨✨
✨✨✨✨✨✨✨✨
✨✨✨✨✨✨✨✨
🚀🚀🚀🚀🚀🚀🚀🚀
✨✨✨✨✨✨✨✨
✨✨✨✨✨✨✨✨
✨✨✨✨✨✨✨✨

<!-- == export: test_exporter / end == -->
