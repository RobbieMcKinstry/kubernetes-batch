. Commit and redact password

. Test out UUID generator

. Write a wrapper for command execution

type Command struct {
    Pre func(*Command) (in string) {}
    Post func(*Command, out string) (ok bool)
    exec.Cmd
}
