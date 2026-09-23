# fish integration and completion for qsoku (docs/reference/cli.md
# "Shell integration" / "Shell completion").
#
# Install with:  qsoku .shell fish | source

function qsoku
    set -l f (mktemp)
    set -lx QSOKU_CWD_FILE $f
    command qsoku $argv
    set -l status_ $status
    set -l dir (cat $f)
    rm -f $f
    if test -n "$dir"; and test "$dir" != "$PWD"
        cd $dir
    end
    return $status_
end

function __qsoku_names
    command qsoku .names 2>/dev/null
end

# Only the first word after "qsoku " is completed: the defined names (from
# `qsoku .names` on every Tab press, never a list embedded here, so
# completion always matches the qsokufile actually in use), and the fixed
# set of management commands (a qsokufile entry can never define one --
# docs/reference/qsokufile.md "Names").
complete -c qsoku -f -n 'test (count (commandline -opc)) -eq 1' -a '(__qsoku_names)'
complete -c qsoku -f -n 'test (count (commandline -opc)) -eq 1' -a '.init .add .rm .list .names .edit .where .version .shell .help'
