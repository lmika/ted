# ted - Terminal Table Editor

TED is an editor for editing tabular data, like CSVs.  It's like a spreadsheet without the calculation logic.  It's inspired by Vim and was developed to easily update CSV files without leaving the terminal.

This project is still in it's early phases, and on a bit of a slow burn in terms of development.  I tend to only work on it when I use it, which is not very often.  However, it's still in a form of maintenance so feel free to use it and post feedback.

## Downloading

For those with Go:

```
go get github.com/lmika/ted
```

## Usage

```
ted [FLAGS] FILE
```

Flags:

- `-c <codec>` the format that the file is in.  Either `csv` or `tsv` files are supported.  Default is `csv`

File can either be a new file, or an existing file.

TED is similar to Vim in that it is modal.  After opening a file, the editor starts off in view mode, which permits navigating around.

## Keyboard Keys

Moving around:

| Key        | Modifier   | `i` | `j` | `k` | `l` |
|:-----------|:-----------|:-----|:-----|:-----|:-----|
| Movement by single cell | (none) | Up | Left | Down | Right |
| Movement by a page | Shift | Up 25 rows | Left 15 cells | Down 25 rows | Right 15 cells |
| Movement to boundary | Ctrl | Top row | Leftmost cell | Bottom row | Rightmost cell |

You can also use the arrows to move by a single cell.

Editing:

| Key        | Action              |
|:-----------|:--------------------|
| `e`        | Edit cell value    |
| `r`        | Replace cell value  |
| `a`        | Insert row below cursor and edit value |
| `D`        | Delete current row |

Others:

| Key        | Action              |
|:-----------|:--------------------|
| `{`        | Reduce cell width    |
| `}`        | Increase cell width  |
| `/`        | Search for cell matching regular expression |
| `n`        | Find next cell matching search |
| `y`        | Copy cell value |
 | `p`       | Paste cell value |
| `:`        | Enter command |

## Commands

Commands can be entered by pressing `:` and typing in the command or alias.

| Command               | Alias      | Description             |
|:----------------------|:-----------|:------------------------|
| `save`                | `w`        | Save the current file. |
| `quit`                | `q`        | Quit the application without saving changes. |
| `save-and-quit`       | `wq`       | Save the current file and quit the application. |
| `open-down`           |            | Insert a new row below the currently selected row. |
| `open-right`          |            | Insert a new column to the right of the currently selected column. |
| `delete-row`          |            | Delete the currently selected row. |
| `delete-column`       |            | Delete the currently selected column. |

