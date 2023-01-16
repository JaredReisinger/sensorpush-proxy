# Contributing to sensorpush-proxy

:+1::tada: First off, thank you for taking the time to contribute! :tada::+1:

The following is a set of guidelines and expectations for contributing to sensorpush-proxy. Note that these are _guidelines_, not set-in-stone rules. Use your best judgment, and also feel free to propose changes to this document itself in a code request, or by opening an issue.

#### Table of contents

<!-- TOC depthfrom:2 depthto:3 orderedlist:false updateonsave:true -->

- [Code of conduct](#code-of-conduct)
- [I just have a question!](#i-just-have-a-question)
- [How can I contribute?](#how-can-i-contribute)
    - [Reporting bugs](#reporting-bugs)
    - [Suggesting enhancements](#suggesting-enhancements)
    - [Pull requests](#pull-requests)

<!-- /TOC -->

## Code of conduct

This project and everyone participating in it is goverened by the [code of conduct](./CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code.

## I just have a question!

If you have a question, or found what looks like a bug, please feel free to [open a new GitHub issue](https://github.com/JaredReisinger/sensorpush-proxy/issues?q=is%3Aissue) for it. Even better, check to see if there’s already been an issue that addresses the question first.

## How can I contribute?

### Reporting bugs

If you’ve run across what feels like a bug, first [check to see if there’s already an open issue](https://github.com/JaredReisinger/sensorpush-proxy/issues) about it. If there is, you can add any additional details you think would helpful in tracking down the issue. If there is a **closed** issue that seems to be the same thing that you’re seeing, please create a new issue and reference the old one rather than re-opening the closed issue.

When describing the bug, please be as detailed as possible:

- **Use a clear and descriptive title** to help identify the problem.
- **Include the version number of sensorpush-proxy** that you’re using.
- **Provide a description or exceprt of the code** if you can.
- **Include the error message** if there was one.
- **Include a screenshot or animated GIF** if it’s a rendering issue.

### Suggesting enhancements

Much as with bugs, please first [check to see if there’s already an open issue](https://github.com/JaredReisinger/sensorpush-proxy/issues) that describes the enhancement you’d like to see. If there is, you can add any additional specifics that would help the enhancement meet your use case.

When suggesting a new enhancement, please be as detailed as possible:

- **Use a clear and descriptive title** to help identify the suggestion.
- **Provide a description of the problem you’re trying to solve** rather than simply suggesting the solution outright. In some cases, there might be a different approach altogether that actually fulfills your needs more directly.
- **Explain why this problem might be of general interest** and not an isolated one-off solution to your specific problem. (If you feel it’s more of a one-off situation, you should still feel free to open an issue so we can discuss possible solutions!)
- **Include the version number of sensorpush-proxy** that you’re using.

### Pull requests

I’m a big fan of automating as much as possible; this means that where other projects might have prescriptive style guides you need to follow, most of sensorpush-proxy’s requirements are automatically applied by tooling (see [A note on tooling, below](#a-note-on-tooling)). Between `.editorconfig`, `.prettierrc.yaml`, `.stylelintrc.yaml`, `.eslintrc.yaml` and `.commitlint.yaml`, you are practically _forced_ into using the project’s standard style. Most of the time, your editor (VSCode, Atom, etc.) will simply take care of this for you, and you won’t even need to think about it.

When you first clone/fork the repo, make sure that you have [`task`](https://taskfile.dev) 3.15 or newer (the project uses the `ROOT_DIR` special variable added in 3.15), and then run:

```shell
task prepare
```

to set up all the tooling (e.g. git hooks).

Your commit message should follow the [Conventional Commits](https://www.conventionalcommits.org) standard. There’s tooling for this, too; `git commit` _should_ run an interactive command-line that helps you create the correct format for the commit message if you’ve done the `task prepare` step. (But honestly, if the commit message isn’t in the“standard” form, don’t worry about it too much… I can adjust the PR to ensure that the right things happen.)  You can use `task commit` to force the interactive command to run, regardless.

#### A note on tooling

I have gotten _**very**_ accustomed to Node/npm tooling like `husky`, `commitizen`, `commitlint`, `semantic-release`, etc. This tooling is mature and full-featured and comes effectively “for free” as soon as you `npm install` in a cloned project. The tooling for Go doesn’t seem quite as mature, but I’m attempting to mirror it. One possibility would be to just use a bogus `package.json` and use npm packages, but that seems less “pure” than an all-Go solution. There’s also the language-agnostic-but-really-python `pre-commit` project, but it still suffers from the non-native bootstrapping problem.

The solution I’m using at present is to find as many all-Go equivalents that I can, and—aside from the initial installation of `task`—use a single `task prepare` as a mirror of what `npm install` would do. (But where many npm `devDependencies` would auto-init during `npm install`, I have to manually perform the tooling setup during `task prepare` myself.)

The good news is that I’m able to get to about XX% of the experience I like as a developer. The `task prepare` step sets up a [Go husky workalike](https://github.com/automation-co/husky), and those hooks then leverage the [Taskfile](./Taskfile.yml) to run the pre-commit and commit message hooks. On the semantic-release front, I don’t _really_ need to run that locally—as much as local dry-runs can be helpful—I really only need it to run as a part of CI, and I can do that directly as a GitHub action as long as the release configuration is in the project.
