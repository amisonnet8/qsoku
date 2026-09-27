# PowerShell (7+, "pwsh") integration and completion for qsoku
# (docs/reference/cli.md "Shell integration" / "Shell completion").
#
# Install with (in $PROFILE):  Invoke-Expression (& qsoku .shell pwsh | Out-String)
#
# pwsh is only ever the *caller's* shell here: qsokufile commands themselves
# always run under sh (CLAUDE.md "実行は常にsh"), on every OS including
# Windows (where that means sh.exe from Git for Windows on PATH).

function qsoku {
    $f = [System.IO.Path]::GetTempFileName()
    $real = Get-Command qsoku -CommandType Application | Select-Object -First 1
    try {
        $env:QSOKU_CWD_FILE = $f
        & $real.Source @args
        $qsoku_status = $LASTEXITCODE
    } finally {
        Remove-Item Env:\QSOKU_CWD_FILE -ErrorAction SilentlyContinue
    }

    $dir = $null
    if (Test-Path -LiteralPath $f) {
        $dir = (Get-Content -LiteralPath $f -ErrorAction SilentlyContinue | Select-Object -First 1)
        Remove-Item -LiteralPath $f -ErrorAction SilentlyContinue
    }
    if ($dir -and $dir -ne $PWD.Path) {
        Set-Location -LiteralPath $dir
    }

    $global:LASTEXITCODE = $qsoku_status
}

# Only the first word after "qsoku " is completed: the defined names (from
# `qsoku .names` on every Tab press, never a list embedded here, so
# completion always matches the qsokufile actually in use), and the fixed
# set of management commands (a qsokufile entry can never define one --
# docs/reference/qsokufile.md "Names").
Register-ArgumentCompleter -CommandName qsoku -ScriptBlock {
    param($wordToComplete, $commandAst, $cursorPosition)

    if ($commandAst.CommandElements.Count -gt 2) {
        return
    }

    $names = (& qsoku .names 2>$null) -split "`n" | Where-Object { $_ }
    $candidates = $names + @('.init', '.add', '.rm', '.list', '.names', '.edit', '.where', '.version', '.shell', '.help')
    $candidates | Where-Object { $_ -like "$wordToComplete*" } | ForEach-Object {
        [System.Management.Automation.CompletionResult]::new($_, $_, 'ParameterValue', $_)
    }
}
