#compdef qsoku
# zsh integration and completion for qsoku (docs/reference/cli.md
# "Shell integration" / "Shell completion").
#
# Install with:  eval "$(qsoku .shell zsh)"

qsoku() {
  local f; f=$(mktemp)
  QSOKU_CWD_FILE=$f command qsoku "$@"
  # zsh reserves the name "status" as a synonym for $? and will not let a
  # local variable shadow it, unlike bash; qsoku_status avoids that.
  local qsoku_status=$?
  local dir; dir=$(cat "$f"); rm -f "$f"
  [ -n "$dir" ] && [ "$dir" != "$PWD" ] && cd "$dir"
  return $qsoku_status
}

# Only the first word after "qsoku " is completed: the defined names, and
# the fixed set of management commands (a qsokufile entry can never define
# one, so it never collides -- docs/reference/qsokufile.md "Names"). The
# names come from `qsoku .names` on every Tab press, never from a list
# embedded here, so completion always matches the qsokufile actually in use.
_qsoku() {
  if (( CURRENT != 2 )); then
    return
  fi
  local -a names
  names=(${(f)"$(command qsoku .names 2>/dev/null)"})
  compadd -- $names .init .add .rm .list .names .edit .where .version .shell .help
}

if [ "$funcstack[1]" = "_qsoku" ]; then
  _qsoku "$@"
else
  compdef _qsoku qsoku
fi
