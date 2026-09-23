# bash integration and completion for qsoku (docs/reference/cli.md
# "Shell integration" / "Shell completion").
#
# Install with:  eval "$(qsoku .shell bash)"

qsoku() {
  local f; f=$(mktemp)
  QSOKU_CWD_FILE=$f command qsoku "$@"
  local status=$?
  local dir; dir=$(cat "$f"); rm -f "$f"
  # shellcheck disable=SC2164 # a failed cd here would be a rare race (dir
  # was just read from a real pwd); exiting the user's shell over it would
  # be worse than silently not moving
  [ -n "$dir" ] && [ "$dir" != "$PWD" ] && cd "$dir"
  return $status
}

# Only the first word after "qsoku " is completed: the defined names, and
# the fixed set of management commands (a qsokufile entry can never define
# one, so it never collides -- docs/reference/qsokufile.md "Names"). The
# names come from `qsoku .names` on every Tab press, never from a list
# embedded here, so completion always matches the qsokufile actually in use.
_qsoku() {
  COMPREPLY=()
  if [ "$COMP_CWORD" -ne 1 ]; then
    return
  fi
  local names value
  names=$(command qsoku .names 2>/dev/null)
  while IFS= read -r value; do
    COMPREPLY+=("$value")
  done < <(compgen -W "$names .init .add .rm .list .names .edit .where .version .shell .help" -- "${COMP_WORDS[1]}")
}
complete -F _qsoku qsoku
