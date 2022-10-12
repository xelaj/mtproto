# Contributing to TL

**english** [русский][index_ru]

🌚🌚 First of all, thanks for your helping! 🌝🌝

This page briefly describes the process of developing both the specific MTProto
package and all Xelaj projects. if you read all these rules, you will be the
best helper in the wild west!

## Code of conduct

We all want to make other people happy! We believe that you are a good guy, but
please, just in case, read our [code of conduct][CoC]. They will help you
understand what ideals we adhere to, among other things, you will be even more
cool!

By joining our community, you automatically agree to our rules _(even if you
have not read them!)_. and if you saw their violation somewhere, write to
<rcooper.xelaj@protnmail.com>, we will help!

## I don't want to read anything, I have a question!

> **Just remind:** you just don’t need to ask anything right to the issues,
> okay? just not necessary. you will quickly solve your problem if you find the
> answer below

We have the official Xelaj chat in Telegram: [@xelaj_developers][telegram_chat].
In this chat you can promptly clarify the information of interest to you.

Also, github create discussions page! Ask in discussions, that's why this
feature released. It's like stack overflow for **any** repo like this one!

And we also actually want to do [FAQ][gh_discussions_faq], but we don’t know
what questions to write there, so, if you are reading this, probably write while
in the Telegram, we'll figure it out :)

## What do I need to know before I can help?

Please, sign ALL of your commits. It's really really really important! Guide,
how to sign your commits is [here][signing_commits]. Unsigned commits will not
be merged

## And how can I help?

### For example, report a bug.

#### before reporting a bug:

* Look for issues with a bug / bug label, it is likely that it has already been
  reported.
* **even if you found issue**: describe your situation in the comments of issue,
  attach logs, backup database, just do not duplicate issues.

### You can still offer a new feature:

We love to add new features! Use the New Feature issues template and fill in all
the fields. Attaching labels is also very important!

### and you can immediately offer a pull request!

Here it is up to you, the only thing is: we are more willing to take pull
requests based on a specific issue (i.e., created pull request based on issue
\#100500 or something like this) This will help us understand what problem your
request solves.

## Styleguide

### commit comments

We love [Conventional Commits][conventional_commits], so, please, use
recommendations from that page.

#### Emojis in commit messages

One important change of conventional commits in our projects: **use emojis,
please!** They looks awesome, especially in git logs, in github ui, everywhere!
So, unlike `feat: yakedey-yakedey`, please name commits like
`✨ yakedey-yakedey`. Want prefix? `✨ (abcd): ooga booga` is right for you.
**Please add one space after emoji!** Some terminals renders emojis by 2 columns

Note also, that prefix **MUST** only code of any [project in github][gh_project]
(In description, you can see `Code: ABCD`, so use it. it's case insensitive).
This is real important to automatically add tasks to boards. **But it's not
important!** If you just want add scope — you're welcome, but please, don't add
unknown before scope

Don't add semicolon after emoji without scope: `✏️ fixed typo` looks better than
`✏️ : fixed typo`.

Commit messages like `👷 fixed #1337 issue` are NOT allowed. Git repository must
not be dependent by github.

**There is a list of features and emojis right for you:**

<!-- markdownlint-disable MD013 -->
|Feature    |Emoji| Desc |
| :-------: | :-: | ---- |
|`feat`     | ✨  | Added some new feature, which wasn't before |
|`fix`      | 🐛  | Bug fixes |
|`deprecate`| 🗑  | API Deprecation (but with saving backwards compatibility). Deprecation changes **MUST** be splitted to separate pull request. </br> See? Emojis are shorter! |
|`docs`     | 📖  | Documentation improvements |
|`typo`     | ✏️  | Fixed typos, everywhere. In docs, in tests, in file names, in variables. Breaking changes allowed, but note this change in commit body, please |
|`style`    | 💎  | Style improvements of code, without change of logic e.g. formatting of code. If you changed code logic (even internal logic with backward compatibility), please, use 📦 feature |
|`refactor` | 📦  | Code Refactoring. |
|`perf`     | 🚀  | Performance improvements (new benchmarks counts too) |
|`test`     | ✅  | New tests, test fixes, etc. Benchmarks are not included |
|`ci`       | 👷  | Building scripts, CI/CD stuff, etc. Any config changes for build systems, use also this prefix |
|`merge`    | 🔀  | Merge commits. Do not use it manually. |
|`chore`    | 🎫  | Chore stuff (add new dependencies, update links in docs, change year of license header). Adding contributors is also chore. |
|`egg`      | 🥚  | Some tweaks for easter eggs in code. These changes won't be showed in any changelogs |
|`revert`   | 🔙  | Revert any stuff, **MUST NOT BE USED** except in too bad cases! When you use this feature, you MUST add `Reverts: <commit hash>` commit tag </br> _(any pull requests which contains commits with this feature will be automatically closed)_ |
<!-- markdownlint-enable MD013 -->
#### Do's and dont's

<!-- markdownlint-disable MD013 -->
|                              |❌ Don't         | ✅ Do              |
| ---------------------------- | --------------- | ----------------- |
|do not write what commit does |`commit adds ...`|`added support ...`|
| do not write **who** made a commit | `I changed ...` </br>`our team worked for a long time and created ...` | `Added support of ...`|
| write briefly (no more than 60 characters), everything else - in the description after two (2) new lines | `fool, i forgot to delete ...` | `removed ...`|
<!-- markdownlint-enable MD013 -->

### Error types format

Unlike official golang convention

Please, use format `ErrXxx` even for structs. It's important, cause in godoc all
errors are tighten to single place. Also, it's better to see unified view of
comparison for error types and error constants, e.g.

```go
if errors.Is(err, ErrSomething) {
    ...
}

if e := ErrStructed; errors.As(err, &e) {
    ...
}
```

So, please, use this naming. I know, you can hate us, and hate for obvious
reasons, but even in go different teams can have different styleguides.

Most perfect way for structs is to use generic function like this:

```go
func as118[T error](err error) (T, bool) {
    var converted T
    return converted, errors.As(err, &converted)
}

...

if err, ok := as118[ErrStructed](err); ok {
    ...
}
```

And [it's working!][go_err_example]

[conventional_commits]: https://www.conventionalcommits.org/en/v1.0.0/
[signing_commits]:      https://docs.gitlab.com/ee/user/project/repository/gpg_signed_commits/
[go_err_example]:       https://go.dev/play/p/JUW68wmHxwc

<!-- localizations -->
[index_ru]: https://github.com/xelaj/mtproto/blob/-/docs/ru_RU/CONTRIBUTING.md

<!-- project links -->
[telegram_chat]:      https://t.me/xelaj_developers
[gh_project]:         https://github.com/xelaj/mtproto/projects
[gh_discussions_faq]: https://github.com/xelaj/mtproto/discussions/categories/q-a
[CoC]:                https://github.com/xelaj/mtproto/blob/-/.github/CODE_OF_CONDUCT.md
