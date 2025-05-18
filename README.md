# Prunehild

A file pruning utility

Prunehild takes a list of files, and based on a given schedule, it will
prune this list based on their date. Different actions can be executed
for the pruned files, the most common action would probably be "delete".

The schedule will decide which ones to keep on which ones to delete.

The idea is to use this for example in a directory that is full of
backup files like

```
backup-sqldump-1.gz
backup-sqldump-2.gz
backup-sqldump-3.gz
backup-sqldump-4.gz
[...]
```

and remove some of them so the size of the directory stays
approximately constant. In this example, you could keep putting sql
backup files to the directory and use prunehild to prune them.

The names of the files do not matter, only the date is relevant.

## Schedules

In a schedule you can define how many files you want to keep per
interval. A schedule definition of:

```json
"myschedule" : [ {"Day":1}, {"Week":2}, {"Month":1}, {"Long":1} ],
```

would ensure one file that is at most a day old is kept, two files
at least a week old, one at least a month old and one will be kept
for eternity.

Let's demonstrate this with an example:

...

## Actions

Possible actions you can run on the pruned files:

### List

Print the names of the obsoleted files to stdout. Use this as a dry run
option or pipe the output to another command.

### Delete

Removes the files from disk. Use with care.

### Move

Move the files into another directory.

### Compress

Compress the files, retaining the modification date.