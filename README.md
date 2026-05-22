# kp

`kp` is a macOS-first utility CLI for prompt paste workflows.

## Install

```bash
brew install jamesonstone/tap/kp
```

## Usage

```bash
kp --version
kp prompt list
kp prompt clarify
kp prompt clarify --copy
kp prompt clarify --print
kp prompt new my-prompt
kp prompt edit clarify
kp prompt rm my-prompt
```

## skhd examples

```bash
cmd + alt - i : kp prompt instructions
cmd + alt - c : kp prompt clarify
```

## Secure Input troubleshooting

If paste fails, find the Secure Input holder:

```bash
ioreg -l -w 0 | grep SecureInputPID
```

## Exit codes

- `0`: success
- `1`: user error
- `2`: clipboard/paste system error
- `3`: config/IO error
- `130`: user cancelled
