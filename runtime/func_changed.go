package runtime

import (
	"fmt"
	"os/exec"
	"otsukai"
	"otsukai/parser"
	"otsukai/runtime/context"
	re "otsukai/runtime/errors"
	"otsukai/runtime/helpers"
	"otsukai/runtime/value"
)

// https://git-scm.com/docs/git-rev-parse#Documentation/git-rev-parse.txt-codeHEADcode
const LAST_COMMIT = "HEAD^1"      // names the commit on which you based the changes in the working tree.
const FETCH_COMMIT = "FETCH_HEAD" // records the branch which you fetched from a remote repository with your last git fetch invocation.
const BEFORE_MERGE = "ORIG_HEAD"  // is created by commands that move your HEAD in a drastic way (git am, git merge, git rebase, git reset), to record the position of the HEAD before their operation, so that you can easily change the tip of the branch back to the state before you ran them.
const AFTER_MERGE = "MERGE_HEAD"  // records the commit(s) which you are merging into your branch when you run git merge.

func getTargetCommit(arguments []parser.Argument) (string, error) {
	commit := helpers.GetNamedArgument(arguments, "commit_from")
	if commit == nil {
		return LAST_COMMIT, nil
	}

	if commit.Expression.ValueExpression == nil || commit.Expression.ValueExpression.Value.HashSymbol == nil {
		return LAST_COMMIT, re.RUNTIME_ERROR
	}

	symbol := commit.Expression.ValueExpression.Value.HashSymbol.Identifier
	switch symbol {
	case "last_commit":
		return LAST_COMMIT, nil

	case "fetch_commit":
		return FETCH_COMMIT, nil

	case "before_merge":
		return BEFORE_MERGE, nil

	case "after_merge":
		return AFTER_MERGE, nil

	default:
		return LAST_COMMIT, re.RUNTIME_ERROR
	}
}

func getTargetFile(arguments []parser.Argument) (*string, error) {
	path := helpers.GetNamedArgument(arguments, "path")
	if path == nil {
		return nil, nil
	}

	if path.Expression.ValueExpression == nil || path.Expression.ValueExpression.Value.Literal == nil {
		return nil, re.RUNTIME_ERROR
	}

	return path.Expression.ValueExpression.Value.Literal.String, nil
}

func hasChangedInGitHistory(path string, revision string) (bool, error) {
	cmd := exec.Command("git", "diff-index", "--quiet", revision, "--", path)
	cmd.Run()

	exit := cmd.ProcessState.ExitCode()
	if exit == 0 {
		return false, nil
	}

	if exit == 1 {
		return true, nil
	}

	return false, re.RUNTIME_ERROR
}

func InvokeChanged(ctx context.IContext, arguments []parser.Argument) (value.IValueObject, error) {
	phase := ctx.GetPhase()
	if phase != PHASE_RUN {
		return nil, re.RUNTIME_ERROR
	}

	scope := ctx.GetContextFlag()
	if scope&context.CONTEXT_TASK != context.CONTEXT_TASK {
		otsukai.Errf("invalid context")
		return nil, re.SYNTAX_ERROR
	}
	if len(arguments) != 2 {
		otsukai.Errf("set method must be one argument")
		return nil, re.SYNTAX_ERROR
	}

	commit, err := getTargetCommit(arguments)
	if err != nil {
		return nil, err
	}

	path, err := getTargetFile(arguments)
	if err != nil {
		return nil, err
	}

	val, err := hasChangedInGitHistory(*path, commit)
	fmt.Println(val)

	return value.BooleanValueObject{Val: val}, nil
}
